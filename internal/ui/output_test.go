package ui

import (
	"strings"
	"testing"

	"github.com/ovitorvalente/git-loom/internal/app"
	domaincommit "github.com/ovitorvalente/git-loom/internal/domain/commit"
)

func TestFormatCommitResult(t *testing.T) {
	t.Parallel()

	result := app.CommitResult{
		Message: "feat(core): add commit formatter",
		Commit: domaincommit.Model{
			Type:        domaincommit.TypeFeat,
			Scope:       "core",
			Description: "add commit formatter",
			Body:        "keep output easy to read",
		},
	}

	formattedOutput := FormatCommitResult(result)

	expectedParts := []string{
		"commit gerado",
		"tipo: feat",
		"escopo: core",
		"descricao: add commit formatter",
		"mensagem: feat(core): add commit formatter",
		"detalhes:",
		"keep output easy to read",
	}

	for _, expectedPart := range expectedParts {
		if !strings.Contains(formattedOutput, expectedPart) {
			t.Fatalf("expected output to contain %q, got %q", expectedPart, formattedOutput)
		}
	}
}

func TestFormatCommitResultWithoutOptionalFields(t *testing.T) {
	t.Parallel()

	result := app.CommitResult{
		Message: "chore: atualizar mudancas do repositorio",
		Commit: domaincommit.Model{
			Type:        domaincommit.TypeChore,
			Description: "atualizar mudancas do repositorio",
		},
	}

	formattedOutput := FormatCommitResult(result)

	if !strings.Contains(formattedOutput, "escopo: -") {
		t.Fatalf("expected empty scope placeholder, got %q", formattedOutput)
	}
	if strings.Contains(formattedOutput, "detalhes:") {
		t.Fatalf("did not expect body line, got %q", formattedOutput)
	}
}

func TestFormatCommitPlan(t *testing.T) {
	t.Parallel()

	result := app.CommitResult{
		Message: "feat(cli): adicionar commit\n\n- adicionado commit",
		Commit: domaincommit.Model{
			Type:        domaincommit.TypeFeat,
			Scope:       "cli",
			Description: "adicionar commit",
		},
		Paths: []string{"internal/cli/commit.go"},
	}

	formattedOutput := FormatCommitPlan(1, 2, result)
	if !strings.Contains(formattedOutput, "bloco 1/2") {
		t.Fatalf("unexpected output: %q", formattedOutput)
	}
	if !strings.Contains(formattedOutput, "arquivos:") {
		t.Fatalf("unexpected output: %q", formattedOutput)
	}
}

func TestFormatChangedFiles(t *testing.T) {
	t.Parallel()

	formattedOutput := FormatChangedFiles([]string{"internal/app/commit_service.go"})
	if !strings.Contains(formattedOutput, "arquivos em changes") {
		t.Fatalf("unexpected output: %q", formattedOutput)
	}
}
