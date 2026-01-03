package cmd

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/m-uesaka/cli-record/internal/models"
	"github.com/m-uesaka/cli-record/internal/storage"
)

func TestRunRemove_EntryNotFound(t *testing.T) {
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

	// Try to remove non-existent entry
	cmd := removeCmd
	removeForce = true // Skip confirmation
	err = runRemove(cmd, []string{"nonexistent-id"})
	
	if err == nil {
		t.Error("Expected error for non-existent entry, got nil")
	}
}

func TestRunRemove_WithForce(t *testing.T) {
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

	entry := &models.TimeEntry{
		TaskName:  "Test Task",
		StartTime: time.Now(),
		Tags:      []string{"test"},
	}
	if err := store.SaveEntry(entry); err != nil {
		t.Fatalf("Failed to add entry: %v", err)
	}

	// Remove with force flag
	cmd := removeCmd
	removeForce = true
	err = runRemove(cmd, []string{entry.ID})
	
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Verify entry is removed
	_, err = store.GetEntry(entry.ID)
	if err == nil {
		t.Error("Entry should have been removed")
	}
}

func TestFormatTagsOrNone(t *testing.T) {
	tests := []struct {
		name     string
		tags     []string
		expected string
	}{
		{
			name:     "no tags",
			tags:     []string{},
			expected: "(none)",
		},
		{
			name:     "nil tags",
			tags:     nil,
			expected: "(none)",
		},
		{
			name:     "single tag",
			tags:     []string{"work"},
			expected: "work",
		},
		{
			name:     "multiple tags",
			tags:     []string{"work", "meeting", "client"},
			expected: "work, meeting, client",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatTagsOrNone(tt.tags)
			if result != tt.expected {
				t.Errorf("formatTagsOrNone(%v) = %s, want %s", tt.tags, result, tt.expected)
			}
		})
	}
}
