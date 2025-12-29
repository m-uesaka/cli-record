package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// AutocompleteModel provides text input with autocomplete suggestions
type AutocompleteModel struct {
	Input         textinput.Model
	Suggestions   []string
	AllSuggestions []string
	SelectedIndex int
	ShowSuggestions bool
	Title         string
	Label         string
	Submitted     bool
	Cancelled     bool
}

// NewAutocomplete creates a new autocomplete input
func NewAutocomplete(title, label, placeholder string, suggestions []string) AutocompleteModel {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.Focus()
	ti.CharLimit = 200
	ti.Width = 50

	return AutocompleteModel{
		Input:          ti,
		AllSuggestions: suggestions,
		Suggestions:    []string{},
		SelectedIndex:  0,
		ShowSuggestions: false,
		Title:          title,
		Label:          label,
	}
}

func (m AutocompleteModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m AutocompleteModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.Cancelled = true
			return m, tea.Quit

		case "enter":
			// If a suggestion is selected, use it
			if m.ShowSuggestions && len(m.Suggestions) > 0 && m.SelectedIndex >= 0 && m.SelectedIndex < len(m.Suggestions) {
				m.Input.SetValue(m.Suggestions[m.SelectedIndex])
				m.ShowSuggestions = false
				m.Suggestions = []string{}
				return m, nil
			}
			m.Submitted = true
			return m, tea.Quit

		case "down":
			if m.ShowSuggestions && len(m.Suggestions) > 0 {
				m.SelectedIndex++
				if m.SelectedIndex >= len(m.Suggestions) {
					m.SelectedIndex = 0
				}
				return m, nil
			}

		case "up":
			if m.ShowSuggestions && len(m.Suggestions) > 0 {
				m.SelectedIndex--
				if m.SelectedIndex < 0 {
					m.SelectedIndex = len(m.Suggestions) - 1
				}
				return m, nil
			}

		case "tab":
			// Auto-complete with first suggestion
			if m.ShowSuggestions && len(m.Suggestions) > 0 {
				m.Input.SetValue(m.Suggestions[0])
				m.ShowSuggestions = false
				m.Suggestions = []string{}
				return m, nil
			}
		}
	}

	var cmd tea.Cmd
	m.Input, cmd = m.Input.Update(msg)

	// Update suggestions based on input
	m.updateSuggestions()

	return m, cmd
}

func (m *AutocompleteModel) updateSuggestions() {
	value := strings.ToLower(strings.TrimSpace(m.Input.Value()))
	if value == "" {
		m.ShowSuggestions = false
		m.Suggestions = []string{}
		m.SelectedIndex = 0
		return
	}

	filtered := []string{}
	for _, suggestion := range m.AllSuggestions {
		if strings.Contains(strings.ToLower(suggestion), value) {
			filtered = append(filtered, suggestion)
		}
	}

	if len(filtered) > 0 {
		m.Suggestions = filtered
		m.ShowSuggestions = true
		if m.SelectedIndex >= len(filtered) {
			m.SelectedIndex = 0
		}
	} else {
		m.ShowSuggestions = false
		m.Suggestions = []string{}
		m.SelectedIndex = 0
	}
}

func (m AutocompleteModel) View() string {
	var b strings.Builder

	if m.Title != "" {
		b.WriteString(m.Title)
		b.WriteString("\n\n")
	}

	if m.Label != "" {
		b.WriteString(m.Label)
		b.WriteString("\n")
	}

	b.WriteString(m.Input.View())
	b.WriteString("\n")

	if m.ShowSuggestions && len(m.Suggestions) > 0 {
		b.WriteString("\nSuggestions:\n")
		maxShow := 5
		for i, suggestion := range m.Suggestions {
			if i >= maxShow {
				b.WriteString("  ...\n")
				break
			}
			if i == m.SelectedIndex {
				b.WriteString("  ▸ ")
			} else {
				b.WriteString("    ")
			}
			b.WriteString(suggestion)
			b.WriteString("\n")
		}
		b.WriteString("\nUse ↑/↓ to select • Tab to autocomplete • Enter to confirm")
	} else {
		b.WriteString("\nEnter to confirm • Esc to cancel")
	}

	return b.String()
}

// GetValue returns the current input value
func (m AutocompleteModel) GetValue() string {
	return strings.TrimSpace(m.Input.Value())
}
