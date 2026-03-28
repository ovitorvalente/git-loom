package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadReturnsDefaultsWhenFileDoesNotExist(t *testing.T) {
	t.Parallel()

	configuration, err := Load(filepath.Join(t.TempDir(), ".gitloom.yaml"))
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if configuration.Commit.Scope != "" {
		t.Fatalf("expected empty default scope, got %q", configuration.Commit.Scope)
	}
	if configuration.CLI.AutoConfirm {
		t.Fatal("expected auto confirm to default to false")
	}
}

func TestLoadParsesSupportedGitloomConfig(t *testing.T) {
	t.Parallel()

	directory := t.TempDir()
	path := filepath.Join(directory, ".gitloom.yaml")
	content := "commit:\n  scope: core\ncli:\n  auto_confirm: true\n"

	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	configuration, err := Load(path)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if configuration.Commit.Scope != "core" {
		t.Fatalf("expected scope core, got %q", configuration.Commit.Scope)
	}
	if !configuration.CLI.AutoConfirm {
		t.Fatal("expected auto confirm to be true")
	}
}
