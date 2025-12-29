package models

import (
	"testing"
	"time"
)

func TestTimeEntry_Duration(t *testing.T) {
	t.Run("completed entry returns correct duration", func(t *testing.T) {
		startTime := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
		endTime := time.Date(2025, 1, 1, 12, 30, 0, 0, time.UTC)
		entry := &TimeEntry{
			ID:        "test-id",
			StartTime: startTime,
			EndTime:   &endTime,
		}

		duration := entry.Duration()
		expected := 2*time.Hour + 30*time.Minute

		if duration != expected {
			t.Errorf("Duration() = %v, want %v", duration, expected)
		}
	})

	t.Run("running entry returns duration since start", func(t *testing.T) {
		startTime := time.Now().Add(-1 * time.Hour)
		entry := &TimeEntry{
			ID:        "test-id",
			StartTime: startTime,
			EndTime:   nil,
		}

		duration := entry.Duration()

		if duration < 59*time.Minute || duration > 61*time.Minute {
			t.Errorf("Duration() = %v, expected approximately 1 hour", duration)
		}
	})
}

func TestTimeEntry_IsRunning(t *testing.T) {
	t.Run("entry with no end time is running", func(t *testing.T) {
		entry := &TimeEntry{
			ID:        "test-id",
			StartTime: time.Now(),
			EndTime:   nil,
		}

		if !entry.IsRunning() {
			t.Error("IsRunning() = false, want true")
		}
	})

	t.Run("entry with end time is not running", func(t *testing.T) {
		endTime := time.Now()
		entry := &TimeEntry{
			ID:        "test-id",
			StartTime: time.Now().Add(-1 * time.Hour),
			EndTime:   &endTime,
		}

		if entry.IsRunning() {
			t.Error("IsRunning() = true, want false")
		}
	})
}
