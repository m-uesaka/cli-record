package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/m-uesaka/cli-record/internal/models"
)

// validateChoice validates that a value is in the list of valid choices
func validateChoice(value string, validChoices []string, fieldName string) error {
	if value == "" {
		return nil // Empty is often valid (means use default)
	}

	for _, choice := range validChoices {
		if value == choice {
			return nil
		}
	}

	return NewErrorWithSuggestion(
		fmt.Errorf("invalid %s: %s", fieldName, value),
		fmt.Sprintf("Valid values are: %s", strings.Join(validChoices, ", ")),
	)
}

// calculatePercentage safely calculates percentage avoiding division by zero
func calculatePercentage(part, total time.Duration) float64 {
	if total <= 0 {
		return 0.0
	}
	return float64(part) / float64(total) * 100
}

// initializeGroups creates a map with pre-initialized groups
func initializeGroups(names []string) map[string]*ReportGroup {
	groups := make(map[string]*ReportGroup, len(names))
	for _, name := range names {
		groups[name] = &ReportGroup{Name: name, Duration: 0, Count: 0}
	}
	return groups
}

// aggregateByKey adds entry duration to the appropriate group
func aggregateByKey(groups map[string]*ReportGroup, key string, entry *models.TimeEntry) {
	if _, exists := groups[key]; !exists {
		groups[key] = &ReportGroup{Name: key}
	}
	groups[key].Duration += entry.Duration()
	groups[key].Count++
}

// getGroupKey extracts the appropriate key for grouping based on groupBy type
func getGroupKey(entry *models.TimeEntry, groupBy string) string {
	switch groupBy {
	case "task":
		return entry.TaskName
	case "day":
		return entry.StartTime.Format("2006-01-02")
	case "week":
		year, week := entry.StartTime.ISOWeek()
		return fmt.Sprintf("%d-W%02d", year, week)
	case "month":
		return entry.StartTime.Format("2006-01")
	case "year":
		return entry.StartTime.Format("2006")
	default:
		return entry.TaskName
	}
}

// distributeEntryAcrossHours distributes an entry's duration across hourly buckets
func distributeEntryAcrossHours(groups map[string]*ReportGroup, entry *models.TimeEntry) {
	duration := entry.Duration()
	startTime := entry.StartTime
	endTime := startTime.Add(duration)

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

// getWeekdayNames returns ordered weekday names starting from Monday
func getWeekdayNames() []string {
	return []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
}

// getMonthNames returns ordered month names
func getMonthNames() []string {
	return []string{"January", "February", "March", "April", "May", "June",
		"July", "August", "September", "October", "November", "December"}
}

// getHourlyKeys returns all 24 hourly time slot keys
func getHourlyKeys() []string {
	keys := make([]string, 24)
	for h := 0; h < 24; h++ {
		keys[h] = fmt.Sprintf("%02d:00-%02d:00", h, (h+1)%24)
	}
	return keys
}

// getDayOfMonthKeys returns all day-of-month keys (1-31)
func getDayOfMonthKeys() []string {
	keys := make([]string, 31)
	for day := 1; day <= 31; day++ {
		keys[day-1] = fmt.Sprintf("Day %02d", day)
	}
	return keys
}

// getWeekdayKey converts time.Weekday to Monday-first weekday name
func getWeekdayKey(weekday time.Weekday) string {
	weekdays := getWeekdayNames()
	return weekdays[(weekday+6)%7] // Adjust so Monday is first
}

// getMonthKey converts time.Month to month name
func getMonthKey(month time.Month) string {
	months := getMonthNames()
	return months[month-1]
}
