package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/m-uesaka/cli-record/internal/models"
)

func TestGenerateGroupedReport(t *testing.T) {
	now := time.Now()
	entries := []*models.TimeEntry{
		{
			ID:        "1",
			TaskName:  "TaskA",
			StartTime: now,
			EndTime:   timePtr(now.Add(1 * time.Hour)),
			Tags:      []string{"tag1", "tag2"},
		},
		{
			ID:        "2",
			TaskName:  "TaskB",
			StartTime: now,
			EndTime:   timePtr(now.Add(30 * time.Minute)),
			Tags:      []string{"tag1"},
		},
		{
			ID:        "3",
			TaskName:  "TaskA",
			StartTime: now,
			EndTime:   timePtr(now.Add(2 * time.Hour)),
			Tags:      []string{},
		},
	}

	t.Run("group by task", func(t *testing.T) {
		report := generateGroupedReport(entries, "task")

		if report.Title != "Grouped by task" {
			t.Errorf("expected title 'Grouped by task', got %s", report.Title)
		}

		if len(report.Groups) != 2 {
			t.Fatalf("expected 2 groups, got %d", len(report.Groups))
		}

		// Find TaskA group
		var taskAGroup *ReportGroup
		for i := range report.Groups {
			if report.Groups[i].Name == "TaskA" {
				taskAGroup = &report.Groups[i]
				break
			}
		}

		if taskAGroup == nil {
			t.Fatal("TaskA group not found")
		}

		if taskAGroup.Duration != 3*time.Hour {
			t.Errorf("expected TaskA duration 3h, got %v", taskAGroup.Duration)
		}

		if taskAGroup.Count != 2 {
			t.Errorf("expected TaskA count 2, got %d", taskAGroup.Count)
		}
	})

	t.Run("group by tag", func(t *testing.T) {
		report := generateGroupedReport(entries, "tag")

		if report.Title != "Grouped by tag" {
			t.Errorf("expected title 'Grouped by tag', got %s", report.Title)
		}

		// Should have: tag1, tag2, (no tags)
		if len(report.Groups) != 3 {
			t.Fatalf("expected 3 groups, got %d", len(report.Groups))
		}

		// Find tag1 group
		var tag1Group *ReportGroup
		for i := range report.Groups {
			if report.Groups[i].Name == "tag1" {
				tag1Group = &report.Groups[i]
				break
			}
		}

		if tag1Group == nil {
			t.Fatal("tag1 group not found")
		}

		// tag1 appears in entries 1 and 2
		if tag1Group.Duration != 90*time.Minute {
			t.Errorf("expected tag1 duration 90m, got %v", tag1Group.Duration)
		}

		if tag1Group.Count != 2 {
			t.Errorf("expected tag1 count 2, got %d", tag1Group.Count)
		}
	})

	t.Run("group by day", func(t *testing.T) {
		report := generateGroupedReport(entries, "day")

		if len(report.Groups) != 1 {
			t.Fatalf("expected 1 group (same day), got %d", len(report.Groups))
		}

		expectedDay := now.Format("2006-01-02")
		if report.Groups[0].Name != expectedDay {
			t.Errorf("expected group name %s, got %s", expectedDay, report.Groups[0].Name)
		}
	})

	t.Run("group by week", func(t *testing.T) {
		report := generateGroupedReport(entries, "week")

		if len(report.Groups) != 1 {
			t.Fatalf("expected 1 group (same week), got %d", len(report.Groups))
		}

		year, week := now.ISOWeek()
		expectedWeek := strings.Contains(report.Groups[0].Name, fmt.Sprintf("%d-W%02d", year, week))
		if !expectedWeek {
			t.Errorf("unexpected week format: %s", report.Groups[0].Name)
		}
	})

	t.Run("group by month", func(t *testing.T) {
		report := generateGroupedReport(entries, "month")

		if len(report.Groups) != 1 {
			t.Fatalf("expected 1 group (same month), got %d", len(report.Groups))
		}

		expectedMonth := now.Format("2006-01")
		if report.Groups[0].Name != expectedMonth {
			t.Errorf("expected group name %s, got %s", expectedMonth, report.Groups[0].Name)
		}
	})

	t.Run("group by year", func(t *testing.T) {
		report := generateGroupedReport(entries, "year")

		if len(report.Groups) != 1 {
			t.Fatalf("expected 1 group (same year), got %d", len(report.Groups))
		}

		expectedYear := now.Format("2006")
		if report.Groups[0].Name != expectedYear {
			t.Errorf("expected group name %s, got %s", expectedYear, report.Groups[0].Name)
		}
	})
}

func TestGenerateHourlyReport(t *testing.T) {
	startTime := time.Date(2024, 1, 1, 14, 30, 0, 0, time.UTC)
	entries := []*models.TimeEntry{
		{
			ID:        "1",
			TaskName:  "Task1",
			StartTime: startTime,
			EndTime:   timePtr(startTime.Add(90 * time.Minute)),
		},
	}

	report := generateHourlyReport(entries)

	if report.Title != "Hourly Breakdown" {
		t.Errorf("expected title 'Hourly Breakdown', got %s", report.Title)
	}

	if len(report.Groups) != 24 {
		t.Errorf("expected 24 hourly groups, got %d", len(report.Groups))
	}

	// Check that entry is distributed correctly
	for _, group := range report.Groups {
		switch group.Name {
		case "14:00-15:00":
			if group.Duration != 30*time.Minute {
				t.Errorf("expected 30m in hour 14, got %v", group.Duration)
			}
		case "15:00-16:00":
			if group.Duration != 60*time.Minute {
				t.Errorf("expected 60m in hour 15, got %v", group.Duration)
			}
		}
	}
}

func TestGenerateWeekdayReport(t *testing.T) {
	// Monday
	monday := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	// Wednesday
	wednesday := time.Date(2024, 1, 3, 10, 0, 0, 0, time.UTC)

	entries := []*models.TimeEntry{
		{
			ID:        "1",
			TaskName:  "Task1",
			StartTime: monday,
			EndTime:   timePtr(monday.Add(1 * time.Hour)),
		},
		{
			ID:        "2",
			TaskName:  "Task2",
			StartTime: wednesday,
			EndTime:   timePtr(wednesday.Add(2 * time.Hour)),
		},
	}

	report := generateWeekdayReport(entries)

	if report.Title != "Weekday Breakdown" {
		t.Errorf("expected title 'Weekday Breakdown', got %s", report.Title)
	}

	if len(report.Groups) != 7 {
		t.Errorf("expected 7 weekday groups, got %d", len(report.Groups))
	}

	for _, group := range report.Groups {
		switch group.Name {
		case "Monday":
			if group.Duration != 1*time.Hour {
				t.Errorf("expected 1h on Monday, got %v", group.Duration)
			}
			if group.Count != 1 {
				t.Errorf("expected 1 entry on Monday, got %d", group.Count)
			}
		case "Wednesday":
			if group.Duration != 2*time.Hour {
				t.Errorf("expected 2h on Wednesday, got %v", group.Duration)
			}
			if group.Count != 1 {
				t.Errorf("expected 1 entry on Wednesday, got %d", group.Count)
			}
		}
	}
}

func TestWeekdayReportCorrectOrder(t *testing.T) {
	// Create entries for each day of the week
	// Starting from Sunday, Jan 7, 2024
	sunday := time.Date(2024, 1, 7, 10, 0, 0, 0, time.UTC)
	entries := []*models.TimeEntry{}
	
	// Create one entry for each day of the week
	for i := 0; i < 7; i++ {
		day := sunday.AddDate(0, 0, i)
		entries = append(entries, &models.TimeEntry{
			ID:        fmt.Sprintf("%d", i),
			TaskName:  "Task",
			StartTime: day,
			EndTime:   timePtr(day.Add(1 * time.Hour)),
		})
	}

	report := generateWeekdayReport(entries)

	// Verify the order is Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday
	expectedOrder := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	
	if len(report.Groups) != 7 {
		t.Fatalf("expected 7 groups, got %d", len(report.Groups))
	}

	for i, expectedName := range expectedOrder {
		if report.Groups[i].Name != expectedName {
			t.Errorf("expected group %d to be %s, got %s", i, expectedName, report.Groups[i].Name)
		}
	}
}

func TestGenerateDayOfMonthReport(t *testing.T) {
	day1 := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	day15 := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	entries := []*models.TimeEntry{
		{
			ID:        "1",
			TaskName:  "Task1",
			StartTime: day1,
			EndTime:   timePtr(day1.Add(1 * time.Hour)),
		},
		{
			ID:        "2",
			TaskName:  "Task2",
			StartTime: day15,
			EndTime:   timePtr(day15.Add(3 * time.Hour)),
		},
	}

	report := generateDayOfMonthReport(entries)

	if report.Title != "Day-of-Month Breakdown" {
		t.Errorf("expected title 'Day-of-Month Breakdown', got %s", report.Title)
	}

	if len(report.Groups) != 31 {
		t.Errorf("expected 31 day groups, got %d", len(report.Groups))
	}

	for _, group := range report.Groups {
		switch group.Name {
		case "Day 01":
			if group.Duration != 1*time.Hour {
				t.Errorf("expected 1h on day 1, got %v", group.Duration)
			}
		case "Day 15":
			if group.Duration != 3*time.Hour {
				t.Errorf("expected 3h on day 15, got %v", group.Duration)
			}
		}
	}
}

func TestGenerateMonthlyReport(t *testing.T) {
	jan := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	mar := time.Date(2024, 3, 15, 10, 0, 0, 0, time.UTC)

	entries := []*models.TimeEntry{
		{
			ID:        "1",
			TaskName:  "Task1",
			StartTime: jan,
			EndTime:   timePtr(jan.Add(2 * time.Hour)),
		},
		{
			ID:        "2",
			TaskName:  "Task2",
			StartTime: mar,
			EndTime:   timePtr(mar.Add(5 * time.Hour)),
		},
	}

	report := generateMonthlyReport(entries)

	if report.Title != "Monthly Breakdown" {
		t.Errorf("expected title 'Monthly Breakdown', got %s", report.Title)
	}

	if len(report.Groups) != 12 {
		t.Errorf("expected 12 month groups, got %d", len(report.Groups))
	}

	for _, group := range report.Groups {
		switch group.Name {
		case "January":
			if group.Duration != 2*time.Hour {
				t.Errorf("expected 2h in January, got %v", group.Duration)
			}
		case "March":
			if group.Duration != 5*time.Hour {
				t.Errorf("expected 5h in March, got %v", group.Duration)
			}
		}
	}
}

func TestCreateReport(t *testing.T) {
	groups := map[string]*ReportGroup{
		"Group1": {Name: "Group1", Duration: 1 * time.Hour, Count: 2},
		"Group2": {Name: "Group2", Duration: 30 * time.Minute, Count: 1},
	}

	report := createReport("Test Report", groups)

	if report.Title != "Test Report" {
		t.Errorf("expected title 'Test Report', got %s", report.Title)
	}

	if len(report.Groups) != 2 {
		t.Errorf("expected 2 groups, got %d", len(report.Groups))
	}

	if report.Total != 90*time.Minute {
		t.Errorf("expected total 90m, got %v", report.Total)
	}

	// Verify groups are sorted by name
	if report.Groups[0].Name > report.Groups[1].Name {
		t.Error("groups are not sorted by name")
	}
}

func TestOutputTableReport(t *testing.T) {
	report := ViewReport{
		Title: "Test Report",
		Groups: []ReportGroup{
			{Name: "Group1", Duration: 1 * time.Hour, Count: 2},
			{Name: "Group2", Duration: 30 * time.Minute, Count: 1},
		},
		Total: 90 * time.Minute,
	}

	var buf bytes.Buffer
	file := os.NewFile(uintptr(0), "/dev/null") // Dummy file
	defer file.Close()

	// Use a buffer instead
	err := outputTableReport(os.Stdout, report)
	if err != nil {
		t.Fatalf("outputTableReport failed: %v", err)
	}

	// Just ensure it doesn't panic/error
	_ = buf.String()
}

func TestOutputCSVReport(t *testing.T) {
	report := ViewReport{
		Title: "Test Report",
		Groups: []ReportGroup{
			{Name: "Group1", Duration: 1 * time.Hour, Count: 2},
			{Name: "Group2", Duration: 30 * time.Minute, Count: 1},
		},
		Total: 90 * time.Minute,
	}

	tmpFile, err := os.CreateTemp("", "test-csv-*.csv")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	err = outputCSVReport(tmpFile, report)
	if err != nil {
		t.Fatalf("outputCSVReport failed: %v", err)
	}

	tmpFile.Seek(0, 0)
	content, _ := os.ReadFile(tmpFile.Name())
	lines := strings.Split(string(content), "\n")

	if len(lines) < 3 {
		t.Errorf("expected at least 3 lines (header + 2 groups), got %d", len(lines))
	}

	// Check header
	if !strings.Contains(lines[0], "Group") {
		t.Error("CSV header missing 'Group' column")
	}
}

func TestOutputJSONReport(t *testing.T) {
	report := ViewReport{
		Title: "Test Report",
		Groups: []ReportGroup{
			{Name: "Group1", Duration: 1 * time.Hour, Count: 2},
			{Name: "Group2", Duration: 30 * time.Minute, Count: 1},
		},
		Total: 90 * time.Minute,
	}

	tmpFile, err := os.CreateTemp("", "test-json-*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	err = outputJSONReport(tmpFile, report)
	if err != nil {
		t.Fatalf("outputJSONReport failed: %v", err)
	}

	tmpFile.Seek(0, 0)
	content, _ := os.ReadFile(tmpFile.Name())

	// Verify it's valid JSON
	var jsonData map[string]interface{}
	if err := json.Unmarshal(content, &jsonData); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	// Verify structure
	if title, ok := jsonData["title"].(string); !ok || title != "Test Report" {
		t.Error("JSON missing or incorrect title")
	}

	if groups, ok := jsonData["groups"].([]interface{}); !ok || len(groups) != 2 {
		t.Error("JSON missing or incorrect groups")
	}

	if total, ok := jsonData["total"].(map[string]interface{}); !ok {
		t.Error("JSON missing total")
	} else {
		if _, ok := total["duration_seconds"]; !ok {
			t.Error("JSON total missing duration_seconds")
		}
	}
}
