package cli

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ovitorvalente/git-loom/internal/app"
	infraai "github.com/ovitorvalente/git-loom/internal/infra/ai"
	infraconfig "github.com/ovitorvalente/git-loom/internal/infra/config"
	infragit "github.com/ovitorvalente/git-loom/internal/infra/git"
	"github.com/ovitorvalente/git-loom/internal/interfaces"
	"github.com/ovitorvalente/git-loom/internal/shared"
	"github.com/ovitorvalente/git-loom/internal/ui"
)

type commitDependencies struct {
	gitRepository interfaces.GitRepository
	aiProvider    interfaces.AIProvider
	config        commitConfig
}

type commitOptions struct {
	dryRun   bool
	yes      bool
	apply    string
	editPlan bool
	review   reviewOptions
}

type commitConfig struct {
	DefaultScope string
	AutoConfirm  bool
}

func newCommitCommand() *cobra.Command {
	configuration, err := infraconfig.Load(".gitloom.yaml")
	if err != nil {
		configuration = infraconfig.DefaultConfig()
	}

	return newCommitCommandWithDependencies(commitDependencies{
		gitRepository: infragit.NewRepository(),
		aiProvider:    infraai.NewNoopProvider(),
		config: commitConfig{
			DefaultScope: configuration.Commit.Scope,
			AutoConfirm:  configuration.CLI.AutoConfirm,
		},
	})
}

func newCommitCommandWithDependencies(dependencies commitDependencies) *cobra.Command {
	options := commitOptions{}
	command := &cobra.Command{
		Use:           "commit",
		Aliases:       []string{"ci"},
		Short:         shared.MessageCommitShort,
		Long:          commitHelpText(),
		Example:       commitExamples(),
		Args:          cobra.NoArgs,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCommitCommand(cmd, dependencies, options)
		},
	}

	command.Flags().BoolVar(&options.dryRun, "dry-run", false, shared.MessageDryRunFlag)
	command.Flags().BoolVar(&options.yes, "yes", false, shared.MessageYesFlag)
	command.Flags().StringVar(&options.apply, "apply", "", "aplica apenas os blocos informados (ex: 1,3-4)")
	command.Flags().BoolVar(&options.editPlan, "edit-plan", false, "abre revisao interativa por bloco (enter/e/s/m/q) antes de confirmar")
	addReviewFlags(command, &options.review, false)
	return command
}

func runCommitCommand(command *cobra.Command, dependencies commitDependencies, options commitOptions) error {
	if options.review.json && !options.dryRun && !shouldSkipConfirmation(options, dependencies.config) {
		return fmt.Errorf("--json no fluxo de commit requer --yes ou auto_confirm; use `gitloom analyze --json` para revisar sem executar")
	}

	service := app.NewCommitService(dependencies.gitRepository, dependencies.aiProvider)
	execution, err := executeReview(command, dependencies, options.review)
	if err != nil {
		return err
	}
	review := execution.review
	renderer := execution.renderer

	if shouldAutoApplySuggestions(options, dependencies.config, review.Suggestions) {
		updatedReview, applyErr := service.ApplySuggestions(review, app.GenerateCommitOptions{
			Scope:             dependencies.config.DefaultScope,
			MaxFilesPerCommit: options.review.maxFiles,
		})
		if applyErr != nil {
			return applyErr
		}
		review = updatedReview
		execution.review = review
	} else if shouldAskToApplySuggestions(options, dependencies.config, review.Suggestions) {
		printSuggestions(command, renderer, review.Suggestions)
		confirmed, confirmErr := ui.ConfirmCommit(command.InOrStdin(), command.OutOrStdout(), shared.MessageApplySuggestions)
		if confirmErr != nil {
			return confirmErr
		}
		if confirmed {
			updatedReview, applyErr := service.ApplySuggestions(review, app.GenerateCommitOptions{
				Scope:             dependencies.config.DefaultScope,
				MaxFilesPerCommit: options.review.maxFiles,
			})
			if applyErr != nil {
				return applyErr
			}
			review = updatedReview
			execution.review = review
		}
	}

	if options.review.strict {
		if strictErr := validateStrictReview(review, options.review.maxFiles); strictErr != nil {
			return strictErr
		}
	}

	selected, parseErr := parseApplySelection(len(review.Plans), options.apply)
	if parseErr != nil {
		return parseErr
	}

	if !options.review.json || options.dryRun || options.review.preview {
		if printErr := printReview(command, execution, options.review); printErr != nil {
			return printErr
		}
	}

	if options.dryRun || options.review.preview {
		return nil
	}

	if options.editPlan && !options.review.json && !shouldSkipConfirmation(options, dependencies.config) && options.apply == "" {
		editedReview, editedSelection, canceled, editErr := reviewPlanActions(command, service, review, options.review)
		if editErr != nil {
			return editErr
		}
		if canceled {
			_, err := fmt.Fprintln(command.OutOrStdout(), shared.MessageCommitCanceled)
			return err
		}
		review = editedReview
		if len(editedSelection) > 0 {
			selected = editedSelection
		}
	}

	if shouldSkipConfirmation(options, dependencies.config) {
		return createPlannedCommits(command, dependencies.gitRepository, renderer, review, options.review.json, true, selected)
	}

	confirmed, err := ui.ConfirmCommit(command.InOrStdin(), command.OutOrStdout(), shared.MessageCommitPlanQuestion)
	if err != nil {
		return err
	}
	if !confirmed {
		_, err = fmt.Fprintln(command.OutOrStdout(), shared.MessageCommitCanceled)
		return err
	}

	return createPlannedCommits(command, dependencies.gitRepository, renderer, review, options.review.json, false, selected)
}

func shouldSkipConfirmation(options commitOptions, configuration commitConfig) bool {
	return options.yes || configuration.AutoConfirm
}

func shouldAutoApplySuggestions(options commitOptions, configuration commitConfig, suggestions []app.CommitSuggestion) bool {
	if !hasAutoApplicableSuggestions(suggestions) {
		return false
	}

	return options.yes || configuration.AutoConfirm
}

func shouldAskToApplySuggestions(options commitOptions, configuration commitConfig, suggestions []app.CommitSuggestion) bool {
	if !hasAutoApplicableSuggestions(suggestions) {
		return false
	}

	return !shouldSkipConfirmation(options, configuration)
}

func createPlannedCommits(command *cobra.Command, gitRepository interfaces.GitRepository, renderer ui.Renderer, review app.CommitReview, asJSON bool, autoApprove bool, selected map[int]bool) error {
	createdCommits := 0
	totalScore := 0
	created := []jsonCreatedCommit{}
	skipped := []jsonSkippedCommit{}
	plans := review.Plans

	if !asJSON && len(plans) > 1 {
		ui.PrintStatus(command.OutOrStdout(), fmt.Sprintf("criando %d commits...", len(plans)))
	}

	for index, plan := range plans {
		if len(selected) > 0 && !selected[index+1] {
			if !asJSON {
				if _, err := fmt.Fprintf(command.OutOrStdout(), shared.MessageIgnoredBlock+"\n", index+1); err != nil {
					return err
				}
			}
			skipped = append(skipped, jsonSkippedCommit{Index: index + 1})
			continue
		}

		confirmed, err := confirmPlannedCommit(command, index+1, len(plans), autoApprove)
		if err != nil {
			return err
		}
		if !confirmed {
			if !asJSON {
				if _, err := fmt.Fprintf(command.OutOrStdout(), shared.MessageIgnoredBlock+"\n", index+1); err != nil {
					return err
				}
			}
			skipped = append(skipped, jsonSkippedCommit{Index: index + 1})
			continue
		}

		if err := gitRepository.CommitPaths(plan.Result.Message, plan.Result.Paths); err != nil {
			return err
		}

		if !asJSON {
			if _, err := fmt.Fprintf(command.OutOrStdout(), shared.MessageCommitCreated+"\n", strings.TrimSpace(plan.Result.Message)); err != nil {
				return err
			}
		}
		createdCommits++
		totalScore += plan.Quality.Score
		created = append(created, jsonCreatedCommit{
			Index:   index + 1,
			Message: strings.TrimSpace(plan.Result.Message),
			Paths:   append([]string(nil), plan.Result.Paths...),
		})
	}

	if createdCommits > 0 {
		summary, err := buildCommitSummary(gitRepository, createdCommits, totalScore)
		if err != nil {
			return err
		}
		if asJSON {
			payload, err := buildJSONCommitExecutionOutput(review, jsonCommitExecutionSummary{
				CreatedCommits: created,
				SkippedCommits: skipped,
				AverageQuality: totalScore / createdCommits,
				Status:         summary.Status,
			})
			if err != nil {
				return err
			}
			if _, err := fmt.Fprintln(command.OutOrStdout(), payload); err != nil {
				return err
			}
		} else {
			if _, err := fmt.Fprintf(command.OutOrStdout(), "\n%s\n", renderer.CommitSummary(summary)); err != nil {
				return err
			}
		}
	} else if asJSON {
		payload, err := buildJSONCommitExecutionOutput(review, jsonCommitExecutionSummary{
			CreatedCommits: created,
			SkippedCommits: skipped,
			Status:         buildWorkingTreeStatus(nil, nil),
		})
		if err != nil {
			return err
		}
		if _, err := fmt.Fprintln(command.OutOrStdout(), payload); err != nil {
			return err
		}
	}

	return nil
}

func hasAutoApplicableSuggestions(suggestions []app.CommitSuggestion) bool {
	for _, suggestion := range suggestions {
		if suggestion.AutoApplicable {
			return true
		}
	}

	return false
}

func printSuggestions(command *cobra.Command, renderer ui.Renderer, suggestions []app.CommitSuggestion) {
	formattedSuggestions := renderer.Suggestions(suggestions)
	if formattedSuggestions == "" {
		return
	}

	_, _ = fmt.Fprintf(command.OutOrStdout(), "%s\n", formattedSuggestions)
}

func validateStrictReview(review app.CommitReview, maxFilesPerCommit int) error {
	if maxFilesPerCommit <= 0 {
		maxFilesPerCommit = 4
	}

	for _, plan := range review.Plans {
		if plan.Quality.Score < 80 {
			return fmt.Errorf("%s: %s", shared.MessageStrictModeFailed, strings.TrimSpace(plan.Result.Message))
		}
		if len(plan.Result.Paths) > maxFilesPerCommit {
			return fmt.Errorf("%s: %s", shared.MessageStrictModeFailed, strings.TrimSpace(plan.Result.Message))
		}
	}

	return nil
}

func prepareCommitPaths(command *cobra.Command, dependencies commitDependencies, options commitOptions) ([]string, error) {
	stagedPaths, err := dependencies.gitRepository.ListStagedFiles()
	if err != nil {
		return nil, err
	}

	changedPaths, err := dependencies.gitRepository.ListChangedFiles()
	if err != nil {
		return nil, err
	}

	if len(changedPaths) > 0 {
		if hasPartiallyStagedFiles(stagedPaths, changedPaths) {
			return nil, errors.New(shared.MessagePartialStage)
		}

		renderer := ui.NewRenderer(ui.RenderOptions{Mode: renderMode(options)})
		if _, err := fmt.Fprintln(command.OutOrStdout(), renderer.ChangedFiles(stagedPaths, changedPaths)); err != nil {
			return nil, err
		}

		confirmed, err := confirmStageChangedFiles(command, dependencies, options)
		if err != nil {
			return nil, err
		}
		if confirmed {
			if err := dependencies.gitRepository.StageFiles(changedPaths); err != nil {
				return nil, err
			}
			stagedPaths = append(stagedPaths, changedPaths...)
		}
	}

	selectedPaths := uniquePaths(stagedPaths)
	if len(selectedPaths) == 0 {
		return nil, app.ErrEmptyDiff
	}

	return selectedPaths, nil
}

func buildCommitSummary(gitRepository interfaces.GitRepository, createdCommits int, totalScore int) (ui.CommitSummary, error) {
	stagedPaths, err := gitRepository.ListStagedFiles()
	if err != nil {
		return ui.CommitSummary{}, err
	}

	changedPaths, err := gitRepository.ListChangedFiles()
	if err != nil {
		return ui.CommitSummary{}, err
	}

	return ui.CommitSummary{
		Created:        createdCommits,
		AverageQuality: totalScore / createdCommits,
		Status:         buildWorkingTreeStatus(stagedPaths, changedPaths),
	}, nil
}

func buildWorkingTreeStatus(stagedPaths []string, changedPaths []string) string {
	switch {
	case len(stagedPaths) == 0 && len(changedPaths) == 0:
		return "working tree limpa"
	case len(stagedPaths) > 0 && len(changedPaths) > 0:
		return "restam mudancas staged e unstaged"
	case len(stagedPaths) > 0:
		return "restam mudancas staged"
	default:
		return "restam mudancas unstaged"
	}
}

func renderMode(options commitOptions) ui.RenderMode {
	if options.review.verbose {
		return ui.RenderModeVerbose
	}

	return ui.RenderModeClean
}

func hasPartiallyStagedFiles(stagedPaths []string, changedPaths []string) bool {
	for _, changedPath := range changedPaths {
		if slices.Contains(stagedPaths, changedPath) {
			return true
		}
	}

	return false
}

func confirmStageChangedFiles(command *cobra.Command, dependencies commitDependencies, options commitOptions) (bool, error) {
	if shouldSkipConfirmation(options, dependencies.config) {
		return true, nil
	}

	return ui.ConfirmCommit(command.InOrStdin(), command.OutOrStdout(), shared.MessageStageChangedQuestion)
}

func confirmPlannedCommit(command *cobra.Command, index int, total int, autoApprove bool) (bool, error) {
	if autoApprove || total == 1 {
		return true, nil
	}

	return ui.ConfirmCommit(
		command.InOrStdin(),
		command.OutOrStdout(),
		fmt.Sprintf(shared.MessageCreateBlockQuestion, index, total),
	)
}

func uniquePaths(paths []string) []string {
	seenPaths := map[string]bool{}
	result := make([]string, 0, len(paths))

	for _, path := range paths {
		if path == "" || seenPaths[path] {
			continue
		}

		seenPaths[path] = true
		result = append(result, path)
	}

	return result
}

func parseApplySelection(total int, raw string) (map[int]bool, error) {
	selection := map[int]bool{}
	normalized := strings.TrimSpace(raw)
	if normalized == "" {
		return selection, nil
	}

	for _, part := range strings.Split(normalized, ",") {
		token := strings.TrimSpace(part)
		if token == "" {
			continue
		}

		if strings.Contains(token, "-") {
			rangeParts := strings.SplitN(token, "-", 2)
			start, err := strconv.Atoi(strings.TrimSpace(rangeParts[0]))
			if err != nil {
				return nil, fmt.Errorf("--apply invalido: %q", token)
			}
			end, err := strconv.Atoi(strings.TrimSpace(rangeParts[1]))
			if err != nil {
				return nil, fmt.Errorf("--apply invalido: %q", token)
			}
			if start <= 0 || end <= 0 || start > end || end > total {
				return nil, fmt.Errorf("--apply fora do intervalo: %q", token)
			}
			for index := start; index <= end; index++ {
				selection[index] = true
			}
			continue
		}

		index, err := strconv.Atoi(token)
		if err != nil {
			return nil, fmt.Errorf("--apply invalido: %q", token)
		}
		if index <= 0 || index > total {
			return nil, fmt.Errorf("--apply fora do intervalo: %q", token)
		}
		selection[index] = true
	}

	return selection, nil
}

func reviewPlanActions(command *cobra.Command, service app.CommitService, review app.CommitReview, options reviewOptions) (app.CommitReview, map[int]bool, bool, error) {
	if len(review.Plans) == 0 {
		return review, nil, false, nil
	}

	selected := map[int]bool{}
	plans := append([]app.CommitPlan(nil), review.Plans...)
	for index := 0; index < len(plans); index++ {
		plan := plans[index]
		action, err := ui.AskInput(command.InOrStdin(), command.OutOrStdout(), fmt.Sprintf("bloco %d/%d [enter=ok|e=editar|s=pular|m=mesclar|q=cancelar]", index+1, len(plans)))
		if err != nil {
			return review, nil, false, err
		}

		switch strings.ToLower(strings.TrimSpace(action)) {
		case "", "ok", "y", "yes":
			selected[index+1] = true
		case "s", "skip":
			continue
		case "q", "cancel":
			return review, nil, true, nil
		case "e", "edit":
			newMessage, askErr := ui.AskInput(command.InOrStdin(), command.OutOrStdout(), "nova mensagem do commit")
			if askErr != nil {
				return review, nil, false, askErr
			}
			if strings.TrimSpace(newMessage) != "" {
				plan.Result.Message = strings.TrimSpace(newMessage)
				plans[index] = plan
			}
			selected[index+1] = true
		case "m", "merge":
			if index+1 >= len(plans) {
				return review, nil, false, fmt.Errorf("nao existe bloco seguinte para mesclar")
			}
			mergedPaths := append([]string(nil), plan.Result.Paths...)
			mergedPaths = append(mergedPaths, plans[index+1].Result.Paths...)
			mergedPlan, mergeErr := service.BuildPlanForPaths(mergedPaths, app.GenerateCommitOptions{
				Scope:             plan.Result.Commit.Scope,
				MaxFilesPerCommit: options.maxFiles,
			})
			if mergeErr != nil {
				return review, nil, false, mergeErr
			}
			plans[index] = mergedPlan
			plans = append(plans[:index+1], plans[index+2:]...)
			selected[index+1] = true
		default:
			return review, nil, false, fmt.Errorf("acao invalida: %q", action)
		}
	}

	review.Plans = plans
	return review, selected, false, nil
}
