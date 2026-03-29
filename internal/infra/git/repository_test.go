package git

import (
	"errors"
	"testing"
)

func TestRepositoryGetDiff(t *testing.T) {
	t.Parallel()

	repository := Repository{
		runCommand: func(name string, args ...string) ([]byte, error) {
			if name != "git" {
				t.Fatalf("expected git command, got %s", name)
			}
			if len(args) != 2 || args[0] != "diff" || args[1] != "--cached" {
				t.Fatalf("unexpected args: %v", args)
			}

			return []byte("diff --git a/file.go b/file.go\n"), nil
		},
	}

	diff, err := repository.GetDiff()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if diff != "diff --git a/file.go b/file.go" {
		t.Fatalf("unexpected diff: %q", diff)
	}
}

func TestRepositoryGetDiffWithPaths(t *testing.T) {
	t.Parallel()

	repository := Repository{
		runCommand: func(name string, args ...string) ([]byte, error) {
			expectedArgs := []string{"diff", "--cached", "--", "internal/cli/commit.go"}
			if name != "git" {
				t.Fatalf("expected git command, got %s", name)
			}
			if len(args) != len(expectedArgs) {
				t.Fatalf("unexpected args: %v", args)
			}
			for index, expectedArg := range expectedArgs {
				if args[index] != expectedArg {
					t.Fatalf("unexpected arg at %d: %v", index, args)
				}
			}

			return []byte("diff --git a/file.go b/file.go\n"), nil
		},
	}

	_, err := repository.GetDiff("internal/cli/commit.go")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestRepositoryCommit(t *testing.T) {
	t.Parallel()

	repository := Repository{
		runCommand: func(name string, args ...string) ([]byte, error) {
			if name != "git" {
				t.Fatalf("expected git command, got %s", name)
			}
			if len(args) != 3 || args[0] != "commit" || args[1] != "-m" || args[2] != "feat: add infra" {
				t.Fatalf("unexpected args: %v", args)
			}

			return []byte(""), nil
		},
	}

	err := repository.Commit("feat: add infra")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestRepositoryCommitPaths(t *testing.T) {
	t.Parallel()

	repository := Repository{
		runCommand: func(name string, args ...string) ([]byte, error) {
			expectedArgs := []string{"commit", "-m", "feat: add infra", "--", "internal/cli/commit.go"}
			if name != "git" {
				t.Fatalf("expected git command, got %s", name)
			}
			if len(args) != len(expectedArgs) {
				t.Fatalf("unexpected args: %v", args)
			}
			for index, expectedArg := range expectedArgs {
				if args[index] != expectedArg {
					t.Fatalf("unexpected arg at %d: %v", index, args)
				}
			}

			return []byte(""), nil
		},
	}

	err := repository.CommitPaths("feat: add infra", []string{"internal/cli/commit.go"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestRepositoryListFiles(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		invocation   func(repository Repository) ([]string, error)
		expectedArgs []string
	}{
		{
			name: "lists staged files",
			invocation: func(repository Repository) ([]string, error) {
				return repository.ListStagedFiles()
			},
			expectedArgs: []string{"diff", "--cached", "--name-only", "--diff-filter=ACMR"},
		},
		{
			name: "lists changed files",
			invocation: func(repository Repository) ([]string, error) {
				return repository.ListChangedFiles()
			},
			expectedArgs: nil,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			callIndex := 0
			repository := Repository{
				runCommand: func(name string, args ...string) ([]byte, error) {
					if name != "git" {
						t.Fatalf("expected git command, got %s", name)
					}
					if testCase.name == "lists changed files" {
						expectedCommands := [][]string{
							{"diff", "--name-only", "--diff-filter=MD"},
							{"ls-files", "--others", "--exclude-standard"},
						}
						expectedArgs := expectedCommands[callIndex]
						callIndex++
						if len(args) != len(expectedArgs) {
							t.Fatalf("unexpected args: %v", args)
						}
						for index, expectedArg := range expectedArgs {
							if args[index] != expectedArg {
								t.Fatalf("unexpected arg at %d: %v", index, args)
							}
						}

						if args[0] == "diff" {
							return []byte("a.go\n"), nil
						}

						return []byte("b.go\n"), nil
					}

					if len(args) != len(testCase.expectedArgs) {
						t.Fatalf("unexpected args: %v", args)
					}
					for index, expectedArg := range testCase.expectedArgs {
						if args[index] != expectedArg {
							t.Fatalf("unexpected arg at %d: %v", index, args)
						}
					}

					return []byte("a.go\nb.go\n"), nil
				},
			}

			files, err := testCase.invocation(repository)
			if err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}
			if len(files) != 2 {
				t.Fatalf("expected two files, got %d", len(files))
			}
		})
	}
}

func TestRepositoryStageFiles(t *testing.T) {
	t.Parallel()

	repository := Repository{
		runCommand: func(name string, args ...string) ([]byte, error) {
			expectedArgs := []string{"add", "-A", "--", "internal/cli/commit.go", "internal/app/commit_service.go"}
			if name != "git" {
				t.Fatalf("expected git command, got %s", name)
			}
			if len(args) != len(expectedArgs) {
				t.Fatalf("unexpected args: %v", args)
			}
			for index, expectedArg := range expectedArgs {
				if args[index] != expectedArg {
					t.Fatalf("unexpected arg at %d: %v", index, args)
				}
			}

			return []byte(""), nil
		},
	}

	err := repository.StageFiles([]string{"internal/cli/commit.go", "internal/app/commit_service.go"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestRepositoryStagesDeletedFiles(t *testing.T) {
	t.Parallel()

	repository := Repository{
		runCommand: func(name string, args ...string) ([]byte, error) {
			expectedArgs := []string{"add", "-A", "--", "configs/default.yaml"}
			if name != "git" {
				t.Fatalf("expected git command, got %s", name)
			}
			if len(args) != len(expectedArgs) {
				t.Fatalf("unexpected args: %v", args)
			}
			for index, expectedArg := range expectedArgs {
				if args[index] != expectedArg {
					t.Fatalf("unexpected arg at %d: %v", index, args)
				}
			}

			return []byte(""), nil
		},
	}

	err := repository.StageFiles([]string{"configs/default.yaml"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestRepositoryCreateBranch(t *testing.T) {
	t.Parallel()

	repository := Repository{
		runCommand: func(name string, args ...string) ([]byte, error) {
			if name != "git" {
				t.Fatalf("expected git command, got %s", name)
			}
			if len(args) != 3 || args[0] != "checkout" || args[1] != "-b" || args[2] != "feat/infra" {
				t.Fatalf("unexpected args: %v", args)
			}

			return []byte(""), nil
		},
	}

	err := repository.CreateBranch("feat/infra")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestRepositoryPropagatesCommandErrors(t *testing.T) {
	t.Parallel()

	expectedError := errors.New("git failed")
	repository := Repository{
		runCommand: func(name string, args ...string) ([]byte, error) {
			return nil, expectedError
		},
	}

	_, err := repository.GetDiff()
	if !errors.Is(err, expectedError) {
		t.Fatalf("expected %v, got %v", expectedError, err)
	}
}
