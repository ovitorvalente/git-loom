package tui

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/ovitorvalente/git-loom/internal/app"
)

type viewState int

const (
	viewAnalyze viewState = iota
	viewReview
)

type CommitToggle struct {
	Plan     app.CommitPlan
	Selected bool
}

type appModel struct {
	commits   []CommitToggle
	spinner   spinner.Model
	state     viewState
	cursor    int
	width     int
	height    int
	ready     bool
	analyzing bool
	confirmed bool
	canceled  bool
}

type AppResult struct {
	Approved  []bool
	Confirmed bool
	Canceled  bool
}

func newAppModel(plans []app.CommitPlan) appModel {
	commits := make([]CommitToggle, len(plans))
	for i, plan := range plans {
		commits[i] = CommitToggle{Plan: plan, Selected: true}
	}

	s := spinner.New()
	s.Spinner = spinner.MiniDot
	s.Style = StyleSpinner

	return appModel{
		state:     viewAnalyze,
		spinner:   s,
		analyzing: true,
		commits:   commits,
	}
}

func (m appModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		func() tea.Msg {
			time.Sleep(900 * time.Millisecond)
			return analyzeDoneMsg{}
		},
	)
}

type analyzeDoneMsg struct{}

func (m appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		return m, nil

	case analyzeDoneMsg:
		m.analyzing = false
		m.state = viewReview
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	if m.analyzing {
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m appModel) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.state == viewReview {
		return m.handleReviewKey(msg)
	}
	return m, nil
}

func (m appModel) handleReviewKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "esc", "q":
		m.canceled = true
		return m, tea.Quit
	case "enter":
		m.confirmed = true
		return m, tea.Quit
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.commits)-1 {
			m.cursor++
		}
	case " ":
		if m.cursor >= 0 && m.cursor < len(m.commits) {
			m.commits[m.cursor].Selected = !m.commits[m.cursor].Selected
		}
	case "a":
		for i := range m.commits {
			m.commits[i].Selected = true
		}
	case "n":
		for i := range m.commits {
			m.commits[i].Selected = false
		}
	}
	return m, nil
}

func (m appModel) View() string {
	if !m.ready {
		return ""
	}

	switch m.state {
	case viewAnalyze:
		return m.viewAnalyze()
	case viewReview:
		return m.viewReview()
	}
	return ""
}

func (m appModel) viewAnalyze() string {
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString("  ")
	b.WriteString(m.spinner.View())
	b.WriteString(" ")
	b.WriteString(StyleMuted.Render("analisando mudancas..."))
	b.WriteString("\n\n")
	b.WriteString("  ")
	b.WriteString(StyleMuted.Render(cycleTips()))
	b.WriteString("\n")

	return b.String()
}

func (m appModel) viewReview() string {
	var b strings.Builder
	width := m.width
	if width == 0 {
		width = 80
	}
	if width > 100 {
		width = 100
	}

	b.WriteString("\n")

	header := fmt.Sprintf("commits planejados (%d)", len(m.commits))
	b.WriteString("  ")
	b.WriteString(Logo())
	b.WriteString(" ")
	b.WriteString(StyleMuted.Render(header))
	b.WriteString("\n")

	selected := 0
	for _, c := range m.commits {
		if c.Selected {
			selected++
		}
	}

	b.WriteString("  ")
	b.WriteString(Divider(width - 4))
	b.WriteString("\n\n")

	for i, commit := range m.commits {
		b.WriteString(m.renderCommitCard(commit, i == m.cursor, width-6))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString("  ")
	b.WriteString(Divider(width - 4))
	b.WriteString("\n")
	b.WriteString(m.renderFooter(width))
	b.WriteString("\n")

	return b.String()
}

func (m appModel) renderCommitCard(commit CommitToggle, isSelected bool, maxWidth int) string {
	plan := commit.Plan
	commitType := string(plan.Result.Commit.Type)
	scope := plan.Result.Commit.Scope
	desc := plan.Result.Commit.Description
	score := plan.Quality.Score

	checkbox := StyleCheckboxEmpty.Render("○")
	if commit.Selected {
		checkbox = StyleCheckbox.Render("●")
	}

	typeStr := StyleForType(commitType).Render(commitType)
	scopeStr := ""
	if scope != "" {
		scopeStr = "(" + StyleScope.Render(scope) + ")"
	}
	scoreStr := ScoreBadge(score)

	header := fmt.Sprintf("%s %s%s %s", checkbox, typeStr, scopeStr, scoreStr)

	fileLines := []string{}
	for j, file := range plan.Result.Paths {
		if j >= 3 {
			remaining := len(plan.Result.Paths) - 3
			fileLines = append(fileLines, StyleMuted.Render(fmt.Sprintf("  +%d mais", remaining)))
			break
		}
		symbol := "~"
		color := StyleWarning
		if j < len(plan.Context.Files) {
			switch plan.Context.Files[j].Status {
			case "adicionado":
				symbol = "+"
				color = StyleSuccess
			case "removido":
				symbol = "-"
				color = StyleDanger
			}
		}
		fileLines = append(fileLines, fmt.Sprintf("  %s %s", color.Render(symbol), StyleBright.Render(file)))
	}

	impactLine := ""
	if plan.Preview.Additions > 0 || plan.Preview.Deletions > 0 {
		impactLine = StyleMuted.Render(fmt.Sprintf("  +%d -%d", plan.Preview.Additions, plan.Preview.Deletions))
	}

	warnings := []string{}
	for _, reason := range plan.Quality.Reasons {
		warnings = append(warnings, StyleWarning.Render("  ⚠ ")+StyleMuted.Render(reason))
	}

	content := header + "\n"
	if desc != "" {
		content += StyleMuted.Render("  "+desc) + "\n"
	}
	content += strings.Join(fileLines, "\n") + "\n"
	if impactLine != "" {
		content += impactLine + "\n"
	}
	if isSelected && len(warnings) > 0 {
		content += "\n" + strings.Join(warnings, "\n") + "\n"
	}

	style := StyleCard
	if isSelected {
		style = StyleCardSelected
	}

	return style.Width(maxWidth).Render(content)
}

func (m appModel) renderFooter(width int) string {
	keys := []struct{ key, desc string }{
		{"↑/↓", "navegar"},
		{"space", "toggle"},
		{"a", "todos"},
		{"n", "nenhum"},
		{"enter", "confirmar"},
		{"esc", "cancelar"},
	}

	parts := []string{}
	for _, k := range keys {
		parts = append(parts, StyleFooterKey.Render(k.key)+StyleFooter.Render(" "+k.desc))
	}

	return "  " + strings.Join(parts, StyleMuted.Render(" · "))
}

func (m appModel) result() AppResult {
	approved := make([]bool, len(m.commits))
	for i, c := range m.commits {
		approved[i] = c.Selected
	}
	return AppResult{
		Approved:  approved,
		Confirmed: m.confirmed,
		Canceled:  m.canceled,
	}
}

func RunCommitTUI(plans []app.CommitPlan) (AppResult, error) {
	if len(plans) == 0 {
		return AppResult{Confirmed: true}, nil
	}

	if !isTerminal(os.Stdin) || !isTerminal(os.Stdout) {
		result := make([]bool, len(plans))
		for i := range result {
			result[i] = true
		}
		return AppResult{Approved: result, Confirmed: true}, nil
	}

	m := newAppModel(plans)
	p := tea.NewProgram(m, tea.WithInput(os.Stdin), tea.WithOutput(os.Stdout))
	finalModel, err := p.Run()
	if err != nil {
		return AppResult{}, err
	}

	if fm, ok := finalModel.(appModel); ok {
		return fm.result(), nil
	}

	return AppResult{Canceled: true}, nil
}
