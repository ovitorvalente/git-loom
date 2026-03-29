# Git Loom - Development

## Visao atual

`gitloom` e um CLI em Go focado em revisar e criar commits semanticos a partir do estado do repositório.

A proposta central do produto e simples:

- revisar antes de executar
- agrupar melhor
- descrever melhor
- automatizar com seguranca

O projeto hoje ja possui:

- comando `commit`
- comando `analyze`
- comando `doctor`
- comando `config init`
- comando `version`
- renderer textual com modos `clean` e `verbose`
- saida `--json` para revisao e execucao
- configuracao local por `.gitloom.yaml`
- build multiplataforma
- release automatizada por GitHub Actions

## Objetivo tecnico

O objetivo nao e apenas gerar mensagens de commit. O objetivo e construir um motor de planejamento de commits:

- detectar grupos coerentes de mudancas
- produzir descricoes melhores
- alertar quando o plano estiver fraco
- permitir automacao segura

## Principios de desenvolvimento

- manter o produto orientado a decisao, nao a "magia"
- manter dominio e UI separados
- nao mover heuristica semantica para o renderer
- preferir saidas deterministicas quando nao houver provider de IA
- manter comandos pequenos e testaveis
- tratar JSON como contrato publico

## Estado atual por area

### CLI

- Cobra como base
- help forte e orientado ao produto
- prompts interativos com `Enter = yes`
- fluxo de commit com revisao antes da execucao
- comando `doctor` para validacao operacional

### App

- `CommitService` coordena planejamento, sugestoes e otimizacao
- agrupamento por contexto semantico
- reconhecimento de temas transversais como `json`, `help`, `doctor`, `config` e `renderer`

### UI

- renderer desacoplado
- modo padrao focado em escaneabilidade
- modo verbose para contexto tecnico
- resumo final com status do working tree

### Infra

- wrapper simples sobre Git
- config local leve
- release e build multiplataforma
- GitHub Actions para tag e release

## Roadmap recomendado

### Curto prazo

- adicionar suporte real a provider de IA configuravel
- melhorar parser de config para um schema mais robusto
- adicionar CI separado para `push` e `pull_request`
- adicionar instalacao oficial para canais como Homebrew e Winget

### Medio prazo

- suporte a monorepo
- presets de estrategia de agrupamento
- exportacao de relatorio de analise
- changelog e release notes a partir do historico de commits

### Longo prazo

- integracao com GitHub/GitLab issues
- adaptadores para frontend/web app
- configuracao compartilhada por time

## Comandos principais

### Revisao segura

```bash
gitloom analyze
gitloom analyze --preview --verbose
gitloom analyze --json
gitloom analyze --optimize
```

### Execucao

```bash
gitloom commit
gitloom commit --dry-run
gitloom commit --yes
gitloom commit --json --yes
```

### Diagnostico

```bash
gitloom doctor
gitloom doctor --json
```

### Configuracao e distribuicao

```bash
gitloom config init
gitloom version
make dist VERSION=v0.1.0
```

## Build e release

### Desenvolvimento local

```bash
make test
make vet
make build
```

### Distribuicao

```bash
make dist VERSION=v0.1.0
```

Gera binarios para:

- macOS amd64 e arm64
- Linux amd64 e arm64
- Windows amd64 e arm64

### GitHub Release

- `tag-release.yml` cria a tag manualmente
- `release.yml` roda em push de tag `v*`
- testes sao executados antes da publicacao
- artefatos e checksums sao anexados na release

## Estrategia de testes

O projeto hoje depende principalmente de:

- testes unitarios por pacote
- testes da CLI com mocks
- testes da infra Git com repositório real temporario

Areas que exigem alto cuidado ao evoluir:

- `internal/app`
- `internal/cli`
- `internal/semantic`
- `internal/infra/git`

## Padrao de contribuicao tecnica

- manter mensagens e exemplos em portugues quando forem parte do produto
- manter nomes de codigo e estrutura em ingles quando fizer sentido tecnico
- adicionar testes junto com mudancas de comportamento
- nao acoplar renderer a regras de negocio
- preservar estabilidade do formato JSON
- manter consistencia entre README, branding e help do CLI

## Riscos principais atuais

- parser de config ainda e simples
- heuristicas semanticas podem exigir calibracao continua
- suporte a arquivos parcialmente staged ainda e bloqueado
- pipeline de release depende do ambiente GitHub, nao so do Makefile local

## Criterio de evolucao

Uma mudanca so deve entrar se melhorar pelo menos um destes pontos:

- qualidade do agrupamento
- clareza do review
- seguranca da automacao
- portabilidade do binario
- capacidade de integracao externa
- legibilidade e confianca do produto como ferramenta real
