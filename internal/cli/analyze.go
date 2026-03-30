package cli

import (
	"github.com/spf13/cobra"

	infraai "github.com/ovitorvalente/git-loom/internal/infra/ai"
	infraconfig "github.com/ovitorvalente/git-loom/internal/infra/config"
	infragit "github.com/ovitorvalente/git-loom/internal/infra/git"
	"github.com/ovitorvalente/git-loom/internal/shared"
)

type analyzeOptions struct {
	review reviewOptions
}

func newAnalyzeCommand() *cobra.Command {
	configuration, err := infraconfig.Load(".gitloom.yaml")
	if err != nil {
		configuration = infraconfig.DefaultConfig()
	}

	return newAnalyzeCommandWithDependencies(commitDependencies{
		gitRepository: infragit.NewRepository(),
		aiProvider:    infraai.NewNoopProvider(),
		config: commitConfig{
			DefaultScope: configuration.Commit.Scope,
			AutoConfirm:  configuration.CLI.AutoConfirm,
		},
	})
}

func newAnalyzeCommandWithDependencies(dependencies commitDependencies) *cobra.Command {
	options := analyzeOptions{}
	command := &cobra.Command{
		Use:           "analyze",
		Aliases:       []string{"plan"},
		Short:         shared.MessageAnalyzeShort,
		Long:          analyzeHelpText(),
		Example:       analyzeExamples(),
		Args:          cobra.NoArgs,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAnalyzeCommand(cmd, dependencies, options)
		},
	}

	addReviewFlags(command, &options.review, true)
	return command
}

func runAnalyzeCommand(command *cobra.Command, dependencies commitDependencies, options analyzeOptions) error {
	execution, err := executeReview(command, dependencies, options.review)
	if err != nil {
		return err
	}

	return printReview(command, execution, options.review)
}

func analyzeHelpText() string {
	return `Analisa o estado atual do repositorio e mostra o plano de commits sem criar commits.

Use este commando para:
  - revisar agrupamento antes de commitar
  - inspecionar score, detalhes e sugestoes
  - exportar a revisao em json para scripts
  - testar o efeito de --optimize antes do fluxo final`
}

func analyzeExamples() string {
	return `  gitloom analyze
  gitloom analyze --preview --verbose
  gitloom analyze --optimize
  gitloom analyze --json`
}
