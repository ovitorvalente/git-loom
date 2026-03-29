package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestRepositoryWithRealGitRepository(t *testing.T) {
	t.Parallel()

	repositoryPath := t.TempDir()
	runGitCommand(t, repositoryPath, "init")
	runGitCommand(t, repositoryPath, "config", "user.name", "Git Loom")
	runGitCommand(t, repositoryPath, "config", "user.email", "gitloom@example.com")

	writeFile(t, repositoryPath, "file-a.txt", "linha a1\n")
	writeFile(t, repositoryPath, "file-b.txt", "linha b1\n")
	runGitCommand(t, repositoryPath, "add", "file-a.txt", "file-b.txt")
	runGitCommand(t, repositoryPath, "commit", "-m", "chore: commit inicial")

	writeFile(t, repositoryPath, "file-a.txt", "linha a1\nlinha a2\n")
	writeFile(t, repositoryPath, "file-b.txt", "linha b1\nlinha b2\n")

	repository := Repository{
		runCommand: func(name string, args ...string) ([]byte, error) {
			command := exec.Command(name, args...)
			command.Dir = repositoryPath
			return command.CombinedOutput()
		},
	}

	if err := repository.StageFiles([]string{"file-a.txt", "file-b.txt"}); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	stagedFiles, err := repository.ListStagedFiles()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !equalStringSlices(stagedFiles, []string{"file-a.txt", "file-b.txt"}) {
		t.Fatalf("unexpected staged files: %v", stagedFiles)
	}

	diff, err := repository.GetDiff("file-a.txt")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !strings.Contains(diff, "diff --git a/file-a.txt b/file-a.txt") {
		t.Fatalf("unexpected diff: %q", diff)
	}

	if err := repository.CommitPaths("feat(core): atualizar file a", []string{"file-a.txt"}); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	headMessage := strings.TrimSpace(runGitCommand(t, repositoryPath, "log", "-1", "--pretty=%s"))
	if headMessage != "feat(core): atualizar file a" {
		t.Fatalf("unexpected head message: %q", headMessage)
	}

	remainingStagedFiles, err := repository.ListStagedFiles()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !equalStringSlices(remainingStagedFiles, []string{"file-b.txt"}) {
		t.Fatalf("unexpected remaining staged files: %v", remainingStagedFiles)
	}
}

func runGitCommand(t *testing.T, repositoryPath string, args ...string) string {
	t.Helper()

	command := exec.Command("git", args...)
	command.Dir = repositoryPath
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s failed: %v\n%s", strings.Join(args, " "), err, string(output))
	}

	return string(output)
}

func writeFile(t *testing.T, repositoryPath string, name string, content string) {
	t.Helper()

	path := filepath.Join(repositoryPath, name)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func equalStringSlices(left []string, right []string) bool {
	if len(left) != len(right) {
		return false
	}

	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}

	return true
}
