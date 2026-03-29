package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ovitorvalente/git-loom/internal/shared"
)

var (
	Version   = "dev"
	GitCommit = "unknown"
	BuildDate = "unknown"
)

func newVersionCommand() *cobra.Command {
	command := &cobra.Command{
		Use:           "version",
		Aliases:       []string{"ver"},
		Short:         shared.MessageVersionShort,
		Args:          cobra.NoArgs,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := fmt.Fprintf(cmd.OutOrStdout(), "gitloom %s\ncommit: %s\nbuild date: %s\n", Version, GitCommit, BuildDate)
			return err
		},
	}

	return command
}
