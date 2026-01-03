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

func TestJSONStorage_DeleteEntry(t *testing.T) {
	storage, cleanup := setupTestStorage(t)
	defer cleanup()

	// Create test entries
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

	t.Run("delete existing entry", func(t *testing.T) {
		err := storage.DeleteEntry("test-1")
		if err != nil {
			t.Fatalf("DeleteEntry() error = %v", err)
		}

		// Verify entry was deleted
		entries, _ := storage.ListEntries()
		if len(entries) != 1 {
			t.Errorf("expected 1 entry after deletion, got %d", len(entries))
		}

		if entries[0].ID != "test-2" {
			t.Errorf("expected remaining entry ID 'test-2', got '%s'", entries[0].ID)
		}
	})

	t.Run("delete non-existent entry", func(t *testing.T) {
		err := storage.DeleteEntry("non-existent")
		if err == nil {
			t.Error("expected error when deleting non-existent entry, got nil")
		}
	})
}

func TestJSONStorage_ArchiveData(t *testing.T) {
	storage, cleanup := setupTestStorage(t)
	defer cleanup()

	// Create test data
	entry := &models.TimeEntry{
		ID:        "test-1",
		StartTime: time.Now(),
		TaskName:  "Test Task",
		Tags:      []string{"tag1"},
	}
	storage.SaveEntry(entry)

	t.Run("archive existing data", func(t *testing.T) {
		tmpDir, _ := os.MkdirTemp("", "cli-record-archive-test-*")
		defer os.RemoveAll(tmpDir)

		archivePath := filepath.Join(tmpDir, "archive.json")
		err := storage.ArchiveData(archivePath)
		if err != nil {
			t.Fatalf("ArchiveData() error = %v", err)
		}

		// Verify archive file was created
		if _, err := os.Stat(archivePath); os.IsNotExist(err) {
			t.Error("archive file was not created")
		}

		// Verify archive contains the data
		archiveStorage := &JSONStorage{filePath: archivePath}
		entries, err := archiveStorage.ListEntries()
		if err != nil {
			t.Fatalf("failed to read archive: %v", err)
		}

		if len(entries) != 1 {
			t.Errorf("expected 1 entry in archive, got %d", len(entries))
		}

		if entries[0].ID != "test-1" {
			t.Errorf("expected entry ID 'test-1', got '%s'", entries[0].ID)
		}
	})

	t.Run("archive to non-existent directory", func(t *testing.T) {
		tmpDir, _ := os.MkdirTemp("", "cli-record-archive-test-*")
		defer os.RemoveAll(tmpDir)

		archivePath := filepath.Join(tmpDir, "subdir", "archive.json")
		err := storage.ArchiveData(archivePath)
		if err != nil {
			t.Fatalf("ArchiveData() should create directory, got error: %v", err)
		}

		// Verify archive was created
		if _, err := os.Stat(archivePath); os.IsNotExist(err) {
			t.Error("archive file was not created in new directory")
		}
	})
}

func TestJSONStorage_RestoreData(t *testing.T) {
	t.Run("restore with merge", func(t *testing.T) {
		storage, cleanup := setupTestStorage(t)
		defer cleanup()

		// Create current data
		currentEntry := &models.TimeEntry{
			ID:        "current-1",
			StartTime: time.Now(),
			TaskName:  "Current Task",
			Tags:      []string{"current"},
		}
		storage.SaveEntry(currentEntry)

		// Create archive with different data
		archiveStorage, archiveCleanup := setupTestStorage(t)
		defer archiveCleanup()

		archiveEntry := &models.TimeEntry{
			ID:        "archive-1",
			StartTime: time.Now(),
			TaskName:  "Archived Task",
			Tags:      []string{"archived"},
		}
		archiveStorage.SaveEntry(archiveEntry)

		// Restore with merge
		err := storage.RestoreData(archiveStorage.filePath, true)
		if err != nil {
			t.Fatalf("RestoreData() error = %v", err)
		}

		// Verify both entries exist
		entries, _ := storage.ListEntries()
		if len(entries) != 2 {
			t.Errorf("expected 2 entries after merge, got %d", len(entries))
		}

		// Verify both IDs are present
		idMap := make(map[string]bool)
		for _, e := range entries {
			idMap[e.ID] = true
		}

		if !idMap["current-1"] || !idMap["archive-1"] {
			t.Error("expected both current and archived entries after merge")
		}
	})

	t.Run("restore with merge - skip duplicates", func(t *testing.T) {
		storage, cleanup := setupTestStorage(t)
		defer cleanup()

		// Create entry with same ID in both
		entry := &models.TimeEntry{
			ID:        "same-id",
			StartTime: time.Now(),
			TaskName:  "Current Task",
			Tags:      []string{"current"},
		}
		storage.SaveEntry(entry)

		// Create archive with same ID
		archiveStorage, archiveCleanup := setupTestStorage(t)
		defer archiveCleanup()

		archiveEntry := &models.TimeEntry{
			ID:        "same-id",
			StartTime: time.Now(),
			TaskName:  "Archived Task",
			Tags:      []string{"archived"},
		}
		archiveStorage.SaveEntry(archiveEntry)

		// Restore with merge
		err := storage.RestoreData(archiveStorage.filePath, true)
		if err != nil {
			t.Fatalf("RestoreData() error = %v", err)
		}

		// Verify only one entry exists (duplicate skipped)
		entries, _ := storage.ListEntries()
		if len(entries) != 1 {
			t.Errorf("expected 1 entry after merge with duplicate, got %d", len(entries))
		}
	})

	t.Run("restore with replace", func(t *testing.T) {
		storage, cleanup := setupTestStorage(t)
		defer cleanup()

		// Create current data
		currentEntry := &models.TimeEntry{
			ID:        "current-1",
			StartTime: time.Now(),
			TaskName:  "Current Task",
			Tags:      []string{"current"},
		}
		storage.SaveEntry(currentEntry)

		// Create archive
		archiveStorage, archiveCleanup := setupTestStorage(t)
		defer archiveCleanup()

		archiveEntry := &models.TimeEntry{
			ID:        "archive-1",
			StartTime: time.Now(),
			TaskName:  "Archived Task",
			Tags:      []string{"archived"},
		}
		archiveStorage.SaveEntry(archiveEntry)

		// Restore with replace
		err := storage.RestoreData(archiveStorage.filePath, false)
		if err != nil {
			t.Fatalf("RestoreData() error = %v", err)
		}

		// Verify only archived entry exists
		entries, _ := storage.ListEntries()
		if len(entries) != 1 {
			t.Errorf("expected 1 entry after replace, got %d", len(entries))
		}

		if entries[0].ID != "archive-1" {
			t.Errorf("expected archived entry, got ID '%s'", entries[0].ID)
		}

		// Verify backup was created
		backupPath := storage.filePath + ".backup"
		if _, err := os.Stat(backupPath); os.IsNotExist(err) {
			t.Error("backup file was not created")
		}
	})

	t.Run("restore from invalid file", func(t *testing.T) {
		storage, cleanup := setupTestStorage(t)
		defer cleanup()

		err := storage.RestoreData("/non/existent/file.json", true)
		if err == nil {
			t.Error("expected error when restoring from non-existent file, got nil")
		}
	})
}

func TestNewJSONStorage(t *testing.T) {
	// This test just ensures NewJSONStorage doesn't panic
	storage, err := NewJSONStorage()
	if err != nil {
		t.Fatalf("NewJSONStorage() error = %v", err)
	}

	if storage == nil {
		t.Error("expected non-nil storage")
	}

	// Verify default path is set
	if storage.filePath == "" {
		t.Error("expected non-empty file path")
	}
}
