package ui

import (
	"fmt"
	"os"
	"strings"
)

type RenderMode string

const (
	RenderModeClean   RenderMode = "clean"
	RenderModeVerbose RenderMode = "verbose"
)

type RenderOptions struct {
	Mode        RenderMode
	ShowPreview bool
}

type Renderer struct {
	options RenderOptions
}

type CommitSummary struct {
	Created        int
	AverageQuality int
	Status         string
}

func NewRenderer(options RenderOptions) Renderer {
	if options.Mode == "" {
		options.Mode = RenderModeClean
	}

	return Renderer{options: options}
}

func (renderer Renderer) sectionTitle(title string) string {
	return colorizeLine(accentColor, title+":")
}

func (renderer Renderer) bulletLine(prefix string, value string) string {
	return colorizeLine(defaultColor, prefix+" "+value)
}

func (renderer Renderer) mode() RenderMode {
	if renderer.options.Mode == "" {
		return RenderModeClean
	}

	return renderer.options.Mode
}

func (renderer Renderer) withPreview() bool {
	return renderer.options.ShowPreview
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

func colorizeLine(color string, line string) string {
	if !useANSIColors() {
		return line
	}

	return "\x1b[" + color + "m" + line + "\x1b[0m"
}

func useANSIColors() bool {
	return os.Getenv("NO_COLOR") == ""
}

func horizontalRule() string {
	return strings.Repeat("─", 60)
}

func splitCommitMessage(message string) (string, string) {
	parts := strings.SplitN(strings.TrimSpace(message), "\n\n", 2)
	if len(parts) == 1 {
		return strings.TrimSpace(parts[0]), ""
	}

	return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
}

func scoreBadge(score int) string {
	color := successColor
	label := "bom"

	switch {
	case score >= 90:
		label = "excelente"
	case score >= 80:
		label = "bom"
	case score >= 70:
		color = warningColor
		label = "aceitavel"
	default:
		color = dangerColor
		label = "critico"
	}

	return colorizeLine(color, fmt.Sprintf("[%d] %s", score, label))
}

func pluralizeCommits(total int) string {
	if total == 1 {
		return "commit criado"
	}

	return "commits criados"
}
