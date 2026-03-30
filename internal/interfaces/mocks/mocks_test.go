package mocks_test

import (
	"errors"
	"testing"

	"github.com/ovitorvalente/git-loom/internal/interfaces"
	"github.com/ovitorvalente/git-loom/internal/interfaces/mocks"
)

func TestGitRepositoryImplementsInterface(t *testing.T) {
	t.Parallel()

	var _ interfaces.GitRepository = &mocks.GitRepository{}
}

func TestGitRepositoryTracksCalls(t *testing.T) {
	t.Parallel()

	expectedError := errors.New("commit failed")
	repository := &mocks.GitRepository{
		GetDiffFunc: func(paths ...string) (string, error) {
			return "diff content", nil
		},
		CommitFunc: func(message string) error {
			if message != "feat(core): add mocks" {
				t.Fatalf("unexpected commit message: %s", message)
			}

			return expectedError
		},
		CreateBranchFunc: func(name string) error {
			if name != "feat/mocks" {
				t.Fatalf("unexpected branch name: %s", name)
			}

			return nil
		},
	}

	diff, err := repository.GetDiff()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if diff != "diff content" {
		t.Fatalf("expected diff content, got %s", diff)
	}
	if len(repository.GetDiffCalls) != 1 {
		t.Fatalf("expected one diff call, got %d", len(repository.GetDiffCalls))
	}

	err = repository.Commit("feat(core): add mocks")
	if !errors.Is(err, expectedError) {
		t.Fatalf("expected commit error, got %v", err)
	}

	err = repository.CreateBranch("feat/mocks")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if len(repository.CommitCalls) != 1 {
		t.Fatalf("expected one commit call, got %d", len(repository.CommitCalls))
	}
	if repository.CommitCalls[0] != "feat(core): add mocks" {
		t.Fatalf("unexpected commit call: %s", repository.CommitCalls[0])
	}
	if len(repository.CreateBranchCalls) != 1 {
		t.Fatalf("expected one branch call, got %d", len(repository.CreateBranchCalls))
	}
	if repository.CreateBranchCalls[0] != "feat/mocks" {
		t.Fatalf("unexpected branch call: %s", repository.CreateBranchCalls[0])
	}
}

func TestAIProviderImplementsInterface(t *testing.T) {
	t.Parallel()

	var _ interfaces.AIProvider = &mocks.AIProvider{}
}

func TestAIProviderTracksCalls(t *testing.T) {
	t.Parallel()

	expectedError := errors.New("provider failed")
	provider := &mocks.AIProvider{
		GenerateCommitFunc: func(diff string) (string, error) {
			if diff != "diff content" {
				t.Fatalf("unexpected diff: %s", diff)
			}

			return "", expectedError
		},
	}

	message, err := provider.GenerateCommit("diff content")
	if message != "" {
		t.Fatalf("expected empty message, got %s", message)
	}
	if !errors.Is(err, expectedError) {
		t.Fatalf("expected provider error, got %v", err)
	}
	if len(provider.GenerateCommitCalls) != 1 {
		t.Fatalf("expected one provider call, got %d", len(provider.GenerateCommitCalls))
	}
	if provider.GenerateCommitCalls[0] != "diff content" {
		t.Fatalf("unexpected provider call: %s", provider.GenerateCommitCalls[0])
	}
}
