package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/ovitorvalente/git-loom/internal/interfaces/mocks"
)

func TestAnalyzeCommandTextOutput(t *testing.T) {
	t.Parallel()

	output := &bytes.Buffer{}
	command := newAnalyzeCommandWithDependencies(commitDependencies{
		gitRepository: &mocks.GitRepository{
			ListStagedFilesFunc: func() ([]string, error) {
				return []string{"internal/cli/commit.go"}, nil
			},
			GetDiffFunc: func(paths ...string) (string, error) {
				return "diff --git a/internal/cli/commit.go b/internal/cli/commit.go\nnew file mode 100644\n", nil
			},
		},
		aiProvider: &mocks.AIProvider{},
	})

	command.SetOut(output)
	command.SetErr(output)
	command.SetArgs([]string{"--preview"})

	err := command.Execute()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !strings.Contains(output.String(), "feat") {
		t.Fatalf("unexpected output: %q", output.String())
	}
	if strings.Contains(output.String(), "commit criado:") {
		t.Fatalf("did not expect commit creation output, got %q", output.String())
	}
}

func TestAnalyzeCommandJSONOutput(t *testing.T) {
	t.Parallel()

	output := &bytes.Buffer{}
	command := newAnalyzeCommandWithDependencies(commitDependencies{
		gitRepository: &mocks.GitRepository{
			ListStagedFilesFunc: func() ([]string, error) {
				return []string{"internal/cli/commit.go"}, nil
			},
			GetDiffFunc: func(paths ...string) (string, error) {
				return "diff --git a/internal/cli/commit.go b/internal/cli/commit.go\nnew file mode 100644\n", nil
			},
		},
		aiProvider: &mocks.AIProvider{},
	})

	command.SetOut(output)
	command.SetErr(output)
	command.SetArgs([]string{"--json"})

	err := command.Execute()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	expectedParts := []string{
		"\"plans\"",
		"\"message\": \"feat(cli): adicionar fluxo de commit",
		"\"type\": \"feat\"",
		"\"summary\"",
	}
	for _, expectedPart := range expectedParts {
		if !strings.Contains(output.String(), expectedPart) {
			t.Fatalf("expected output to contain %q, got %q", expectedPart, output.String())
		}
	}
}

func TestAnalyzeCommandFocusFiltersPaths(t *testing.T) {
	t.Parallel()

	output := &bytes.Buffer{}
	command := newAnalyzeCommandWithDependencies(commitDependencies{
		gitRepository: &mocks.GitRepository{
			ListStagedFilesFunc: func() ([]string, error) {
				return []string{"internal/cli/commit.go", "internal/app/commit_service.go"}, nil
			},
			GetDiffFunc: func(paths ...string) (string, error) {
				lines := []string{}
				for _, path := range paths {
					lines = append(lines, "diff --git a/"+path+" b/"+path, "index 1111111..2222222 100644")
				}
				return strings.Join(lines, "\n"), nil
			},
		},
		aiProvider: &mocks.AIProvider{},
	})

	command.SetOut(output)
	command.SetErr(output)
	command.SetArgs([]string{"--focus", "internal/cli"})

	err := command.Execute()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if strings.Contains(output.String(), "internal/app/commit_service.go") {
		t.Fatalf("expected focus to filter app path, got %q", output.String())
	}
}

func TestAnalyzeCommandExplainShowsGroupingReason(t *testing.T) {
	t.Parallel()

	output := &bytes.Buffer{}
	command := newAnalyzeCommandWithDependencies(commitDependencies{
		gitRepository: &mocks.GitRepository{
			ListStagedFilesFunc: func() ([]string, error) {
				return []string{"internal/cli/commit.go"}, nil
			},
			GetDiffFunc: func(paths ...string) (string, error) {
				return "diff --git a/internal/cli/commit.go b/internal/cli/commit.go\nnew file mode 100644\n", nil
			},
		},
		aiProvider: &mocks.AIProvider{},
	})

	command.SetOut(output)
	command.SetErr(output)
	command.SetArgs([]string{"--explain"})

	err := command.Execute()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !strings.Contains(output.String(), "porque esse agrupamento:") {
		t.Fatalf("expected explain section, got %q", output.String())
	}
}
