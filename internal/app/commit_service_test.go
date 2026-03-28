package app

import (
	"errors"
	"testing"

	domaincommit "github.com/ovitorvalente/git-loom/internal/domain/commit"
	"github.com/ovitorvalente/git-loom/internal/interfaces/mocks"
)

func TestCommitServiceGenerateCommit(t *testing.T) {
	t.Parallel()

	diffError := errors.New("diff failed")
	aiError := errors.New("ai failed")

	testCases := []struct {
		name              string
		gitRepository     *mocks.GitRepository
		aiProvider        *mocks.AIProvider
		expectedType      domaincommit.Type
		expectedMessage   string
		expectedDiff      string
		expectedAIInvokes int
		expectedError     error
	}{
		{
			name: "returns generated message without ai provider",
			gitRepository: &mocks.GitRepository{
				GetDiffFunc: func() (string, error) {
					return "add commit workflow support", nil
				},
			},
			expectedType:      domaincommit.TypeFeat,
			expectedMessage:   "feat: add commit workflow support",
			expectedDiff:      "add commit workflow support",
			expectedAIInvokes: 0,
		},
		{
			name: "returns ai generated message when available",
			gitRepository: &mocks.GitRepository{
				GetDiffFunc: func() (string, error) {
					return "fix generator regression", nil
				},
			},
			aiProvider: &mocks.AIProvider{
				GenerateCommitFunc: func(diff string) (string, error) {
					if diff != "fix generator regression" {
						t.Fatalf("unexpected diff: %s", diff)
					}

					return "fix(core): correct generator regression", nil
				},
			},
			expectedType:      domaincommit.TypeFix,
			expectedMessage:   "fix(core): correct generator regression",
			expectedDiff:      "fix generator regression",
			expectedAIInvokes: 1,
		},
		{
			name: "propagates git diff errors",
			gitRepository: &mocks.GitRepository{
				GetDiffFunc: func() (string, error) {
					return "", diffError
				},
			},
			expectedError: diffError,
		},
		{
			name: "propagates ai provider errors",
			gitRepository: &mocks.GitRepository{
				GetDiffFunc: func() (string, error) {
					return "add commit workflow support", nil
				},
			},
			aiProvider: &mocks.AIProvider{
				GenerateCommitFunc: func(diff string) (string, error) {
					return "", aiError
				},
			},
			expectedError: aiError,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			service := NewCommitService(testCase.gitRepository, testCase.aiProvider)

			result, err := service.GenerateCommit()
			if !errors.Is(err, testCase.expectedError) {
				t.Fatalf("expected error %v, got %v", testCase.expectedError, err)
			}
			if testCase.expectedError != nil {
				return
			}
			if result.Diff != testCase.expectedDiff {
				t.Fatalf("expected diff %q, got %q", testCase.expectedDiff, result.Diff)
			}
			if result.Commit.Type != testCase.expectedType {
				t.Fatalf("expected type %s, got %s", testCase.expectedType, result.Commit.Type)
			}
			if result.Message != testCase.expectedMessage {
				t.Fatalf("expected message %q, got %q", testCase.expectedMessage, result.Message)
			}
			if len(testCase.gitRepository.CommitCalls) != 0 {
				t.Fatalf("expected no commit calls, got %d", len(testCase.gitRepository.CommitCalls))
			}
			if testCase.aiProvider != nil && len(testCase.aiProvider.GenerateCommitCalls) != testCase.expectedAIInvokes {
				t.Fatalf("expected %d ai calls, got %d", testCase.expectedAIInvokes, len(testCase.aiProvider.GenerateCommitCalls))
			}
		})
	}
}
