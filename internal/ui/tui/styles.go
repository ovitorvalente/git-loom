package tui

import "github.com/charmbracelet/lipgloss"

var (
	StyleHeader   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("99"))
	StyleSuccess  = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	StyleWarning  = lipgloss.NewStyle().Foreground(lipgloss.Color("220"))
	StyleDanger   = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	StyleMuted    = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	StyleType     = lipgloss.NewStyle().Bold(true)
	StyleScope    = lipgloss.NewStyle().Foreground(lipgloss.Color("219"))
	StyleSelected = lipgloss.NewStyle().Foreground(lipgloss.Color("170"))
	StyleCheckbox = lipgloss.NewStyle().Foreground(lipgloss.Color("99"))
	StyleSpinner  = lipgloss.NewStyle().Foreground(lipgloss.Color("99"))
	StyleScore    = lipgloss.NewStyle()
)

func TypeColor(commitType string) lipgloss.Style {
	switch commitType {
	case "feat":
		return lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	case "fix":
		return lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	case "refactor":
		return lipgloss.NewStyle().Foreground(lipgloss.Color("99"))
	case "docs":
		return lipgloss.NewStyle().Foreground(lipgloss.Color("34"))
	case "test":
		return lipgloss.NewStyle().Foreground(lipgloss.Color("220"))
	case "chore":
		return lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	default:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	}
}

func ScoreColor(score int) lipgloss.Style {
	switch {
	case score >= 90:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	case score >= 80:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("99"))
	case score >= 70:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("220"))
	default:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	}
}
