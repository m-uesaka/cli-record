package cmd

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/m-uesaka/cli-record/internal/models"
	"github.com/m-uesaka/cli-record/internal/storage"
)

func TestRunShow_EntryNotFound(t *testing.T) {
	// Setup test environment
	tmpDir := t.TempDir()
	dataPath := filepath.Join(tmpDir, "data.json")
	os.Setenv("CLI_RECORD_DATA", dataPath)
	defer os.Unsetenv("CLI_RECORD_DATA")

	// Create empty storage
	_, err := storage.NewJSONStorage()
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	// Try to show non-existent entry
	cmd := showCmd
	err = runShow(cmd, []string{"nonexistent-id"})
	
	if err == nil {
		t.Error("Expected error for non-existent entry, got nil")
	}
}

func TestRunShow_Success(t *testing.T) {
	// Setup test environment
	tmpDir := t.TempDir()
	dataPath := filepath.Join(tmpDir, "data.json")
	os.Setenv("CLI_RECORD_DATA", dataPath)
	defer os.Unsetenv("CLI_RECORD_DATA")

	// Create storage and add entry
	store, err := storage.NewJSONStorage()
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	now := time.Now()
	entry := &models.TimeEntry{
		TaskName:  "Test Task",
		StartTime: now,
		EndTime:   &now,
		Tags:      []string{"test", "demo"},
	}
	if err := store.SaveEntry(entry); err != nil {
		t.Fatalf("Failed to add entry: %v", err)
	}

	// Show the entry
	cmd := showCmd
	err = runShow(cmd, []string{entry.ID})
	
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestRunShow_RunningEntry(t *testing.T) {
	// Setup test environment
	tmpDir := t.TempDir()
	dataPath := filepath.Join(tmpDir, "data.json")
	os.Setenv("CLI_RECORD_DATA", dataPath)
	defer os.Unsetenv("CLI_RECORD_DATA")

	// Create storage and add running entry
	store, err := storage.NewJSONStorage()
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	entry := &models.TimeEntry{
		TaskName:  "Running Task",
		StartTime: time.Now(),
		EndTime:   nil, // Still running
		Tags:      []string{},
	}
	if err := store.SaveEntry(entry); err != nil {
		t.Fatalf("Failed to add entry: %v", err)
	}

	// Show the running entry
	cmd := showCmd
	err = runShow(cmd, []string{entry.ID})
	
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}
