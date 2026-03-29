package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/ovitorvalente/git-loom/internal/interfaces/mocks"
)

func TestDoctorCommandTextOutput(t *testing.T) {
	t.Parallel()

	output := &bytes.Buffer{}
	command := newDoctorCommandWithDependencies(commitDependencies{
		gitRepository: &mocks.GitRepository{
			IsRepositoryFunc: func() (bool, error) { return true, nil },
			ListStagedFilesFunc: func() ([]string, error) {
				return []string{"internal/cli/commit.go"}, nil
			},
			ListChangedFilesFunc: func() ([]string, error) {
				return []string{"internal/ui/renderer.go"}, nil
			},
		},
	})
	command.SetOut(output)
	command.SetErr(output)

	err := command.Execute()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	expectedParts := []string{
		"◆ doctor",
		"repositorio: repositorio git detectado",
		"working-tree:",
		"changes: 1",
		"status geral: warn",
	}
	for _, expectedPart := range expectedParts {
		if !strings.Contains(output.String(), expectedPart) {
			t.Fatalf("expected output to contain %q, got %q", expectedPart, output.String())
		}
	}
}

func TestDoctorCommandJSONOutput(t *testing.T) {
	t.Parallel()

	output := &bytes.Buffer{}
	command := newDoctorCommandWithDependencies(commitDependencies{
		gitRepository: &mocks.GitRepository{
			IsRepositoryFunc: func() (bool, error) { return false, nil },
		},
	})
	command.SetOut(output)
	command.SetErr(output)
	command.SetArgs([]string{"--json"})

	err := command.Execute()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !strings.Contains(output.String(), "\"status\": \"fail\"") {
		t.Fatalf("unexpected output: %q", output.String())
	}
	if !strings.Contains(output.String(), "\"name\": \"repositorio\"") {
		t.Fatalf("unexpected output: %q", output.String())
	}
}
