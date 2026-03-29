# Contributing

Obrigado por contribuir com o `gitloom`.

## Antes de comecar

Requisitos locais:

- Go 1.21+
- Git instalado
- ambiente com suporte a execucao de `go test ./...`

Clone e valide a base:

```bash
git clone https://github.com/ovitorvalente/git-loom.git
cd git-loom
go test ./...
make vet
make build
```

## Principios de contribuicao

- mantenha o produto orientado a decisao, nao a "magia"
- preserve a separacao entre dominio, app, UI e infra
- nao mova heuristica semantica para o renderer
- trate JSON como contrato publico
- adicione testes para toda mudanca de comportamento

## Fluxo recomendado

1. abra uma issue para bugs ou mudancas maiores
2. crie uma branch pequena e focada
3. implemente a mudanca com testes
4. atualize docs e help quando mudar comportamento publico
5. abra o PR com contexto claro

## Checklist de PR

Antes de abrir o PR:

```bash
go test ./...
make vet
```

Valide tambem manualmente quando aplicavel:

- `gitloom help`
- `gitloom commit --dry-run`
- `gitloom analyze`
- `gitloom doctor`

## Mudancas que exigem cuidado extra

- contratos JSON
- heuristicas de agrupamento
- textos e UX do CLI
- integracao com Git
- release/build multiplataforma

## Documentacao

Se sua mudanca altera comportamento publico, revise pelo menos um dos arquivos abaixo:

- `README.md`
- `Branding.md`
- `Development.md`
- `Project-structure.md`

## Issues

Ao relatar um problema, inclua:

- OS
- versao do Go
- versao do `gitloom`
- comandos executados
- saida observada
- saida de `gitloom doctor` quando fizer sentido

## Escopo atual

O projeto ainda nao suporta:

- fluxo automatico com arquivos parcialmente staged
- provider de IA real por padrao

PRs nessas areas sao bem-vindos, mas devem vir com testes e documentacao.
