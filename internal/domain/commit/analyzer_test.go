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
	if !strings.Contains(analysis.Body, "- atualizado commit service") {
		t.Fatalf("expected body to mention updated file, got %q", analysis.Body)
	}
	if !strings.Contains(analysis.Body, "- adicionado commit") {
		t.Fatalf("expected body to mention added file, got %q", analysis.Body)
	}
}
