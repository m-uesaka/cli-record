package cmd

import (
	"fmt"
	"strings"

	"github.com/m-uesaka/cli-record/internal/storage"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show <ID>",
	Short: "Display detailed information about a specific time entry",
	Long:  `Display comprehensive details of a single time entry identified by its unique ID.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runShow,
}

func init() {
	rootCmd.AddCommand(showCmd)
}

func runShow(cmd *cobra.Command, args []string) error {
	entryID := args[0]

	store, err := storage.NewJSONStorage()
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}

	entry, err := store.GetEntry(entryID)
	if err != nil {
		return fmt.Errorf("failed to get entry: %w", err)
	}

	// Display detailed information
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("Time Entry Details")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("ID:         %s\n", entry.ID)
	fmt.Printf("Task Name:  %s\n", entry.TaskName)
	fmt.Printf("Start Time: %s\n", entry.StartTime.Format("2006-01-02 15:04:05 (Mon)"))

	if entry.EndTime != nil {
		fmt.Printf("End Time:   %s\n", entry.EndTime.Format("2006-01-02 15:04:05 (Mon)"))
		fmt.Printf("Duration:   %s\n", formatDuration(entry.Duration()))
	} else {
		fmt.Printf("End Time:   Running\n")
		fmt.Printf("Duration:   %s (still running)\n", formatDuration(entry.Duration()))
	}

	if len(entry.Tags) > 0 {
		fmt.Printf("Tags:       %s\n", strings.Join(entry.Tags, ", "))
	} else {
		fmt.Printf("Tags:       (none)\n")
	}

	fmt.Println(strings.Repeat("=", 60))

	return nil
}
