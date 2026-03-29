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
		ListStagedFilesFunc: func() ([]string, error) {
			return []string{"internal/cli/commit.go"}, nil
		},
		GetDiffFunc: func(paths ...string) (string, error) {
			return "diff --git a/internal/cli/commit.go b/internal/cli/commit.go\nnew file mode 100644\n", nil
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
	if len(gitRepository.CommitPathsCalls) != 0 {
		t.Fatalf("expected no commit calls, got %d", len(gitRepository.CommitPathsCalls))
	}
	if !strings.Contains(output.String(), "◆ feat(cli)") {
		t.Fatalf("unexpected output: %q", output.String())
	}
	if !strings.Contains(output.String(), "mensagem:") {
		t.Fatalf("unexpected output: %q", output.String())
	}
	if !strings.Contains(output.String(), "feat(cli): adicionar fluxo de commit") {
		t.Fatalf("unexpected output: %q", output.String())
	}
	if !strings.Contains(output.String(), "[") {
		t.Fatalf("unexpected output: %q", output.String())
	}
}

func TestCommitCommandStagesChangedFilesWhenConfirmed(t *testing.T) {
	t.Parallel()

	output := &bytes.Buffer{}
	gitRepository := &mocks.GitRepository{
		ListStagedFilesFunc: func() ([]string, error) {
			return []string{"internal/cli/commit.go"}, nil
		},
		ListChangedFilesFunc: func() ([]string, error) {
			return []string{"internal/app/commit_service.go"}, nil
		},
		GetDiffFunc: func(paths ...string) (string, error) {
			if len(paths) == 1 && paths[0] == "internal/app/commit_service.go" {
				return "diff --git a/internal/app/commit_service.go b/internal/app/commit_service.go\nindex 1111111..2222222 100644\n", nil
			}
			return "diff --git a/internal/cli/commit.go b/internal/cli/commit.go\nnew file mode 100644\n", nil
		},
	}
	command := newCommitCommandWithDependencies(commitDependencies{
		gitRepository: gitRepository,
		aiProvider:    &mocks.AIProvider{},
	})

	command.SetOut(output)
	command.SetErr(output)
	command.SetIn(strings.NewReader("y\nn\n"))
	command.SetArgs([]string{"--dry-run"})

	err := command.Execute()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(gitRepository.StageFilesCalls) != 1 {
		t.Fatalf("expected one stage call, got %d", len(gitRepository.StageFilesCalls))
	}
	if !strings.Contains(output.String(), "arquivos alterados") {
		t.Fatalf("unexpected output: %q", output.String())
	}
}

func TestCommitCommandSplitsCommitsInBlocksOfFour(t *testing.T) {
	t.Parallel()

	output := &bytes.Buffer{}
	paths := []string{
		"internal/cli/a.go",
		"internal/cli/b.go",
		"internal/cli/c.go",
		"internal/cli/d.go",
		"internal/cli/e.go",
	}
	gitRepository := &mocks.GitRepository{
		ListChangedFilesFunc: func() ([]string, error) {
			return nil, nil
		},
		ListStagedFilesFunc: func() ([]string, error) {
			return paths, nil
		},
		GetDiffFunc: func(requestedPaths ...string) (string, error) {
			lines := []string{}
			for _, path := range requestedPaths {
				lines = append(lines, "diff --git a/"+path+" b/"+path, "index 1111111..2222222 100644")
			}
			return strings.Join(lines, "\n"), nil
		},
	}
	command := newCommitCommandWithDependencies(commitDependencies{
		gitRepository: gitRepository,
		aiProvider:    &mocks.AIProvider{},
	})

	command.SetOut(output)
	command.SetErr(output)
	command.SetArgs([]string{"--yes"})

	err := command.Execute()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(gitRepository.CommitPathsCalls) != 2 {
		t.Fatalf("expected two commit blocks, got %d", len(gitRepository.CommitPathsCalls))
	}
	if len(gitRepository.CommitPathsCalls[0].Paths) != 4 {
		t.Fatalf("expected first block with four files, got %d", len(gitRepository.CommitPathsCalls[0].Paths))
	}
	if len(gitRepository.CommitPathsCalls[1].Paths) != 1 {
		t.Fatalf("expected second block with one file, got %d", len(gitRepository.CommitPathsCalls[1].Paths))
	}
	if !strings.Contains(output.String(), "✔ 2 commits criados") {
		t.Fatalf("unexpected output: %q", output.String())
	}
	if !strings.Contains(output.String(), "qualidade media:") {
		t.Fatalf("unexpected output: %q", output.String())
	}
}

func TestCommitCommandPreviewShowsDiffImpact(t *testing.T) {
	t.Parallel()

	output := &bytes.Buffer{}
	gitRepository := &mocks.GitRepository{
		ListStagedFilesFunc: func() ([]string, error) {
			return []string{"internal/cli/commit.go"}, nil
		},
		GetDiffFunc: func(paths ...string) (string, error) {
			return strings.Join([]string{
				"diff --git a/internal/cli/commit.go b/internal/cli/commit.go",
				"new file mode 100644",
				"+++ b/internal/cli/commit.go",
				"+package cli",
				"+func run() {}",
			}, "\n"), nil
		},
	}

	command := newCommitCommandWithDependencies(commitDependencies{
		gitRepository: gitRepository,
		aiProvider:    &mocks.AIProvider{},
	})

	command.SetOut(output)
	command.SetErr(output)
	command.SetArgs([]string{"--preview"})

	err := command.Execute()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !strings.Contains(output.String(), "preview:") {
		t.Fatalf("unexpected output: %q", output.String())
	}
	if !strings.Contains(output.String(), "linhas: +2 -0") {
		t.Fatalf("unexpected output: %q", output.String())
	}
}

func TestCommitCommandVerboseShowsExpandedDetails(t *testing.T) {
	t.Parallel()

	output := &bytes.Buffer{}
	gitRepository := &mocks.GitRepository{
		ListStagedFilesFunc: func() ([]string, error) {
			return []string{"internal/cli/commit.go", "internal/ui/renderer.go"}, nil
		},
		GetDiffFunc: func(paths ...string) (string, error) {
			return strings.Join([]string{
				"diff --git a/internal/cli/commit.go b/internal/cli/commit.go",
				"index 1111111..2222222 100644",
				"+++ b/internal/cli/commit.go",
				"+package cli",
				"diff --git a/internal/ui/renderer.go b/internal/ui/renderer.go",
				"new file mode 100644",
				"+++ b/internal/ui/renderer.go",
				"+package ui",
			}, "\n"), nil
		},
	}

	command := newCommitCommandWithDependencies(commitDependencies{
		gitRepository: gitRepository,
		aiProvider:    &mocks.AIProvider{},
	})

	command.SetOut(output)
	command.SetErr(output)
	command.SetArgs([]string{"--dry-run", "--verbose", "--preview"})

	err := command.Execute()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !strings.Contains(output.String(), "verbose:") {
		t.Fatalf("unexpected output: %q", output.String())
	}
	if !strings.Contains(output.String(), "impacto: +2 -0") {
		t.Fatalf("unexpected output: %q", output.String())
	}
	if !strings.Contains(output.String(), "internal/ui/renderer.go") {
		t.Fatalf("unexpected output: %q", output.String())
	}
}

func TestCommitCommandStrictFailsForLowQualityPlan(t *testing.T) {
	t.Parallel()

	command := newCommitCommandWithDependencies(commitDependencies{
		gitRepository: &mocks.GitRepository{
			ListStagedFilesFunc: func() ([]string, error) {
				return []string{"internal/domain/unknown.go"}, nil
			},
			GetDiffFunc: func(paths ...string) (string, error) {
				return "diff --git a/internal/domain/unknown.go b/internal/domain/unknown.go\nindex 1111111..2222222 100644\n", nil
			},
		},
		aiProvider: &mocks.AIProvider{},
	})
	command.SetArgs([]string{"--strict", "--dry-run"})

	err := command.Execute()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "modo estrito falhou") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCommitCommandCancelsWithoutConfirmation(t *testing.T) {
	t.Parallel()

	output := &bytes.Buffer{}
	gitRepository := &mocks.GitRepository{
		ListStagedFilesFunc: func() ([]string, error) {
			return []string{"internal/cli/commit.go"}, nil
		},
		GetDiffFunc: func(paths ...string) (string, error) {
			return "diff --git a/internal/cli/commit.go b/internal/cli/commit.go\nindex 1111111..2222222 100644\n", nil
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
	if len(gitRepository.CommitPathsCalls) != 0 {
		t.Fatalf("expected no commit calls, got %d", len(gitRepository.CommitPathsCalls))
	}
	if !strings.Contains(output.String(), "commit cancelado") {
		t.Fatalf("unexpected output: %q", output.String())
	}
}

func TestCommitCommandPropagatesErrors(t *testing.T) {
	t.Parallel()

	expectedError := errors.New("diff failed")
	command := newCommitCommandWithDependencies(commitDependencies{
		gitRepository: &mocks.GitRepository{
			ListStagedFilesFunc: func() ([]string, error) {
				return []string{"internal/cli/commit.go"}, nil
			},
			GetDiffFunc: func(paths ...string) (string, error) {
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
			ListStagedFilesFunc: func() ([]string, error) {
				return nil, nil
			},
		},
		aiProvider: &mocks.AIProvider{},
	})

	err := command.Execute()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "nenhuma mudanca staged encontrada; execute git add antes de gitloom commit" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCommitCommandRejectsPartiallyStagedFiles(t *testing.T) {
	t.Parallel()

	command := newCommitCommandWithDependencies(commitDependencies{
		gitRepository: &mocks.GitRepository{
			ListStagedFilesFunc: func() ([]string, error) {
				return []string{"internal/cli/commit.go"}, nil
			},
			ListChangedFilesFunc: func() ([]string, error) {
				return []string{"internal/cli/commit.go"}, nil
			},
		},
		aiProvider: &mocks.AIProvider{},
	})

	err := command.Execute()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "arquivos parcialmente staged ainda nao sao suportados neste fluxo automatico; finalize ou descarte as mudancas unstaged antes de continuar" {
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
