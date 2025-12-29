package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// InputFormField represents a single field in a form
type InputFormField struct {
	Label       string
	Placeholder string
	Value       string
	Input       textinput.Model
}

// InputFormModel is a reusable form with multiple text inputs
type InputFormModel struct {
	Fields      []InputFormField
	FocusIndex  int
	Title       string
	HelpText    string
	Submitted   bool
	Cancelled   bool
}

// NewInputForm creates a new input form
func NewInputForm(title string, fields []InputFormField) InputFormModel {
	// Initialize text inputs for each field
	for i := range fields {
		ti := textinput.New()
		ti.Placeholder = fields[i].Placeholder
		ti.CharLimit = 200
		ti.Width = 50
		if fields[i].Value != "" {
			ti.SetValue(fields[i].Value)
		}
		fields[i].Input = ti
	}

	// Focus first field
	if len(fields) > 0 {
		fields[0].Input.Focus()
	}

	return InputFormModel{
		Fields:     fields,
		FocusIndex: 0,
		Title:      title,
		HelpText:   "Press Enter to submit • Tab to switch fields • Esc to cancel",
	}
}

func (m InputFormModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m InputFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.Cancelled = true
			return m, tea.Quit

		case "enter":
			m.Submitted = true
			return m, tea.Quit

		case "tab", "shift+tab", "up", "down":
			// Switch focus between fields
			if msg.String() == "up" || msg.String() == "shift+tab" {
				m.FocusIndex--
			} else {
				m.FocusIndex++
			}

			if m.FocusIndex >= len(m.Fields) {
				m.FocusIndex = 0
			} else if m.FocusIndex < 0 {
				m.FocusIndex = len(m.Fields) - 1
			}

			// Update focus
			for i := range m.Fields {
				if i == m.FocusIndex {
					m.Fields[i].Input.Focus()
				} else {
					m.Fields[i].Input.Blur()
				}
			}

			return m, nil
		}
	}

	// Update the focused field
	var cmd tea.Cmd
	if m.FocusIndex < len(m.Fields) {
		m.Fields[m.FocusIndex].Input, cmd = m.Fields[m.FocusIndex].Input.Update(msg)
	}

	return m, cmd
}

func (m InputFormModel) View() string {
	var b strings.Builder

	if m.Title != "" {
		b.WriteString(m.Title)
		b.WriteString("\n\n")
	}

	for i, field := range m.Fields {
		if field.Label != "" {
			b.WriteString(field.Label)
			b.WriteString("\n")
		}
		b.WriteString(field.Input.View())
		if i < len(m.Fields)-1 {
			b.WriteString("\n\n")
		}
	}

	b.WriteString("\n\n")
	if m.HelpText != "" {
		b.WriteString(m.HelpText)
	}

	return b.String()
}

// GetValues returns the current values of all fields
func (m InputFormModel) GetValues() []string {
	values := make([]string, len(m.Fields))
	for i, field := range m.Fields {
		values[i] = strings.TrimSpace(field.Input.Value())
	}
	return values
}
