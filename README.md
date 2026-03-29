# Git Loom

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square)](https://golang.org)
[![Latest Release](https://img.shields.io/github/v/release/ovitorvalente/git-loom?style=flat-square)](https://github.com/ovitorvalente/git-loom/releases)
[![License](https://img.shields.io/badge/License-MIT-blue.svg?style=flat-square)](LICENSE)

`gitloom` e um CLI em Go para revisar, agrupar e criar commits semanticos com feedback acionavel antes da execucao.

Ele nao tenta esconder o Git. Ele adiciona uma camada de planejamento entre o estado do repositório e a decisao final de commitar.

Em vez de agir como um gerador cego de mensagem, o `gitloom` ajuda voce a decidir melhor antes de executar.

## Quick Start

```bash
go install github.com/ovitorvalente/git-loom/cmd/gitloom@latest
git add .
gitloom doctor
gitloom analyze
gitloom commit
```

Se quiser automacao:

```bash
gitloom analyze --json
gitloom commit --yes --json
```

## O que o projeto faz hoje

- agrupa mudancas em blocos pequenos
- gera mensagens em portugues seguindo Conventional Commits
- calcula score de qualidade por commit planejado
- mostra alertas, detalhes e sugestoes antes de commitar
- suporta revisao sem execucao com `analyze`
- valida o estado do ambiente com `doctor`
- expoe saida JSON para automacao
- suporta configuracao local com `.gitloom.yaml`
- gera binarios para Linux, macOS e Windows

## Comandos principais

### `gitloom commit`

Planeja e cria commits.

Exemplos:

```bash
gitloom commit
gitloom commit --dry-run
gitloom commit --preview
gitloom commit --verbose
gitloom commit --strict
gitloom commit --yes
gitloom commit --json --yes
```

### `gitloom analyze`

Revisa o plano sem criar commits.

```bash
gitloom analyze
gitloom analyze --preview --verbose
gitloom analyze --optimize
gitloom analyze --json
```

### `gitloom doctor`

Valida se o repositório e o ambiente estao prontos para o fluxo.

```bash
gitloom doctor
gitloom doctor --json
```

### Comandos auxiliares

```bash
gitloom config init
gitloom version
gitloom help
gitloom help commit
```

## Filosofia do produto

`gitloom` foi desenhado com quatro principios:

- clareza antes de automacao
- controle antes de "magia"
- saida escaneavel para humanos
- JSON estavel para integracoes

## Exemplo de fluxo

```bash
git add internal/cli/commit.go
gitloom commit
```

Saida esperada:

```text
────────────────────────────────────────────────────────────
◆ feat(cli) [92] excelente

mensagem:
feat(cli): adicionar fluxo de commit

detalhes:
• adiciona comando commit em cli

arquivos:
+12 -3 internal/cli/commit.go

analise:
ok sem alertas relevantes

> criar commits planejados? [Y/n]:
```

## Instalacao

### Go

```bash
go install github.com/ovitorvalente/git-loom/cmd/gitloom@latest
```

### Requisitos

- Go 1.21+
- Git instalado e disponivel no `PATH`
- acesso a um repositorio Git local

### Build local

```bash
git clone https://github.com/ovitorvalente/git-loom.git
cd git-loom
make build
```

O binario local e gerado como:

```bash
./gitloom
```

## Configuracao

Arquivo opcional `.gitloom.yaml`:

```yaml
commit:
  scope: core

cli:
  auto_confirm: false
```

Para gerar o arquivo automaticamente:

```bash
gitloom config init
```

## JSON para automacao

### Revisao

```bash
gitloom analyze --json
gitloom commit --dry-run --json
```

### Execucao

```bash
gitloom commit --yes --json
```

O modo JSON e pensado para CI, scripts e integracoes externas.

## Limitacoes conhecidas

- arquivos parcialmente staged ainda bloqueiam o fluxo automatico
- o parser de configuracao ainda e propositalmente simples
- o provider de IA atual e `noop`; a geracao principal continua heuristica e deterministica

## Build multiplataforma

Para gerar binarios para Linux, macOS e Windows:

```bash
make dist VERSION=v0.1.0
```

O comando gera:

- `.tar.gz` para Linux e macOS
- `.zip` para Windows
- `checksums.txt` com SHA-256

## Plataformas suportadas

Build oficial para:

- macOS amd64
- macOS arm64
- Linux amd64
- Linux arm64
- Windows amd64
- Windows arm64

Observacao:

- o projeto e distribuido como CLI em Go com `CGO_ENABLED=0`
- o comportamento depende da disponibilidade do Git no ambiente do usuario

## Releases GitHub

O projeto suporta dois workflows:

- `.github/workflows/tag-release.yml`
  cria uma tag via `workflow_dispatch`
- `.github/workflows/release.yml`
  roda em `push` de tags `v*`, executa testes, gera artefatos multiplataforma e publica a release

Fluxo recomendado:

1. abra o workflow `tag-release`
2. informe uma versao como `v0.1.0`
3. a action cria a tag
4. o push da tag dispara o workflow `release`
5. a release e publicada com binarios e checksums

## Desenvolvimento

```bash
make test
make vet
make build
make dist VERSION=v0.1.0
```

Fluxo recomendado para desenvolvimento local:

```bash
git clone https://github.com/ovitorvalente/git-loom.git
cd git-loom
go test ./...
make build
./gitloom help
```

## Testes

O projeto possui:

- testes unitarios por camada
- testes da CLI com mocks
- testes de integracao do wrapper Git com repositório temporario real

## Contribuicao

Contribuicoes sao bem-vindas.

Antes de abrir PR:

1. rode `go test ./...`
2. rode `make vet`
3. valide manualmente o fluxo que voce alterou
4. mantenha o README e o help do CLI alinhados quando mudar comportamento publico

Guia rapido:

- abra uma issue para bugs ou propostas maiores
- prefira PRs pequenas e focadas
- inclua testes para mudancas de comportamento
- preserve a estabilidade do formato JSON

Veja tambem [CONTRIBUTING.md](CONTRIBUTING.md).

Documentos de governanca do projeto:

- [CONTRIBUTING.md](CONTRIBUTING.md)
- [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md)
- [SECURITY.md](SECURITY.md)

## Como reportar problemas

Ao abrir uma issue, inclua:

- sistema operacional
- versao do Go
- versao do `gitloom`
- comando executado
- saida observada
- estado do repositório quando relevante

Se o problema for de planejamento de commits, tente anexar:

- lista de arquivos staged
- lista de arquivos em changes
- saida de `gitloom doctor`
- saida de `gitloom analyze --json`

## Roadmap util para contribuidores

Areas com maior potencial de contribuicao:

- provider de IA configuravel
- parser de configuracao mais robusto
- suporte a instalacao por Homebrew e Winget
- CI para `push` e `pull_request`
- melhor agrupamento para monorepo

## GitHub repository metadata

Descricao curta recomendada para o repositório:

`CLI em Go para revisar e criar commits semanticos com agrupamento inteligente, score de qualidade e saida JSON.`

Topics sugeridos:

- `git`
- `cli`
- `golang`
- `conventional-commits`
- `developer-tools`
- `automation`
- `productivity`
- `semantic-commits`

## Documentacao complementar

- [Branding.md](Branding.md)
- [CONTRIBUTING.md](CONTRIBUTING.md)
- [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md)
- [Development.md](Development.md)
- [Project-structure.md](Project-structure.md)
- [SECURITY.md](SECURITY.md)
- [frontend.md](frontend.md)

## Licenca

MIT. Veja [LICENSE](LICENSE).
