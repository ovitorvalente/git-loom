package ui

import (
	"fmt"
	"strings"

	"github.com/ovitorvalente/git-loom/internal/app"
	"github.com/ovitorvalente/git-loom/internal/semantic"
)

func (renderer Renderer) CommitPlan(index int, total int, plan app.CommitPlan) string {
	feedback := app.BuildCommitFeedback(plan)
	subject, body := splitCommitMessage(plan.Result.Message)

	lines := []string{
		renderer.renderCompactHeader(plan, subject),
	}

	if renderer.mode() == RenderModeVerbose {
		lines = append(lines, "", renderer.renderVerboseHeader(index, total))
	}

	lines = append(lines, "")

	impactLine := renderer.renderImpact(plan)
	if impactLine != "" {
		lines = append(lines, impactLine, "")
	}

	if renderer.mode() == RenderModeVerbose {
		detailsBlock := renderer.renderDetails(body)
		if detailsBlock != "" {
			lines = append(lines, detailsBlock, "")
		}
	}

	analysisBlock := renderer.renderAnalysis(feedback.Highlights)
	if analysisBlock != "" {
		lines = append(lines, renderer.sectionTitle("analise"), analysisBlock)
	}

	suggestionsBlock := renderer.renderSuggestionsList(feedback.Suggestions)
	if suggestionsBlock != "" {
		if analysisBlock != "" {
			lines = append(lines, "")
		}
		lines = append(lines, renderer.sectionTitle("sugestoes"), suggestionsBlock)
	}

	if renderer.mode() == RenderModeVerbose {
		criteriaBlock := renderer.renderCriteria(plan.Quality.Criteria)
		if criteriaBlock != "" {
			lines = append(lines, "", renderer.sectionTitle("criterios"), criteriaBlock)
		}
	}

	if renderer.withPreview() {
		preview := renderer.renderPreview(plan.Preview)
		if preview != "" {
			lines = append(lines, "", renderer.sectionTitle("preview"), preview)
		}
	}

	return strings.Join(lines, "\n")
}

func (renderer Renderer) renderCompactHeader(plan app.CommitPlan, subject string) string {
	commitType := string(plan.Result.Commit.Type)
	if commitType == "" {
		commitType = "chore"
	}

	scope := strings.TrimSpace(plan.Result.Commit.Scope)
	score := plan.Quality.Score
	badge := scoreEmoji(score)

	typeColor := typeColorFor(commitType)
	typeStr := colorizeText(typeColor, commitType)

	scopeStr := ""
	if scope != "" {
		scopeStr = colorizeText(scopeValueColor, scope)
	}

	scoreStr := colorizeText(scoreColor(score), fmt.Sprintf("%d", score))

	if scopeStr != "" {
		return fmt.Sprintf("◆ %s(%s) [%s] %s", typeStr, scopeStr, scoreStr, badge)
	}

	return fmt.Sprintf("◆ %s [%s] %s", typeStr, scoreStr, badge)
}

func (renderer Renderer) renderVerboseHeader(index int, total int) string {
	return colorizeText(borderColor, fmt.Sprintf("commit %d/%d", index, total))
}

func (renderer Renderer) renderImpact(plan app.CommitPlan) string {
	files := plan.Context.Files
	if len(files) == 0 {
		for _, path := range plan.Result.Paths {
			files = append(files, semantic.ChangedFile{Path: path, Status: "atualizado"})
		}
	}

	if renderer.mode() != RenderModeVerbose && len(files) > 3 {
		files = files[:3]
	}

	lines := make([]string, 0, len(files))
	for _, file := range files {
		lines = append(lines, renderer.renderFileLine(file))
	}

	if renderer.mode() != RenderModeVerbose && len(plan.Context.Files) > len(files) {
		lines = append(lines, colorizeText(borderColor, fmt.Sprintf("  +%d arquivo(s)", len(plan.Context.Files)-len(files))))
	}

	if plan.Preview.Additions > 0 || plan.Preview.Deletions > 0 {
		impact := fmt.Sprintf("+%d -%d", plan.Preview.Additions, plan.Preview.Deletions)
		lines = append(lines, colorizeText(mutedColor, "  "+impact))
	}

	return strings.Join(lines, "\n")
}

func (renderer Renderer) renderCriteria(criteria []semantic.QualityCriteria) string {
	if len(criteria) == 0 {
		return ""
	}

	lines := make([]string, 0, len(criteria))
	for _, c := range criteria {
		if c.Passed {
			lines = append(lines, colorizeText(successColor, "  ✔ "+c.Name))
		} else if c.Warning {
			msg := c.Name
			if c.Message != "" {
				msg += ": " + c.Message
			}
			lines = append(lines, colorizeText(warningColor, "  ⚠ "+msg))
		} else {
			lines = append(lines, colorizeText(dangerColor, "  ✗ "+c.Name))
		}
	}

	return strings.Join(lines, "\n")
}

func (renderer Renderer) renderFiles(plan app.CommitPlan) string {
	files := plan.Context.Files
	if len(files) == 0 {
		for _, path := range plan.Result.Paths {
			files = append(files, semantic.ChangedFile{Path: path, Status: "atualizado"})
		}
	}

	if len(files) == 0 {
		return colorizeLine(defaultColor, "  nenhum arquivo")
	}

	if renderer.mode() != RenderModeVerbose && len(files) > 3 {
		files = files[:3]
	}

	lines := make([]string, 0, len(files))
	for _, file := range files {
		lines = append(lines, renderer.renderFileLine(file))
	}

	if renderer.mode() != RenderModeVerbose && len(plan.Context.Files) > len(files) {
		lines = append(lines, colorizeLine(borderColor, fmt.Sprintf("  +%d arquivo(s)", len(plan.Context.Files)-len(files))))
	}

	return strings.Join(lines, "\n")
}

func (renderer Renderer) renderAnalysis(highlights []string) string {
	if len(highlights) == 0 {
		return ""
	}

	lines := make([]string, 0, len(highlights))
	for _, highlight := range highlights {
		lines = append(lines, colorizeLine(warningColor, "  ⚠ "+highlight))
	}

	return strings.Join(lines, "\n")
}

func (renderer Renderer) renderSuggestionsList(suggestions []string) string {
	if len(suggestions) == 0 {
		return ""
	}

	lines := make([]string, 0, len(suggestions))
	for _, suggestion := range suggestions {
		lines = append(lines, colorizeLine(mutatedColor, "  → "+suggestion))
	}

	return strings.Join(lines, "\n")
}

func (renderer Renderer) renderDetails(body string) string {
	if strings.TrimSpace(body) == "" {
		return ""
	}

	lines := []string{}
	for _, line := range strings.Split(body, "\n") {
		trimmed := strings.TrimSpace(strings.TrimPrefix(line, "-"))
		if trimmed == "" {
			continue
		}
		lines = append(lines, colorizeLine(defaultColor, "  • "+trimmed))
	}

	if len(lines) == 0 {
		return ""
	}

	if renderer.mode() != RenderModeVerbose && len(lines) > 3 {
		lines = lines[:3]
	}

	return strings.Join(lines, "\n")
}

func (renderer Renderer) renderPreview(preview app.CommitPreview) string {
	if preview.FilesChanged == 0 {
		return ""
	}

	lines := []string{
		colorizeLine(defaultColor, fmt.Sprintf("  %d arquivo(s): +%d -%d",
			preview.FilesChanged, preview.Additions, preview.Deletions)),
	}

	return strings.Join(lines, "\n")
}

func (renderer Renderer) renderFileLine(file semantic.ChangedFile) string {
	symbol := "~"
	color := statusUpdateColor

	switch strings.TrimSpace(strings.ToLower(file.Status)) {
	case "adicionado":
		symbol = "+"
		color = statusAddColor
	case "removido":
		symbol = "-"
		color = statusRemoveColor
	}

	return "  " + colorizeText(color, symbol) + " " + colorizeText(defaultColor, file.Path)
}

func scoreEmoji(score int) string {
	switch {
	case score >= 90:
		return colorizeText(successColor, "●")
	case score >= 80:
		return colorizeText(accentColor, "●")
	case score >= 70:
		return colorizeText(warningColor, "◐")
	default:
		return colorizeText(dangerColor, "○")
	}
}

func scoreColor(score int) string {
	switch {
	case score >= 90:
		return successColor
	case score >= 80:
		return accentColor
	case score >= 70:
		return warningColor
	default:
		return dangerColor
	}
}

func typeColorFor(commitType string) string {
	switch commitType {
	case "feat":
		return successColor
	case "fix":
		return dangerColor
	case "refactor":
		return accentColor
	case "docs":
		return infoColor
	case "test":
		return warningColor
	case "chore":
		return borderColor
	default:
		return defaultColor
	}
}

func wrapLines(lines []string, width int) []string {
	if width <= 0 {
		return lines
	}

	wrapped := []string{}
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			if len(wrapped) == 0 || wrapped[len(wrapped)-1] == "" {
				continue
			}
			wrapped = append(wrapped, "")
			continue
		}

		current := ""
		for _, word := range strings.Fields(trimmed) {
			candidate := word
			if current != "" {
				candidate = current + " " + word
			}
			if len(candidate) <= width {
				current = candidate
				continue
			}
			if current != "" {
				wrapped = append(wrapped, current)
			}
			current = word
		}

		if current != "" {
			wrapped = append(wrapped, current)
		}
	}

	if len(wrapped) == 0 {
		return []string{""}
	}

	return wrapped
}
