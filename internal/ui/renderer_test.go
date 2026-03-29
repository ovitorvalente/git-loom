package ui

import (
	"strings"
	"testing"

	"github.com/ovitorvalente/git-loom/internal/app"
	domaincommit "github.com/ovitorvalente/git-loom/internal/domain/commit"
	"github.com/ovitorvalente/git-loom/internal/semantic"
)

func TestRendererCommitPlanCleanMode(t *testing.T) {
	t.Parallel()

	renderer := NewRenderer(RenderOptions{})
	plan := app.CommitPlan{
		Result: app.CommitResult{
			Message: "docs(gitignore): atualizar regras do .gitignore\n\n- remove entradas antigas",
			Commit: domaincommit.Model{
				Type:        domaincommit.TypeDocs,
				Scope:       "gitignore",
				Description: "atualizar regras do .gitignore",
			},
			Paths: []string{".gitignore"},
		},
		Preview: semantic.CommitPreview{
			FilesChanged: 1,
			Additions:    4,
			Deletions:    1,
		},
		Quality: semantic.CommitQuality{
			Score:   85,
			Reasons: []string{"descricao ainda esta generica"},
		},
		Context: semantic.CommitContext{
			Files: []semantic.ChangedFile{{Path: ".gitignore", Status: "atualizado"}},
			Tags:  []string{"config"},
		},
	}

	output := renderer.CommitPlan(1, 1, plan)

	expectedParts := []string{
		"◆ docs(gitignore) [85] bom",
		"mensagem:",
		"docs(gitignore): atualizar regras do .gitignore",
		"detalhes:",
		"remove entradas antigas",
		"arquivos:",
		"+4 -1 .gitignore",
		"analise:",
		"escopo generico: gitignore",
		"sugestoes:",
		"→ config",
		"→ repo",
	}

	for _, expectedPart := range expectedParts {
		if !strings.Contains(output, expectedPart) {
			t.Fatalf("expected output to contain %q, got %q", expectedPart, output)
		}
	}
}

func TestRendererCommitPlanVerboseMode(t *testing.T) {
	t.Parallel()

	renderer := NewRenderer(RenderOptions{
		Mode:        RenderModeVerbose,
		ShowPreview: true,
	})
	plan := app.CommitPlan{
		Result: app.CommitResult{
			Message: "feat(cli): adicionar fluxo de commit\n\n- adiciona renderer novo",
			Commit: domaincommit.Model{
				Type:        domaincommit.TypeFeat,
				Scope:       "cli",
				Description: "adicionar fluxo de commit",
			},
			Paths: []string{"internal/cli/commit.go", "internal/ui/renderer.go"},
		},
		Preview: semantic.CommitPreview{
			FilesChanged: 2,
			Additions:    12,
			Deletions:    3,
		},
		Quality: semantic.CommitQuality{Score: 92},
		Context: semantic.CommitContext{
			Files: []semantic.ChangedFile{
				{Path: "internal/cli/commit.go", Status: "atualizado"},
				{Path: "internal/ui/renderer.go", Status: "adicionado"},
			},
		},
	}

	output := renderer.CommitPlan(1, 2, plan)

	expectedParts := []string{
		"◆ feat(cli) 1/2 [92] excelente",
		"detalhes:",
		"+12 -3 internal/cli/commit.go",
		"+12 -3 internal/ui/renderer.go",
		"verbose:",
		"mensagem final: feat(cli): adicionar fluxo de commit",
		"impacto: +12 -3",
		"detalhe: adiciona renderer novo",
	}

	for _, expectedPart := range expectedParts {
		if !strings.Contains(output, expectedPart) {
			t.Fatalf("expected output to contain %q, got %q", expectedPart, output)
		}
	}
}

func TestRendererChangedFiles(t *testing.T) {
	t.Parallel()

	renderer := NewRenderer(RenderOptions{})
	output := renderer.ChangedFiles(
		[]string{"internal/cli/commit.go"},
		[]string{"internal/app/commit_service.go", "internal/ui/renderer.go"},
	)

	if !strings.Contains(output, "arquivos alterados") {
		t.Fatalf("unexpected output: %q", output)
	}
	if !strings.Contains(output, "staged:") {
		t.Fatalf("unexpected output: %q", output)
	}
	if !strings.Contains(output, "changes (2):") {
		t.Fatalf("unexpected output: %q", output)
	}
	if !strings.Contains(output, "internal/app/commit_service.go") {
		t.Fatalf("unexpected output: %q", output)
	}
}

func TestRendererSuggestions(t *testing.T) {
	t.Parallel()

	renderer := NewRenderer(RenderOptions{})
	output := renderer.Suggestions([]app.CommitSuggestion{
		{Message: "agrupar arquivos de config em um unico bloco"},
	})

	if !strings.Contains(output, "◆ sugestoes") {
		t.Fatalf("unexpected output: %q", output)
	}
	if !strings.Contains(output, "agrupar arquivos de config em um unico bloco") {
		t.Fatalf("unexpected output: %q", output)
	}
}

func TestRendererCommitSummary(t *testing.T) {
	t.Parallel()

	renderer := NewRenderer(RenderOptions{})
	output := renderer.CommitSummary(CommitSummary{
		Created:        1,
		AverageQuality: 85,
		Status:         "working tree limpa",
	})

	expectedParts := []string{
		"✔ 1 commit criado",
		"qualidade media: 85",
		"status: working tree limpa",
		"☕ ate a proxima",
	}

	for _, expectedPart := range expectedParts {
		if !strings.Contains(output, expectedPart) {
			t.Fatalf("expected output to contain %q, got %q", expectedPart, output)
		}
	}
}
