package ui

import (
	"fmt"
	"strings"

	"github.com/ovitorvalente/git-loom/internal/app"
)

func (renderer Renderer) ChangedFiles(stagedPaths []string, changedPaths []string) string {
	lines := []string{
		colorizeLine(emphasisColor, "arquivos detectados no working tree"),
	}

	if len(stagedPaths) > 0 {
		lines = append(lines, "", renderer.sectionTitle("staged"))
		for _, path := range stagedPaths {
			lines = append(lines, "  "+colorizeText(statusAddColor, "+")+" "+colorizeText(defaultColor, path))
		}
	}

	if len(changedPaths) > 0 {
		lines = append(lines, "", renderer.sectionTitle(fmt.Sprintf("changes (%d)", len(changedPaths))))
		for _, path := range changedPaths {
			lines = append(lines, "  "+colorizeText(statusUpdateColor, "~")+" "+colorizeText(defaultColor, path))
		}
	}

	return strings.Join(lines, "\n")
}

func (renderer Renderer) Suggestions(suggestions []app.CommitSuggestion) string {
	if len(suggestions) == 0 {
		return ""
	}

	lines := []string{
		colorizeLine(emphasisColor, "sugestoes encontradas"),
	}

	for _, suggestion := range suggestions {
		lines = append(lines, colorizeLine(mutatedColor, "  • "+suggestion.Message))
	}

	return strings.Join(lines, "\n")
}

func (renderer Renderer) CommitSummary(summary CommitSummary) string {
	created := summary.Created
	label := pluralizeCommits(created)
	qualityStr := scoreBadge(summary.AverageQuality)

	lines := []string{
		colorizeLine(successColor, fmt.Sprintf("✔ %d %s", created, label)),
		"",
		colorizeLine(defaultColor, "qualidade media: "+qualityStr),
		colorizeLine(defaultColor, "status: "+strings.TrimSpace(summary.Status)),
	}

	return strings.Join(lines, "\n")
}

func (renderer Renderer) FinalPreview(plans []app.CommitPlan) string {
	if len(plans) == 0 {
		return ""
	}

	commitLabel := "commit"
	if len(plans) > 1 {
		commitLabel = "commits"
	}

	lines := []string{
		colorizeLine(emphasisColor, fmt.Sprintf("resumo: %d %s sera criado", len(plans), commitLabel)),
		"",
	}

	for _, plan := range plans {
		subject, _ := splitCommitMessage(plan.Result.Message)
		commitType := string(plan.Result.Commit.Type)
		scope := strings.TrimSpace(plan.Result.Commit.Scope)

		header := commitType
		if scope != "" {
			header += "(" + scope + ")"
		}

		lines = append(lines,
			colorizeLine(typeColorFor(commitType), header),
			colorizeLine(defaultColor, "→ "+subject),
			"",
		)
	}

	return strings.Join(lines, "\n")
}

type CommitSummary struct {
	Created        int
	AverageQuality int
	Status         string
}
