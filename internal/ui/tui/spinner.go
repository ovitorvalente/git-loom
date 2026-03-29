package tui

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type spinnerDoneMsg struct{ err error }
type spinnerMinTimeDoneMsg struct{}

const minSpinnerDuration = 800 * time.Millisecond

func RunWithSpinner(msg string, fn func() error) error {
	if !isTerminal(os.Stdout) {
		fmt.Fprintf(os.Stderr, "  ⠋ %s\n", msg)
		err := fn()
		if err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "\r  ✔ %s\n", msg)
		return nil
	}

	m := newSpinnerModel(msg, fn)
	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return err
	}

	if fm, ok := finalModel.(spinnerModel); ok {
		return fm.err
	}

	return nil
}

type spinnerModel struct {
	spinner     spinner.Model
	message     string
	fn          func() error
	err         error
	done        bool
	fnDone      bool
	minTimeDone bool
}

func newSpinnerModel(msg string, fn func() error) spinnerModel {
	s := spinner.New()
	s.Spinner = spinner.MiniDot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("99"))

	return spinnerModel{
		spinner: s,
		message: msg,
		fn:      fn,
	}
}

func (m spinnerModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		func() tea.Msg {
			err := m.fn()
			return spinnerDoneMsg{err: err}
		},
		func() tea.Msg {
			time.Sleep(minSpinnerDuration)
			return spinnerMinTimeDoneMsg{}
		},
	)
}

func (m spinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinnerDoneMsg:
		m.fnDone = true
		m.err = msg.err
		if m.minTimeDone {
			m.done = true
			return m, tea.Quit
		}
		return m, nil
	case spinnerMinTimeDoneMsg:
		m.minTimeDone = true
		if m.fnDone {
			m.done = true
			return m, tea.Quit
		}
		return m, nil
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			m.done = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m spinnerModel) View() string {
	if m.done {
		if m.err != nil {
			return StyleDanger.Render(fmt.Sprintf("\r  ✗ %s: %v\n", m.message, m.err))
		}
		return StyleSuccess.Render(fmt.Sprintf("\r  ✔ %s\n", m.message))
	}
	return fmt.Sprintf("\r  %s %s", m.spinner.View(), StyleMuted.Render(m.message))
}
