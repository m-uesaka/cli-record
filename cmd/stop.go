package cmd

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/m-uesaka/cli-record/internal/storage"
	"github.com/m-uesaka/cli-record/internal/tui"
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

func promptTaskName() (string, error) {
	fields := []tui.InputFormField{
		{
			Label:       "Task Name:",
			Placeholder: "Enter task name",
		},
	}

	form := tui.NewInputForm("Task Name Required", fields)
	form.HelpText = "The task name was not provided when starting.\nPlease enter it now.\n\n" + form.HelpText

	p := tea.NewProgram(form)
	m, err := p.Run()
	if err != nil {
		return "", fmt.Errorf("error running prompt: %w", err)
	}

	model := m.(tui.InputFormModel)
	if model.Cancelled {
		return "", fmt.Errorf("cancelled by user")
	}

	values := model.GetValues()
	taskName := values[0]

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
