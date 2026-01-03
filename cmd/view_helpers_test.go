package cmd

import (
	"testing"
	"time"

	"github.com/m-uesaka/cli-record/internal/models"
)

func TestValidateChoice(t *testing.T) {
	tests := []struct {
		name          string
		value         string
		validChoices  []string
		fieldName     string
		expectError   bool
	}{
		{
			name:         "valid choice",
			value:        "table",
			validChoices: []string{"table", "csv", "json"},
			fieldName:    "format",
			expectError:  false,
		},
		{
			name:         "empty value - allowed",
			value:        "",
			validChoices: []string{"table", "csv", "json"},
			fieldName:    "format",
			expectError:  false,
		},
		{
			name:         "invalid choice",
			value:        "xml",
			validChoices: []string{"table", "csv", "json"},
			fieldName:    "format",
			expectError:  true,
		},
		{
			name:         "case sensitive",
			value:        "Table",
			validChoices: []string{"table", "csv", "json"},
			fieldName:    "format",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateChoice(tt.value, tt.validChoices, tt.fieldName)
			if tt.expectError && err == nil {
				t.Error("expected error but got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestCalculatePercentage(t *testing.T) {
	tests := []struct {
		name     string
		part     time.Duration
		total    time.Duration
		expected float64
	}{
		{
			name:     "normal calculation",
			part:     30 * time.Minute,
			total:    60 * time.Minute,
			expected: 50.0,
		},
		{
			name:     "zero total - avoid division by zero",
			part:     30 * time.Minute,
			total:    0,
			expected: 0.0,
		},
		{
			name:     "zero part",
			part:     0,
			total:    60 * time.Minute,
			expected: 0.0,
		},
		{
			name:     "both zero",
			part:     0,
			total:    0,
			expected: 0.0,
		},
		{
			name:     "100 percent",
			part:     1 * time.Hour,
			total:    1 * time.Hour,
			expected: 100.0,
		},
		{
			name:     "small percentage",
			part:     1 * time.Minute,
			total:    100 * time.Minute,
			expected: 1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculatePercentage(tt.part, tt.total)
			if result != tt.expected {
				t.Errorf("expected %.2f, got %.2f", tt.expected, result)
			}
		})
	}
}

func TestInitializeGroups(t *testing.T) {
	tests := []struct {
		name     string
		names    []string
		expected int
	}{
		{
			name:     "empty list",
			names:    []string{},
			expected: 0,
		},
		{
			name:     "single group",
			names:    []string{"Group1"},
			expected: 1,
		},
		{
			name:     "multiple groups",
			names:    []string{"Group1", "Group2", "Group3"},
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			groups := initializeGroups(tt.names)
			if len(groups) != tt.expected {
				t.Errorf("expected %d groups, got %d", tt.expected, len(groups))
			}

			for _, name := range tt.names {
				group, exists := groups[name]
				if !exists {
					t.Errorf("expected group %s to exist", name)
				}
				if group.Name != name {
					t.Errorf("expected group name %s, got %s", name, group.Name)
				}
				if group.Duration != 0 {
					t.Errorf("expected duration to be 0, got %v", group.Duration)
				}
				if group.Count != 0 {
					t.Errorf("expected count to be 0, got %d", group.Count)
				}
			}
		})
	}
}

func TestAggregateByKey(t *testing.T) {
	entry := &models.TimeEntry{
		ID:        "test-1",
		TaskName:  "Test Task",
		StartTime: time.Now(),
		EndTime:   timePtr(time.Now().Add(1 * time.Hour)),
	}

	t.Run("add to new group", func(t *testing.T) {
		groups := make(map[string]*ReportGroup)
		aggregateByKey(groups, "key1", entry)

		if len(groups) != 1 {
			t.Fatalf("expected 1 group, got %d", len(groups))
		}

		group := groups["key1"]
		if group.Name != "key1" {
			t.Errorf("expected name 'key1', got %s", group.Name)
		}
		// Allow for small nanosecond differences
		if group.Duration < 59*time.Minute || group.Duration > 61*time.Minute {
			t.Errorf("expected duration ~1h, got %v", group.Duration)
		}
		if group.Count != 1 {
			t.Errorf("expected count 1, got %d", group.Count)
		}
	})

	t.Run("add to existing group", func(t *testing.T) {
		groups := make(map[string]*ReportGroup)
		groups["key1"] = &ReportGroup{Name: "key1", Duration: 30 * time.Minute, Count: 1}

		aggregateByKey(groups, "key1", entry)

		group := groups["key1"]
		// Allow for small nanosecond differences
		if group.Duration < 89*time.Minute || group.Duration > 91*time.Minute {
			t.Errorf("expected duration ~90m, got %v", group.Duration)
		}
		if group.Count != 2 {
			t.Errorf("expected count 2, got %d", group.Count)
		}
	})
}

func TestGetGroupKey(t *testing.T) {
	testTime := time.Date(2024, 3, 15, 14, 30, 0, 0, time.UTC)
	entry := &models.TimeEntry{
		ID:        "test-1",
		TaskName:  "Test Task",
		StartTime: testTime,
	}

	tests := []struct {
		name     string
		groupBy  string
		expected string
	}{
		{
			name:     "group by task",
			groupBy:  "task",
			expected: "Test Task",
		},
		{
			name:     "group by day",
			groupBy:  "day",
			expected: "2024-03-15",
		},
		{
			name:     "group by week",
			groupBy:  "week",
			expected: "2024-W11",
		},
		{
			name:     "group by month",
			groupBy:  "month",
			expected: "2024-03",
		},
		{
			name:     "group by year",
			groupBy:  "year",
			expected: "2024",
		},
		{
			name:     "invalid groupBy defaults to task",
			groupBy:  "invalid",
			expected: "Test Task",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getGroupKey(entry, tt.groupBy)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestDistributeEntryAcrossHours(t *testing.T) {
	t.Run("entry within single hour", func(t *testing.T) {
		startTime := time.Date(2024, 1, 1, 14, 10, 0, 0, time.UTC)
		entry := &models.TimeEntry{
			ID:        "test-1",
			TaskName:  "Test",
			StartTime: startTime,
			EndTime:   timePtr(startTime.Add(30 * time.Minute)),
		}

		groups := initializeGroups(getHourlyKeys())
		distributeEntryAcrossHours(groups, entry)

		key := "14:00-15:00"
		if groups[key].Duration != 30*time.Minute {
			t.Errorf("expected 30m in hour 14, got %v", groups[key].Duration)
		}
		if groups[key].Count != 1 {
			t.Errorf("expected count 1, got %d", groups[key].Count)
		}
	})

	t.Run("entry spanning multiple hours", func(t *testing.T) {
		startTime := time.Date(2024, 1, 1, 14, 45, 0, 0, time.UTC)
		entry := &models.TimeEntry{
			ID:        "test-1",
			TaskName:  "Test",
			StartTime: startTime,
			EndTime:   timePtr(startTime.Add(90 * time.Minute)),
		}

		groups := initializeGroups(getHourlyKeys())
		distributeEntryAcrossHours(groups, entry)

		// Hour 14: 15 minutes (14:45-15:00)
		if groups["14:00-15:00"].Duration != 15*time.Minute {
			t.Errorf("expected 15m in hour 14, got %v", groups["14:00-15:00"].Duration)
		}

		// Hour 15: 60 minutes (15:00-16:00)
		if groups["15:00-16:00"].Duration != 60*time.Minute {
			t.Errorf("expected 60m in hour 15, got %v", groups["15:00-16:00"].Duration)
		}

		// Hour 16: 15 minutes (16:00-16:15)
		if groups["16:00-17:00"].Duration != 15*time.Minute {
			t.Errorf("expected 15m in hour 16, got %v", groups["16:00-17:00"].Duration)
		}
	})

	t.Run("entry spanning midnight", func(t *testing.T) {
		startTime := time.Date(2024, 1, 1, 23, 30, 0, 0, time.UTC)
		entry := &models.TimeEntry{
			ID:        "test-1",
			TaskName:  "Test",
			StartTime: startTime,
			EndTime:   timePtr(startTime.Add(90 * time.Minute)),
		}

		groups := initializeGroups(getHourlyKeys())
		distributeEntryAcrossHours(groups, entry)

		// Hour 23: 30 minutes
		if groups["23:00-00:00"].Duration != 30*time.Minute {
			t.Errorf("expected 30m in hour 23, got %v", groups["23:00-00:00"].Duration)
		}

		// Hour 0: 60 minutes
		if groups["00:00-01:00"].Duration != 60*time.Minute {
			t.Errorf("expected 60m in hour 0, got %v", groups["00:00-01:00"].Duration)
		}
	})
}

func TestGetWeekdayNames(t *testing.T) {
	names := getWeekdayNames()

	if len(names) != 7 {
		t.Errorf("expected 7 weekday names, got %d", len(names))
	}

	if names[0] != "Monday" {
		t.Errorf("expected first day to be Monday, got %s", names[0])
	}

	if names[6] != "Sunday" {
		t.Errorf("expected last day to be Sunday, got %s", names[6])
	}
}

func TestGetMonthNames(t *testing.T) {
	names := getMonthNames()

	if len(names) != 12 {
		t.Errorf("expected 12 month names, got %d", len(names))
	}

	if names[0] != "January" {
		t.Errorf("expected first month to be January, got %s", names[0])
	}

	if names[11] != "December" {
		t.Errorf("expected last month to be December, got %s", names[11])
	}
}

func TestGetHourlyKeys(t *testing.T) {
	keys := getHourlyKeys()

	if len(keys) != 24 {
		t.Errorf("expected 24 hourly keys, got %d", len(keys))
	}

	if keys[0] != "00:00-01:00" {
		t.Errorf("expected first key to be '00:00-01:00', got %s", keys[0])
	}

	if keys[23] != "23:00-00:00" {
		t.Errorf("expected last key to be '23:00-00:00', got %s", keys[23])
	}
}

func TestGetDayOfMonthKeys(t *testing.T) {
	keys := getDayOfMonthKeys()

	if len(keys) != 31 {
		t.Errorf("expected 31 day keys, got %d", len(keys))
	}

	if keys[0] != "Day 01" {
		t.Errorf("expected first key to be 'Day 01', got %s", keys[0])
	}

	if keys[30] != "Day 31" {
		t.Errorf("expected last key to be 'Day 31', got %s", keys[30])
	}
}

func TestGetWeekdayKey(t *testing.T) {
	tests := []struct {
		weekday  time.Weekday
		expected string
	}{
		{time.Monday, "Monday"},
		{time.Tuesday, "Tuesday"},
		{time.Wednesday, "Wednesday"},
		{time.Thursday, "Thursday"},
		{time.Friday, "Friday"},
		{time.Saturday, "Saturday"},
		{time.Sunday, "Sunday"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := getWeekdayKey(tt.weekday)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestGetMonthKey(t *testing.T) {
	tests := []struct {
		month    time.Month
		expected string
	}{
		{time.January, "January"},
		{time.February, "February"},
		{time.March, "March"},
		{time.April, "April"},
		{time.May, "May"},
		{time.June, "June"},
		{time.July, "July"},
		{time.August, "August"},
		{time.September, "September"},
		{time.October, "October"},
		{time.November, "November"},
		{time.December, "December"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := getMonthKey(tt.month)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func timePtr(t time.Time) *time.Time {
	return &t
}
