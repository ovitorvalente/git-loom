package commit

import (
	"strings"
	"testing"
)

func TestAnalyzeDiff(t *testing.T) {
	t.Parallel()

	diff := strings.Join([]string{
		"diff --git a/internal/app/commit_service.go b/internal/app/commit_service.go",
		"index 1111111..2222222 100644",
		"--- a/internal/app/commit_service.go",
		"+++ b/internal/app/commit_service.go",
		"@@ -1,3 +1,3 @@",
		"diff --git a/internal/cli/commit.go b/internal/cli/commit.go",
		"new file mode 100644",
		"--- /dev/null",
		"+++ b/internal/cli/commit.go",
	}, "\n")

	analysis := AnalyzeDiff(diff, TypeFix)

	if analysis.Scope != "app" {
		t.Fatalf("expected scope app, got %q", analysis.Scope)
	}
	if analysis.Description != "corrigir app" {
		t.Fatalf("expected structured description, got %q", analysis.Description)
	}
	if !strings.Contains(analysis.Body, "- atualiza commit service em app") {
		t.Fatalf("expected body to mention updated file, got %q", analysis.Body)
	}
	if !strings.Contains(analysis.Body, "- adiciona comando commit em cli") { //nolint:misspell
		t.Fatalf("expected body to mention added file, got %q", analysis.Body)
	}
}

func TestAnalyzeDiffForTests(t *testing.T) {
	t.Parallel()

	diff := strings.Join([]string{
		"diff --git a/internal/domain/commit/analyzer_test.go b/internal/domain/commit/analyzer_test.go",
		"index 1111111..2222222 100644",
	}, "\n")

	analysis := AnalyzeDiff(diff, TypeTest)

	if analysis.Scope != "commit" {
		t.Fatalf("expected scope commit, got %q", analysis.Scope)
	}
	if analysis.Description != "ajustar testes de analyzer" {
		t.Fatalf("expected structured description, got %q", analysis.Description)
	}
	if !strings.Contains(analysis.Body, "- atualiza testes de analyzer em domain/commit") {
		t.Fatalf("expected detailed body, got %q", analysis.Body)
	}
}

func TestAnalyzeDiffAddsSemanticHintsFromPatch(t *testing.T) {
	t.Parallel()

	diff := strings.Join([]string{
		"diff --git a/internal/cli/commit.go b/internal/cli/commit.go",
		"index 1111111..2222222 100644",
		"--- a/internal/cli/commit.go",
		"+++ b/internal/cli/commit.go",
		"@@ -10,6 +10,12 @@",
		"+func buildJSONExecutionSummary() {}",
		"+command.Flags().BoolVar(&options.json, \"json\", false, \"render json\")",
		"+Use: \"commit\"",
	}, "\n")

	analysis := AnalyzeDiff(diff, TypeFeat)

	if analysis.Description != "atualizar saida json do cli" {
		t.Fatalf("expected semantic description, got %q", analysis.Description)
	}
	if !strings.Contains(analysis.Body, "- ajusta funcoes: buildJSONExecutionSummary") {
		t.Fatalf("expected function hint in body, got %q", analysis.Body)
	}
	if !strings.Contains(analysis.Body, "- ajusta flags: --json") {
		t.Fatalf("expected flag hint in body, got %q", analysis.Body)
	}
	if !strings.Contains(analysis.Body, "- ajusta comandos: commit") {
		t.Fatalf("expected command hint in body, got %q", analysis.Body)
	}
}
