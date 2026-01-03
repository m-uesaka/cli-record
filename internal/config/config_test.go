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
