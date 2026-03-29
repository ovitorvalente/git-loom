package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestVersionCommand(t *testing.T) {
	t.Parallel()

	output := &bytes.Buffer{}
	command := newVersionCommand()
	command.SetOut(output)
	command.SetErr(output)

	previousVersion := Version
	previousCommit := GitCommit
	previousBuildDate := BuildDate
	Version = "1.2.3"
	GitCommit = "abc123"
	BuildDate = "2026-03-28"
	defer func() {
		Version = previousVersion
		GitCommit = previousCommit
		BuildDate = previousBuildDate
	}()

	err := command.Execute()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	expectedParts := []string{"gitloom 1.2.3", "commit: abc123", "build date: 2026-03-28"}
	for _, expectedPart := range expectedParts {
		if !strings.Contains(output.String(), expectedPart) {
			t.Fatalf("expected output to contain %q, got %q", expectedPart, output.String())
		}
	}
}
