package config

import "testing"

func TestDefaultConfig(t *testing.T) {
	t.Parallel()

	configuration := DefaultConfig()
	if configuration.Commit.Scope != "" {
		t.Fatalf("expected empty scope, got %q", configuration.Commit.Scope)
	}
	if configuration.CLI.AutoConfirm {
		t.Fatal("expected auto confirm to default to false")
	}
}
