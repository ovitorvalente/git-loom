package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

var (
	colorPrimary = lipgloss.Color("99")
	colorSuccess = lipgloss.Color("42")
	colorWarning = lipgloss.Color("220")
	colorDanger  = lipgloss.Color("196")
	colorMuted   = lipgloss.Color("245")
	colorAccent  = lipgloss.Color("170")
	colorScope   = lipgloss.Color("219")
	colorBright  = lipgloss.Color("252")
	colorBorder  = lipgloss.Color("240")
)

var (
	StyleHeader = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorPrimary)

	StyleSuccess = lipgloss.NewStyle().
			Foreground(colorSuccess)

	StyleWarning = lipgloss.NewStyle().
			Foreground(colorWarning)

	StyleDanger = lipgloss.NewStyle().
			Foreground(colorDanger)

	StyleMuted = lipgloss.NewStyle().
			Foreground(colorMuted)

	StyleAccent = lipgloss.NewStyle().
			Foreground(colorAccent)

	StyleBright = lipgloss.NewStyle().
			Foreground(colorBright)

	StyleBorder = lipgloss.NewStyle().
			Foreground(colorBorder)

	StyleTypeFeat = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorSuccess)

	StyleTypeFix = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorDanger)

	StyleTypeRefactor = lipgloss.NewStyle().
				Bold(true).
				Foreground(colorPrimary)

	StyleTypeDocs = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("34"))

	StyleTypeTest = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorWarning)

	StyleTypeChore = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorMuted)

	StyleScope = lipgloss.NewStyle().
			Foreground(colorScope)

	StyleCard = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorBorder).
			Padding(0, 1)

	StyleCardSelected = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(colorPrimary).
				Padding(0, 1)

	StyleFooter = lipgloss.NewStyle().
			Foreground(colorMuted)

	StyleFooterKey = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorBright)

	StyleSpinner = lipgloss.NewStyle().
			Foreground(colorPrimary)

	StyleCheckbox = lipgloss.NewStyle().
			Foreground(colorSuccess)

	StyleCheckboxEmpty = lipgloss.NewStyle().
				Foreground(colorMuted)

	StyleScoreGood = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorSuccess)

	StyleScoreOk = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorPrimary)

	StyleScoreWarn = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorWarning)

	StyleScoreBad = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorDanger)
)

func StyleForType(commitType string) lipgloss.Style {
	switch commitType {
	case "feat":
		return StyleTypeFeat
	case "fix":
		return StyleTypeFix
	case "refactor":
		return StyleTypeRefactor
	case "docs":
		return StyleTypeDocs
	case "test":
		return StyleTypeTest
	default:
		return StyleTypeChore
	}
}

func StyleForScore(score int) lipgloss.Style {
	switch {
	case score >= 90:
		return StyleScoreGood
	case score >= 80:
		return StyleScoreOk
	case score >= 70:
		return StyleScoreWarn
	default:
		return StyleScoreBad
	}
}

func ScoreBadge(score int) string {
	var label string
	switch {
	case score >= 90:
		label = "excelente"
	case score >= 80:
		label = "bom"
	case score >= 70:
		label = "aceitavel"
	default:
		label = "critico"
	}
	return StyleForScore(score).Render(fmt.Sprintf("%d", score)) + StyleMuted.Render(" "+label)
}

func Divider(width int) string {
	return StyleMuted.Render(strings.Repeat("─", width))
}

func Logo() string {
	return StyleHeader.Render("◆ gitloom")
}

var tips = []string{
	"use --preview para ver o diff antes de commitar",
	"use --strict para falhar em commits de baixa qualidade",
	"use --yes para pular confirmacoes",
	"commits sao agrupados por tipo e escopo automaticamente",
	"o score de qualidade detecta escopos genericos",
	"commits com 4+ arquivos sao divididos em blocos",
	"use --verbose para ver criterios de qualidade",
	"gl commit eh o mesmo que gitloom commit",
}

func cycleTips() string {
	idx := int(time.Now().Unix()/5) % len(tips)
	return tips[idx]
}
