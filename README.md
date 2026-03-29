# Git Loom

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg?style=flat-square)](LICENSE)

`gitloom` é um CLI em Go para sugerir e criar commits semânticos a partir das mudanças do repositório.

Esta primeira versão pública é intencionalmente pequena: o foco atual é o comando `commit`.

## Escopo da v0.1.0

- gera mensagens em português seguindo Conventional Commits
- agrupa arquivos em blocos de no máximo 4 arquivos por commit
- considera arquivos já staged
- detecta arquivos modificados fora do stage e pergunta se devem ser adicionados
- mostra um plano visual antes de criar os commits
- calcula score de qualidade por commit planejado
- sugere melhorias de agrupamento e descrição antes de commitar
- suporta `--dry-run` e `--yes`
- suporta `--preview` e `--strict`
- suporta configuração mínima por `.gitloom.yaml`

## Limitações conhecidas

- o projeto ainda não suporta fluxo automático com arquivos parcialmente staged
- comandos de branch, analyze e config ainda não fazem parte desta versão
- o provider de IA atual é `noop`; a geração é heurística e determinística

## Instalação

### Go

```bash
go install github.com/ovitorvalente/git-loom/cmd/gitloom@latest
```

### Build local

```bash
git clone https://github.com/ovitorvalente/git-loom.git
cd git-loom
go build -o gitloom ./cmd/gitloom
```

## Uso

```bash
git add .
gitloom commit
```

Para apenas revisar o plano:

```bash
gitloom commit --dry-run
```

Para criar os commits sem confirmação:

```bash
gitloom commit --yes
```

Para revisar impacto e diff resumido sem commitar:

```bash
gitloom commit --preview
```

Para falhar quando o plano estiver fraco semanticamente:

```bash
gitloom commit --strict
```

## Exemplo de fluxo

```bash
git add internal/cli/commit.go
gitloom commit
```

Saída esperada:

```text
◆ commit 1/1  [feat]  (qualidade: 92/100)
  resumo: 1 arquivo(s) | 1 detalhe(s) | impacto: +12/-3
✦ commit gerado
  • tipo: feat
  • escopo: cli
  • intencao: deixar o fluxo mais claro e previsivel
  • descricao: adicionar fluxo de commit
  › mensagem: feat(cli): adicionar fluxo de commit
  ✧ detalhes:
    - adiciona comando commit em cli
  ✧ arquivos:
    ◦ internal/cli/commit.go

> criar commits planejados? [y/N]:
```

## Configuração

Arquivo opcional `.gitloom.yaml` na raiz do repositório:

```yaml
commit:
  scope: core

cli:
  auto_confirm: false
```

## Regras de commit

O projeto segue as regras descritas em [.rules/commit.md](.rules/commit.md):

- mensagens em pt-BR
- cabeçalho semântico
- no máximo 72 caracteres na linha principal
- detalhes em blocos com bullets
- no máximo 4 arquivos por commit
- score mínimo no modo estrito

## Desenvolvimento

```bash
make test
make vet
make build
```

## Testes

O projeto possui:

- testes unitários por camada
- testes de integração do repositório Git com um repositório temporário real

## Licença

MIT. Veja [LICENSE](LICENSE).
