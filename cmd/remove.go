package cmd

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/m-uesaka/cli-record/internal/models"
	"github.com/m-uesaka/cli-record/internal/storage"
	"github.com/m-uesaka/cli-record/internal/tui"
	"github.com/spf13/cobra"
)

var removeForce bool

var removeCmd = &cobra.Command{
	Use:     "remove <ID|prefix>",
	Aliases: []string{"rm"},
	Short:   "Remove a time entry",
	Long:    `Permanently delete a time entry from the database. You can specify the full ID or a unique prefix. Prompts for confirmation by default.`,
	Args:    cobra.ExactArgs(1),
	RunE:    runRemove,
}

func init() {
	rootCmd.AddCommand(removeCmd)
	removeCmd.Flags().BoolVarP(&removeForce, "force", "f", false, "Skip confirmation prompt")
}

func runRemove(cmd *cobra.Command, args []string) error {
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

	// If not using --force, show confirmation
	if !removeForce {
		confirmed, err := showRemoveConfirmation(entry)
		if err != nil {
			return err
		}
		if !confirmed {
			fmt.Println("Deletion cancelled.")
			return nil
		}
	}

	// Delete the entry using the actual ID
	if err := store.DeleteEntry(entry.ID); err != nil {
		return fmt.Errorf("failed to delete entry: %w", err)
	}

	// Display success message
	fmt.Println("✓ Entry removed successfully")
	fmt.Printf("  Task: %s\n", entry.TaskName)
	fmt.Printf("  Duration: %s\n", formatDuration(entry.Duration()))

	return nil
}

func showRemoveConfirmation(entry *models.TimeEntry) (bool, error) {
	message := fmt.Sprintf(`Are you sure you want to remove this entry?

Task:     %s
Started:  %s
Duration: %s
Tags:     %s

This action cannot be undone.`,
		entry.TaskName,
		entry.StartTime.Format("2006-01-02 15:04:05"),
		formatDuration(entry.Duration()),
		formatTagsOrNone(entry.Tags))

	confirm := tui.NewConfirm(message)

	p := tea.NewProgram(confirm)
	m, err := p.Run()
	if err != nil {
		return false, fmt.Errorf("error running confirmation: %w", err)
	}

	model := m.(tui.ConfirmModel)
	return model.IsConfirmed(), nil
}

func formatTagsOrNone(tags []string) string {
	if len(tags) == 0 {
		return "(none)"
	}
	return strings.Join(tags, ", ")
}
