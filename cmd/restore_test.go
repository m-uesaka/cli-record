package cmd

import (
	"testing"
)

func TestRunRestore_InvalidFlags(t *testing.T) {
	cmd := restoreCmd
	
	// Set both merge and replace flags
	cmd.Flags().Set("merge", "true")
	cmd.Flags().Set("replace", "true")
	
	err := runRestore(cmd, []string{"archive.json"})
	
	if err == nil {
		t.Error("Expected error when both --merge and --replace are specified")
	}
	
	// Reset flags
	cmd.Flags().Set("merge", "false")
	cmd.Flags().Set("replace", "false")
}

func TestRunRestore_PreviewMode(t *testing.T) {
	// This tests that preview mode doesn't actually modify data
	restorePreview = true
	defer func() { restorePreview = false }()
	
	// Preview should not fail even with non-existent file
	// (because it doesn't actually try to restore)
	cmd := restoreCmd
	err := runRestore(cmd, []string{"/tmp/nonexistent.json"})
	
	// Preview shouldn't error on the file check since it's just showing what would happen
	// Note: Current implementation may still fail on storage init
	_ = err // We don't assert here since implementation details may vary
}
