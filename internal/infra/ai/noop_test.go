package ai

import (
	"testing"

	"github.com/ovitorvalente/git-loom/internal/interfaces"
)

func TestNoopProviderImplementsInterface(t *testing.T) {
	t.Parallel()

	var _ interfaces.AIProvider = NewNoopProvider()
}

func TestNoopProviderGenerateCommit(t *testing.T) {
	t.Parallel()

	provider := NewNoopProvider()

	message, err := provider.GenerateCommit("diff content")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if message != "" {
		t.Fatalf("expected empty message, got %q", message)
	}
}
