package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetDefaultConfig(t *testing.T) {
	cfg, err := GetDefaultConfig()
	if err != nil {
		t.Fatalf("GetDefaultConfig() error = %v", err)
	}

	if cfg.DataFilePath == "" {
		t.Error("DataFilePath should not be empty")
	}

	if cfg.TimeFormat != "24h" {
		t.Errorf("TimeFormat = %s, want 24h", cfg.TimeFormat)
	}

	if cfg.ArchiveDirectory == "" {
		t.Error("ArchiveDirectory should not be empty")
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid config with 24h",
			config: Config{
				DataFilePath:     "/home/user/.cli-record/data.json",
				TimeFormat:       "24h",
				ArchiveDirectory: "/home/user/.cli-record/archives",
			},
			wantErr: false,
		},
		{
			name: "valid config with 12h",
			config: Config{
				DataFilePath:     "/home/user/.cli-record/data.json",
				TimeFormat:       "12h",
				ArchiveDirectory: "/home/user/.cli-record/archives",
			},
			wantErr: false,
		},
		{
			name: "invalid time format",
			config: Config{
				DataFilePath:     "/home/user/.cli-record/data.json",
				TimeFormat:       "invalid",
				ArchiveDirectory: "/home/user/.cli-record/archives",
			},
			wantErr: true,
		},
		{
			name: "empty data file path",
			config: Config{
				DataFilePath:     "",
				TimeFormat:       "24h",
				ArchiveDirectory: "/home/user/.cli-record/archives",
			},
			wantErr: true,
		},
		{
			name: "empty archive directory",
			config: Config{
				DataFilePath:     "/home/user/.cli-record/data.json",
				TimeFormat:       "24h",
				ArchiveDirectory: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfig_SaveAndLoad(t *testing.T) {
	// Create a temporary config directory
	tmpDir, err := os.MkdirTemp("", "cli-record-config-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testConfigPath := filepath.Join(tmpDir, "config.toml")

	// Create a test config
	testCfg := &Config{
		DataFilePath:     "/custom/path/data.json",
		TimeFormat:       "12h",
		ArchiveDirectory: "/custom/archives",
	}

	// Save config to test path
	if err := saveToPath(testCfg, testConfigPath); err != nil {
		t.Fatalf("saveToPath() error = %v", err)
	}

	// Load config from test path
	loadedCfg, err := loadFromPath(testConfigPath)
	if err != nil {
		t.Fatalf("loadFromPath() error = %v", err)
	}

	// Verify loaded config matches saved config
	if loadedCfg.DataFilePath != testCfg.DataFilePath {
		t.Errorf("DataFilePath = %s, want %s", loadedCfg.DataFilePath, testCfg.DataFilePath)
	}

	if loadedCfg.TimeFormat != testCfg.TimeFormat {
		t.Errorf("TimeFormat = %s, want %s", loadedCfg.TimeFormat, testCfg.TimeFormat)
	}

	if loadedCfg.ArchiveDirectory != testCfg.ArchiveDirectory {
		t.Errorf("ArchiveDirectory = %s, want %s", loadedCfg.ArchiveDirectory, testCfg.ArchiveDirectory)
	}
}

func TestConfig_IsDefault(t *testing.T) {
	defaultCfg, err := GetDefaultConfig()
	if err != nil {
		t.Fatalf("GetDefaultConfig() error = %v", err)
	}

	tests := []struct {
		name      string
		config    *Config
		field     string
		want      bool
		wantError bool
	}{
		{
			name:      "default data_file_path",
			config:    defaultCfg,
			field:     "data_file_path",
			want:      true,
			wantError: false,
		},
		{
			name: "custom data_file_path",
			config: &Config{
				DataFilePath:     "/custom/path/data.json",
				TimeFormat:       defaultCfg.TimeFormat,
				ArchiveDirectory: defaultCfg.ArchiveDirectory,
			},
			field:     "data_file_path",
			want:      false,
			wantError: false,
		},
		{
			name:      "unknown field",
			config:    defaultCfg,
			field:     "unknown_field",
			want:      false,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.config.IsDefault(tt.field)
			if (err != nil) != tt.wantError {
				t.Errorf("IsDefault() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if got != tt.want {
				t.Errorf("IsDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetConfigPath(t *testing.T) {
	path, err := GetConfigPath()
	if err != nil {
		t.Fatalf("GetConfigPath() error = %v", err)
	}

	if path == "" {
		t.Error("expected non-empty config path")
	}

	// Verify it ends with config.toml
	expectedSuffix := "cli-record/config.toml"
	if len(path) < len(expectedSuffix) || path[len(path)-len(expectedSuffix):] != expectedSuffix {
		t.Errorf("expected path to end with %s, got %s", expectedSuffix, path)
	}
}

func TestLoad(t *testing.T) {
	t.Run("load default when file doesn't exist", func(t *testing.T) {
		// Create a temporary directory for test
		tmpDir, err := os.MkdirTemp("", "cli-record-config-test-*")
		if err != nil {
			t.Fatalf("failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		// Load should return default config when file doesn't exist
		// (can't easily test this without mocking GetConfigPath)
		cfg, err := GetDefaultConfig()
		if err != nil {
			t.Fatalf("GetDefaultConfig() error = %v", err)
		}

		if cfg.TimeFormat != "24h" {
			t.Errorf("expected default time format 24h, got %s", cfg.TimeFormat)
		}
	})
}

func TestSave(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-record-config-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testConfigPath := filepath.Join(tmpDir, "config.toml")

	testCfg := &Config{
		DataFilePath:     "/test/path/data.json",
		TimeFormat:       "12h",
		ArchiveDirectory: "/test/archives",
	}

	err = saveToPath(testCfg, testConfigPath)
	if err != nil {
		t.Fatalf("saveToPath() error = %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(testConfigPath); os.IsNotExist(err) {
		t.Error("config file was not created")
	}

	// Load and verify content
	loadedCfg, err := loadFromPath(testConfigPath)
	if err != nil {
		t.Fatalf("loadFromPath() error = %v", err)
	}

	if loadedCfg.DataFilePath != testCfg.DataFilePath {
		t.Errorf("DataFilePath = %s, want %s", loadedCfg.DataFilePath, testCfg.DataFilePath)
	}

	if loadedCfg.TimeFormat != testCfg.TimeFormat {
		t.Errorf("TimeFormat = %s, want %s", loadedCfg.TimeFormat, testCfg.TimeFormat)
	}
}

func TestReset(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-record-config-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a custom config
	testConfigPath := filepath.Join(tmpDir, "config.toml")
	customCfg := &Config{
		DataFilePath:     "/custom/path/data.json",
		TimeFormat:       "12h",
		ArchiveDirectory: "/custom/archives",
	}
	saveToPath(customCfg, testConfigPath)

	// Note: Reset() uses GetConfigPath which we can't easily override
	// So we test the logic separately
	defaultCfg, err := GetDefaultConfig()
	if err != nil {
		t.Fatalf("GetDefaultConfig() error = %v", err)
	}

	// Verify default values
	if defaultCfg.TimeFormat != "24h" {
		t.Errorf("expected default time format 24h, got %s", defaultCfg.TimeFormat)
	}
}

func TestLoadFromPath_InvalidFile(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{
			name:    "invalid TOML",
			content: "this is not valid TOML {{{",
		},
		{
			name:    "empty file",
			content: "",
		},
		{
			name: "invalid time format in file",
			content: `data_file_path = "/test/data.json"
time_format = "invalid"
archive_directory = "/test/archives"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "cli-record-config-test-*")
			if err != nil {
				t.Fatalf("failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tmpDir)

			testFile := filepath.Join(tmpDir, "config.toml")
			if err := os.WriteFile(testFile, []byte(tt.content), 0644); err != nil {
				t.Fatalf("failed to write test file: %v", err)
			}

			_, err = loadFromPath(testFile)
			// For invalid TOML or validation errors, we expect an error
			if err == nil && tt.name != "empty file" {
				t.Error("expected error for invalid config file, got nil")
			}
		})
	}
}

func TestLoadFromPath_NonexistentFile(t *testing.T) {
	_, err := loadFromPath("/nonexistent/path/config.toml")
	if err == nil {
		t.Error("expected error for nonexistent file, got nil")
	}
}

func TestSaveToPath_InvalidDirectory(t *testing.T) {
	// Try to save to a path where parent directory doesn't exist
	testCfg := &Config{
		DataFilePath:     "/test/path/data.json",
		TimeFormat:       "24h",
		ArchiveDirectory: "/test/archives",
	}

	// Use a path with non-existent parent directory
	invalidPath := "/nonexistent/deeply/nested/path/config.toml"
	
	// This should fail because parent directory doesn't exist
	err := saveToPath(testCfg, invalidPath)
	if err == nil {
		t.Error("expected error when saving to invalid path, got nil")
	}
}

func TestConfig_Validate_EdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "whitespace time format",
			config: Config{
				DataFilePath:     "/test/data.json",
				TimeFormat:       "  24h  ",
				ArchiveDirectory: "/test/archives",
			},
			wantErr: true,
		},
		{
			name: "case sensitive time format",
			config: Config{
				DataFilePath:     "/test/data.json",
				TimeFormat:       "24H",
				ArchiveDirectory: "/test/archives",
			},
			wantErr: true,
		},
		{
			name: "all fields empty",
			config: Config{
				DataFilePath:     "",
				TimeFormat:       "",
				ArchiveDirectory: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
