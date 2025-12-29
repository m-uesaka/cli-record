package integration

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/m-uesaka/cli-record/internal/models"
	"github.com/m-uesaka/cli-record/internal/storage"
)

// TestFullWorkflow tests a complete workflow: start -> stop -> list
func TestFullWorkflow(t *testing.T) {
	// Use default storage for now
	store, err := storage.NewJSONStorage()
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}

	// Get current count of entries
	initialEntries, err := store.ListEntries()
	if err != nil {
		t.Fatalf("failed to list entries: %v", err)
	}
	initialCount := len(initialEntries)

	// Start a time entry
	entry1 := &models.TimeEntry{
		ID:        uuid.New().String(),
		StartTime: time.Now(),
		TaskName:  "Integration Test Task",
		Tags:      []string{"testing", "integration"},
	}

	if err := store.SaveEntry(entry1); err != nil {
		t.Fatalf("failed to save entry: %v", err)
	}

	// Verify entry was saved
	entries, err := store.ListEntries()
	if err != nil {
		t.Fatalf("failed to list entries: %v", err)
	}
	if len(entries) != initialCount+1 {
		t.Errorf("expected %d entries, got %d", initialCount+1, len(entries))
	}

	// Verify entry is running
	running, err := store.GetRunningEntry()
	if err != nil {
		t.Fatalf("failed to get running entry: %v", err)
	}
	if running == nil {
		t.Fatal("expected running entry, got nil")
	}

	// Stop the entry
	time.Sleep(100 * time.Millisecond)
	endTime := time.Now()
	running.EndTime = &endTime

	if err := store.UpdateEntry(running); err != nil {
		t.Fatalf("failed to update entry: %v", err)
	}

	// Verify it's stopped
	updatedEntry, err := store.GetEntry(running.ID)
	if err != nil {
		t.Fatalf("failed to get entry: %v", err)
	}
	if updatedEntry.IsRunning() {
		t.Error("expected entry to be stopped")
	}

	t.Logf("Successfully completed workflow test with %d entries", len(entries))
}

// TestEdgeCases tests various edge cases
func TestEdgeCases(t *testing.T) {
	store, err := storage.NewJSONStorage()
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}

	t.Run("entry without task name", func(t *testing.T) {
		entry := &models.TimeEntry{
			ID:        uuid.New().String(),
			StartTime: time.Now(),
			TaskName:  "",
			Tags:      []string{"test"},
		}

		if err := store.SaveEntry(entry); err != nil {
			t.Fatalf("failed to save entry with empty task name: %v", err)
		}

		retrieved, err := store.GetEntry(entry.ID)
		if err != nil {
			t.Fatalf("failed to retrieve entry: %v", err)
		}

		if retrieved.TaskName != "" {
			t.Errorf("expected empty task name, got %s", retrieved.TaskName)
		}
	})

	t.Run("entry without tags", func(t *testing.T) {
		entry := &models.TimeEntry{
			ID:        uuid.New().String(),
			StartTime: time.Now(),
			TaskName:  "No Tags",
			Tags:      []string{},
		}

		if err := store.SaveEntry(entry); err != nil {
			t.Fatalf("failed to save entry without tags: %v", err)
		}

		retrieved, err := store.GetEntry(entry.ID)
		if err != nil {
			t.Fatalf("failed to retrieve entry: %v", err)
		}

		if len(retrieved.Tags) != 0 {
			t.Errorf("expected no tags, got %v", retrieved.Tags)
		}
	})

	t.Run("duration calculation for running entry", func(t *testing.T) {
		entry := &models.TimeEntry{
			ID:        uuid.New().String(),
			StartTime: time.Now().Add(-5 * time.Minute),
			TaskName:  "Running Entry",
			Tags:      []string{},
		}

		duration := entry.Duration()
		if duration < 4*time.Minute || duration > 6*time.Minute {
			t.Errorf("expected duration around 5 minutes, got %v", duration)
		}
	})
}
