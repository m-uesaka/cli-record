package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/m-uesaka/cli-record/internal/models"
	"github.com/m-uesaka/cli-record/internal/storage"
	"github.com/spf13/cobra"
)

var (
	listFrom    string
	listTo      string
	listTags    string
	listTask    string
	listGroupBy string
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "Display a list of recorded time entries",
	Long: `Display a list of recorded time entries with optional filtering by date range, tags, and task names.
Also displays total time spent with optional grouping.`,
	RunE: runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVar(&listFrom, "from", "", "Start date (YYYY-MM-DD or YYYY-MM-DD HH:MM:SS)")
	listCmd.Flags().StringVar(&listTo, "to", "", "End date (YYYY-MM-DD or YYYY-MM-DD HH:MM:SS)")
	listCmd.Flags().StringVar(&listTags, "tags", "", "Filter by tags (comma-separated)")
	listCmd.Flags().StringVar(&listTask, "task", "", "Filter by task name (partial matching)")
	listCmd.Flags().StringVar(&listGroupBy, "group-by", "", "Group entries (task, tag, date)")
}

func runList(cmd *cobra.Command, args []string) error {
	store, err := storage.NewJSONStorage()
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}

	entries, err := store.ListEntries()
	if err != nil {
		return fmt.Errorf("failed to list entries: %w", err)
	}

	// Apply filters
	entries, err = filterEntries(entries)
	if err != nil {
		return err
	}

	if len(entries) == 0 {
		fmt.Println("No entries found.")
		return nil
	}

	// Display entries
	if listGroupBy != "" {
		displayGroupedEntries(entries, listGroupBy)
	} else {
		displayEntries(entries)
	}

	return nil
}

func filterEntries(entries []*models.TimeEntry) ([]*models.TimeEntry, error) {
	var filtered []*models.TimeEntry

	// Parse date filters
	var fromTime, toTime *time.Time
	var err error

	if listFrom != "" {
		t, err := parseDateTime(listFrom)
		if err != nil {
			return nil, fmt.Errorf("invalid --from date: %w", err)
		}
		fromTime = &t
	}

	if listTo != "" {
		t, err := parseDateTime(listTo)
		if err != nil {
			return nil, fmt.Errorf("invalid --to date: %w", err)
		}
		toTime = &t
	}

	// Parse tag filters
	filterTags := parseTags(listTags)

	for _, entry := range entries {
		// Date range filter
		if fromTime != nil && entry.StartTime.Before(*fromTime) {
			continue
		}
		if toTime != nil && entry.StartTime.After(*toTime) {
			continue
		}

		// Tags filter
		if len(filterTags) > 0 {
			hasTag := false
			for _, filterTag := range filterTags {
				for _, entryTag := range entry.Tags {
					if strings.EqualFold(entryTag, filterTag) {
						hasTag = true
						break
					}
				}
				if hasTag {
					break
				}
			}
			if !hasTag {
				continue
			}
		}

		// Task name filter
		if listTask != "" {
			if !strings.Contains(strings.ToLower(entry.TaskName), strings.ToLower(listTask)) {
				continue
			}
		}

		filtered = append(filtered, entry)
	}

	return filtered, err
}

func displayEntries(entries []*models.TimeEntry) {
	fmt.Println("Time Entries")
	fmt.Println(strings.Repeat("=", 100))
	fmt.Printf("%-36s %-20s %-20s %-12s %-20s %s\n", 
		"ID", "Start Time", "End Time", "Duration", "Task", "Tags")
	fmt.Println(strings.Repeat("-", 100))

	var totalDuration time.Duration

	for _, entry := range entries {
		endTimeStr := "Running"
		if entry.EndTime != nil {
			endTimeStr = entry.EndTime.Format("2006-01-02 15:04:05")
		}

		duration := entry.Duration()
		totalDuration += duration

		tagsStr := strings.Join(entry.Tags, ", ")
		if tagsStr == "" {
			tagsStr = "-"
		}

		// Truncate ID for display
		shortID := entry.ID
		if len(shortID) > 8 {
			shortID = shortID[:8] + "..."
		}

		fmt.Printf("%-36s %-20s %-20s %-12s %-20s %s\n",
			shortID,
			entry.StartTime.Format("2006-01-02 15:04:05"),
			endTimeStr,
			formatDuration(duration),
			truncate(entry.TaskName, 20),
			truncate(tagsStr, 20))
	}

	fmt.Println(strings.Repeat("=", 100))
	fmt.Printf("Total: %s (%d entries)\n", formatDuration(totalDuration), len(entries))
}

func displayGroupedEntries(entries []*models.TimeEntry, groupBy string) {
	groups := make(map[string][]*models.TimeEntry)

	for _, entry := range entries {
		var key string
		switch groupBy {
		case "task":
			key = entry.TaskName
		case "tag":
			if len(entry.Tags) == 0 {
				groups["(no tags)"] = append(groups["(no tags)"], entry)
			} else {
				for _, tag := range entry.Tags {
					groups[tag] = append(groups[tag], entry)
				}
			}
			continue
		case "date":
			key = entry.StartTime.Format("2006-01-02")
		default:
			fmt.Printf("Invalid group-by value: %s (use: task, tag, date)\n", groupBy)
			return
		}
		groups[key] = append(groups[key], entry)
	}

	fmt.Printf("Time Entries (Grouped by %s)\n", groupBy)
	fmt.Println(strings.Repeat("=", 80))

	var grandTotal time.Duration

	for groupName, groupEntries := range groups {
		var groupTotal time.Duration
		for _, entry := range groupEntries {
			groupTotal += entry.Duration()
		}
		grandTotal += groupTotal

		fmt.Printf("\n%s: %s (%d entries)\n", groupName, formatDuration(groupTotal), len(groupEntries))
		fmt.Println(strings.Repeat("-", 80))

		for _, entry := range groupEntries {
			endTimeStr := "Running"
			if entry.EndTime != nil {
				endTimeStr = entry.EndTime.Format("15:04:05")
			}

			fmt.Printf("  %-20s - %-10s  %-12s  %s\n",
				entry.StartTime.Format("2006-01-02 15:04:05"),
				endTimeStr,
				formatDuration(entry.Duration()),
				truncate(entry.TaskName, 40))
		}
	}

	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("Grand Total: %s\n", formatDuration(grandTotal))
}

func parseDateTime(dateStr string) (time.Time, error) {
	// Try parsing with time first
	layouts := []string{
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, layout := range layouts {
		t, err := time.Parse(layout, dateStr)
		if err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("invalid date format (use YYYY-MM-DD or YYYY-MM-DD HH:MM:SS)")
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
