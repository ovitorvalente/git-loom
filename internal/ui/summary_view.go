package ui

import (
	"fmt"
	"strings"

	"github.com/ovitorvalente/git-loom/internal/app"
)

func (renderer Renderer) ChangedFiles(stagedPaths []string, changedPaths []string) string {
	lines := []string{
		colorizeLine(borderColor, horizontalRule()),
		colorizeLine(headerColor, "◆ arquivos alterados"),
	}

	if len(stagedPaths) > 0 {
		lines = append(lines, "", renderer.sectionTitle("staged"))
		for _, path := range stagedPaths {
			lines = append(lines, renderer.bulletLine("•", path))
		}
	}

	if len(changedPaths) > 0 {
		lines = append(lines, "", renderer.sectionTitle("changes"))
		for _, path := range changedPaths {
			lines = append(lines, renderer.bulletLine("•", path))
		}
	}

	return strings.Join(lines, "\n")
}

func (renderer Renderer) Suggestions(suggestions []app.CommitSuggestion) string {
	if len(suggestions) == 0 {
		return ""
	}

	lines := []string{
		colorizeLine(borderColor, horizontalRule()),
		colorizeLine(headerColor, "◆ sugestoes"),
	}

	for _, suggestion := range suggestions {
		lines = append(lines, renderer.bulletLine("→", suggestion.Message))
	}

	return strings.Join(lines, "\n")
}

func (renderer Renderer) CommitSummary(summary CommitSummary) string {
	lines := []string{
		colorizeLine(borderColor, horizontalRule()),
		colorizeLine(successColor, fmt.Sprintf("✔ %d %s", summary.Created, pluralizeCommits(summary.Created))),
		colorizeLine(defaultColor, fmt.Sprintf("qualidade media: %d", summary.AverageQuality)),
		colorizeLine(defaultColor, "status: "+strings.TrimSpace(summary.Status)),
	}

	return strings.Join(lines, "\n")
}
