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
	previousDirectory, err := os.Getwd()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	defer func() {
		_ = os.Chdir(previousDirectory)
	}()
	if err := os.Chdir(directory); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	output := &bytes.Buffer{}
	command := newConfigCommand()
	command.SetOut(output)
	command.SetErr(output)
	command.SetArgs([]string{"init"})

	err = command.Execute()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	content, err := os.ReadFile(filepath.Join(directory, ".gitloom.yaml"))
	if err != nil {
		t.Fatalf("expected config file to be created, got %v", err)
	}
	if !strings.Contains(string(content), "auto_confirm: false") {
		t.Fatalf("unexpected config content: %q", string(content))
	}
}

func TestConfigInitFailsWhenFileExistsWithoutForce(t *testing.T) {
	directory := t.TempDir()
	previousDirectory, err := os.Getwd()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	defer func() {
		_ = os.Chdir(previousDirectory)
	}()
	if err := os.Chdir(directory); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if err := os.WriteFile(".gitloom.yaml", []byte("commit:\n  scope: core\n"), 0o600); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	command := newConfigCommand()
	command.SetArgs([]string{"init"})

	err = command.Execute()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), ".gitloom.yaml ja existe") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestConfigInitOverwritesWithForce(t *testing.T) {
	directory := t.TempDir()
	previousDirectory, err := os.Getwd()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	defer func() {
		_ = os.Chdir(previousDirectory)
	}()
	if err := os.Chdir(directory); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if err := os.WriteFile(".gitloom.yaml", []byte("old"), 0o600); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	command := newConfigCommand()
	command.SetArgs([]string{"init", "--force"})

	err = command.Execute()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	content, err := os.ReadFile(".gitloom.yaml")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if string(content) == "old" {
		t.Fatalf("expected config file to be overwritten, got %q", string(content))
	}
}
