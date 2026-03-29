package cli

import (
	"fmt"
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
	dryRun bool
	yes    bool
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
		Use:   "commit",
		Short: shared.MessageCommitShort,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCommitCommand(cmd, dependencies, options)
		},
	}

	command.Flags().BoolVar(&options.dryRun, "dry-run", false, shared.MessageDryRunFlag)
	command.Flags().BoolVar(&options.yes, "yes", false, shared.MessageYesFlag)
	return command
}

func runCommitCommand(command *cobra.Command, dependencies commitDependencies, options commitOptions) error {
	service := app.NewCommitService(dependencies.gitRepository, dependencies.aiProvider)
	selectedPaths, err := prepareCommitPaths(command, dependencies, options)
	if err != nil {
		return err
	}

	plans, err := service.PlanCommits(selectedPaths, app.GenerateCommitOptions{
		Scope: dependencies.config.DefaultScope,
	})
	if err != nil {
		return err
	}

	for index, plan := range plans {
		formattedOutput := ui.FormatCommitPlan(index+1, len(plans), plan.Result)
		if _, err := fmt.Fprintf(command.OutOrStdout(), "%s\n", formattedOutput); err != nil {
			return err
		}
	}

	if options.dryRun {
		return nil
	}

	if shouldSkipConfirmation(options, dependencies.config) {
		return createPlannedCommits(command, dependencies.gitRepository, plans, true)
	}

	confirmed, err := ui.ConfirmCommit(command.InOrStdin(), command.OutOrStdout(), shared.MessageCommitPlanQuestion)
	if err != nil {
		return err
	}
	if !confirmed {
		_, err = fmt.Fprintln(command.OutOrStdout(), shared.MessageCommitCanceled)
		return err
	}

	return createPlannedCommits(command, dependencies.gitRepository, plans, false)
}

func shouldSkipConfirmation(options commitOptions, configuration commitConfig) bool {
	return options.yes || configuration.AutoConfirm
}

func createPlannedCommits(command *cobra.Command, gitRepository interfaces.GitRepository, plans []app.CommitPlan, autoApprove bool) error {
	for index, plan := range plans {
		confirmed, err := confirmPlannedCommit(command, index+1, len(plans), autoApprove)
		if err != nil {
			return err
		}
		if !confirmed {
			if _, err := fmt.Fprintf(command.OutOrStdout(), shared.MessageIgnoredBlock+"\n", index+1); err != nil {
				return err
			}
			continue
		}

		if err := gitRepository.CommitPaths(plan.Result.Message, plan.Result.Paths); err != nil {
			return err
		}

		if _, err := fmt.Fprintf(command.OutOrStdout(), shared.MessageCommitCreated+"\n", strings.TrimSpace(plan.Result.Message)); err != nil {
			return err
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
		if _, err := fmt.Fprintln(command.OutOrStdout(), ui.FormatChangedFiles(changedPaths)); err != nil {
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
