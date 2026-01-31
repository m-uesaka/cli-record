package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

// Config represents the application configuration
type Config struct {
	DataFilePath     string            `toml:"data_file_path"`
	TimeFormat       string            `toml:"time_format"`
	ArchiveDirectory string            `toml:"archive_directory"`
	GroupColors      map[string]string `toml:"group_colors,omitempty"`
}

// Default configuration values
const (
	DefaultTimeFormat = "24h"
)

// Valid ANSI color names
var validColors = map[string]bool{
	"black":   true,
	"red":     true,
	"green":   true,
	"yellow":  true,
	"blue":    true,
	"magenta": true,
	"cyan":    true,
	"white":   true,
	"gray":    true,
	"bright-red":     true,
	"bright-green":   true,
	"bright-yellow":  true,
	"bright-blue":    true,
	"bright-magenta": true,
	"bright-cyan":    true,
	"bright-white":   true,
}

// IsValidColor checks if a color name is valid
func IsValidColor(color string) bool {
	_, ok := validColors[color]
	return ok
}

// GetDefaultConfig returns the default configuration
func GetDefaultConfig() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	return &Config{
		DataFilePath:     filepath.Join(homeDir, ".cli-record", "data.json"),
		TimeFormat:       DefaultTimeFormat,
		ArchiveDirectory: filepath.Join(homeDir, ".cli-record", "archives"),
		GroupColors:      make(map[string]string),
	}, nil
}

// GetConfigPath returns the path to the configuration file
func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".config", "cli-record")
	return filepath.Join(configDir, "config.toml"), nil
}

// Load loads the configuration from file, or returns default if file doesn't exist
func Load() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	// If config file doesn't exist, return default config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return GetDefaultConfig()
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse TOML
	var cfg Config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &cfg, nil
}

// Save saves the configuration to file
func (c *Config) Save() error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	// Create config directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Validate before saving
	if err := c.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Marshal to TOML
	data, err := toml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Validate validates the configuration values
func (c *Config) Validate() error {
	// Validate time format
	if c.TimeFormat != "12h" && c.TimeFormat != "24h" {
		return fmt.Errorf("time_format must be either '12h' or '24h', got '%s'", c.TimeFormat)
	}

	// Validate data file path is not empty
	if c.DataFilePath == "" {
		return fmt.Errorf("data_file_path cannot be empty")
	}

	// Validate archive directory is not empty
	if c.ArchiveDirectory == "" {
		return fmt.Errorf("archive_directory cannot be empty")
	}

	// Validate group colors
	for group, color := range c.GroupColors {
		if !IsValidColor(color) {
			return fmt.Errorf("invalid color '%s' for group '%s'", color, group)
		}
	}

	return nil
}

// IsDefault checks if a field is using the default value
func (c *Config) IsDefault(field string) (bool, error) {
	defaultCfg, err := GetDefaultConfig()
	if err != nil {
		return false, err
	}

	switch field {
	case "data_file_path":
		return c.DataFilePath == defaultCfg.DataFilePath, nil
	case "time_format":
		return c.TimeFormat == defaultCfg.TimeFormat, nil
	case "archive_directory":
		return c.ArchiveDirectory == defaultCfg.ArchiveDirectory, nil
	default:
		return false, fmt.Errorf("unknown field: %s", field)
	}
}

// Reset resets the configuration to default values
func Reset() error {
	defaultCfg, err := GetDefaultConfig()
	if err != nil {
		return err
	}

	return defaultCfg.Save()
}

// SetGroupColor sets the color for a specific group
func (c *Config) SetGroupColor(group, color string) error {
	if !IsValidColor(color) {
		return fmt.Errorf("invalid color '%s'", color)
	}
	
	if c.GroupColors == nil {
		c.GroupColors = make(map[string]string)
	}
	
	c.GroupColors[group] = color
	return nil
}

// GetGroupColor returns the color for a specific group, or empty string if not set
func (c *Config) GetGroupColor(group string) string {
	if c.GroupColors == nil {
		return ""
	}
	return c.GroupColors[group]
}

// RemoveGroupColor removes the color setting for a specific group
func (c *Config) RemoveGroupColor(group string) {
	if c.GroupColors != nil {
		delete(c.GroupColors, group)
	}
}

// Helper functions for testing

func loadFromPath(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &cfg, nil
}

func saveToPath(c *Config, path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := c.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	data, err := toml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
