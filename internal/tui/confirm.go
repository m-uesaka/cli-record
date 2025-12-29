package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// ConfirmModel is a simple yes/no confirmation dialog
type ConfirmModel struct {
	Message   string
	Yes       bool
	Confirmed bool
	Cancelled bool
}

// NewConfirm creates a new confirmation dialog
func NewConfirm(message string) ConfirmModel {
	return ConfirmModel{
		Message: message,
		Yes:     true, // Default to Yes
	}
}

func (m ConfirmModel) Init() tea.Cmd {
	return nil
}

func (m ConfirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "n", "N":
			m.Cancelled = true
			return m, tea.Quit

		case "enter", "y", "Y":
			m.Confirmed = true
			return m, tea.Quit

		case "left", "h":
			m.Yes = true
			return m, nil

		case "right", "l":
			m.Yes = false
			return m, nil
		}
	}

	return m, nil
}

func (m ConfirmModel) View() string {
	s := m.Message + "\n\n"

	if m.Yes {
		s += "[ Yes ]  No\n"
	} else {
		s += " Yes  [ No ]\n"
	}

	s += "\nUse ←/→ to select • Enter/Y to confirm • N/Esc to cancel"

	return s
}

// IsConfirmed returns true if user confirmed
func (m ConfirmModel) IsConfirmed() bool {
	return m.Confirmed && m.Yes
}
