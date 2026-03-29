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
		renderer.renderHeader(index, total, plan, subject),
		"",
		renderer.renderFiles(plan),
		"",
		renderer.sectionTitle("Sugestao de commit"),
		styledPanel(renderer.renderMessagePanel(plan, subject)),
		"",
		renderer.renderMeta(plan),
	}

	detailsBlock := renderer.renderDetails(body)
	if detailsBlock != "" {
		lines = append(lines, "", renderer.sectionTitle("Detalhes"), detailsBlock)
	}

	analysisBlock := renderer.renderAnalysis(feedback.Highlights)
	if analysisBlock != "" {
		lines = append(lines, "", renderer.sectionTitle("Analise"), analysisBlock)
	}

	suggestionsBlock := renderer.renderSuggestionsList(feedback.Suggestions)
	if suggestionsBlock != "" {
		lines = append(lines, "", renderer.sectionTitle("Sugestoes"), suggestionsBlock)
	}

	if renderer.mode() == RenderModeVerbose {
		verboseBlock := renderer.renderVerbose(plan, subject, body)
		if verboseBlock != "" {
			lines = append(lines, "", renderer.sectionTitle("Verbose"), verboseBlock)
		}
	} else if renderer.withPreview() {
		preview := renderer.renderPreview(plan.Preview)
		if preview != "" {
			lines = append(lines, "", renderer.sectionTitle("Preview"), preview)
		}
	}

	return strings.Join(lines, "\n")
}

func (renderer Renderer) renderHeader(index int, total int, plan app.CommitPlan, subject string) string {
	fileCount := len(plan.Result.Paths)
	label := "arquivo modificado"
	if fileCount != 1 {
		label = "arquivos modificados"
	}

	header := fmt.Sprintf("Analisando %d %s...", fileCount, label)
	if total > 1 {
		header = fmt.Sprintf("%s %s", header, colorizeText(borderColor, fmt.Sprintf("[%d/%d]", index, total)))
	}

	return colorizeLine(emphasisColor, header)
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
	files := plan.Context.Files
	if len(files) == 0 {
		for _, path := range plan.Result.Paths {
			files = append(files, semantic.ChangedFile{Path: path, Status: "atualizado"})
		}
	}

	if len(files) == 0 {
		return colorizeLine(defaultColor, "  nenhum arquivo identificado")
	}

	if renderer.mode() != RenderModeVerbose && len(files) > 3 {
		files = files[:3]
	}

	lines := make([]string, 0, len(files))
	for _, file := range files {
		lines = append(lines, renderer.renderFileLine(file))
	}

	if renderer.mode() != RenderModeVerbose && len(plan.Context.Files) > len(files) {
		lines = append(lines, colorizeLine(borderColor, fmt.Sprintf("│  +%d arquivo(s) adicional(is)", len(plan.Context.Files)-len(files))))
	}

	return strings.Join(lines, "\n")
}

func (renderer Renderer) renderAnalysis(highlights []string) string {
	if len(highlights) == 0 {
		return colorizeLine(successColor, "• sem alertas relevantes")
	}

	lines := make([]string, 0, len(highlights))
	for _, highlight := range highlights {
		lines = append(lines, colorizeLine(warningColor, "• "+highlight))
	}

	return strings.Join(lines, "\n")
}

func (renderer Renderer) renderSuggestionsList(suggestions []string) string {
	if len(suggestions) == 0 {
		return ""
	}

	lines := make([]string, 0, len(suggestions))
	for _, suggestion := range suggestions {
		lines = append(lines, colorizeLine(mutatedColor, "• "+suggestion))
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
	message := strings.TrimSpace(strings.ReplaceAll(plan.Result.Message, "\n\n", " | "))
	lines := []string{
		renderer.bulletLine("•", "mensagem final: "+message),
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

func (renderer Renderer) renderMessagePanel(plan app.CommitPlan, subject string) []string {
	lines := []string{renderer.renderMessage(plan, subject)}

	if renderer.mode() == RenderModeVerbose {
		for _, line := range strings.Split(strings.TrimSpace(plan.Result.Commit.Body), "\n") {
			trimmed := strings.TrimSpace(strings.TrimPrefix(line, "-"))
			if trimmed != "" {
				lines = append(lines, trimmed)
			}
		}
	}

	return wrapLines(lines, 54)
}

func (renderer Renderer) renderMeta(plan app.CommitPlan) string {
	score := colorizeText(labelColor, "Score:") + " " + scoreText(plan.Quality.Score)
	commitType := string(plan.Result.Commit.Type)
	if commitType == "" {
		commitType = "chore"
	}

	parts := []string{
		score,
		colorizeText(labelColor, "Tipo:") + " " + colorizeText(typeValueColor, commitType),
	}

	scope := strings.TrimSpace(plan.Result.Commit.Scope)
	if scope != "" {
		parts = append(parts, colorizeText(labelColor, "Escopo:")+" "+colorizeText(scopeValueColor, scope))
	}

	return strings.Join(parts, "   ")
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

	return colorizeText(panelBorderColor, "│")+" "+colorizeText(color, symbol)+" "+colorizeText(defaultColor, file.Path)
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
