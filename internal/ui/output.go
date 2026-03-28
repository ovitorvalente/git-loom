package ui

import (
	"fmt"
	"strings"

	"github.com/ovitorvalente/git-loom/internal/app"
)

func FormatCommitResult(result app.CommitResult) string {
	header, body := splitCommitMessage(result.Message)
	if body == "" {
		body = strings.TrimSpace(result.Commit.Body)
	}

	lines := []string{
		"commit gerado",
		fmt.Sprintf("  tipo: %s", result.Commit.Type),
		fmt.Sprintf("  escopo: %s", formatValue(result.Commit.Scope)),
		fmt.Sprintf("  descricao: %s", formatValue(result.Commit.Description)),
		fmt.Sprintf("  mensagem: %s", header),
	}

	if body != "" {
		lines = append(lines, fmt.Sprintf("  detalhes:\n%s", indentBody(body)))
	}

	return strings.Join(lines, "\n")
}

func FormatCommitPlan(index int, total int, result app.CommitResult) string {
	lines := []string{
		fmt.Sprintf("bloco %d/%d", index, total),
		FormatCommitResult(result),
	}

	if len(result.Paths) > 0 {
		lines = append(lines, fmt.Sprintf("  arquivos:\n%s", indentBody(formatPaths(result.Paths))))
	}

	return strings.Join(lines, "\n")
}

func FormatChangedFiles(paths []string) string {
	return "arquivos em changes:\n" + indentBody(formatPaths(paths))
}

func formatValue(value string) string {
	trimmedValue := strings.TrimSpace(value)
	if trimmedValue == "" {
		return "-"
	}

	return trimmedValue
}

func splitCommitMessage(message string) (string, string) {
	parts := strings.SplitN(strings.TrimSpace(message), "\n\n", 2)
	if len(parts) == 1 {
		return parts[0], ""
	}

	return parts[0], parts[1]
}

func indentBody(body string) string {
	lines := strings.Split(body, "\n")
	for index, line := range lines {
		lines[index] = "    " + line
	}

	return strings.Join(lines, "\n")
}

func formatPaths(paths []string) string {
	lines := make([]string, 0, len(paths))
	for _, path := range paths {
		lines = append(lines, "- "+path)
	}

	return strings.Join(lines, "\n")
}
