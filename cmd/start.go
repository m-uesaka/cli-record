package cmd

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
	"github.com/m-uesaka/cli-record/internal/models"
	"github.com/m-uesaka/cli-record/internal/storage"
	"github.com/m-uesaka/cli-record/internal/tui"
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

func runPrompt(taskName string, tags []string, existingTags []string) (string, []string, error) {
	// Use the new reusable form component
	fields := []tui.InputFormField{
		{
			Label:       "Task Name:",
			Placeholder: "Task name (optional, can be set when stopping)",
			Value:       taskName,
		},
		{
			Label:       "Tags:",
			Placeholder: "Comma-separated tags (optional)",
			Value:       strings.Join(tags, ", "),
		},
	}

	form := tui.NewInputForm("Start Time Entry", fields)
	if len(existingTags) > 0 {
		form.HelpText = fmt.Sprintf("Existing tags: %s\n\n%s", 
			strings.Join(existingTags, ", "), form.HelpText)
	}

	p := tea.NewProgram(form)
	m, err := p.Run()
	if err != nil {
		return "", nil, fmt.Errorf("error running prompt: %w", err)
	}

	model := m.(tui.InputFormModel)
	if model.Cancelled {
		return "", nil, fmt.Errorf("cancelled by user")
	}

	values := model.GetValues()
	resultTaskName := values[0]
	resultTags := parseTags(values[1])

	return resultTaskName, resultTags, nil
}
