package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/ovitorvalente/git-loom/internal/app"
	"github.com/ovitorvalente/git-loom/internal/semantic"
	"github.com/ovitorvalente/git-loom/internal/shared"
)

func FormatCommitResult(result app.CommitResult) string {
	header, body := splitCommitMessage(result.Message)
	if body == "" {
		body = strings.TrimSpace(result.Commit.Body)
	}

	lines := []string{
		colorizeLine(accentColor, "✦ "+shared.MessageCommitGenerated),
		formatMetaLine("•", shared.MessageTypeLabel, fmt.Sprint(result.Commit.Type), typeColor(result.Commit.Type)),
		formatMetaLine("•", shared.MessageScopeLabel, formatValue(result.Commit.Scope), mutatedColor),
		formatMetaLine("•", shared.MessageIntentLabel, formatValue(result.Commit.Intent), accentColor),
		formatMetaLine("•", shared.MessageDescriptionLabel, formatValue(result.Commit.Description), defaultColor),
		formatMetaLine("›", shared.MessageHeaderLabel, header, headerColor),
	}

	if body != "" {
		lines = append(lines, colorizeLine(accentColor, fmt.Sprintf("  ✧ %s:", shared.MessageDetailsLabel)))
		lines = append(lines, indentBody(body))
	}

	return strings.Join(lines, "\n")
}

func FormatCommitPlan(index int, total int, plan app.CommitPlan, showPreview bool) string {
	scoreSummary := fmt.Sprintf("%s: %d/100", shared.MessageQualityLabel, plan.Quality.Score)
	lines := []string{
		colorizeLine(borderColor, horizontalRule()),
		colorizeLine(accentColor, fmt.Sprintf("◆ %s  [%s]  (%s)", fmt.Sprintf(shared.MessageCommitLabel, index, total), plan.Result.Commit.Type, scoreSummary)),
		colorizeLine(mutatedColor, "  "+formatPlanSummary(plan)),
		FormatCommitResult(plan.Result),
	}

	if showPreview {
		lines = append(lines, colorizeLine(accentColor, fmt.Sprintf("  ✧ %s:", shared.MessagePreviewLabel)))
		lines = append(lines, indentBody(formatPreview(plan.Preview)))
	}

	if len(plan.Quality.Reasons) > 0 {
		lines = append(lines, colorizeLine(accentColor, fmt.Sprintf("  ✧ %s:", shared.MessageQualityLabel)))
		lines = append(lines, indentBody(formatReasons(plan.Quality.Reasons)))
	}

	if len(plan.Result.Paths) > 0 {
		lines = append(lines, colorizeLine(accentColor, fmt.Sprintf("  ✧ %s:", shared.MessageFilesLabel)))
		lines = append(lines, indentBody(formatPaths(plan.Result.Paths)))
	}

	return strings.Join(lines, "\n")
}

func FormatChangedFiles(paths []string) string {
	lines := []string{
		colorizeLine(borderColor, horizontalRule()),
		colorizeLine(accentColor, "↺ "+shared.MessageChangedFiles),
		colorizeLine(mutatedColor, fmt.Sprintf("  total: %d arquivo(s) modificado(s)", len(paths))),
		indentBody(formatPaths(paths)),
	}

	return strings.Join(lines, "\n")
}

func FormatCommitConclusion() string {
	lines := []string{
		colorizeLine(borderColor, horizontalRule()),
		colorizeLine(successColor, "✔ "+shared.MessageCommitFinished),
		colorizeLine(accentColor, "☕ "+shared.MessageCommitFarewell),
	}

	return strings.Join(lines, "\n")
}

func FormatSuggestions(suggestions []app.CommitSuggestion) string {
	if len(suggestions) == 0 {
		return ""
	}

	lines := []string{
		colorizeLine(borderColor, horizontalRule()),
		colorizeLine(accentColor, "✦ "+shared.MessageSuggestionsLabel),
	}

	for index, suggestion := range suggestions {
		lines = append(lines, colorizeLine(defaultColor, fmt.Sprintf("  %d. %s", index+1, suggestion.Message)))
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
		lines[index] = colorizeLine(defaultColor, "    "+strings.TrimSpace(line))
	}

	return strings.Join(lines, "\n")
}

func formatPaths(paths []string) string {
	lines := make([]string, 0, len(paths))
	for _, path := range paths {
		lines = append(lines, "◦ "+path)
	}

	return strings.Join(lines, "\n")
}

const (
	borderColor  = "90"
	accentColor  = "36"
	headerColor  = "33"
	defaultColor = "37"
	mutatedColor = "94"
	successColor = "32"
	warningColor = "33"
	dangerColor  = "31"
)

func formatPlanSummary(plan app.CommitPlan) string {
	detailsCount := countDetails(plan.Result.Commit.Body)
	return fmt.Sprintf("resumo: %d arquivo(s) | %d detalhe(s) | %s: +%d/-%d", len(plan.Result.Paths), detailsCount, shared.MessageImpactLabel, plan.Preview.Additions, plan.Preview.Deletions)
}

func countDetails(body string) int {
	trimmedBody := strings.TrimSpace(body)
	if trimmedBody == "" {
		return 0
	}

	return len(strings.Split(trimmedBody, "\n"))
}

func horizontalRule() string {
	return strings.Repeat("─", 60)
}

func typeColor(commitType any) string {
	switch fmt.Sprint(commitType) {
	case "feat":
		return successColor
	case "fix":
		return dangerColor
	case "refactor":
		return accentColor
	case "docs":
		return warningColor
	case "test":
		return mutatedColor
	default:
		return defaultColor
	}
}

func colorizeLine(color string, line string) string {
	if !useANSIColors() {
		return line
	}

	return "\x1b[" + color + "m" + line + "\x1b[0m"
}

func useANSIColors() bool {
	return os.Getenv("NO_COLOR") == ""
}

func formatMetaLine(icon string, label string, value string, color string) string {
	return colorizeLine(color, fmt.Sprintf("  %s %s: %s", icon, label, value))
}

func formatPreview(preview semantic.CommitPreview) string {
	return fmt.Sprintf("%d arquivo(s) alterado(s)\n+%d linha(s)\n-%d linha(s)", preview.FilesChanged, preview.Additions, preview.Deletions)
}

func formatReasons(reasons []string) string {
	lines := make([]string, 0, len(reasons))
	for _, reason := range reasons {
		lines = append(lines, "- "+reason)
	}

	return strings.Join(lines, "\n")
}
