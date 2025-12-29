package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/m-uesaka/cli-record/internal/models"
)

type Storage interface {
	SaveEntry(entry *models.TimeEntry) error
	GetEntry(id string) (*models.TimeEntry, error)
	ListEntries() ([]*models.TimeEntry, error)
	UpdateEntry(entry *models.TimeEntry) error
	GetRunningEntry() (*models.TimeEntry, error)
	ListTags() ([]string, error)
}

type JSONStorage struct {
	filePath string
	mu       sync.RWMutex
}

func NewJSONStorage() (*JSONStorage, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	dataDir := filepath.Join(homeDir, ".cli-record")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	return &JSONStorage{
		filePath: filepath.Join(dataDir, "data.json"),
	}, nil
}

func (s *JSONStorage) SaveEntry(entry *models.TimeEntry) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entries, err := s.loadEntries()
	if err != nil {
		return err
	}

	entries = append(entries, entry)
	return s.saveEntries(entries)
}

func (s *JSONStorage) GetEntry(id string) (*models.TimeEntry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entries, err := s.loadEntries()
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.ID == id {
			return entry, nil
		}
	}

	return nil, fmt.Errorf("entry with ID %s not found", id)
}

func (s *JSONStorage) ListEntries() ([]*models.TimeEntry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.loadEntries()
}

func (s *JSONStorage) UpdateEntry(entry *models.TimeEntry) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entries, err := s.loadEntries()
	if err != nil {
		return err
	}

	for i, e := range entries {
		if e.ID == entry.ID {
			entries[i] = entry
			return s.saveEntries(entries)
		}
	}

	return fmt.Errorf("entry with ID %s not found", entry.ID)
}

func (s *JSONStorage) GetRunningEntry() (*models.TimeEntry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entries, err := s.loadEntries()
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsRunning() {
			return entry, nil
		}
	}

	return nil, nil
}

func (s *JSONStorage) ListTags() ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entries, err := s.loadEntries()
	if err != nil {
		return nil, err
	}

	tagMap := make(map[string]bool)
	for _, entry := range entries {
		for _, tag := range entry.Tags {
			tagMap[tag] = true
		}
	}

	tags := make([]string, 0, len(tagMap))
	for tag := range tagMap {
		tags = append(tags, tag)
	}

	return tags, nil
}

func (s *JSONStorage) loadEntries() ([]*models.TimeEntry, error) {
	if _, err := os.Stat(s.filePath); os.IsNotExist(err) {
		return []*models.TimeEntry{}, nil
	}

	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	if len(data) == 0 {
		return []*models.TimeEntry{}, nil
	}

	var entries []*models.TimeEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return entries, nil
}

func (s *JSONStorage) saveEntries(entries []*models.TimeEntry) error {
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	if err := os.WriteFile(s.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
