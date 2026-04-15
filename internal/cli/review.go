package cli

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ovitorvalente/git-loom/internal/app"
	"github.com/ovitorvalente/git-loom/internal/shared"
	"github.com/ovitorvalente/git-loom/internal/ui"
)

type reviewOptions struct {
	preview  bool
	strict   bool
	verbose  bool
	json     bool
	optimize bool
	explain  bool
	focus    string
	maxFiles int
}

type reviewExecution struct {
	renderer ui.Renderer
	review   app.CommitReview
}

type jsonReviewOutput struct {
	Plans       []jsonPlanOutput         `json:"plans"`
	Summary     *jsonReviewSummaryOutput `json:"summary,omitempty"`
	Suggestions []app.CommitSuggestion   `json:"suggestions,omitempty"`
}

type jsonCommitExecutionOutput struct {
	Summary     *jsonReviewSummaryOutput   `json:"summary,omitempty"`
	Plans       []jsonPlanOutput           `json:"plans"`
	Suggestions []app.CommitSuggestion     `json:"suggestions,omitempty"`
	Execution   jsonCommitExecutionSummary `json:"execution"`
}

type jsonPlanOutput struct {
	Message     string             `json:"message"`
	Type        string             `json:"type"`
	Scope       string             `json:"scope,omitempty"`
	Description string             `json:"description,omitempty"`
	Feedback    app.CommitFeedback `json:"feedback"`
	Paths       []string           `json:"paths"`
	Quality     app.CommitQuality  `json:"quality"`
	Preview     app.CommitPreview  `json:"preview"`
	Index       int                `json:"index"`
	Total       int                `json:"total"`
}

type jsonReviewSummaryOutput struct {
	PlannedCommits  int `json:"planned_commits"`
	AverageQuality  int `json:"average_quality"`
	ChangedFiles    int `json:"changed_files"`
	SuggestionCount int `json:"suggestion_count"`
}

type jsonCommitExecutionSummary struct {
	Status         string              `json:"status"`
	CreatedCommits []jsonCreatedCommit `json:"created_commits,omitempty"`
	SkippedCommits []jsonSkippedCommit `json:"skipped_commits,omitempty"`
	AverageQuality int                 `json:"average_quality"`
}

type jsonCreatedCommit struct {
	Message string   `json:"message"`
	Paths   []string `json:"paths"`
	Index   int      `json:"index"`
}

type jsonSkippedCommit struct {
	Index int `json:"index"`
}

func executeReview(command *cobra.Command, dependencies commitDependencies, options reviewOptions) (reviewExecution, error) {
	if options.maxFiles <= 0 {
		return reviewExecution{}, fmt.Errorf("--max-files-per-commit deve ser maior que zero")
	}

	service := app.NewCommitService(dependencies.gitRepository, dependencies.aiProvider)
	renderer := ui.NewRenderer(ui.RenderOptions{
		Mode:        renderModeFromReview(options),
		ShowPreview: options.preview,
		ShowExplain: options.explain,
	})

	selectedPaths, err := prepareCommitPaths(command, dependencies, commitOptions{
		yes: false,
		review: reviewOptions{
			preview: options.preview,
			verbose: options.verbose,
		},
	})
	if err != nil {
		return reviewExecution{}, err
	}
	selectedPaths = filterPathsByFocus(selectedPaths, options.focus)
	if len(selectedPaths) == 0 {
		return reviewExecution{}, app.ErrEmptyDiff
	}

	review, err := service.PlanCommits(selectedPaths, app.GenerateCommitOptions{
		Scope:             dependencies.config.DefaultScope,
		MaxFilesPerCommit: options.maxFiles,
	})
	if err != nil {
		return reviewExecution{}, err
	}

	if options.optimize {
		review, err = service.ApplySuggestions(review, app.GenerateCommitOptions{
			Scope:             dependencies.config.DefaultScope,
			MaxFilesPerCommit: options.maxFiles,
		})
		if err != nil {
			return reviewExecution{}, err
		}
	}

	if options.strict {
		if err := validateStrictReview(review, options.maxFiles); err != nil {
			return reviewExecution{}, err
		}
	}

	return reviewExecution{
		renderer: renderer,
		review:   review,
	}, nil
}

func printReview(command *cobra.Command, execution reviewExecution, options reviewOptions) error {
	if options.json {
		payload, err := buildJSONReviewOutput(execution.review)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(command.OutOrStdout(), payload)
		return err
	}

	for index, plan := range execution.review.Plans {
		if _, err := fmt.Fprintf(command.OutOrStdout(), "%s\n", execution.renderer.CommitPlan(index+1, len(execution.review.Plans), plan)); err != nil {
			return err
		}
	}

	if len(execution.review.Suggestions) > 0 {
		printSuggestions(command, execution.renderer, execution.review.Suggestions)
	}

	if len(execution.review.Plans) > 1 {
		preview := execution.renderer.FinalPreview(execution.review.Plans)
		if preview != "" {
			fmt.Fprintf(command.OutOrStdout(), "\n%s\n", preview)
		}
	}

	return nil
}

func buildJSONReviewOutput(review app.CommitReview) (string, error) {
	output := jsonReviewOutput{
		Plans:       buildJSONPlanOutputs(review),
		Suggestions: review.Suggestions,
		Summary:     buildJSONReviewSummary(review),
	}

	formatted, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return "", err
	}

	return string(formatted), nil
}

func buildJSONCommitExecutionOutput(review app.CommitReview, execution jsonCommitExecutionSummary) (string, error) {
	output := jsonCommitExecutionOutput{
		Plans:       buildJSONPlanOutputs(review),
		Suggestions: review.Suggestions,
		Summary:     buildJSONReviewSummary(review),
		Execution:   execution,
	}

	formatted, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return "", err
	}

	return string(formatted), nil
}

func buildJSONPlanOutputs(review app.CommitReview) []jsonPlanOutput {
	plans := make([]jsonPlanOutput, 0, len(review.Plans))
	for index, plan := range review.Plans {
		plans = append(plans, jsonPlanOutput{
			Index:       index + 1,
			Total:       len(review.Plans),
			Message:     strings.TrimSpace(plan.Result.Message),
			Type:        string(plan.Result.Commit.Type),
			Scope:       strings.TrimSpace(plan.Result.Commit.Scope),
			Description: strings.TrimSpace(plan.Result.Commit.Description),
			Paths:       append([]string(nil), plan.Result.Paths...),
			Preview:     plan.Preview,
			Quality:     plan.Quality,
			Feedback:    app.BuildCommitFeedback(plan),
		})
	}
	return plans
}

func buildJSONReviewSummary(review app.CommitReview) *jsonReviewSummaryOutput {
	if len(review.Plans) == 0 {
		return nil
	}

	totalScore := 0
	totalFiles := 0
	for _, plan := range review.Plans {
		totalScore += plan.Quality.Score
		totalFiles += len(plan.Result.Paths)
	}

	return &jsonReviewSummaryOutput{
		PlannedCommits:  len(review.Plans),
		AverageQuality:  totalScore / len(review.Plans),
		ChangedFiles:    totalFiles,
		SuggestionCount: len(review.Suggestions),
	}
}

func renderModeFromReview(options reviewOptions) ui.RenderMode {
	if options.verbose {
		return ui.RenderModeVerbose
	}

	return ui.RenderModeClean
}

func addReviewFlags(command *cobra.Command, options *reviewOptions, includeOptimize bool) {
	command.Flags().BoolVar(&options.preview, "preview", false, shared.MessagePreviewFlag)
	command.Flags().BoolVar(&options.strict, "strict", false, shared.MessageStrictFlag)
	command.Flags().BoolVar(&options.verbose, "verbose", false, shared.MessageVerboseFlag)
	command.Flags().BoolVar(&options.json, "json", false, shared.MessageJSONFlag)
	command.Flags().BoolVar(&options.explain, "explain", false, "explica as razoes de agrupamento e score de cada commit planejado")
	command.Flags().StringVar(&options.focus, "focus", "", "filtra o planejamento para caminhos/areas que contenham o termo informado")
	command.Flags().IntVar(&options.maxFiles, "max-files-per-commit", 4, "define limite maximo de arquivos por commit planejado")
	if includeOptimize {
		command.Flags().BoolVar(&options.optimize, "optimize", false, shared.MessageOptimizeFlag)
	}
}

func filterPathsByFocus(paths []string, focus string) []string {
	normalizedFocus := strings.ToLower(strings.TrimSpace(focus))
	if normalizedFocus == "" {
		return append([]string(nil), paths...)
	}

	filteredPaths := make([]string, 0, len(paths))
	for _, path := range paths {
		if strings.Contains(strings.ToLower(path), normalizedFocus) {
			filteredPaths = append(filteredPaths, path)
		}
	}

	return filteredPaths
}
