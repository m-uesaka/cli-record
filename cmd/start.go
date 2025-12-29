package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
	"github.com/m-uesaka/cli-record/internal/models"
	"github.com/m-uesaka/cli-record/internal/storage"
	"github.com/spf13/cobra"
)

var (
	startTaskName string
	startTags     string
)

var startCmd = &cobra.Command{
	Use:     "start",
	Aliases: []string{"st"},
	Short:   "Start recording time for a new task",
	Long: `Start recording time for a new task. You can optionally provide task name and tags.
If not provided, you will be prompted to enter them interactively.`,
	RunE: runStart,
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringVarP(&startTaskName, "task", "t", "", "Task name (optional)")
	startCmd.Flags().StringVar(&startTags, "tags", "", "Comma-separated tags (optional)")
}

func runStart(cmd *cobra.Command, args []string) error {
	store, err := storage.NewJSONStorage()
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}

	// Check for running entry
	runningEntry, err := store.GetRunningEntry()
	if err != nil {
		return fmt.Errorf("failed to check for running entry: %w", err)
	}

	if runningEntry != nil {
		return fmt.Errorf("a time entry is already running (ID: %s, Task: %s). Please stop it first", 
			runningEntry.ID, runningEntry.TaskName)
	}

	// Get task name and tags
	taskName := startTaskName
	var tags []string

	if startTags != "" {
		tags = parseTags(startTags)
	}

	// If task name or tags not provided, use interactive prompt
	if taskName == "" || len(tags) == 0 {
		existingTags, err := store.ListTags()
		if err != nil {
			return fmt.Errorf("failed to load existing tags: %w", err)
		}

		promptTaskName, promptTags, err := runPrompt(taskName, tags, existingTags)
		if err != nil {
			return err
		}

		if taskName == "" {
			taskName = promptTaskName
		}
		if len(tags) == 0 {
			tags = promptTags
		}
	}

	// Create new entry
	entry := &models.TimeEntry{
		ID:        uuid.New().String(),
		StartTime: time.Now(),
		TaskName:  taskName,
		Tags:      tags,
	}

	// Save entry
	if err := store.SaveEntry(entry); err != nil {
		return fmt.Errorf("failed to save entry: %w", err)
	}

	// Display success message
	fmt.Printf("✓ Started tracking time\n")
	fmt.Printf("  ID: %s\n", entry.ID)
	if entry.TaskName != "" {
		fmt.Printf("  Task: %s\n", entry.TaskName)
	} else {
		fmt.Printf("  Task: (will be set when stopping)\n")
	}
	if len(entry.Tags) > 0 {
		fmt.Printf("  Tags: %s\n", strings.Join(entry.Tags, ", "))
	}

	return nil
}

func parseTags(tagStr string) []string {
	if tagStr == "" {
		return []string{}
	}

	parts := strings.Split(tagStr, ",")
	tags := make([]string, 0, len(parts))
	for _, tag := range parts {
		trimmed := strings.TrimSpace(tag)
		if trimmed != "" {
			tags = append(tags, trimmed)
		}
	}
	return tags
}

type promptModel struct {
	taskInput textinput.Model
	tagsInput textinput.Model
	focusIndex int
	existingTags []string
	err error
}

func initialPromptModel(taskName string, tags []string, existingTags []string) promptModel {
	ti := textinput.New()
	ti.Placeholder = "Task name (optional, can be set when stopping)"
	ti.Focus()
	ti.CharLimit = 200
	ti.Width = 50
	if taskName != "" {
		ti.SetValue(taskName)
	}

	tagsInput := textinput.New()
	tagsInput.Placeholder = "Tags (comma-separated, optional)"
	tagsInput.CharLimit = 200
	tagsInput.Width = 50
	if len(tags) > 0 {
		tagsInput.SetValue(strings.Join(tags, ", "))
	}

	return promptModel{
		taskInput: ti,
		tagsInput: tagsInput,
		focusIndex: 0,
		existingTags: existingTags,
	}
}

func (m promptModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m promptModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "enter":
			return m, tea.Quit

		case "tab", "shift+tab", "up", "down":
			if msg.String() == "up" || msg.String() == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > 1 {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = 1
			}

			if m.focusIndex == 0 {
				m.taskInput.Focus()
				m.tagsInput.Blur()
			} else {
				m.taskInput.Blur()
				m.tagsInput.Focus()
			}

			return m, nil
		}
	}

	var cmd tea.Cmd
	if m.focusIndex == 0 {
		m.taskInput, cmd = m.taskInput.Update(msg)
	} else {
		m.tagsInput, cmd = m.tagsInput.Update(msg)
	}

	return m, cmd
}

func (m promptModel) View() string {
	var b strings.Builder

	b.WriteString("Start Time Entry\n\n")
	b.WriteString(m.taskInput.View())
	b.WriteString("\n\n")
	b.WriteString(m.tagsInput.View())
	b.WriteString("\n\n")

	if len(m.existingTags) > 0 {
		b.WriteString(fmt.Sprintf("Existing tags: %s\n\n", strings.Join(m.existingTags, ", ")))
	}

	b.WriteString("Press Enter to start • Tab to switch fields • Esc to cancel")

	return b.String()
}

func runPrompt(taskName string, tags []string, existingTags []string) (string, []string, error) {
	p := tea.NewProgram(initialPromptModel(taskName, tags, existingTags))
	m, err := p.Run()
	if err != nil {
		return "", nil, fmt.Errorf("error running prompt: %w", err)
	}

	model := m.(promptModel)
	resultTaskName := strings.TrimSpace(model.taskInput.Value())
	resultTags := parseTags(model.tagsInput.Value())

	return resultTaskName, resultTags, nil
}
