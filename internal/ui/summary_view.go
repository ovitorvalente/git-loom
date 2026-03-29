package ui

import (
	"fmt"
	"strings"

	"github.com/ovitorvalente/git-loom/internal/app"
	"github.com/ovitorvalente/git-loom/internal/shared"
)

func (renderer Renderer) ChangedFiles(stagedPaths []string, changedPaths []string) string {
	lines := []string{
		colorizeLine(emphasisColor, "Arquivos detectados no working tree"),
	}

	if len(stagedPaths) > 0 {
		lines = append(lines, "", renderer.sectionTitle("Staged"))
		for _, path := range stagedPaths {
			lines = append(lines, colorizeText(panelBorderColor, "│")+" "+colorizeText(statusAddColor, "+")+" "+colorizeText(defaultColor, path))
		}
	}

	if len(changedPaths) > 0 {
		lines = append(lines, "", renderer.sectionTitle(fmt.Sprintf("Changes (%d)", len(changedPaths))))
		for _, path := range changedPaths {
			lines = append(lines, colorizeText(panelBorderColor, "│")+" "+colorizeText(statusUpdateColor, "~")+" "+colorizeText(defaultColor, path))
		}
	}

	return strings.Join(lines, "\n")
}

func (renderer Renderer) Suggestions(suggestions []app.CommitSuggestion) string {
	if len(suggestions) == 0 {
		return ""
	}

	lines := []string{
		colorizeLine(emphasisColor, "Sugestoes encontradas"),
	}

	for _, suggestion := range suggestions {
		lines = append(lines, colorizeLine(mutatedColor, "• "+suggestion.Message))
	}

	return strings.Join(lines, "\n")
}

func (renderer Renderer) CommitSummary(summary CommitSummary) string {
	lines := []string{
		colorizeLine(successColor, fmt.Sprintf("Commit concluido: %d %s", summary.Created, pluralizeCommits(summary.Created))),
		colorizeLine(defaultColor, fmt.Sprintf("Qualidade media: %d", summary.AverageQuality)),
		colorizeLine(defaultColor, "Status: "+strings.TrimSpace(summary.Status)),
		colorizeLine(accentColor, shared.MessageCommitFarewell),
	}

	return strings.Join(lines, "\n")
}
