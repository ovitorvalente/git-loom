# Git Loom - Project Structure

## Visao geral

Estrutura atual do projeto:

```text
git-loom/
├── .github/
│   └── workflows/
│       ├── release.yml
│       └── tag-release.yml
├── cmd/
│   └── gitloom/
│       └── main.go
├── internal/
│   ├── app/
│   │   ├── commit_feedback.go
│   │   ├── commit_service.go
│   │   └── commit_service_test.go
│   ├── cli/
│   │   ├── analyze.go
│   │   ├── analyze_test.go
│   │   ├── commit.go
│   │   ├── commit_test.go
│   │   ├── config.go
│   │   ├── config_test.go
│   │   ├── doctor.go
│   │   ├── doctor_test.go
│   │   ├── help.go
│   │   ├── help_test.go
│   │   ├── review.go
│   │   ├── root.go
│   │   ├── version.go
│   │   └── version_test.go
│   ├── domain/
│   │   ├── branch/
│   │   ├── commit/
│   │   └── shared/
│   ├── infra/
│   │   ├── ai/
│   │   ├── config/
│   │   ├── git/
│   │   └── system/
│   ├── interfaces/
│   │   ├── mocks/
│   │   ├── ai.go
│   │   ├── git.go
│   │   └── logger.go
│   ├── semantic/
│   │   ├── commit_scorer.go
│   │   ├── context.go
│   │   ├── intent_detector.go
│   │   ├── scope_normalizer.go
│   │   └── semantic_test.go
│   ├── shared/
│   │   └── messages.go
│   └── ui/
│       ├── commit_view.go
│       ├── prompts.go
│       ├── prompts_test.go
│       ├── renderer.go
│       ├── renderer_test.go
│       └── summary_view.go
├── pkg/
│   └── gitloom/
├── .rules/
├── Makefile
├── README.md
├── Branding.md
├── Development.md
└── Project-structure.md
```

## Responsabilidade por camada

### `cmd/gitloom`

Ponto de entrada do binario.

- inicializa o CLI
- delega a execucao para `internal/cli`

### `internal/cli`

Camada de interface com o usuario.

Responsavel por:

- comandos Cobra
- flags e help
- prompts de confirmacao
- serializacao JSON de comandos
- orquestracao do fluxo de review e execucao

Comandos atuais:

- `commit`
- `analyze`
- `doctor`
- `config init`
- `version`

### `internal/app`

Camada de aplicacao.

Responsavel por:

- montar plano de commits
- aplicar sugestoes de otimizacao
- consolidar feedback acionavel

Aqui ficam as regras de orquestracao do caso de uso, sem detalhes de UI.

### `internal/domain`

Camada de dominio puro.

#### `internal/domain/commit`

- classificacao de tipo
- analise de diff
- geracao de mensagem
- regras de descricao e corpo

#### `internal/domain/branch`

- regras de nomenclatura de branch

#### `internal/domain/shared`

- tipos compartilhados do dominio

### `internal/semantic`

Motor semantico complementar.

Responsavel por:

- normalizacao de escopo
- deteccao de intencao
- score de qualidade
- contexto semantico do diff
- agrupamento por tema, area e sinal transversal

Essa camada influencia diretamente a qualidade do planner.

### `internal/ui`

Camada de apresentacao textual.

Responsavel por:

- renderer do review
- secoes visuais de commit
- resumo final
- prompts interativos

Nao deve conter regra de negocio.

### `internal/infra`

Adaptadores para ambiente externo.

#### `internal/infra/git`

- wrapper sobre comandos Git
- leitura de diff
- leitura de staged/changes
- commit por paths
- deteccao de repositório

#### `internal/infra/config`

- leitura de `.gitloom.yaml`
- renderizacao de config padrao

#### `internal/infra/ai`

- provider atual `noop`
- ponto de extensao para provedores reais

#### `internal/infra/system`

- utilitarios de sistema

### `internal/interfaces`

Contratos entre camadas e mocks para testes.

Principal uso:

- `GitRepository`
- `AIProvider`

## Arquivos raiz importantes

### `Makefile`

Targets principais:

- `make test`
- `make vet`
- `make build`
- `make dist VERSION=vX.Y.Z`
- `make clean`

### `README.md`

Documentacao principal do produto.

### `Branding.md`

Diretrizes de identidade e posicionamento.

### `Development.md`

Guia de desenvolvimento, principios e roadmap.

### `.github/workflows`

Automacao de release:

- `tag-release.yml`
- `release.yml`

## Convencoes atuais

- codigo Go dentro de `internal/`
- `cmd/` apenas como entrypoint
- testes ao lado dos arquivos de implementacao
- JSON tratado como contrato publico da CLI
- documentacao operacional na raiz do projeto

## Areas onde o projeto ainda pode crescer

- `pkg/gitloom/`
  Hoje praticamente vazio. Pode virar superficie publica se o motor precisar ser reutilizado fora do binario.

- `configs/`
  Pode concentrar exemplos e presets de configuracao no futuro.

- `scripts/`
  Pode receber scripts de release, verificacoes locais e setup de dev.
