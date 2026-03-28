# 🤖 CLAUDE.md

## 📌 Project Overview

**Project Name:** gitloom
**Language:** Go (Golang)
**CLI Framework:** Cobra
**Architecture:** Clean Architecture + Domain-Driven Design (DDD)
**Development Approach:** Test-Driven Development (TDD)

**Purpose:**
gitloom is a CLI tool that automates Git workflows by analyzing code changes, generating structured commit messages, managing branches, and enforcing conventions.

---

# 🧠 Architectural Vision

gitloom is not just a CLI tool.

> It is a **Git Automation Engine**, designed to be:

- Modular
- Extensible
- Testable
- Deterministic
- Provider-agnostic (Git, AI, etc.)

---

# 🧱 Architecture Overview

The system follows a **layered architecture with strict boundaries**:

```
┌──────────────────────────────┐
│           CLI Layer          │  → cmd/ + cli/
└──────────────┬───────────────┘
               ↓
┌──────────────────────────────┐
│      Application Layer       │  → internal/app/
│   (Use Cases / Orchestration)│
└──────────────┬───────────────┘
               ↓
┌──────────────────────────────┐
│        Domain Layer          │  → internal/domain/
│   (Business Rules / Logic)   │
└──────────────┬───────────────┘
               ↓
┌──────────────────────────────┐
│    Infrastructure Layer      │  → internal/infra/
│ (Git, AI, FS, External APIs) │
└──────────────────────────────┘
```

---

# 📂 Project Structure (Official)

```
gitloom/
├── cmd/
│   └── gitloom/
│       └── main.go

├── internal/

│   ├── app/                      # Use cases (orchestration layer)
│   │   ├── commit_service.go
│   │   ├── branch_service.go
│   │   └── workflow_service.go

│   ├── domain/                   # Core business logic (pure)
│   │   ├── commit/
│   │   │   ├── analyzer.go
│   │   │   ├── classifier.go
│   │   │   ├── generator.go
│   │   │   └── model.go
│   │   │
│   │   ├── branch/
│   │   │   ├── namer.go
│   │   │   └── model.go
│   │   │
│   │   └── shared/
│   │       └── types.go

│   ├── interfaces/               # Contracts (abstractions)
│   │   ├── git.go
│   │   ├── ai.go
│   │   └── logger.go

│   ├── infra/                    # External implementations
│   │   ├── git/
│   │   │   ├── client.go
│   │   │   ├── repository.go
│   │   │   └── staging.go
│   │   │
│   │   ├── ai/
│   │   │   ├── provider.go
│   │   │   ├── ollama.go
│   │   │   └── noop.go
│   │   │
│   │   ├── config/
│   │   │   ├── loader.go
│   │   │   └── schema.go
│   │   │
│   │   └── system/
│   │       ├── exec.go
│   │       └── fs.go

│   ├── cli/                      # CLI commands (Cobra)
│   │   ├── root.go
│   │   ├── commit.go
│   │   ├── branch.go
│   │   └── init.go

│   └── ui/                       # User interaction
│       ├── prompts.go
│       └── output.go

├── pkg/                          # Public API (future SDK)
│   └── gitloom/
│       └── client.go

├── configs/
│   └── default.yaml

├── scripts/
│   └── install.sh

├── .gitloom.yaml
├── go.mod
├── Makefile
```

---

# 🔄 Execution Flow (Commit Example)

```
CLI (commit command)
↓
CommitService (app)
↓
GitRepository.GetDiff()
↓
CommitAnalyzer (domain)
↓
CommitClassifier (domain)
↓
CommitGenerator (domain)
↓
AIProvider (optional)
↓
User confirmation (ui)
↓
GitRepository.Commit()
```

---

# 🧠 Domain Layer Rules (STRICT)

- MUST be pure (no IO, no Git, no HTTP)
- MUST be deterministic
- MUST be fully testable

### Responsibilities:

- Analyze diff
- Classify commit type
- Generate commit message
- Define domain models

---

## Example:

```go
func ClassifyCommit(diff string) CommitType
```

---

# ⚙️ Application Layer (Use Cases)

Acts as the **orchestrator**.

### Responsibilities:

- Coordinate domain + infra
- Execute workflows
- Handle business flows

---

## Example:

```go
type CommitService struct {
    git GitRepository
    ai  AIProvider
}
```

---

# 🔌 Interfaces Layer (Contracts)

All dependencies must be abstracted.

---

## Git Interface:

```go
type GitRepository interface {
    GetDiff() (string, error)
    Commit(message string) error
    CreateBranch(name string) error
}
```

---

## AI Interface:

```go
type AIProvider interface {
    GenerateCommit(diff string) (string, error)
}
```

---

# 🏗️ Infrastructure Layer

Implements external concerns:

- Git CLI
- AI (Ollama, future providers)
- File system
- Config

---

### Rules:

- Must implement interfaces
- Must NOT contain business logic
- Must be replaceable

---

# 🖥️ CLI Layer (Cobra)

Responsibilities:

- Parse flags
- Call application services
- Render output

---

### ❌ Forbidden:

- Business logic
- Direct git operations

---

# 🎨 UI Layer

Handles:

- Prompts
- Output formatting

Must be independent of business logic.

---

# 🧩 Dependency Injection

Dependencies must be injected manually.

Example:

```go
gitRepo := git.NewRepository()
ai := ai.NewNoopProvider()

service := app.NewCommitService(gitRepo, ai)
```

---

# 🧼 Clean Code Rules

- Functions ≤ 30 lines
- No hidden side effects
- Explicit naming
- No magic numbers
- Early returns preferred

---

# 🧱 SOLID Principles

- S → One responsibility per module
- O → Extend via interfaces
- L → Replace implementations safely
- I → Small interfaces
- D → Depend on abstractions

---

# 🧪 Testing Strategy (TDD)

### Approach:

1. Write failing test
2. Implement minimal code
3. Refactor

---

## Test Types:

### Unit Tests

- Domain layer
- Pure logic

### Application Tests

- Use mocks

### Integration Tests

- Infra layer

---

## Rules:

- No real Git calls in unit tests
- Use mocks for interfaces

---

# 🧾 Configuration

File: `.gitloom.yaml`

Rules:

- Always validate
- Provide defaults
- Never crash if missing

---

# 🧠 Commit Convention

Based on Conventional Commits

---

## Format:

```
type(scope): description

- detail 1
- detail 2
```

---

# 🌿 Branch Naming

```
<type>/<description>
```

---

# 🔍 Diff Analysis Rules

- Detect intent from code changes
- Use keyword heuristics
- Group related changes

---

# 🧩 Commit Splitting

If multiple logical changes are found:

→ Generate multiple commits

---

# 🚀 Extensibility Design

System must support:

- AI providers
- Git providers
- Plugin system (future)

---

# ⚠️ Error Handling

- Never ignore errors
- Always wrap:

```go
fmt.Errorf("context: %w", err)
```

---

# 📊 Logging

- Structured logs preferred
- Avoid noisy output

---

# 🔐 Safety Rules

- Never commit without confirmation (unless flag)
- Never rewrite history
- Never delete branches automatically

---

# 🚫 Anti-Patterns

- ❌ Business logic in CLI
- ❌ Domain using infra
- ❌ Large functions (>100 lines)
- ❌ Tight coupling

---

# 🚀 Performance Guidelines

- Prefer simplicity first
- Avoid premature optimization
- Use streaming for large diffs

---

# 🧬 Future Roadmap

- AI commit enhancement
- PR description generation
- GitHub integration
- Plugin system

---

# 🧑‍💻 Developer Workflow

1. Create branch
2. Write test (TDD)
3. Implement feature
4. Refactor
5. Validate
6. Open PR

---

# ✅ Definition of Done

- Code follows architecture
- Tests implemented
- No lint errors
- CLI is consistent
- Code is readable

---

# 🧠 Final Rule

> If the architecture is violated, the code is invalid — even if it works.
