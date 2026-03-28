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
				Description: "atualizar arquivos gerados",
			},
			expected: "chore: atualizar arquivos gerados",
		},
		{
			name: "returns error for empty description",
			model: Model{
				Type:  TypeFeat,
				Scope: "core",
			},
			expectedErr: ErrEmptyDescription,
		},
		{
			name: "limits header length",
			model: Model{
				Type:        TypeFeat,
				Scope:       "commit",
				Description: "adicionar suporte para gerar mensagens estruturadas a partir de arquivos alterados",
			},
			expected: "feat(commit): adicionar suporte para gerar mensagens estruturadas a",
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
