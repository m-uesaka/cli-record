package cmd

import (
	"fmt"
	"time"

	"github.com/m-uesaka/cli-record/internal/storage"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the current time tracking status",
	Long:  `Show the current time tracking status. If a time entry is running, displays task details. If not, shows an appropriate message.`,
	RunE:  runStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func runStatus(cmd *cobra.Command, args []string) error {
	store, err := storage.NewJSONStorage()
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}

	// Check for running entry
	runningEntry, err := store.GetRunningEntry()
	if err != nil {
		return fmt.Errorf("failed to check for running entry: %w", err)
	}

	if runningEntry == nil {
		// No entry running
		fmt.Println("⏸️  No time entry is currently running")
		fmt.Println()
		fmt.Println("Use 'cli-record start' to begin tracking time.")
		return nil
	}

	// Entry is running
	elapsed := time.Since(runningEntry.StartTime)
	
	fmt.Println("⏱️  Currently tracking time")
	fmt.Println()
	fmt.Printf("Task:     %s\n", runningEntry.TaskName)
	if runningEntry.TaskName == "" {
		fmt.Printf("Task:     (not set yet)\n")
	}
	fmt.Printf("Started:  %s (%s ago)\n", 
		runningEntry.StartTime.Format("2006-01-02 15:04:05"),
		formatElapsedTime(elapsed))
	fmt.Printf("Duration: %s\n", formatDuration(elapsed))
	if len(runningEntry.Tags) > 0 {
		fmt.Printf("Tags:     %s\n", formatTags(runningEntry.Tags))
	} else {
		fmt.Printf("Tags:     (none)\n")
	}

	return nil
}

func formatElapsedTime(d time.Duration) string {
	d = d.Round(time.Second)
	
	if d < time.Minute {
		return "just now"
	}
	
	if d < time.Hour {
		minutes := int(d.Minutes())
		if minutes == 1 {
			return "1 minute"
		}
		return fmt.Sprintf("%d minutes", minutes)
	}
	
	if d < 24*time.Hour {
		hours := int(d.Hours())
		minutes := int(d.Minutes()) % 60
		if hours == 1 {
			if minutes == 0 {
				return "1 hour"
			}
			return fmt.Sprintf("1 hour %d minutes", minutes)
		}
		if minutes == 0 {
			return fmt.Sprintf("%d hours", hours)
		}
		return fmt.Sprintf("%d hours %d minutes", hours, minutes)
	}
	
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	if days == 1 {
		if hours == 0 {
			return "1 day"
		}
		return fmt.Sprintf("1 day %d hours", hours)
	}
	if hours == 0 {
		return fmt.Sprintf("%d days", days)
	}
	return fmt.Sprintf("%d days %d hours", days, hours)
}

func formatTags(tags []string) string {
	if len(tags) == 0 {
		return ""
	}
	result := ""
	for i, tag := range tags {
		if i > 0 {
			result += ", "
		}
		result += tag
	}
	return result
}
