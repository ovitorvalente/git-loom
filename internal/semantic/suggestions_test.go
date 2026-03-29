package semantic

import "testing"

func TestSuggestScope(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name               string
		scope              string
		files              []ChangedFile
		expectGeneric      bool
		expectAlternatives int
	}{
		{
			name:  "non-generic scope returns as-is",
			scope: "config",
			files: []ChangedFile{{Path: ".gitignore"}},
		},
		{
			name:               "generic scope gitignore suggests alternatives",
			scope:              "gitignore",
			files:              []ChangedFile{{Path: ".gitignore"}},
			expectGeneric:      true,
			expectAlternatives: 1,
		},
		{
			name:               "empty scope is generic",
			scope:              "",
			files:              []ChangedFile{{Path: "internal/cli/commit.go"}},
			expectGeneric:      true,
			expectAlternatives: 1,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result := SuggestScope(tc.scope, tc.files)
			if result.Generic != tc.expectGeneric {
				t.Fatalf("expected Generic=%v, got %v", tc.expectGeneric, result.Generic)
			}
			if tc.expectAlternatives > 0 && len(result.Alternatives) < tc.expectAlternatives {
				t.Fatalf("expected at least %d alternatives, got %d", tc.expectAlternatives, len(result.Alternatives))
			}
		})
	}
}

func TestIsGenericScope(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scope    string
		expected bool
	}{
		{"", true},
		{"gitignore", true},
		{"core", true},
		{"cli", true},
		{"config", false},
		{"commit_service", false},
		{"renderer", false},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scope, func(t *testing.T) {
			t.Parallel()

			if got := IsGenericScope(tc.scope); got != tc.expected {
				t.Fatalf("IsGenericScope(%q) = %v, want %v", tc.scope, got, tc.expected)
			}
		})
	}
}

func TestIsGenericDescription(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		desc     string
		expected bool
	}{
		{"atualizar", true},
		{"ajustar", true},
		{"atualizar projeto", true},
		{"adicionar suporte a retry", false},
		{"corrigir race condition no worker", false},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			if got := IsGenericDescription(tc.desc); got != tc.expected {
				t.Fatalf("IsGenericDescription(%q) = %v, want %v", tc.desc, got, tc.expected)
			}
		})
	}
}
