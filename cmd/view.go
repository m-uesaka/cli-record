package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/m-uesaka/cli-record/internal/models"
	"github.com/m-uesaka/cli-record/internal/storage"
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
	validFormats := []string{"table", "csv", "json"}
	isValidFormat := false
	for _, f := range validFormats {
		if viewFormat == f {
			isValidFormat = true
			break
		}
	}
	if !isValidFormat {
		return NewErrorWithSuggestion(
			fmt.Errorf("invalid format: %s", viewFormat),
			fmt.Sprintf("Valid formats are: %s", strings.Join(validFormats, ", ")),
		)
	}

	// Validate groupBy
	if viewBy != "" {
		validGroupBy := []string{"task", "tag", "day", "week", "month", "year"}
		isValid := false
		for _, g := range validGroupBy {
			if viewBy == g {
				isValid = true
				break
			}
		}
		if !isValid {
			return NewErrorWithSuggestion(
				fmt.Errorf("invalid --by value: %s", viewBy),
				fmt.Sprintf("Valid values are: %s", strings.Join(validGroupBy, ", ")),
			)
		}
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
		var key string
		switch groupBy {
		case "task":
			key = entry.TaskName
		case "tag":
			if len(entry.Tags) == 0 {
				key = "(no tags)"
			} else {
				for _, tag := range entry.Tags {
					if _, exists := groups[tag]; !exists {
						groups[tag] = &ReportGroup{Name: tag}
					}
					groups[tag].Duration += entry.Duration()
					groups[tag].Count++
				}
				continue
			}
		case "day":
			key = entry.StartTime.Format("2006-01-02")
		case "week":
			year, week := entry.StartTime.ISOWeek()
			key = fmt.Sprintf("%d-W%02d", year, week)
		case "month":
			key = entry.StartTime.Format("2006-01")
		case "year":
			key = entry.StartTime.Format("2006")
		default:
			key = entry.TaskName
		}

		if groupBy != "tag" {
			if _, exists := groups[key]; !exists {
				groups[key] = &ReportGroup{Name: key}
			}
			groups[key].Duration += entry.Duration()
			groups[key].Count++
		}
	}

	return createReport(fmt.Sprintf("Grouped by %s", groupBy), groups)
}

func generateHourlyReport(entries []*models.TimeEntry) ViewReport {
	groups := make(map[string]*ReportGroup)

	// Initialize all hours
	for h := 0; h < 24; h++ {
		key := fmt.Sprintf("%02d:00-%02d:00", h, (h+1)%24)
		groups[key] = &ReportGroup{Name: key, Duration: 0, Count: 0}
	}

	for _, entry := range entries {
		duration := entry.Duration()
		startTime := entry.StartTime
		endTime := startTime.Add(duration)

		// Distribute duration across hours
		current := startTime
		for current.Before(endTime) {
			hour := current.Hour()
			nextHour := current.Truncate(time.Hour).Add(time.Hour)
			if nextHour.After(endTime) {
				nextHour = endTime
			}

			segmentDuration := nextHour.Sub(current)
			key := fmt.Sprintf("%02d:00-%02d:00", hour, (hour+1)%24)
			groups[key].Duration += segmentDuration
			if segmentDuration > 0 {
				groups[key].Count++
			}

			current = nextHour
		}
	}

	return createReport("Hourly Breakdown", groups)
}

func generateWeekdayReport(entries []*models.TimeEntry) ViewReport {
	groups := make(map[string]*ReportGroup)
	weekdays := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}

	for _, day := range weekdays {
		groups[day] = &ReportGroup{Name: day, Duration: 0, Count: 0}
	}

	for _, entry := range entries {
		weekday := entry.StartTime.Weekday()
		day := weekdays[(weekday+6)%7] // Adjust so Monday is first
		groups[day].Duration += entry.Duration()
		groups[day].Count++
	}

	return createReport("Weekday Breakdown", groups)
}

func generateDayOfMonthReport(entries []*models.TimeEntry) ViewReport {
	groups := make(map[string]*ReportGroup)

	for day := 1; day <= 31; day++ {
		key := fmt.Sprintf("Day %02d", day)
		groups[key] = &ReportGroup{Name: key, Duration: 0, Count: 0}
	}

	for _, entry := range entries {
		day := entry.StartTime.Day()
		key := fmt.Sprintf("Day %02d", day)
		groups[key].Duration += entry.Duration()
		groups[key].Count++
	}

	return createReport("Day-of-Month Breakdown", groups)
}

func generateMonthlyReport(entries []*models.TimeEntry) ViewReport {
	groups := make(map[string]*ReportGroup)
	months := []string{"January", "February", "March", "April", "May", "June",
		"July", "August", "September", "October", "November", "December"}

	for _, month := range months {
		groups[month] = &ReportGroup{Name: month, Duration: 0, Count: 0}
	}

	for _, entry := range entries {
		month := months[entry.StartTime.Month()-1]
		groups[month].Duration += entry.Duration()
		groups[month].Count++
	}

	return createReport("Monthly Breakdown", groups)
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
	fmt.Fprintln(output, report.Title)
	fmt.Fprintln(output, strings.Repeat("=", 80))
	fmt.Fprintf(output, "%-40s %15s %12s %10s\n", "Group", "Duration", "Entries", "Percentage")
	fmt.Fprintln(output, strings.Repeat("-", 80))

	for _, group := range report.Groups {
		percentage := 0.0
		if report.Total > 0 {
			percentage = float64(group.Duration) / float64(report.Total) * 100
		}
		fmt.Fprintf(output, "%-40s %15s %12d %9.1f%%\n",
			truncate(group.Name, 40),
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
		percentage := 0.0
		if report.Total > 0 {
			percentage = float64(group.Duration) / float64(report.Total) * 100
		}
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
		percentage := 0.0
		if report.Total > 0 {
			percentage = float64(group.Duration) / float64(report.Total) * 100
		}
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
