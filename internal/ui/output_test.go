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
		Message: "feat(core): adicionar formatter de commit",
		Commit: domaincommit.Model{
			Type:        domaincommit.TypeFeat,
			Scope:       "core",
			Intent:      "deixar o fluxo mais claro",
			Description: "adicionar formatter de commit",
			Body:        "manter o output facil de ler",
		},
	}

	formattedOutput := FormatCommitResult(result)

	expectedParts := []string{
		"commit gerado",
		"tipo: feat",
		"escopo: core",
		"intencao: deixar o fluxo mais claro",
		"descricao: adicionar formatter de commit",
		"mensagem: feat(core): adicionar formatter de commit",
		"detalhes:",
		"manter o output facil de ler",
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

	plan := app.CommitPlan{
		Result: app.CommitResult{
			Message: "feat(cli): adicionar fluxo de commit\n\n- adiciona comando commit em cli",
			Commit: domaincommit.Model{
				Type:        domaincommit.TypeFeat,
				Scope:       "cli",
				Intent:      "deixar o fluxo mais claro",
				Description: "adicionar fluxo de commit",
			},
			Paths: []string{"internal/cli/commit.go"},
		},
		Quality: app.CommitQuality{Score: 91},
		Preview: app.CommitPreview{
			FilesChanged: 1,
			Additions:    12,
			Deletions:    3,
		},
	}

	formattedOutput := FormatCommitPlan(1, 2, plan, true)
	if !strings.Contains(formattedOutput, "commit 1/2") {
		t.Fatalf("unexpected output: %q", formattedOutput)
	}
	if !strings.Contains(formattedOutput, "arquivos:") {
		t.Fatalf("unexpected output: %q", formattedOutput)
	}
	if !strings.Contains(formattedOutput, "qualidade: 91/100") {
		t.Fatalf("unexpected output: %q", formattedOutput)
	}
	if !strings.Contains(formattedOutput, "preview:") {
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

func TestFormatSuggestions(t *testing.T) {
	t.Parallel()

	formattedOutput := FormatSuggestions([]app.CommitSuggestion{
		{Message: "melhorar descricao do commit 1"},
	})
	if !strings.Contains(formattedOutput, "sugestoes") {
		t.Fatalf("unexpected output: %q", formattedOutput)
	}
}
