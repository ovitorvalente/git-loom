package app

import (
	"errors"
	"strings"
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
		options           GenerateCommitOptions
		expectedType      domaincommit.Type
		expectedMessage   string
		expectedDiff      string
		expectedAIInvokes int
		expectedError     error
	}{
		{
			name: "returns generated message without ai provider",
			gitRepository: &mocks.GitRepository{
				GetDiffFunc: func(paths ...string) (string, error) {
					return "diff --git a/internal/cli/commit.go b/internal/cli/commit.go\nnew file mode 100644\n", nil
				},
			},
			expectedType:      domaincommit.TypeFeat,
			expectedMessage:   "feat(cli): adicionar fluxo de commit\n\n- adiciona comando commit em cli",
			expectedDiff:      "diff --git a/internal/cli/commit.go b/internal/cli/commit.go\nnew file mode 100644\n",
			expectedAIInvokes: 0,
		},
		{
			name: "applies configured scope",
			gitRepository: &mocks.GitRepository{
				GetDiffFunc: func(paths ...string) (string, error) {
					return "diff --git a/internal/app/commit_service.go b/internal/app/commit_service.go\nindex 1111111..2222222 100644\n", nil
				},
			},
			options: GenerateCommitOptions{
				Scope: "core",
			},
			expectedType:      domaincommit.TypeRefactor,
			expectedMessage:   "refactor(core): refinar planejamento de commits\n\n- atualiza commit service em app",
			expectedDiff:      "diff --git a/internal/app/commit_service.go b/internal/app/commit_service.go\nindex 1111111..2222222 100644\n",
			expectedAIInvokes: 0,
		},
		{
			name: "returns ai generated message when available",
			gitRepository: &mocks.GitRepository{
				GetDiffFunc: func(paths ...string) (string, error) {
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
				GetDiffFunc: func(paths ...string) (string, error) {
					return "", diffError
				},
			},
			expectedError: diffError,
		},
		{
			name: "returns error when diff is empty",
			gitRepository: &mocks.GitRepository{
				GetDiffFunc: func(paths ...string) (string, error) {
					return " \n\t", nil
				},
			},
			expectedError: ErrEmptyDiff,
		},
		{
			name: "propagates ai provider errors",
			gitRepository: &mocks.GitRepository{
				GetDiffFunc: func(paths ...string) (string, error) {
					return "diff --git a/internal/app/commit_service.go b/internal/app/commit_service.go\nindex 1111111..2222222 100644\n", nil
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

			result, err := service.GenerateCommit(testCase.options)
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

func TestCommitServicePlanCommits(t *testing.T) {
	t.Parallel()

	service := NewCommitService(&mocks.GitRepository{
		GetDiffFunc: func(paths ...string) (string, error) {
			lines := []string{}
			for _, path := range paths {
				lines = append(lines, "diff --git a/"+path+" b/"+path, "index 1111111..2222222 100644")
			}
			return strings.Join(lines, "\n"), nil
		},
	}, &mocks.AIProvider{})

	review, err := service.PlanCommits([]string{
		"internal/cli/a.go",
		"internal/cli/b.go",
		"internal/cli/c.go",
		"internal/cli/d.go",
		"internal/cli/e.go",
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(review.Plans) != 2 {
		t.Fatalf("expected two plans, got %d", len(review.Plans))
	}
	if len(review.Plans[0].Result.Paths) != 4 {
		t.Fatalf("expected first plan with four files, got %d", len(review.Plans[0].Result.Paths))
	}
	if len(review.Plans[1].Result.Paths) != 1 {
		t.Fatalf("expected second plan with one file, got %d", len(review.Plans[1].Result.Paths))
	}
	if review.Plans[0].Quality.Score == 0 {
		t.Fatal("expected quality score to be calculated")
	}
}
