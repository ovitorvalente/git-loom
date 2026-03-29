package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/ovitorvalente/git-loom/internal/interfaces/mocks"
)

func TestRootCommandHelp(t *testing.T) {
	t.Parallel()

	output := &bytes.Buffer{}
	command := newRootCommand()
	command.SetOut(output)
	command.SetErr(output)
	command.SetArgs([]string{"--help"})

	err := command.Execute()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	expectedParts := []string{
		"Git Loom automatiza commits semanticos com revisao antes de executar.",
		"Comandos:",
		"commit",
		"analyze",
		"config",
		"version",
		"Exemplos:",
		"gitloom help commit",
	}

	for _, expectedPart := range expectedParts {
		if !strings.Contains(output.String(), expectedPart) {
			t.Fatalf("expected help output to contain %q, got %q", expectedPart, output.String())
		}
	}
}

func TestCommitCommandHelp(t *testing.T) {
	t.Parallel()

	output := &bytes.Buffer{}
	command := newCommitCommandWithDependencies(commitDependencies{
		gitRepository: &mocks.GitRepository{},
		aiProvider:    &mocks.AIProvider{},
	})
	command.SetOut(output)
	command.SetErr(output)
	command.SetArgs([]string{"--help"})

	err := command.Execute()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	expectedParts := []string{
		"Planeja e cria commits semanticos a partir do estado atual do repositorio.",
		"--dry-run",
		"--preview",
		"--strict",
		"--verbose",
		"--json",
		"Config:",
		".gitloom.yaml",
	}

	for _, expectedPart := range expectedParts {
		if !strings.Contains(output.String(), expectedPart) {
			t.Fatalf("expected help output to contain %q, got %q", expectedPart, output.String())
		}
	}
}

func TestRootCommandRegistersCommitAlias(t *testing.T) {
	t.Parallel()

	command := newRootCommand()
	subCommand, _, err := command.Find([]string{"ci"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if subCommand == nil || subCommand.Name() != "commit" {
		t.Fatal("expected ci alias to resolve to commit")
	}
}
