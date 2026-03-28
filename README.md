# 🔧 Git Loom

> **Automação Inteligente de Git Workflow** — De Dev para Dev

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg?style=flat-square)](LICENSE)
[![GitHub Release](https://img.shields.io/github/v/release/seu-usuario/gitcraft?style=flat-square)](https://github.com/seu-usuario/gitcraft/releases)
[![GitHub Stars](https://img.shields.io/github/stars/seu-usuario/gitcraft?style=flat-square)](https://github.com/seu-usuario/gitcraft)

**Git Loom** é um **CLI moderno** que automatiza completamente seu workflow Git. Gera commits semânticos, cria branches automáticas com nomes profissionais, segue padrões de **Conventional Commits**, e estrutura seu versionamento de forma inteligente.

Feito **de dev para dev** — sem burocracias, sem complicações.

---

## ✨ Funcionalidades

- 🚀 **Commits Automáticos Semânticos** - Analisa suas mudanças e gera mensagens estruturadas
- 🌿 **Branches Inteligentes** - Cria branch names baseadas em tipo de mudança (feat, fix, docs, etc)
- 📝 **Padrão Conventional Commits** - Segue `type(scope): description` automaticamente
- 🧱 **Estrutura de Blocos** - Agrupa commits logicamente em blocos temáticos
- 📊 **Análise de Mudanças** - Detecta tipos de alteração (adição, modificação, remoção)
- ⚙️ **Configurável** - Templates customizáveis para seu workflow específico
- 🎨 **UI Interativa** - Prompts elegantes com TUI usando Bubble Tea
- 📦 **Zero Config** - Funciona pronto para usar, configure se quiser

---

## 🚀 Quick Start

### Instalação

#### Homebrew (macOS)

```bash
brew install seu-usuario/gitcraft/gitcraft
```

#### Go

```bash
go install github.com/seu-usuario/gitcraft@latest
```

#### Build from Source

```bash
git clone https://github.com/seu-usuario/gitcraft.git
cd gitcraft
go build -o gitcraft ./cmd/gitcraft
sudo mv gitcraft /usr/local/bin/
```

### Inicialização

```bash
# Inicie o GitCraft em seu repositório
gitcraft init

# Você será guiado através de uma configuração rápida
```

### Uso Básico

```bash
# Crie um commit automático analisando suas mudanças
gitcraft commit

# Crie uma branch semântica
gitcraft branch feat "adicionar autenticação"
# Resultado: feat/adicionar-autenticacao

# Analise mudanças pendentes
gitcraft analyze

# Veja suas configurações
gitcraft config
```

---

## 📖 Documentação Completa

### Comandos Principais

#### `gitcraft init`

Inicializa o GitCraft no seu repositório. Cria um arquivo `.gitcraft.yml` com configurações padrão.

```bash
gitcraft init
```

#### `gitcraft commit`

Analisa suas staged changes e abre um workflow interativo para criar um commit semântico.

```bash
# Modo interativo (recomendado)
gitcraft commit

# Com flags
gitcraft commit --type feat --scope auth --message "adicionar login com Google"
gitcraft commit -t fix -s api -m "corrigir timeout de requisição"
```

**Tipos de Commit Suportados:**

- `feat` - Nova feature
- `fix` - Correção de bug
- `docs` - Mudanças em documentação
- `style` - Formatação, sem mudanças lógicas
- `refactor` - Refatoração de código
- `perf` - Melhorias de performance
- `test` - Adição/modificação de testes
- `ci` - Mudanças em CI/CD
- `chore` - Outras mudanças (deps, config)

#### `gitcraft branch`

Cria uma nova branch com nome semântico baseada em padrões.

```bash
# Modo interativo
gitcraft branch

# Com argumentos diretos
gitcraft branch feat "autenticação JWT"
gitcraft branch bugfix "memory leak em listeners"
gitcraft branch docs "atualizar API docs"

# Resultado:
# feat/autenticacao-jwt
# bugfix/memory-leak-em-listeners
# docs/atualizar-api-docs
```

#### `gitcraft analyze`

Analisa alterações no repositório e sugere tipos de commit.

```bash
gitcraft analyze

# Output exemplo:
# 📊 Análise de Mudanças
# ├─ 3 arquivos modificados
# ├─ 2 arquivos adicionados
# ├─ 1 arquivo removido
# ├─ Tipo sugerido: feat
# └─ Escopo sugerido: core
```

#### `gitcraft config`

Gerencia configurações do GitCraft.

```bash
gitcraft config set author-name "Seu Nome"
gitcraft config set author-email "seu@email.com"
gitcraft config get author-name
gitcraft config list
```

---

## ⚙️ Configuração

GitCraft utiliza um arquivo `.gitcraft.yml` na raiz do repositório:

```yaml
# .gitcraft.yml
project:
  name: meu-projeto
  description: Descrição do projeto

commit:
  pattern: conventional
  max-length: 72
  scopes:
    - auth
    - api
    - core
    - ui
    - docs

branch:
  auto-prefix: true
  lowercase: true
  replace-spaces: "-"

author:
  name: "Seu Nome"
  email: "seu@email.com"

ai:
  enabled: false
  provider: openai
  api-key: ""
```

---

## 🔄 Workflow Exemplo

```bash
# 1. Criar uma nova feature
gitcraft branch feat "sistema de cache"
# ✓ Branch criada: feat/sistema-de-cache

git checkout feat/sistema-de-cache

# 2. Fazer mudanças
echo "cache implementation" > cache.go
git add cache.go

# 3. Criar commit automático
gitcraft commit
# ? Qual é o tipo de commit? → feat
# ? Qual é o escopo? → core
# ? Descrição breve? → implementar cache distribuído
# ? Body (opcional)? → Adiciona suporte a Redis e Memcached
# ✓ Commit: feat(core): implementar cache distribuído

# 4. Analizar e revisar
gitcraft analyze

# 5. Push para remoto
git push origin feat/sistema-de-cache
```

---

## 🎨 Padrão de Commits (Conventional Commits)

GitCraft segue o padrão **Conventional Commits v1.0.0**:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Exemplos

```
feat(auth): adicionar autenticação OAuth2

Implementa login via Google e GitHub.
Adiciona refresh token automático.

Closes #123
```

```
fix(api): corrigir tratamento de erros 500

Anteriormente, erros internos não eram logados corretamente.
Agora todos os erros são capturados e armazenados no ELK.

BREAKING CHANGE: Error response format changed
```

```
docs: atualizar README com exemplos
```

---

## 📦 Estrutura do Projeto

```
gitcraft/
├── cmd/
│   └── gitcraft/
│       └── main.go              # Entry point
├── internal/
│   ├── commit/
│   │   ├── analyzer.go          # Analisa mudanças
│   │   ├── generator.go         # Gera commits
│   │   └── message.go           # Estrutura de mensagem
│   ├── branch/
│   │   ├── creator.go           # Cria branches
│   │   └── namer.go             # Nomenclatura semântica
│   ├── config/
│   │   ├── loader.go            # Carrega config
│   │   └── defaults.go          # Configurações padrão
│   ├── git/
│   │   ├── repository.go        # Operações Git
│   │   └── staging.go           # Gerencia staging area
│   └── ui/
│       ├── prompts.go           # Prompts interativos
│       └── output.go            # Formatação de output
├── pkg/
│   └── types/
│       └── commit.go            # Tipos de commit
├── go.mod
├── go.sum
├── Makefile
├── LICENSE                       # MIT License
├── README.md
├── CONTRIBUTING.md
└── .github/
    ├── workflows/
    │   ├── ci.yml               # CI/CD pipeline
    │   └── release.yml          # Release automation
    └── ISSUE_TEMPLATE/
        ├── bug_report.md
        └── feature_request.md
```

---

## 🤝 Contribuindo

Adoramos contribuições! Antes de começar, leia [CONTRIBUTING.md](CONTRIBUTING.md).

### Setup para Desenvolvimento

```bash
# Clone o repositório
git clone https://github.com/seu-usuario/gitcraft.git
cd gitcraft

# Instale dependências
go mod download

# Build
make build

# Testes
make test

# Lint
make lint
```

### Reportar Bugs

Use [GitHub Issues](https://github.com/seu-usuario/gitcraft/issues) com a template de bug report.

### Sugerir Features

Abra uma [Discussion](https://github.com/seu-usuario/gitcraft/discussions) ou [Issue](https://github.com/seu-usuario/gitcraft/issues) com o label `enhancement`.

---

## 📄 Licença

GitCraft é licenciado sob a [MIT License](LICENSE).

---

## 🙏 Agradecimentos

Inspirado por projetos incríveis como:

- [Commitizen](https://commitizen-tools.github.io/commitizen/)
- [conventional-commits](https://www.conventionalcommits.org/)
- [Semantic Versioning](https://semver.org/)

---

## 📞 Suporte

- 📖 **Documentação**: [docs/](docs/)
- 💬 **Discussions**: [GitHub Discussions](https://github.com/seu-usuario/gitcraft/discussions)
- 🐛 **Issues**: [GitHub Issues](https://github.com/seu-usuario/gitcraft/issues)
- 🐦 **Twitter**: [@seu-usuario](https://twitter.com/seu-usuario)

---

<div align="center">

**[⬆ back to top](#gitcraft)**

Made with ❤️ by developers, for developers.

</div>
