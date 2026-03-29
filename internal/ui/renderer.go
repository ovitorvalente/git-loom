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

func NewRenderer(options RenderOptions) Renderer {
	if options.Mode == "" {
		options.Mode = RenderModeClean
	}

	return Renderer{options: options}
}

func (renderer Renderer) sectionTitle(title string) string {
	return colorizeLine(accentColor, title+":")
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
	borderColor       = "90"
	accentColor       = "36"
	headerColor       = "33"
	defaultColor      = "37"
	mutatedColor      = "94"
	mutedColor        = "90"
	successColor      = "32"
	warningColor      = "33"
	dangerColor       = "31"
	emphasisColor     = "96"
	infoColor         = "34"
	magentaColor      = "35"
	panelBackground   = "48;5;53"
	panelBorderColor  = "38;5;203"
	panelTextColor    = "38;5;230"
	mutedCapsColor    = "37"
	statusAddColor    = "32"
	statusUpdateColor = "33"
	statusRemoveColor = "31"
	statusPromptColor = "36"
	labelColor        = "38;5;109"
	typeValueColor    = "38;5;45"
	scopeValueColor   = "38;5;219"
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

func colorizeText(color string, value string) string {
	if !useANSIColors() || strings.TrimSpace(color) == "" {
		return value
	}

	return "\x1b[" + color + "m" + value + "\x1b[0m"
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
