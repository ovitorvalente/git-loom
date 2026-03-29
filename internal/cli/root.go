package cli

import (
	"github.com/spf13/cobra"

	"github.com/ovitorvalente/git-loom/internal/shared"
)

func newRootCommand() *cobra.Command {
	command := &cobra.Command{
		Use:           "gitloom",
		Short:         shared.MessageRootShort,
		Long:          rootHelpText(),
		Example:       rootExamples(),
		Aliases:       []string{"gl", "loom"},
		SilenceUsage:  true,
		SilenceErrors: true,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}
	command.SetHelpTemplate(rootHelpTemplate)

	command.AddCommand(newCommitCommand())
	command.AddCommand(newAnalyzeCommand())
	command.AddCommand(newConfigCommand())
	command.AddCommand(newDoctorCommand())
	command.AddCommand(newVersionCommand())
	return command
}

func Execute() error {
	return newRootCommand().Execute()
}
