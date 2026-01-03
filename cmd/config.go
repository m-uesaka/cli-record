package cmd

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/m-uesaka/cli-record/internal/config"
	"github.com/m-uesaka/cli-record/internal/tui"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage application configuration",
	Long:  `Manage application configuration including data file location, time format, and archive directory.`,
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Display current configuration",
	Long:  `Display all configuration values and indicate which values are using defaults.`,
	RunE:  runConfigShow,
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Long: `Set a configuration value. Valid keys:
  - data-path: Set data file location
  - time-format: Set time format (12h or 24h)
  - archive-dir: Set archive directory location`,
	Args: cobra.ExactArgs(2),
	RunE: runConfigSet,
}

var configResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset configuration to defaults",
	Long:  `Reset all configuration values to their defaults. This action requires confirmation.`,
	RunE:  runConfigReset,
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configResetCmd)
}

func runConfigShow(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	fmt.Println("Current Configuration")
	fmt.Println("====================")
	fmt.Println()

	// Data file path
	isDefault, _ := cfg.IsDefault("data_file_path")
	if isDefault {
		fmt.Printf("data-path:    %s (default)\n", cfg.DataFilePath)
	} else {
		fmt.Printf("data-path:    %s\n", cfg.DataFilePath)
	}

	// Time format
	isDefault, _ = cfg.IsDefault("time_format")
	if isDefault {
		fmt.Printf("time-format:  %s (default)\n", cfg.TimeFormat)
	} else {
		fmt.Printf("time-format:  %s\n", cfg.TimeFormat)
	}

	// Archive directory
	isDefault, _ = cfg.IsDefault("archive_directory")
	if isDefault {
		fmt.Printf("archive-dir:  %s (default)\n", cfg.ArchiveDirectory)
	} else {
		fmt.Printf("archive-dir:  %s\n", cfg.ArchiveDirectory)
	}

	fmt.Println()
	configPath, _ := config.GetConfigPath()
	fmt.Printf("Configuration file: %s\n", configPath)

	return nil
}

func runConfigSet(cmd *cobra.Command, args []string) error {
	key := args[0]
	value := args[1]

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Update configuration based on key
	switch key {
	case "data-path":
		cfg.DataFilePath = value
	case "time-format":
		if value != "12h" && value != "24h" {
			return NewErrorWithSuggestion(
				fmt.Errorf("invalid time format: %s", value),
				"Valid values are: 12h, 24h",
			)
		}
		cfg.TimeFormat = value
	case "archive-dir":
		cfg.ArchiveDirectory = value
	default:
		return NewErrorWithSuggestion(
			fmt.Errorf("unknown configuration key: %s", key),
			"Valid keys are: data-path, time-format, archive-dir",
		)
	}

	// Save configuration
	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	fmt.Printf("✓ Configuration updated\n")
	fmt.Printf("  %s = %s\n", key, value)

	return nil
}

func runConfigReset(cmd *cobra.Command, args []string) error {
	// Prompt for confirmation
	confirmed, err := showConfigResetConfirmation()
	if err != nil {
		return err
	}

	if !confirmed {
		fmt.Println("Reset cancelled.")
		return nil
	}

	// Reset configuration
	if err := config.Reset(); err != nil {
		return fmt.Errorf("failed to reset configuration: %w", err)
	}

	fmt.Println("✓ Configuration reset to defaults")

	// Show the default configuration
	defaultCfg, _ := config.GetDefaultConfig()
	fmt.Println()
	fmt.Printf("data-path:   %s\n", defaultCfg.DataFilePath)
	fmt.Printf("time-format: %s\n", defaultCfg.TimeFormat)
	fmt.Printf("archive-dir: %s\n", defaultCfg.ArchiveDirectory)

	return nil
}

func showConfigResetConfirmation() (bool, error) {
	message := `Are you sure you want to reset configuration to defaults?

This will overwrite your current configuration.

This action cannot be undone.`

	confirm := tui.NewConfirm(message)

	p := tea.NewProgram(confirm)
	m, err := p.Run()
	if err != nil {
		return false, fmt.Errorf("error running confirmation: %w", err)
	}

	model := m.(tui.ConfirmModel)
	return model.IsConfirmed(), nil
}
