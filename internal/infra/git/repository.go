package git

import (
	"os/exec"
	"strings"
)

type commandRunner func(name string, args ...string) ([]byte, error)

type Repository struct {
	runCommand commandRunner
}

func NewRepository() Repository {
	return Repository{
		runCommand: executeCommand,
	}
}

func (repository Repository) GetDiff(paths ...string) (string, error) {
	args := []string{"diff", "--cached"}
	if len(paths) > 0 {
		args = append(args, "--")
		args = append(args, paths...)
	}

	output, err := repository.run("git", args...)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

func (repository Repository) ListStagedFiles() ([]string, error) {
	return repository.listFiles("diff", "--cached", "--name-only", "--diff-filter=ACMR")
}

func (repository Repository) ListChangedFiles() ([]string, error) {
	return repository.listFiles("diff", "--name-only", "--diff-filter=M")
}

func (repository Repository) StageFiles(paths []string) error {
	if len(paths) == 0 {
		return nil
	}

	args := []string{"add"}
	args = append(args, paths...)
	_, err := repository.run("git", args...)
	return err
}

func (repository Repository) Commit(message string) error {
	_, err := repository.run("git", "commit", "-m", message)
	return err
}

func (repository Repository) CommitPaths(message string, paths []string) error {
	args := []string{"commit", "-m", message}
	if len(paths) > 0 {
		args = append(args, "--")
		args = append(args, paths...)
	}

	_, err := repository.run("git", args...)
	return err
}

func (repository Repository) CreateBranch(name string) error {
	_, err := repository.run("git", "checkout", "-b", name)
	return err
}

func (repository Repository) listFiles(args ...string) ([]string, error) {
	output, err := repository.run("git", args...)
	if err != nil {
		return nil, err
	}

	return splitLines(string(output)), nil
}

func (repository Repository) run(name string, args ...string) ([]byte, error) {
	if repository.runCommand != nil {
		return repository.runCommand(name, args...)
	}

	return executeCommand(name, args...)
}

func executeCommand(name string, args ...string) ([]byte, error) {
	command := exec.Command(name, args...)
	return command.CombinedOutput()
}

func splitLines(content string) []string {
	lines := strings.Split(strings.TrimSpace(content), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return nil
	}

	result := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" {
			continue
		}
		result = append(result, trimmedLine)
	}

	return result
}
