package cli

import (
	"github.com/spf13/cobra"

	"github.com/ovitorvalente/git-loom/internal/shared"
)

func newRootCommand() *cobra.Command {
	command := &cobra.Command{
		Use:           "gitloom",
		Short:         shared.MessageRootShort,
		SilenceUsage:  true,
		SilenceErrors: true,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}

	command.AddCommand(newCommitCommand())
	return command
}

func Execute() error {
	return newRootCommand().Execute()
}
