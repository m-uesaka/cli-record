package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/m-uesaka/cli-record/internal/config"
)

func TestRunConfigSet_InvalidTimeFormat(t *testing.T) {
	// Setup test environment
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	os.Setenv("CLI_RECORD_CONFIG", configPath)
	defer os.Unsetenv("CLI_RECORD_CONFIG")

	// Initialize with default config
	cfg, err := config.GetDefaultConfig()
	if err != nil {
		t.Fatalf("Failed to get default config: %v", err)
	}
	if err := cfg.Save(); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	cmd := configSetCmd
	cmd.SetArgs([]string{"time-format", "invalid"})

	err = cmd.Execute()
	if err == nil {
		t.Error("Expected error for invalid time format, got nil")
	}
}

func TestRunConfigSet_InvalidKey(t *testing.T) {
	// Setup test environment
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	os.Setenv("CLI_RECORD_CONFIG", configPath)
	defer os.Unsetenv("CLI_RECORD_CONFIG")

	// Initialize with default config
	cfg, err := config.GetDefaultConfig()
	if err != nil {
		t.Fatalf("Failed to get default config: %v", err)
	}
	if err := cfg.Save(); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	cmd := configSetCmd
	cmd.SetArgs([]string{"invalid-key", "value"})

	err = cmd.Execute()
	if err == nil {
		t.Error("Expected error for invalid config key, got nil")
	}
}

func TestRunConfigSet_ValidTimeFormat(t *testing.T) {
	// Setup test environment
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	os.Setenv("CLI_RECORD_CONFIG", configPath)
	defer os.Unsetenv("CLI_RECORD_CONFIG")

	// Initialize with default config
	cfg, err := config.GetDefaultConfig()
	if err != nil {
		t.Fatalf("Failed to get default config: %v", err)
	}
	if err := cfg.Save(); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	tests := []string{"12h", "24h"}
	for _, format := range tests {
		t.Run(format, func(t *testing.T) {
			cmd := configSetCmd
			cmd.SetArgs([]string{"time-format", format})

			err := runConfigSet(cmd, []string{"time-format", format})
			if err != nil {
				t.Errorf("Expected no error for valid time format %s, got: %v", format, err)
			}

			// Verify the config was updated
			loadedCfg, err := config.Load()
			if err != nil {
				t.Fatalf("Failed to load config: %v", err)
			}
			if loadedCfg.TimeFormat != format {
				t.Errorf("TimeFormat = %s, want %s", loadedCfg.TimeFormat, format)
			}
		})
	}
}

func TestRunConfigSet_DataPath(t *testing.T) {
	// Setup test environment
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	os.Setenv("CLI_RECORD_CONFIG", configPath)
	defer os.Unsetenv("CLI_RECORD_CONFIG")

	// Initialize with default config
	cfg, err := config.GetDefaultConfig()
	if err != nil {
		t.Fatalf("Failed to get default config: %v", err)
	}
	if err := cfg.Save(); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	newPath := "/tmp/test/data.json"
	cmd := configSetCmd
	cmd.SetArgs([]string{"data-path", newPath})

	err = runConfigSet(cmd, []string{"data-path", newPath})
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// Verify the config was updated
	loadedCfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	if loadedCfg.DataFilePath != newPath {
		t.Errorf("DataFilePath = %s, want %s", loadedCfg.DataFilePath, newPath)
	}
}

func TestRunConfigSet_ArchiveDir(t *testing.T) {
	// Setup test environment
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	os.Setenv("CLI_RECORD_CONFIG", configPath)
	defer os.Unsetenv("CLI_RECORD_CONFIG")

	// Initialize with default config
	cfg, err := config.GetDefaultConfig()
	if err != nil {
		t.Fatalf("Failed to get default config: %v", err)
	}
	if err := cfg.Save(); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	newPath := "/tmp/test/archives"
	cmd := configSetCmd
	cmd.SetArgs([]string{"archive-dir", newPath})

	err = runConfigSet(cmd, []string{"archive-dir", newPath})
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// Verify the config was updated
	loadedCfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	if loadedCfg.ArchiveDirectory != newPath {
		t.Errorf("ArchiveDirectory = %s, want %s", loadedCfg.ArchiveDirectory, newPath)
	}
}
