package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/m-uesaka/cli-record/internal/config"
	"github.com/m-uesaka/cli-record/internal/models"
	"github.com/m-uesaka/cli-record/internal/storage"
	"github.com/m-uesaka/cli-record/internal/tui"
	"github.com/spf13/cobra"
)

var (
	viewFrom            string
	viewTo              string
	viewTags            string
	viewTask            string
	viewBy              string
	viewByHour          bool
	viewByWeekday       bool
	viewByDayOfMonth    bool
	viewByMonth         bool
	viewFormat          string
	viewOutput          string
)

var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "View aggregated time reports in various formats",
	Long: `Generate comprehensive reports of time spent on tasks with various grouping and filtering options.
Supports multiple view types and export formats.`,
	RunE: runView,
}

func init() {
	rootCmd.AddCommand(viewCmd)
	viewCmd.Flags().StringVar(&viewFrom, "from", "", "Start date (YYYY-MM-DD or YYYY-MM-DD HH:MM:SS)")
	viewCmd.Flags().StringVar(&viewTo, "to", "", "End date (YYYY-MM-DD or YYYY-MM-DD HH:MM:SS)")
	viewCmd.Flags().StringVar(&viewTags, "tags", "", "Filter by tags (comma-separated)")
	viewCmd.Flags().StringVar(&viewTask, "task", "", "Filter by task name (partial matching)")
	viewCmd.Flags().StringVar(&viewBy, "by", "task", "Primary grouping (task, tag, day, week, month, year)")
	viewCmd.Flags().BoolVar(&viewByHour, "view-by-hour", false, "Show hourly breakdown")
	viewCmd.Flags().BoolVar(&viewByWeekday, "view-by-weekday", false, "Show weekday breakdown")
	viewCmd.Flags().BoolVar(&viewByDayOfMonth, "view-by-day-of-month", false, "Show day-of-month breakdown")
	viewCmd.Flags().BoolVar(&viewByMonth, "view-by-month", false, "Show monthly breakdown")
	viewCmd.Flags().StringVar(&viewFormat, "format", "table", "Output format (table, csv, json)")
	viewCmd.Flags().StringVar(&viewOutput, "output", "", "Output file path (default: stdout)")
}

func runView(cmd *cobra.Command, args []string) error {
	// Validate date formats
	if err := ValidateDateFormat(viewFrom); err != nil {
		return err
	}
	if err := ValidateDateFormat(viewTo); err != nil {
		return err
	}

	// Validate format
	if err := validateChoice(viewFormat, []string{"table", "csv", "json"}, "format"); err != nil {
		return err
	}

	// Validate groupBy
	if err := validateChoice(viewBy, []string{"task", "tag", "day", "week", "month", "year"}, "--by value"); err != nil {
		return err
	}

	store, err := storage.NewJSONStorage()
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}

	entries, err := store.ListEntries()
	if err != nil {
		return fmt.Errorf("failed to list entries: %w", err)
	}

	// Apply filters (reuse from list command)
	listFrom = viewFrom
	listTo = viewTo
	listTags = viewTags
	listTask = viewTask
	entries, err = filterEntries(entries)
	if err != nil {
		return err
	}

	if len(entries) == 0 {
		fmt.Println("No entries found.")
		return nil
	}

	// Determine which view to use
	var report ViewReport
	if viewByHour {
		report = generateHourlyReport(entries)
	} else if viewByWeekday {
		report = generateWeekdayReport(entries)
	} else if viewByDayOfMonth {
		report = generateDayOfMonthReport(entries)
	} else if viewByMonth {
		report = generateMonthlyReport(entries)
	} else {
		report = generateGroupedReport(entries, viewBy)
	}

	// Output report
	return outputReport(report)
}

type ViewReport struct {
	Title  string
	Groups []ReportGroup
	Total  time.Duration
}

type ReportGroup struct {
	Name     string
	Duration time.Duration
	Count    int
}

func generateGroupedReport(entries []*models.TimeEntry, groupBy string) ViewReport {
	groups := make(map[string]*ReportGroup)

	for _, entry := range entries {
		if groupBy == "tag" {
			// Special handling for tags - one entry can belong to multiple groups
			if len(entry.Tags) == 0 {
				aggregateByKey(groups, "(no tags)", entry)
			} else {
				for _, tag := range entry.Tags {
					aggregateByKey(groups, tag, entry)
				}
			}
		} else {
			// Standard grouping - one entry belongs to one group
			key := getGroupKey(entry, groupBy)
			aggregateByKey(groups, key, entry)
		}
	}

	return createReport(fmt.Sprintf("Grouped by %s", groupBy), groups)
}

func generateHourlyReport(entries []*models.TimeEntry) ViewReport {
	hourlyKeys := getHourlyKeys()
	groups := initializeGroups(hourlyKeys)

	for _, entry := range entries {
		distributeEntryAcrossHours(groups, entry)
	}

	return createOrderedReport("Hourly Breakdown", groups, hourlyKeys)
}

func generateWeekdayReport(entries []*models.TimeEntry) ViewReport {
	weekdayNames := getWeekdayNames()
	groups := initializeGroups(weekdayNames)

	for _, entry := range entries {
		key := getWeekdayKey(entry.StartTime.Weekday())
		groups[key].Duration += entry.Duration()
		groups[key].Count++
	}

	return createOrderedReport("Weekday Breakdown", groups, weekdayNames)
}

func generateDayOfMonthReport(entries []*models.TimeEntry) ViewReport {
	dayKeys := getDayOfMonthKeys()
	groups := initializeGroups(dayKeys)

	for _, entry := range entries {
		key := fmt.Sprintf("Day %02d", entry.StartTime.Day())
		groups[key].Duration += entry.Duration()
		groups[key].Count++
	}

	return createOrderedReport("Day-of-Month Breakdown", groups, dayKeys)
}

func generateMonthlyReport(entries []*models.TimeEntry) ViewReport {
	monthNames := getMonthNames()
	groups := initializeGroups(monthNames)

	for _, entry := range entries {
		key := getMonthKey(entry.StartTime.Month())
		groups[key].Duration += entry.Duration()
		groups[key].Count++
	}

	return createOrderedReport("Monthly Breakdown", groups, monthNames)
}

func createReport(title string, groupMap map[string]*ReportGroup) ViewReport {
	report := ViewReport{Title: title, Groups: make([]ReportGroup, 0, len(groupMap))}

	for _, group := range groupMap {
		report.Groups = append(report.Groups, *group)
		report.Total += group.Duration
	}

	// Sort groups by name
	sort.Slice(report.Groups, func(i, j int) bool {
		return report.Groups[i].Name < report.Groups[j].Name
	})

	return report
}

// createOrderedReport creates a report with groups in the specified order
func createOrderedReport(title string, groupMap map[string]*ReportGroup, order []string) ViewReport {
	report := ViewReport{Title: title, Groups: make([]ReportGroup, 0, len(order))}

	for _, name := range order {
		if group, exists := groupMap[name]; exists {
			report.Groups = append(report.Groups, *group)
			report.Total += group.Duration
		}
	}

	return report
}

func outputReport(report ViewReport) error {
	var output *os.File
	var err error

	if viewOutput != "" {
		output, err = os.Create(viewOutput)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer output.Close()
	} else {
		output = os.Stdout
	}

	switch viewFormat {
	case "table":
		return outputTableReport(output, report)
	case "csv":
		return outputCSVReport(output, report)
	case "json":
		return outputJSONReport(output, report)
	default:
		return fmt.Errorf("invalid format: %s (use: table, csv, json)", viewFormat)
	}
}

func outputTableReport(output *os.File, report ViewReport) error {
	// Load configuration to get group colors
	cfg, _ := config.Load()
	
	// Only apply colors when outputting to terminal (stdout)
	useColors := output == os.Stdout && cfg != nil
	
	fmt.Fprintln(output, report.Title)
	fmt.Fprintln(output, strings.Repeat("=", 80))
	fmt.Fprintf(output, "%-40s %15s %12s %10s\n", "Group", "Duration", "Entries", "Percentage")
	fmt.Fprintln(output, strings.Repeat("-", 80))

	for _, group := range report.Groups {
		percentage := calculatePercentage(group.Duration, report.Total)
		groupName := truncate(group.Name, 40)
		
		// Apply color if configured
		if useColors {
			if color := cfg.GetGroupColor(group.Name); color != "" {
				groupName = tui.Colorize(groupName, color)
			}
		}
		
		fmt.Fprintf(output, "%-40s %15s %12d %9.1f%%\n",
			groupName,
			formatDuration(group.Duration),
			group.Count,
			percentage)
	}

	fmt.Fprintln(output, strings.Repeat("=", 80))
	fmt.Fprintf(output, "Total: %s (%d groups)\n", formatDuration(report.Total), len(report.Groups))

	return nil
}

func outputCSVReport(output *os.File, report ViewReport) error {
	writer := csv.NewWriter(output)
	defer writer.Flush()

	// Write header
	if err := writer.Write([]string{"Group", "Duration (seconds)", "Entries", "Percentage"}); err != nil {
		return err
	}

	// Write data
	for _, group := range report.Groups {
		percentage := calculatePercentage(group.Duration, report.Total)
		if err := writer.Write([]string{
			group.Name,
			fmt.Sprintf("%.0f", group.Duration.Seconds()),
			fmt.Sprintf("%d", group.Count),
			fmt.Sprintf("%.1f", percentage),
		}); err != nil {
			return err
		}
	}

	return nil
}

func outputJSONReport(output *os.File, report ViewReport) error {
	type JSONGroup struct {
		Name              string  `json:"group"`
		DurationSeconds   float64 `json:"duration_seconds"`
		FormattedDuration string  `json:"formatted_duration"`
		Entries           int     `json:"entries"`
		Percentage        float64 `json:"percentage"`
	}

	type JSONReport struct {
		Title  string      `json:"title"`
		Groups []JSONGroup `json:"groups"`
		Total  struct {
			DurationSeconds   float64 `json:"duration_seconds"`
			FormattedDuration string  `json:"formatted_duration"`
		} `json:"total"`
	}

	jsonReport := JSONReport{
		Title:  report.Title,
		Groups: make([]JSONGroup, 0, len(report.Groups)),
	}

	for _, group := range report.Groups {
		percentage := calculatePercentage(group.Duration, report.Total)
		jsonReport.Groups = append(jsonReport.Groups, JSONGroup{
			Name:              group.Name,
			DurationSeconds:   group.Duration.Seconds(),
			FormattedDuration: formatDuration(group.Duration),
			Entries:           group.Count,
			Percentage:        percentage,
		})
	}

	jsonReport.Total.DurationSeconds = report.Total.Seconds()
	jsonReport.Total.FormattedDuration = formatDuration(report.Total)

	encoder := json.NewEncoder(output)
	encoder.SetIndent("", "  ")
	return encoder.Encode(jsonReport)
}
