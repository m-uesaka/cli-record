package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/m-uesaka/cli-record/internal/models"
)

func setupTestStorage(t *testing.T) (*JSONStorage, func()) {
	t.Helper()

	tmpDir, err := os.MkdirTemp("", "cli-record-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	storage := &JSONStorage{
		filePath: filepath.Join(tmpDir, "data.json"),
	}

	cleanup := func() {
		os.RemoveAll(tmpDir)
	}

	return storage, cleanup
}

func TestJSONStorage_SaveEntry(t *testing.T) {
	storage, cleanup := setupTestStorage(t)
	defer cleanup()

	entry := &models.TimeEntry{
		ID:        "test-1",
		StartTime: time.Now(),
		TaskName:  "Test Task",
		Tags:      []string{"tag1", "tag2"},
	}

	err := storage.SaveEntry(entry)
	if err != nil {
		t.Fatalf("SaveEntry() error = %v", err)
	}

	entries, err := storage.ListEntries()
	if err != nil {
		t.Fatalf("ListEntries() error = %v", err)
	}

	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	if entries[0].ID != entry.ID {
		t.Errorf("expected ID %s, got %s", entry.ID, entries[0].ID)
	}
}

func TestJSONStorage_GetEntry(t *testing.T) {
	storage, cleanup := setupTestStorage(t)
	defer cleanup()

	entry := &models.TimeEntry{
		ID:        "test-1",
		StartTime: time.Now(),
		TaskName:  "Test Task",
		Tags:      []string{"tag1"},
	}

	storage.SaveEntry(entry)

	t.Run("existing entry", func(t *testing.T) {
		retrieved, err := storage.GetEntry("test-1")
		if err != nil {
			t.Fatalf("GetEntry() error = %v", err)
		}

		if retrieved.ID != entry.ID {
			t.Errorf("expected ID %s, got %s", entry.ID, retrieved.ID)
		}
	})

	t.Run("non-existent entry", func(t *testing.T) {
		_, err := storage.GetEntry("non-existent")
		if err == nil {
			t.Error("expected error for non-existent entry, got nil")
		}
	})
}

func TestJSONStorage_ListEntries(t *testing.T) {
	storage, cleanup := setupTestStorage(t)
	defer cleanup()

	t.Run("empty storage", func(t *testing.T) {
		entries, err := storage.ListEntries()
		if err != nil {
			t.Fatalf("ListEntries() error = %v", err)
		}

		if len(entries) != 0 {
			t.Errorf("expected 0 entries, got %d", len(entries))
		}
	})

	t.Run("multiple entries", func(t *testing.T) {
		entry1 := &models.TimeEntry{
			ID:        "test-1",
			StartTime: time.Now(),
			TaskName:  "Task 1",
			Tags:      []string{"tag1"},
		}
		entry2 := &models.TimeEntry{
			ID:        "test-2",
			StartTime: time.Now(),
			TaskName:  "Task 2",
			Tags:      []string{"tag2"},
		}

		storage.SaveEntry(entry1)
		storage.SaveEntry(entry2)

		entries, err := storage.ListEntries()
		if err != nil {
			t.Fatalf("ListEntries() error = %v", err)
		}

		if len(entries) != 2 {
			t.Errorf("expected 2 entries, got %d", len(entries))
		}
	})
}

func TestJSONStorage_UpdateEntry(t *testing.T) {
	storage, cleanup := setupTestStorage(t)
	defer cleanup()

	entry := &models.TimeEntry{
		ID:        "test-1",
		StartTime: time.Now(),
		TaskName:  "Original Task",
		Tags:      []string{"tag1"},
	}

	storage.SaveEntry(entry)

	t.Run("update existing entry", func(t *testing.T) {
		endTime := time.Now()
		entry.EndTime = &endTime
		entry.TaskName = "Updated Task"

		err := storage.UpdateEntry(entry)
		if err != nil {
			t.Fatalf("UpdateEntry() error = %v", err)
		}

		retrieved, err := storage.GetEntry("test-1")
		if err != nil {
			t.Fatalf("GetEntry() error = %v", err)
		}

		if retrieved.TaskName != "Updated Task" {
			t.Errorf("expected TaskName 'Updated Task', got '%s'", retrieved.TaskName)
		}

		if retrieved.EndTime == nil {
			t.Error("expected EndTime to be set, got nil")
		}
	})

	t.Run("update non-existent entry", func(t *testing.T) {
		nonExistent := &models.TimeEntry{
			ID:        "non-existent",
			StartTime: time.Now(),
			TaskName:  "Test",
		}

		err := storage.UpdateEntry(nonExistent)
		if err == nil {
			t.Error("expected error for non-existent entry, got nil")
		}
	})
}

func TestJSONStorage_GetRunningEntry(t *testing.T) {
	storage, cleanup := setupTestStorage(t)
	defer cleanup()

	t.Run("no running entry", func(t *testing.T) {
		running, err := storage.GetRunningEntry()
		if err != nil {
			t.Fatalf("GetRunningEntry() error = %v", err)
		}

		if running != nil {
			t.Error("expected nil for no running entry, got entry")
		}
	})

	t.Run("with running entry", func(t *testing.T) {
		runningEntry := &models.TimeEntry{
			ID:        "running-1",
			StartTime: time.Now(),
			TaskName:  "Running Task",
			Tags:      []string{},
		}

		storage.SaveEntry(runningEntry)

		retrieved, err := storage.GetRunningEntry()
		if err != nil {
			t.Fatalf("GetRunningEntry() error = %v", err)
		}

		if retrieved == nil {
			t.Fatal("expected running entry, got nil")
		}

		if retrieved.ID != "running-1" {
			t.Errorf("expected ID 'running-1', got '%s'", retrieved.ID)
		}
	})

	t.Run("only completed entries", func(t *testing.T) {
		storage2, cleanup2 := setupTestStorage(t)
		defer cleanup2()

		endTime := time.Now()
		completedEntry := &models.TimeEntry{
			ID:        "completed-1",
			StartTime: time.Now().Add(-1 * time.Hour),
			EndTime:   &endTime,
			TaskName:  "Completed Task",
		}

		storage2.SaveEntry(completedEntry)

		running, err := storage2.GetRunningEntry()
		if err != nil {
			t.Fatalf("GetRunningEntry() error = %v", err)
		}

		if running != nil {
			t.Error("expected nil when only completed entries exist, got entry")
		}
	})
}

func TestJSONStorage_ListTags(t *testing.T) {
	storage, cleanup := setupTestStorage(t)
	defer cleanup()

	t.Run("no tags", func(t *testing.T) {
		tags, err := storage.ListTags()
		if err != nil {
			t.Fatalf("ListTags() error = %v", err)
		}

		if len(tags) != 0 {
			t.Errorf("expected 0 tags, got %d", len(tags))
		}
	})

	t.Run("multiple unique tags", func(t *testing.T) {
		entry1 := &models.TimeEntry{
			ID:        "test-1",
			StartTime: time.Now(),
			TaskName:  "Task 1",
			Tags:      []string{"tag1", "tag2"},
		}
		entry2 := &models.TimeEntry{
			ID:        "test-2",
			StartTime: time.Now(),
			TaskName:  "Task 2",
			Tags:      []string{"tag2", "tag3"},
		}

		storage.SaveEntry(entry1)
		storage.SaveEntry(entry2)

		tags, err := storage.ListTags()
		if err != nil {
			t.Fatalf("ListTags() error = %v", err)
		}

		if len(tags) != 3 {
			t.Errorf("expected 3 unique tags, got %d", len(tags))
		}

		tagMap := make(map[string]bool)
		for _, tag := range tags {
			tagMap[tag] = true
		}

		expectedTags := []string{"tag1", "tag2", "tag3"}
		for _, expected := range expectedTags {
			if !tagMap[expected] {
				t.Errorf("expected tag '%s' not found", expected)
			}
		}
	})
}

func TestJSONStorage_ConcurrentAccess(t *testing.T) {
	storage, cleanup := setupTestStorage(t)
	defer cleanup()

	done := make(chan bool)

	for i := 0; i < 10; i++ {
		go func(id int) {
			entry := &models.TimeEntry{
				ID:        fmt.Sprintf("concurrent-%d", id),
				StartTime: time.Now(),
				TaskName:  fmt.Sprintf("Task %d", id),
				Tags:      []string{fmt.Sprintf("tag%d", id)},
			}
			storage.SaveEntry(entry)
			done <- true
		}(i)
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	entries, err := storage.ListEntries()
	if err != nil {
		t.Fatalf("ListEntries() error = %v", err)
	}

	if len(entries) != 10 {
		t.Errorf("expected 10 entries after concurrent writes, got %d", len(entries))
	}
}
