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
		"generated commit",
		"type: feat",
		"scope: core",
		"description: add commit formatter",
		"message: feat(core): add commit formatter",
		"body: keep output easy to read",
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
		Message: "chore: update repository changes",
		Commit: domaincommit.Model{
			Type:        domaincommit.TypeChore,
			Description: "update repository changes",
		},
	}

	formattedOutput := FormatCommitResult(result)

	if !strings.Contains(formattedOutput, "scope: -") {
		t.Fatalf("expected empty scope placeholder, got %q", formattedOutput)
	}
	if strings.Contains(formattedOutput, "body:") {
		t.Fatalf("did not expect body line, got %q", formattedOutput)
	}
}
