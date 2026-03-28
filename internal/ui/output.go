package ui

import (
	"fmt"
	"strings"

	"github.com/ovitorvalente/git-loom/internal/app"
)

func FormatCommitResult(result app.CommitResult) string {
	lines := []string{
		"generated commit",
		fmt.Sprintf("  type: %s", result.Commit.Type),
		fmt.Sprintf("  scope: %s", formatValue(result.Commit.Scope)),
		fmt.Sprintf("  description: %s", formatValue(result.Commit.Description)),
		fmt.Sprintf("  message: %s", result.Message),
	}

	if strings.TrimSpace(result.Commit.Body) != "" {
		lines = append(lines, fmt.Sprintf("  body: %s", result.Commit.Body))
	}

	return strings.Join(lines, "\n")
}

func formatValue(value string) string {
	trimmedValue := strings.TrimSpace(value)
	if trimmedValue == "" {
		return "-"
	}

	return trimmedValue
}
