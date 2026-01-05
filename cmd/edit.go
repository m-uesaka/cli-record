package cmd

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/m-uesaka/cli-record/internal/models"
	"github.com/m-uesaka/cli-record/internal/storage"
	"github.com/m-uesaka/cli-record/internal/tui"
	"github.com/spf13/cobra"
)

var (
	editTask        string
	editTags        string
	editStart       string
	editEnd         string
	editInteractive bool
)

var editCmd = &cobra.Command{
	Use:   "edit <ID|prefix>",
	Short: "Edit an existing time entry",
	Long:  `Edit an existing time entry's details including start time, end time, task name, and tags. You can specify the full ID or a unique prefix.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runEdit,
}

func init() {
	rootCmd.AddCommand(editCmd)
	editCmd.Flags().StringVarP(&editTask, "task", "t", "", "New task name")
	editCmd.Flags().StringVar(&editTags, "tags", "", "New tags (comma-separated)")
	editCmd.Flags().StringVar(&editStart, "start", "", "New start time (YYYY-MM-DD HH:MM:SS)")
	editCmd.Flags().StringVar(&editEnd, "end", "", "New end time (YYYY-MM-DD HH:MM:SS)")
	editCmd.Flags().BoolVarP(&editInteractive, "interactive", "i", true, "Use interactive mode")
}

func runEdit(cmd *cobra.Command, args []string) error {
	entryID := args[0]

	store, err := storage.NewJSONStorage()
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}

	// Try to get entry by prefix first, fall back to exact ID
	entry, err := store.GetEntryByPrefix(entryID)
	if err != nil {
		// If prefix fails, try exact ID
		entry, err = store.GetEntry(entryID)
		if err != nil {
			return fmt.Errorf("failed to get entry: %w", err)
		}
	}

	// Check if using command-line flags or interactive mode
	hasFlags := cmd.Flags().Changed("task") || cmd.Flags().Changed("tags") ||
		cmd.Flags().Changed("start") || cmd.Flags().Changed("end")

	if hasFlags {
		// Update using command-line flags
		if cmd.Flags().Changed("task") {
			entry.TaskName = editTask
		}
		if cmd.Flags().Changed("tags") {
			entry.Tags = parseTags(editTags)
		}
		if cmd.Flags().Changed("start") {
			startTime, err := parseEditDateTime(editStart)
			if err != nil {
				return NewErrorWithSuggestion(
					fmt.Errorf("invalid start time: %s", editStart),
					"Use format: YYYY-MM-DD HH:MM:SS\nExample: 2025-12-29 14:30:00",
				)
			}
			entry.StartTime = startTime
		}
		if cmd.Flags().Changed("end") {
			endTime, err := parseEditDateTime(editEnd)
			if err != nil {
				return NewErrorWithSuggestion(
					fmt.Errorf("invalid end time: %s", editEnd),
					"Use format: YYYY-MM-DD HH:MM:SS\nExample: 2025-12-29 16:30:00",
				)
			}
			entry.EndTime = &endTime
		}
	} else if editInteractive {
		// Use interactive mode
		updatedEntry, err := runEditPrompt(entry)
		if err != nil {
			return err
		}
		entry = updatedEntry
	}

	// Validate: end time must be after start time
	if entry.EndTime != nil && entry.EndTime.Before(entry.StartTime) {
		return fmt.Errorf("end time must be after start time")
	}

	// Update entry
	if err := store.UpdateEntry(entry); err != nil {
		return fmt.Errorf("failed to update entry: %w", err)
	}

	// Display success message
	fmt.Println("✓ Entry updated successfully")
	fmt.Printf("  Task: %s\n", entry.TaskName)
	fmt.Printf("  Start: %s\n", entry.StartTime.Format("2006-01-02 15:04:05"))
	if entry.EndTime != nil {
		fmt.Printf("  End: %s\n", entry.EndTime.Format("2006-01-02 15:04:05"))
		fmt.Printf("  Duration: %s\n", formatDuration(entry.Duration()))
	} else {
		fmt.Printf("  Status: Running\n")
	}
	if len(entry.Tags) > 0 {
		fmt.Printf("  Tags: %s\n", strings.Join(entry.Tags, ", "))
	}

	return nil
}

func runEditPrompt(entry *models.TimeEntry) (*models.TimeEntry, error) {
	endTimeStr := ""
	if entry.EndTime != nil {
		endTimeStr = entry.EndTime.Format("2006-01-02 15:04:05")
	}

	fields := []tui.InputFormField{
		{
			Label:       "Task Name:",
			Placeholder: "Enter task name",
			Value:       entry.TaskName,
		},
		{
			Label:       "Tags (comma-separated):",
			Placeholder: "Enter tags",
			Value:       strings.Join(entry.Tags, ", "),
		},
		{
			Label:       "Start Time (YYYY-MM-DD HH:MM:SS):",
			Placeholder: "Enter start time",
			Value:       entry.StartTime.Format("2006-01-02 15:04:05"),
		},
		{
			Label:       "End Time (YYYY-MM-DD HH:MM:SS):",
			Placeholder: "Leave empty if still running",
			Value:       endTimeStr,
		},
	}

	form := tui.NewInputForm("Edit Time Entry", fields)
	form.HelpText = "Press Enter to save • Tab to switch fields • Esc to cancel"

	p := tea.NewProgram(form)
	m, err := p.Run()
	if err != nil {
		return nil, fmt.Errorf("error running prompt: %w", err)
	}

	model := m.(tui.InputFormModel)
	if model.Cancelled {
		return nil, fmt.Errorf("edit cancelled by user")
	}

	values := model.GetValues()

	// Parse the values
	entry.TaskName = values[0]
	entry.Tags = parseTags(values[1])

	startTime, err := parseEditDateTime(values[2])
	if err != nil {
		return nil, fmt.Errorf("invalid start time format: %w", err)
	}
	entry.StartTime = startTime

	if values[3] != "" {
		endTime, err := parseEditDateTime(values[3])
		if err != nil {
			return nil, fmt.Errorf("invalid end time format: %w", err)
		}
		entry.EndTime = &endTime
	} else {
		entry.EndTime = nil
	}

	return entry, nil
}

func parseEditDateTime(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, fmt.Errorf("date string is empty")
	}

	layout := "2006-01-02 15:04:05"
	t, err := time.Parse(layout, dateStr)
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}
