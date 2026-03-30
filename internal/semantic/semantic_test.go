package semantic

import "testing"

func TestNormalizeScope(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		path          string
		expectedScope string
	}{
		{path: "README.md", expectedScope: "readme"},
		{path: "go.mod", expectedScope: "deps"},
		{path: "Makefile", expectedScope: "build"},
		{path: "internal/cli/commit.go", expectedScope: "cli"},
		{path: "internal/ui/output.go", expectedScope: "ui"},
		{path: "internal/domain/commit/analyzer.go", expectedScope: "commit"},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.path, func(t *testing.T) {
			t.Parallel()

			scope := NormalizeScope(testCase.path)
			if scope != testCase.expectedScope {
				t.Fatalf("expected %q, got %q", testCase.expectedScope, scope)
			}
		})
	}
}

func TestDetectIntent(t *testing.T) {
	t.Parallel()

	context := CommitContext{
		Files: []ChangedFile{{Path: "README.md", Status: "atualizado"}},
		Diff:  "diff --git a/README.md b/README.md\n+gitloom commit\n",
		Tags:  []string{"readme", "commit", "cli"},
	}

	intent := DetectIntent("docs", context)
	if intent.Scope != "readme" {
		t.Fatalf("expected readme scope, got %q", intent.Scope)
	}
	if intent.Description != "atualizar instrucoes de uso do cli commit" {
		t.Fatalf("unexpected description: %q", intent.Description)
	}
	if intent.Intent == "" {
		t.Fatal("expected non-empty intent")
	}
}

func TestScoreCommit(t *testing.T) {
	t.Parallel()

	context := CommitContext{
		Files: []ChangedFile{
			{Path: "go.mod", Status: "atualizado"},
			{Path: "go.sum", Status: "atualizado"},
		},
	}

	quality := ScoreCommit(ChangeIntent{
		Type:        "chore",
		Scope:       "deps",
		Description: "atualizar dependencias do projeto", //nolint:misspell
	}, context)

	if quality.Score <= 70 {
		t.Fatalf("expected a strong score, got %d", quality.Score)
	}
}
