package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ovitorvalente/git-loom/internal/app"
	infraai "github.com/ovitorvalente/git-loom/internal/infra/ai"
	infragit "github.com/ovitorvalente/git-loom/internal/infra/git"
	"github.com/ovitorvalente/git-loom/internal/interfaces"
)

type commitDependencies struct {
	gitRepository interfaces.GitRepository
	aiProvider    interfaces.AIProvider
}

type commitOptions struct {
	dryRun bool
	yes    bool
}

func newCommitCommand() *cobra.Command {
	return newCommitCommandWithDependencies(commitDependencies{
		gitRepository: infragit.NewRepository(),
		aiProvider:    infraai.NewNoopProvider(),
	})
}

func newCommitCommandWithDependencies(dependencies commitDependencies) *cobra.Command {
	options := commitOptions{}
	command := &cobra.Command{
		Use:   "commit",
		Short: "Generate a commit message from staged changes",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCommitCommand(cmd, dependencies, options)
		},
	}

	command.Flags().BoolVar(&options.dryRun, "dry-run", false, "show the generated message without committing")
	command.Flags().BoolVar(&options.yes, "yes", false, "create the commit without confirmation")
	return command
}

func runCommitCommand(command *cobra.Command, dependencies commitDependencies, options commitOptions) error {
	service := app.NewCommitService(dependencies.gitRepository, dependencies.aiProvider)
	result, err := service.GenerateCommit()
	if err != nil {
		return err
	}

	if _, err := fmt.Fprintf(command.OutOrStdout(), "commit message: %s\n", result.Message); err != nil {
		return err
	}

	if options.dryRun || !options.yes {
		return nil
	}

	if err := dependencies.gitRepository.Commit(result.Message); err != nil {
		return err
	}

	_, err = fmt.Fprintf(command.OutOrStdout(), "commit created\n")
	return err
}
