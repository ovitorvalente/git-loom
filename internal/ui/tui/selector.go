package tui

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/ovitorvalente/git-loom/internal/app"
)

type SelectorResult struct {
	Approved  []bool
	Confirmed bool
	Cancelled bool
}

type commitItem struct {
	index       int
	commitType  string
	scope       string
	description string
	score       int
	paths       []string
	selected    bool
}

func (i commitItem) Title() string { return i.FilterValue() }
func (i commitItem) Description() string {
	return fmt.Sprintf("%d arquivo(s) — qualidade %d", len(i.paths), i.score)
}
func (i commitItem) FilterValue() string {
	scope := ""
	if i.scope != "" {
		scope = "(" + i.scope + ")"
	}
	return fmt.Sprintf("%s%s: %s", i.commitType, scope, i.description)
}

type selectorModel struct {
	list      list.Model
	items     []commitItem
	confirmed bool
	cancelled bool
	width     int
	height    int
}

func newSelectorModel(plans []app.CommitPlan) selectorModel {
	items := make([]list.Item, len(plans))
	commitItems := make([]commitItem, len(plans))

	for i, plan := range plans {
		ci := commitItem{
			index:       i,
			commitType:  string(plan.Result.Commit.Type),
			scope:       plan.Result.Commit.Scope,
			description: plan.Result.Commit.Description,
			score:       plan.Quality.Score,
			paths:       plan.Result.Paths,
			selected:    true,
		}
		items[i] = ci
		commitItems[i] = ci
	}

	delegate := newCommitDelegate(commitItems)
	l := list.New(items, delegate, 0, 0)
	l.Title = fmt.Sprintf("commits planejados (%d)", len(plans))
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.DisableQuitKeybindings()

	return selectorModel{
		list:  l,
		items: commitItems,
	}
}

func (m selectorModel) Init() tea.Cmd {
	return nil
}

func (m selectorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height - 4)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			m.cancelled = true
			return m, tea.Quit
		case "enter":
			m.confirmed = true
			return m, tea.Quit
		case " ":
			idx := m.list.Index()
			if idx >= 0 && idx < len(m.items) {
				m.items[idx].selected = !m.items[idx].selected
				m.list.SetItem(idx, m.items[idx])
			}
			return m, nil
		case "a":
			for i := range m.items {
				m.items[i].selected = true
				m.list.SetItem(i, m.items[i])
			}
			return m, nil
		case "n":
			for i := range m.items {
				m.items[i].selected = false
				m.list.SetItem(i, m.items[i])
			}
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m selectorModel) View() string {
	selectedCount := 0
	for _, item := range m.items {
		if item.selected {
			selectedCount++
		}
	}

	help := StyleMuted.Render(fmt.Sprintf(
		"  %d/%d selecionados — enter confirmar • space toggle • a todos • nenhum • esc cancelar",
		selectedCount, len(m.items),
	))

	return "\n" + m.list.View() + "\n\n" + help + "\n"
}

func (m selectorModel) result() SelectorResult {
	approved := make([]bool, len(m.items))
	for i, item := range m.items {
		approved[i] = item.selected
	}
	return SelectorResult{
		Approved:  approved,
		Confirmed: m.confirmed,
		Cancelled: m.cancelled,
	}
}

type commitDelegate struct {
	items map[int]commitItem
}

func newCommitDelegate(items []commitItem) commitDelegate {
	d := commitDelegate{items: make(map[int]commitItem)}
	for _, item := range items {
		d.items[item.index] = item
	}
	return d
}

func (d commitDelegate) Height() int                             { return 2 }
func (d commitDelegate) Spacing() int                            { return 1 }
func (d commitDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d commitDelegate) Render(w io.Writer, m list.Model, index int, li list.Item) {
	item, ok := li.(commitItem)
	if !ok {
		return
	}

	isSelected := index == m.Index()
	checkbox := "☐"
	if item.selected {
		checkbox = StyleCheckbox.Render("☑")
	} else {
		checkbox = StyleMuted.Render("☐")
	}

	scope := ""
	if item.scope != "" {
		scope = "(" + StyleScope.Render(item.scope) + ")"
	}

	typeStr := TypeColor(item.commitType).Render(item.commitType)
	scoreStr := ScoreColor(item.score).Render(fmt.Sprintf("[%d]", item.score))

	line := fmt.Sprintf("%s %s%s: %s %s",
		checkbox,
		typeStr,
		scope,
		item.description,
		scoreStr,
	)

	if isSelected {
		fmt.Fprintf(w, "  %s", StyleSelected.Render("▸")+" "+line)
	} else {
		fmt.Fprintf(w, "  %s", " "+line)
	}
}

func SelectCommits(plans []app.CommitPlan) (SelectorResult, error) {
	if len(plans) == 0 {
		return SelectorResult{Confirmed: true}, nil
	}

	if !isTerminal(os.Stdin) || !isTerminal(os.Stdout) {
		result := make([]bool, len(plans))
		for i := range result {
			result[i] = true
		}
		return SelectorResult{Approved: result, Confirmed: true}, nil
	}

	m := newSelectorModel(plans)
	p := tea.NewProgram(m, tea.WithInput(os.Stdin), tea.WithOutput(os.Stdout))
	finalModel, err := p.Run()
	if err != nil {
		return SelectorResult{}, err
	}

	if fm, ok := finalModel.(selectorModel); ok {
		return fm.result(), nil
	}

	return SelectorResult{Cancelled: true}, nil
}

func formatCommitLine(item commitItem) string {
	scope := ""
	if item.scope != "" {
		scope = "(" + item.scope + ")"
	}
	return strings.TrimSpace(fmt.Sprintf("%s%s: %s", item.commitType, scope, item.description))
}
