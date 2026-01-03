package cmd

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/m-uesaka/cli-record/internal/storage"
	"github.com/m-uesaka/cli-record/internal/tui"
	"github.com/spf13/cobra"
)

var (
	restoreMerge   bool
	restoreReplace bool
	restorePreview bool
)

var restoreCmd = &cobra.Command{
	Use:   "restore <archive-file>",
	Short: "Restore data from an archived file",
	Long: `Restore data from a specified archived file. By default, merges with existing data.
Use --replace to replace all existing data (creates backup first).
Use --preview to see what would be restored without applying changes.`,
	Args: cobra.ExactArgs(1),
	RunE: runRestore,
}

func init() {
	rootCmd.AddCommand(restoreCmd)
	restoreCmd.Flags().BoolVar(&restoreMerge, "merge", true, "Merge with existing data (default)")
	restoreCmd.Flags().BoolVar(&restoreReplace, "replace", false, "Replace existing data")
	restoreCmd.Flags().BoolVar(&restorePreview, "preview", false, "Preview without applying")
}

func runRestore(cmd *cobra.Command, args []string) error {
	archiveFile := args[0]

	// If both merge and replace are specified, return error
	if cmd.Flags().Changed("merge") && cmd.Flags().Changed("replace") {
		return fmt.Errorf("cannot specify both --merge and --replace")
	}

	// Determine merge mode
	merge := !restoreReplace

	store, err := storage.NewJSONStorage()
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}

	if restorePreview {
		return previewRestore(store, archiveFile, merge)
	}

	// Show confirmation for replace mode
	if restoreReplace {
		confirmed, err := showRestoreConfirmation(archiveFile, merge)
		if err != nil {
			return err
		}
		if !confirmed {
			fmt.Println("Restore cancelled.")
			return nil
		}
	}

	// Restore data
	if err := store.RestoreData(archiveFile, merge); err != nil {
		return fmt.Errorf("failed to restore data: %w", err)
	}

	// Display success message
	fmt.Println("✓ Data restored successfully")
	if merge {
		fmt.Println("  Mode: Merged with existing data")
	} else {
		fmt.Println("  Mode: Replaced existing data (backup created)")
	}
	fmt.Printf("  Source: %s\n", archiveFile)

	// Show entry count
	entries, _ := store.ListEntries()
	fmt.Printf("  Total entries: %d\n", len(entries))

	return nil
}

func previewRestore(store storage.Storage, archiveFile string, merge bool) error {
	// This is a simplified preview - in a production system, you'd want to
	// parse the archive file and show what would be added/updated
	fmt.Println("Preview Mode")
	fmt.Println("============")
	fmt.Println()
	fmt.Printf("Archive file: %s\n", archiveFile)
	if merge {
		fmt.Println("Mode: Merge (add new entries, skip duplicates)")
	} else {
		fmt.Println("Mode: Replace (backup current data, load archived data)")
	}
	fmt.Println()
	fmt.Println("Note: This is a preview. No changes will be made.")
	fmt.Println("Run without --preview to apply the restore.")

	return nil
}

func showRestoreConfirmation(archiveFile string, merge bool) (bool, error) {
	var modeStr string
	if merge {
		modeStr = "merge with existing data"
	} else {
		modeStr = "REPLACE ALL existing data"
	}

	message := fmt.Sprintf(`Are you sure you want to restore from this archive?

Archive file: %s
Mode: %s

`, archiveFile, modeStr)

	if !merge {
		message += `WARNING: This will replace all your current data!
A backup will be created automatically.

`
	}

	message += `This action cannot be undone.`

	confirm := tui.NewConfirm(message)

	p := tea.NewProgram(confirm)
	m, err := p.Run()
	if err != nil {
		return false, fmt.Errorf("error running confirmation: %w", err)
	}

	model := m.(tui.ConfirmModel)
	return model.IsConfirmed(), nil
}
