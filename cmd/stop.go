package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/m-uesaka/cli-record/internal/storage"
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:     "stop",
	Aliases: []string{"sp"},
	Short:   "Stop the currently running time entry",
	Long:    `Stop the currently running time entry and save the recorded time.`,
	RunE:    runStop,
}

func init() {
	rootCmd.AddCommand(stopCmd)
}

func runStop(cmd *cobra.Command, args []string) error {
	store, err := storage.NewJSONStorage()
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}

	// Get running entry
	runningEntry, err := store.GetRunningEntry()
	if err != nil {
		return fmt.Errorf("failed to check for running entry: %w", err)
	}

	if runningEntry == nil {
		return fmt.Errorf("no time entry is currently running")
	}

	// If task name is missing, prompt for it
	if runningEntry.TaskName == "" {
		taskName, err := promptTaskName()
		if err != nil {
			return err
		}
		runningEntry.TaskName = taskName
	}

	// Set end time
	now := time.Now()
	runningEntry.EndTime = &now

	// Update entry
	if err := store.UpdateEntry(runningEntry); err != nil {
		return fmt.Errorf("failed to update entry: %w", err)
	}

	// Display summary
	duration := runningEntry.Duration()
	fmt.Printf("✓ Stopped tracking time\n")
	fmt.Printf("  Task: %s\n", runningEntry.TaskName)
	fmt.Printf("  Duration: %s\n", formatDuration(duration))
	if len(runningEntry.Tags) > 0 {
		fmt.Printf("  Tags: %s\n", strings.Join(runningEntry.Tags, ", "))
	}

	return nil
}

type taskNameModel struct {
	input textinput.Model
}

func initialTaskNameModel() taskNameModel {
	ti := textinput.New()
	ti.Placeholder = "Enter task name"
	ti.Focus()
	ti.CharLimit = 200
	ti.Width = 50

	return taskNameModel{
		input: ti,
	}
}

func (m taskNameModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m taskNameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m taskNameModel) View() string {
	var b strings.Builder

	b.WriteString("Task Name Required\n\n")
	b.WriteString("The task name was not provided when starting.\n")
	b.WriteString("Please enter it now:\n\n")
	b.WriteString(m.input.View())
	b.WriteString("\n\nPress Enter to confirm • Esc to cancel")

	return b.String()
}

func promptTaskName() (string, error) {
	p := tea.NewProgram(initialTaskNameModel())
	m, err := p.Run()
	if err != nil {
		return "", fmt.Errorf("error running prompt: %w", err)
	}

	model := m.(taskNameModel)
	taskName := strings.TrimSpace(model.input.Value())

	if taskName == "" {
		return "", fmt.Errorf("task name is required")
	}

	return taskName, nil
}

func formatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second

	if h > 0 {
		return fmt.Sprintf("%dh %dm %ds", h, m, s)
	} else if m > 0 {
		return fmt.Sprintf("%dm %ds", m, s)
	}
	return fmt.Sprintf("%ds", s)
}
