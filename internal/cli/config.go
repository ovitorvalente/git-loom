package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	infraconfig "github.com/ovitorvalente/git-loom/internal/infra/config"
	"github.com/ovitorvalente/git-loom/internal/shared"
)

type configInitOptions struct {
	force bool
}

func newConfigCommand() *cobra.Command {
	command := &cobra.Command{
		Use:           "config",
		Short:         shared.MessageConfigShort,
		Long:          configHelpText(),
		Example:       configExamples(),
		Args:          cobra.NoArgs,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	command.AddCommand(newConfigInitCommand())
	return command
}

func newConfigInitCommand() *cobra.Command {
	options := configInitOptions{}
	command := &cobra.Command{
		Use:           "init",
		Short:         shared.MessageConfigInitShort,
		Args:          cobra.NoArgs,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigInitCommand(cmd, options)
		},
	}

	command.Flags().BoolVar(&options.force, "force", false, shared.MessageForceFlag)
	return command
}

func runConfigInitCommand(command *cobra.Command, options configInitOptions) error {
	path := filepath.Join(".", ".gitloom.yaml")
	if !options.force {
		if _, err := os.Stat(path); err == nil {
			return errors.New(shared.MessageConfigExists)
		} else if !os.IsNotExist(err) {
			return err
		}
	}

	if err := os.WriteFile(path, []byte(infraconfig.RenderDefaultConfig()), 0o600); err != nil {
		return err
	}

	_, err := fmt.Fprintf(command.OutOrStdout(), shared.MessageConfigCreated+"\n", path)
	return err
}

func configHelpText() string {
	return `Gerencia a configuracao local do gitloom.

Use "config init" para criar um arquivo inicial com as opcoes suportadas hoje.`
}

func configExamples() string {
	return `  gitloom config init
  gitloom config init --force`
}
