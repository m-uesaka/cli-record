package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var Version = "0.1.0"

var rootCmd = &cobra.Command{
	Use:   "cli-record",
	Short: "A simple CLI time tracking application",
	Long: `cli-record is a terminal-based time tracking application that helps you record and analyze time spent on various tasks.

Features:
  • Start and stop time tracking with task names and tags
  • List and filter time entries
  • Generate detailed reports with various grouping options
  • Export data to CSV or JSON formats
  • All data stored locally in JSON format

Examples:
  # Start tracking time
  cli-record start --task "Write documentation" --tags "docs,writing"

  # Stop the current timer
  cli-record stop

  # List all entries
  cli-record list

  # View report grouped by task
  cli-record view --by task

  # Export weekly report to CSV
  cli-record view --by week --format csv --output report.csv`,
	Version: Version,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetVersionTemplate("cli-record version {{.Version}}\n")
}
