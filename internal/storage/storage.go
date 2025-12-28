package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/m-uesaka/cli-record/internal/models"
)

type Storage struct {
	filePath string
}

func NewStorage() (*Storage, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	dataDir := filepath.Join(homeDir, ".cli-record")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	return &Storage{
		filePath: filepath.Join(dataDir, "timeentries.json"),
	}, nil
}

func (s *Storage) LoadEntries() ([]models.TimeEntry, error) {
	if _, err := os.Stat(s.filePath); os.IsNotExist(err) {
		return []models.TimeEntry{}, nil
	}

	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var entries []models.TimeEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return entries, nil
}

func (s *Storage) SaveEntries(entries []models.TimeEntry) error {
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	if err := os.WriteFile(s.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
