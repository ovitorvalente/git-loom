package cli

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/ovitorvalente/git-loom/internal/interfaces/mocks"
)

func TestCommitCommandDryRun(t *testing.T) {
	t.Parallel()

	output := &bytes.Buffer{}
	gitRepository := &mocks.GitRepository{
		GetDiffFunc: func() (string, error) {
			return "add commit command", nil
		},
	}
	command := newCommitCommandWithDependencies(commitDependencies{
		gitRepository: gitRepository,
		aiProvider:    &mocks.AIProvider{},
	})

	command.SetOut(output)
	command.SetErr(output)
	command.SetArgs([]string{"--dry-run"})

	err := command.Execute()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(gitRepository.CommitCalls) != 0 {
		t.Fatalf("expected no commit calls, got %d", len(gitRepository.CommitCalls))
	}
	if !strings.Contains(output.String(), "message: feat: add commit command") {
		t.Fatalf("unexpected output: %q", output.String())
	}
}

func TestCommitCommandUsesConfiguredScope(t *testing.T) {
	t.Parallel()

	output := &bytes.Buffer{}
	gitRepository := &mocks.GitRepository{
		GetDiffFunc: func() (string, error) {
			return "add config support", nil
		},
	}
	command := newCommitCommandWithDependencies(commitDependencies{
		gitRepository: gitRepository,
		aiProvider:    &mocks.AIProvider{},
		config: commitConfig{
			DefaultScope: "core",
		},
	})

	command.SetOut(output)
	command.SetErr(output)
	command.SetArgs([]string{"--dry-run"})

	err := command.Execute()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !strings.Contains(output.String(), "message: feat(core): add config support") {
		t.Fatalf("unexpected output: %q", output.String())
	}
}

func TestCommitCommandYesFlag(t *testing.T) {
	t.Parallel()

	output := &bytes.Buffer{}
	gitRepository := &mocks.GitRepository{
		GetDiffFunc: func() (string, error) {
			return "fix commit execution", nil
		},
	}
	command := newCommitCommandWithDependencies(commitDependencies{
		gitRepository: gitRepository,
		aiProvider:    &mocks.AIProvider{},
	})

	command.SetOut(output)
	command.SetErr(output)
	command.SetIn(strings.NewReader("y\n"))
	command.SetArgs([]string{"--yes"})

	err := command.Execute()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(gitRepository.CommitCalls) != 1 {
		t.Fatalf("expected one commit call, got %d", len(gitRepository.CommitCalls))
	}
	if gitRepository.CommitCalls[0] != "fix: fix commit execution" {
		t.Fatalf("unexpected commit message: %q", gitRepository.CommitCalls[0])
	}
	if !strings.Contains(output.String(), "commit created") {
		t.Fatalf("unexpected output: %q", output.String())
	}
}

func TestCommitCommandConfirmsBeforeCommit(t *testing.T) {
	t.Parallel()

	output := &bytes.Buffer{}
	gitRepository := &mocks.GitRepository{
		GetDiffFunc: func() (string, error) {
			return "fix confirm flow", nil
		},
	}
	command := newCommitCommandWithDependencies(commitDependencies{
		gitRepository: gitRepository,
		aiProvider:    &mocks.AIProvider{},
	})

	command.SetOut(output)
	command.SetErr(output)
	command.SetIn(strings.NewReader("y\n"))

	err := command.Execute()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(gitRepository.CommitCalls) != 1 {
		t.Fatalf("expected one commit call, got %d", len(gitRepository.CommitCalls))
	}
	if !strings.Contains(output.String(), "create commit? [y/N]: ") {
		t.Fatalf("unexpected output: %q", output.String())
	}
}

func TestCommitCommandCancelsWithoutConfirmation(t *testing.T) {
	t.Parallel()

	output := &bytes.Buffer{}
	gitRepository := &mocks.GitRepository{
		GetDiffFunc: func() (string, error) {
			return "fix confirm flow", nil
		},
	}
	command := newCommitCommandWithDependencies(commitDependencies{
		gitRepository: gitRepository,
		aiProvider:    &mocks.AIProvider{},
	})

	command.SetOut(output)
	command.SetErr(output)
	command.SetIn(strings.NewReader("n\n"))

	err := command.Execute()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(gitRepository.CommitCalls) != 0 {
		t.Fatalf("expected no commit calls, got %d", len(gitRepository.CommitCalls))
	}
	if !strings.Contains(output.String(), "commit canceled") {
		t.Fatalf("unexpected output: %q", output.String())
	}
}

func TestCommitCommandPropagatesErrors(t *testing.T) {
	t.Parallel()

	expectedError := errors.New("diff failed")
	command := newCommitCommandWithDependencies(commitDependencies{
		gitRepository: &mocks.GitRepository{
			GetDiffFunc: func() (string, error) {
				return "", expectedError
			},
		},
		aiProvider: &mocks.AIProvider{},
	})

	err := command.Execute()
	if !errors.Is(err, expectedError) {
		t.Fatalf("expected %v, got %v", expectedError, err)
	}
}

func TestCommitCommandShowsHelpfulEmptyDiffError(t *testing.T) {
	t.Parallel()

	command := newCommitCommandWithDependencies(commitDependencies{
		gitRepository: &mocks.GitRepository{
			GetDiffFunc: func() (string, error) {
				return " ", nil
			},
		},
		aiProvider: &mocks.AIProvider{},
	})

	err := command.Execute()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "no staged changes found; run git add before gitloom commit" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRootCommandRegistersCommit(t *testing.T) {
	t.Parallel()

	command := newRootCommand()
	subCommand, _, err := command.Find([]string{"commit"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if subCommand == nil || subCommand.Name() != "commit" {
		t.Fatal("expected commit subcommand to be registered")
	}
}
