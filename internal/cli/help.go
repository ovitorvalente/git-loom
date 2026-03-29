package cli

const rootHelpTemplate = `{{with or .Long .Short}}{{. | trimTrailingWhitespaces}}{{end}}

Uso:
  {{.UseLine}}

{{if .HasAvailableSubCommands}}Comandos:
{{range .Commands}}{{if (and .IsAvailableCommand (not .IsAdditionalHelpTopicCommand))}}  {{rpad .Name .NamePadding }} {{.Short}}
{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}
Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Flags globais:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .Example}}

Exemplos:
{{.Example | trimTrailingWhitespaces}}{{end}}
`

func rootHelpText() string {
	return `Git Loom automatiza commits semanticos com revisao antes de executar.

O CLI hoje e focado no fluxo de commit, com agrupamento de arquivos, analise semantica,
score de qualidade, sugestoes de melhoria e confirmacao interativa.`
}

func rootExamples() string {
	return `  gitloom help
  gitloom help commit
  gitloom analyze --json
  gitloom config init
  gitloom doctor
  gitloom version
  gitloom commit
  gitloom commit --dry-run --preview
  gitloom commit --yes --verbose`
}

func commitHelpText() string {
	return `Planeja e cria commits semanticos a partir do estado atual do repositorio.

O comando:
  - le arquivos staged
  - detecta arquivos em changes e oferece adicionar ao stage
  - agrupa mudancas relacionadas em blocos pequenos
  - gera mensagem semantica em portugues
  - mostra score, detalhes, analise e sugestoes antes de commitar
  - confirma cada bloco, ou executa direto com --yes`
}

func commitExamples() string {
	return `  git add .
  gitloom commit

  gitloom commit --dry-run
  gitloom commit --preview
  gitloom commit --verbose
  gitloom commit --json --dry-run
  gitloom commit --strict
  gitloom commit --yes

Config:
  .gitloom.yaml

  commit:
    scope: core

  cli:
    auto_confirm: false`
}
