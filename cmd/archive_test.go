package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFormatFileSize(t *testing.T) {
	tests := []struct {
		name     string
		size     int64
		expected string
	}{
		{
			name:     "bytes",
			size:     100,
			expected: "100 B",
		},
		{
			name:     "kilobytes",
			size:     1024,
			expected: "1.0 KB",
		},
		{
			name:     "megabytes",
			size:     1024 * 1024,
			expected: "1.0 MB",
		},
		{
			name:     "gigabytes",
			size:     1024 * 1024 * 1024,
			expected: "1.0 GB",
		},
		{
			name:     "zero",
			size:     0,
			expected: "0 B",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatFileSize(tt.size)
			if result != tt.expected {
				t.Errorf("formatFileSize(%d) = %s, want %s", tt.size, result, tt.expected)
			}
		})
	}
}

func TestListArchives_NoDirectory(t *testing.T) {
	// Test with non-existent directory
	nonExistentDir := filepath.Join(t.TempDir(), "nonexistent")
	
	err := listArchives(nonExistentDir)
	if err != nil {
		t.Errorf("listArchives should not return error for non-existent directory, got: %v", err)
	}
}

func TestListArchives_EmptyDirectory(t *testing.T) {
	// Test with empty directory
	emptyDir := t.TempDir()
	
	err := listArchives(emptyDir)
	if err != nil {
		t.Errorf("listArchives should not return error for empty directory, got: %v", err)
	}
}

func TestListArchives_WithFiles(t *testing.T) {
	// Create temporary directory with some archive files
	tmpDir := t.TempDir()
	
	// Create test files
	testFiles := []string{
		"data-2025-01-01-120000.json",
		"data-2025-01-02-120000.json",
		"readme.txt", // Should be ignored
	}
	
	for _, file := range testFiles {
		f, err := os.Create(filepath.Join(tmpDir, file))
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
		f.WriteString("{}")
		f.Close()
	}
	
	err := listArchives(tmpDir)
	if err != nil {
		t.Errorf("listArchives failed: %v", err)
	}
}

func TestListArchives_InvalidDirectory(t *testing.T) {
	// Create a file instead of directory
	tmpFile := filepath.Join(t.TempDir(), "file")
	f, err := os.Create(tmpFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	f.Close()
	
	// Try to list archives from a file (not a directory)
	err = listArchives(tmpFile)
	if err == nil {
		t.Error("listArchives should return error when path is not a directory")
	}
}
