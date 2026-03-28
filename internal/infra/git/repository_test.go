package git

import (
	"errors"
	"testing"
)

func TestRepositoryGetDiff(t *testing.T) {
	t.Parallel()

	repository := Repository{
		runCommand: func(name string, args ...string) ([]byte, error) {
			if name != "git" {
				t.Fatalf("expected git command, got %s", name)
			}
			if len(args) != 2 || args[0] != "diff" || args[1] != "--cached" {
				t.Fatalf("unexpected args: %v", args)
			}

			return []byte("diff --git a/file.go b/file.go\n"), nil
		},
	}

	diff, err := repository.GetDiff()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if diff != "diff --git a/file.go b/file.go" {
		t.Fatalf("unexpected diff: %q", diff)
	}
}

func TestRepositoryCommit(t *testing.T) {
	t.Parallel()

	repository := Repository{
		runCommand: func(name string, args ...string) ([]byte, error) {
			if name != "git" {
				t.Fatalf("expected git command, got %s", name)
			}
			if len(args) != 3 || args[0] != "commit" || args[1] != "-m" || args[2] != "feat: add infra" {
				t.Fatalf("unexpected args: %v", args)
			}

			return []byte(""), nil
		},
	}

	err := repository.Commit("feat: add infra")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestRepositoryCreateBranch(t *testing.T) {
	t.Parallel()

	repository := Repository{
		runCommand: func(name string, args ...string) ([]byte, error) {
			if name != "git" {
				t.Fatalf("expected git command, got %s", name)
			}
			if len(args) != 3 || args[0] != "checkout" || args[1] != "-b" || args[2] != "feat/infra" {
				t.Fatalf("unexpected args: %v", args)
			}

			return []byte(""), nil
		},
	}

	err := repository.CreateBranch("feat/infra")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestRepositoryPropagatesCommandErrors(t *testing.T) {
	t.Parallel()

	expectedError := errors.New("git failed")
	repository := Repository{
		runCommand: func(name string, args ...string) ([]byte, error) {
			return nil, expectedError
		},
	}

	_, err := repository.GetDiff()
	if !errors.Is(err, expectedError) {
		t.Fatalf("expected %v, got %v", expectedError, err)
	}
}
