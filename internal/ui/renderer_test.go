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
		"Analisando 1 arquivo modificado...",
		"Sugestao de commit:",
		"docs(gitignore): atualizar regras do .gitignore",
		"Score: 85/100",
		"Tipo: docs",
		"Escopo: gitignore",
		"Detalhes:",
		"remove entradas antigas",
		"│ ~ .gitignore",
		"Analise:",
		"escopo generico: gitignore",
		"Sugestoes:",
		"• config",
		"• repo",
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
		"Analisando 2 arquivos modificados... [1/2]",
		"Detalhes:",
		"│ ~ internal/cli/commit.go",
		"│ + internal/ui/renderer.go",
		"Verbose:",
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

	if !strings.Contains(output, "Arquivos detectados no working tree") {
		t.Fatalf("unexpected output: %q", output)
	}
	if !strings.Contains(output, "Staged:") {
		t.Fatalf("unexpected output: %q", output)
	}
	if !strings.Contains(output, "Changes (2):") {
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

	if !strings.Contains(output, "Sugestoes encontradas") {
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
		"Commit concluido: 1 commit criado",
		"Qualidade media: 85",
		"Status: working tree limpa",
		"ate a proxima",
	}

	for _, expectedPart := range expectedParts {
		if !strings.Contains(output, expectedPart) {
			t.Fatalf("expected output to contain %q, got %q", expectedPart, output)
		}
	}
}
