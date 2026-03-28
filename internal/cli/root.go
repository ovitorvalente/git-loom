package cli

import "github.com/spf13/cobra"

func newRootCommand() *cobra.Command {
	command := &cobra.Command{
		Use:           "gitloom",
		Short:         "Automatiza fluxos de trabalho com Git",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	command.AddCommand(newCommitCommand())
	return command
}

func Execute() error {
	return newRootCommand().Execute()
}
