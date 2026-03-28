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

func (repository Repository) GetDiff() (string, error) {
	output, err := repository.run("git", "diff", "--cached")
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

func (repository Repository) Commit(message string) error {
	_, err := repository.run("git", "commit", "-m", message)
	return err
}

func (repository Repository) CreateBranch(name string) error {
	_, err := repository.run("git", "checkout", "-b", name)
	return err
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
