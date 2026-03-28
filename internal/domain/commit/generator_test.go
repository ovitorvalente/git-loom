package commit

import (
	"errors"
	"testing"
)

func TestGenerateMessage(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		model       Model
		expected    string
		expectedErr error
	}{
		{
			name: "generates message with scope and body",
			model: Model{
				Type:        TypeFeat,
				Scope:       "core",
				Description: "add commit generator",
				Body:        "keep message formatting deterministic",
			},
			expected: "feat(core): add commit generator\n\nkeep message formatting deterministic",
		},
		{
			name: "generates message without scope",
			model: Model{
				Type:        TypeFix,
				Description: "handle empty diff",
			},
			expected: "fix: handle empty diff",
		},
		{
			name: "defaults unknown type to chore",
			model: Model{
				Type:        Type("unknown"),
				Description: "update generated files",
			},
			expected: "chore: update generated files",
		},
		{
			name: "returns error for empty description",
			model: Model{
				Type:  TypeFeat,
				Scope: "core",
			},
			expectedErr: ErrEmptyDescription,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			message, err := GenerateMessage(testCase.model)
			if !errors.Is(err, testCase.expectedErr) {
				t.Fatalf("expected error %v, got %v", testCase.expectedErr, err)
			}
			if message != testCase.expected {
				t.Fatalf("expected %q, got %q", testCase.expected, message)
			}
		})
	}
}
