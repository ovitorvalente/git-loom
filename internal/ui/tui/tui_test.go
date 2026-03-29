package tui

import (
	"testing"

	"github.com/ovitorvalente/git-loom/internal/app"
	domaincommit "github.com/ovitorvalente/git-loom/internal/domain/commit"
	"github.com/ovitorvalente/git-loom/internal/semantic"
)

func TestRunWithSpinnerExecutesFunction(t *testing.T) {
	t.Parallel()

	executed := false
	err := RunWithSpinner("testando", func() error {
		executed = true
		return nil
	})

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !executed {
		t.Fatal("expected function to be executed")
	}
}

func TestRunWithSpinnerPropagatesError(t *testing.T) {
	t.Parallel()

	expectedErr := &testError{"spinner error"}
	err := RunWithSpinner("testando", func() error {
		return expectedErr
	})

	if err != expectedErr {
		t.Fatalf("expected %v, got %v", expectedErr, err)
	}
}

type testError struct{ msg string }

func (e *testError) Error() string { return e.msg }

func TestSelectCommitsEmptyPlans(t *testing.T) {
	t.Parallel()

	result, err := SelectCommits(nil)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !result.Confirmed {
		t.Fatal("expected confirmed for empty plans")
	}
}

func TestSelectorModelResult(t *testing.T) {
	t.Parallel()

	plans := []app.CommitPlan{
		{
			Result: app.CommitResult{
				Message: "feat(cli): adicionar comando",
				Commit: domaincommit.Model{
					Type:        domaincommit.TypeFeat,
					Scope:       "cli",
					Description: "adicionar comando",
				},
				Paths: []string{"internal/cli/commit.go"},
			},
			Quality: semantic.CommitQuality{Score: 95},
		},
		{
			Result: app.CommitResult{
				Message: "docs(readme): atualizar",
				Commit: domaincommit.Model{
					Type:        domaincommit.TypeDocs,
					Scope:       "readme",
					Description: "atualizar",
				},
				Paths: []string{"README.md"},
			},
			Quality: semantic.CommitQuality{Score: 80},
		},
	}

	m := newSelectorModel(plans)

	if len(m.items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(m.items))
	}
	if m.items[0].commitType != "feat" {
		t.Fatalf("expected feat, got %s", m.items[0].commitType)
	}
	if m.items[1].commitType != "docs" {
		t.Fatalf("expected docs, got %s", m.items[1].commitType)
	}

	m.items[0].selected = true
	m.items[1].selected = false

	result := m.result()
	if len(result.Approved) != 2 {
		t.Fatalf("expected 2 approvals, got %d", len(result.Approved))
	}
	if !result.Approved[0] {
		t.Fatal("expected first item to be approved")
	}
	if result.Approved[1] {
		t.Fatal("expected second item to not be approved")
	}
}

func TestTypeColorDoesNotPanic(t *testing.T) {
	t.Parallel()

	types := []string{"feat", "fix", "refactor", "docs", "test", "chore", "unknown"}
	for _, tp := range types {
		_ = TypeColor(tp)
	}
}

func TestScoreColorDoesNotPanic(t *testing.T) {
	t.Parallel()

	scores := []int{95, 85, 75, 50}
	for _, s := range scores {
		_ = ScoreColor(s)
	}
}
