package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestConfigInitCreatesDefaultFile(t *testing.T) {
	directory := t.TempDir()
	previousDirectory, getCwdErr := os.Getwd()
	if getCwdErr != nil {
		t.Fatalf("expected nil error, got %v", getCwdErr)
	}
	defer func() {
		_ = os.Chdir(previousDirectory)
	}()
	if chdirErr := os.Chdir(directory); chdirErr != nil {
		t.Fatalf("expected nil error, got %v", chdirErr)
	}

	output := &bytes.Buffer{}
	command := newConfigCommand()
	command.SetOut(output)
	command.SetErr(output)
	command.SetArgs([]string{"init"})

	execErr := command.Execute()
	if execErr != nil {
		t.Fatalf("expected nil error, got %v", execErr)
	}

	content, readErr := os.ReadFile(filepath.Join(directory, ".gitloom.yaml"))
	if readErr != nil {
		t.Fatalf("expected config file to be created, got %v", readErr)
	}
	if !strings.Contains(string(content), "auto_confirm: false") {
		t.Fatalf("unexpected config content: %q", string(content))
	}
}

func TestConfigInitFailsWhenFileExistsWithoutForce(t *testing.T) {
	directory := t.TempDir()
	previousDirectory, getCwdErr := os.Getwd()
	if getCwdErr != nil {
		t.Fatalf("expected nil error, got %v", getCwdErr)
	}
	defer func() {
		_ = os.Chdir(previousDirectory)
	}()
	if chdirErr := os.Chdir(directory); chdirErr != nil {
		t.Fatalf("expected nil error, got %v", chdirErr)
	}

	if writeErr := os.WriteFile(".gitloom.yaml", []byte("commit:\n  scope: core\n"), 0o600); writeErr != nil {
		t.Fatalf("expected nil error, got %v", writeErr)
	}

	command := newConfigCommand()
	command.SetArgs([]string{"init"})

	execErr := command.Execute()
	if execErr == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(execErr.Error(), ".gitloom.yaml ja existe") {
		t.Fatalf("unexpected error: %v", execErr)
	}
}

func TestConfigInitOverwritesWithForce(t *testing.T) {
	directory := t.TempDir()
	previousDirectory, getCwdErr := os.Getwd()
	if getCwdErr != nil {
		t.Fatalf("expected nil error, got %v", getCwdErr)
	}
	defer func() {
		_ = os.Chdir(previousDirectory)
	}()
	if chdirErr := os.Chdir(directory); chdirErr != nil {
		t.Fatalf("expected nil error, got %v", chdirErr)
	}

	if writeErr := os.WriteFile(".gitloom.yaml", []byte("old"), 0o600); writeErr != nil {
		t.Fatalf("expected nil error, got %v", writeErr)
	}

	command := newConfigCommand()
	command.SetArgs([]string{"init", "--force"})

	execErr := command.Execute()
	if execErr != nil {
		t.Fatalf("expected nil error, got %v", execErr)
	}

	content, readErr := os.ReadFile(".gitloom.yaml")
	if readErr != nil {
		t.Fatalf("expected nil error, got %v", readErr)
	}
	if string(content) == "old" {
		t.Fatalf("expected config file to be overwritten, got %q", string(content))
	}
}
