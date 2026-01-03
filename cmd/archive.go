package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/m-uesaka/cli-record/internal/config"
	"github.com/m-uesaka/cli-record/internal/storage"
	"github.com/spf13/cobra"
)

var (
	archiveOutput string
	archiveList   bool
)

var archiveCmd = &cobra.Command{
	Use:   "archive",
	Short: "Archive the current data file",
	Long:  `Create a timestamped backup of the current data file in the archive directory.`,
	RunE:  runArchive,
}

func init() {
	rootCmd.AddCommand(archiveCmd)
	archiveCmd.Flags().StringVarP(&archiveOutput, "output", "o", "", "Custom archive file name")
	archiveCmd.Flags().BoolVarP(&archiveList, "list", "l", false, "List available archives")
}

func runArchive(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// If list flag is set, list archives
	if archiveList {
		return listArchives(cfg.ArchiveDirectory)
	}

	store, err := storage.NewJSONStorage()
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}

	// Generate archive filename
	var archivePath string
	if archiveOutput != "" {
		archivePath = filepath.Join(cfg.ArchiveDirectory, archiveOutput)
	} else {
		timestamp := time.Now().Format("2006-01-02-150405")
		archivePath = filepath.Join(cfg.ArchiveDirectory, fmt.Sprintf("data-%s.json", timestamp))
	}

	// Create archive
	if err := store.ArchiveData(archivePath); err != nil {
		return fmt.Errorf("failed to create archive: %w", err)
	}

	// Get file size
	fileInfo, err := os.Stat(archivePath)
	if err != nil {
		return fmt.Errorf("failed to get archive info: %w", err)
	}

	fmt.Println("✓ Archive created successfully")
	fmt.Printf("  Location: %s\n", archivePath)
	fmt.Printf("  Size: %s\n", formatFileSize(fileInfo.Size()))

	return nil
}

func listArchives(archiveDir string) error {
	// Check if archive directory exists
	if _, err := os.Stat(archiveDir); os.IsNotExist(err) {
		fmt.Println("No archives found.")
		fmt.Printf("Archive directory: %s (does not exist)\n", archiveDir)
		return nil
	}

	// Read directory
	entries, err := os.ReadDir(archiveDir)
	if err != nil {
		return fmt.Errorf("failed to read archive directory: %w", err)
	}

	// Filter JSON files
	var archives []os.DirEntry
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".json" {
			archives = append(archives, entry)
		}
	}

	if len(archives) == 0 {
		fmt.Println("No archives found.")
		fmt.Printf("Archive directory: %s\n", archiveDir)
		return nil
	}

	// Sort by name (which includes timestamp) in descending order
	sort.Slice(archives, func(i, j int) bool {
		return archives[i].Name() > archives[j].Name()
	})

	fmt.Printf("Available Archives (%d)\n", len(archives))
	fmt.Println("====================")
	fmt.Println()

	for _, entry := range archives {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		fullPath := filepath.Join(archiveDir, entry.Name())
		fmt.Printf("  %s\n", entry.Name())
		fmt.Printf("    Date: %s\n", info.ModTime().Format("2006-01-02 15:04:05"))
		fmt.Printf("    Size: %s\n", formatFileSize(info.Size()))
		fmt.Printf("    Path: %s\n", fullPath)
		fmt.Println()
	}

	return nil
}

func formatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}
