package ui

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ovitorvalente/git-loom/internal/app"
)

func (renderer Renderer) CommitPlan(index int, total int, plan app.CommitPlan) string {
	feedback := app.BuildCommitFeedback(plan)
	subject, body := splitCommitMessage(plan.Result.Message)
	lines := []string{
		colorizeLine(borderColor, horizontalRule()),
		renderer.renderHeader(index, total, plan, subject),
		"",
		renderer.sectionTitle("mensagem"),
		colorizeLine(defaultColor, renderer.renderMessage(plan, subject)),
	}

	detailsBlock := renderer.renderDetails(body)
	if detailsBlock != "" {
		lines = append(lines, "", renderer.sectionTitle("detalhes"), detailsBlock)
	}

	filesBlock := renderer.renderFiles(plan)
	if filesBlock != "" {
		lines = append(lines, "", renderer.sectionTitle("arquivos"), filesBlock)
	}

	analysisBlock := renderer.renderAnalysis(feedback.Highlights)
	if analysisBlock != "" {
		lines = append(lines, "", renderer.sectionTitle("analise"), analysisBlock)
	}

	suggestionsBlock := renderer.renderSuggestionsList(feedback.Suggestions)
	if suggestionsBlock != "" {
		lines = append(lines, "", renderer.sectionTitle("sugestoes"), suggestionsBlock)
	}

	if renderer.mode() == RenderModeVerbose {
		verboseBlock := renderer.renderVerbose(plan, subject, body)
		if verboseBlock != "" {
			lines = append(lines, "", renderer.sectionTitle("verbose"), verboseBlock)
		}
	} else if renderer.withPreview() {
		preview := renderer.renderPreview(plan.Preview)
		if preview != "" {
			lines = append(lines, "", renderer.sectionTitle("preview"), preview)
		}
	}

	return strings.Join(lines, "\n")
}

func (renderer Renderer) renderHeader(index int, total int, plan app.CommitPlan, subject string) string {
	header := formatCommitLabel(plan, subject)
	if total > 1 {
		header = fmt.Sprintf("%s %d/%d", header, index, total)
	}

	return colorizeLine(headerColor, "◆ "+header+" "+scoreBadge(plan.Quality.Score))
}

func formatCommitLabel(plan app.CommitPlan, subject string) string {
	scope := strings.TrimSpace(plan.Result.Commit.Scope)
	commitType := strings.TrimSpace(string(plan.Result.Commit.Type))
	if commitType == "" {
		commitType = "chore"
	}

	if scope == "" {
		return commitType
	}

	return fmt.Sprintf("%s(%s)", commitType, scope)
}

func (renderer Renderer) renderMessage(plan app.CommitPlan, subject string) string {
	if subject != "" {
		return subject
	}

	description := strings.TrimSpace(plan.Result.Commit.Description)
	if description != "" {
		return description
	}

	return strings.TrimSpace(plan.Result.Message)
}

func (renderer Renderer) renderFiles(plan app.CommitPlan) string {
	lines := []string{}
	for _, path := range plan.Result.Paths {
		lines = append(lines, colorizeLine(defaultColor, fmt.Sprintf("+%d -%d %s", plan.Preview.Additions, plan.Preview.Deletions, path)))
		if renderer.mode() != RenderModeVerbose {
			break
		}
	}

	if len(plan.Result.Paths) > 1 && renderer.mode() != RenderModeVerbose {
		lines[0] = colorizeLine(defaultColor, fmt.Sprintf("+%d -%d %s", plan.Preview.Additions, plan.Preview.Deletions, renderer.primaryFileLabel(plan.Result.Paths)))
	}

	return strings.Join(lines, "\n")
}

func (renderer Renderer) primaryFileLabel(paths []string) string {
	switch len(paths) {
	case 0:
		return "sem arquivos"
	case 1:
		return paths[0]
	default:
		base := filepath.Base(paths[0])
		return fmt.Sprintf("%s +%d", base, len(paths)-1)
	}
}

func (renderer Renderer) renderAnalysis(highlights []string) string {
	if len(highlights) == 0 {
		return colorizeLine(successColor, "ok sem alertas relevantes")
	}

	lines := make([]string, 0, len(highlights))
	for _, highlight := range highlights {
		lines = append(lines, colorizeLine(warningColor, "⚠ "+highlight))
	}

	return strings.Join(lines, "\n")
}

func (renderer Renderer) renderSuggestionsList(suggestions []string) string {
	if len(suggestions) == 0 {
		return ""
	}

	lines := make([]string, 0, len(suggestions))
	for _, suggestion := range suggestions {
		lines = append(lines, colorizeLine(mutatedColor, "→ "+suggestion))
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
		lines = append(lines, renderer.bulletLine("•", trimmed))
	}

	if len(lines) == 0 {
		return ""
	}

	if renderer.mode() != RenderModeVerbose && len(lines) > 3 {
		lines = lines[:3]
	}

	return strings.Join(lines, "\n")
}

func (renderer Renderer) renderVerbose(plan app.CommitPlan, subject string, body string) string {
	lines := []string{
		renderer.bulletLine("•", "mensagem final: "+strings.TrimSpace(plan.Result.Message)),
		renderer.bulletLine("•", fmt.Sprintf("arquivos alterados: %d", plan.Preview.FilesChanged)),
	}

	if renderer.withPreview() {
		lines = append(lines, renderer.bulletLine("•", fmt.Sprintf("impacto: +%d -%d", plan.Preview.Additions, plan.Preview.Deletions)))
	}

	if body != "" {
		for _, line := range strings.Split(body, "\n") {
			trimmed := strings.TrimSpace(strings.TrimPrefix(line, "-"))
			if trimmed == "" {
				continue
			}
			lines = append(lines, renderer.bulletLine("•", "detalhe: "+trimmed))
		}
	}

	return strings.Join(lines, "\n")
}

func (renderer Renderer) renderPreview(preview app.CommitPreview) string {
	return strings.Join([]string{
		renderer.bulletLine("•", fmt.Sprintf("arquivos: %d", preview.FilesChanged)),
		renderer.bulletLine("•", fmt.Sprintf("linhas: +%d -%d", preview.Additions, preview.Deletions)),
	}, "\n")
}
