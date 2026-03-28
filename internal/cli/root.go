package cli

import "github.com/spf13/cobra"

func newRootCommand() *cobra.Command {
	return &cobra.Command{
		Use:           "gitloom",
		Short:         "Automate Git workflows",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
}

func Execute() error {
	return newRootCommand().Execute()
}
