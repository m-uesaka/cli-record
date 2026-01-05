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
	GetEntryByPrefix(prefix string) (*models.TimeEntry, error)
	ListEntries() ([]*models.TimeEntry, error)
	UpdateEntry(entry *models.TimeEntry) error
	GetRunningEntry() (*models.TimeEntry, error)
	ListTags() ([]string, error)
	DeleteEntry(id string) error
	DeleteEntryByPrefix(prefix string) error
	ArchiveData(archivePath string) error
	RestoreData(archivePath string, merge bool) error
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

	if entry == nil {
		return fmt.Errorf("entry cannot be nil")
	}

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

	if id == "" {
		return nil, fmt.Errorf("entry ID cannot be empty")
	}

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

func (s *JSONStorage) GetEntryByPrefix(prefix string) (*models.TimeEntry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if prefix == "" {
		return nil, fmt.Errorf("prefix cannot be empty")
	}

	entries, err := s.loadEntries()
	if err != nil {
		return nil, err
	}

	var matches []*models.TimeEntry
	for _, entry := range entries {
		if len(entry.ID) >= len(prefix) && entry.ID[:len(prefix)] == prefix {
			matches = append(matches, entry)
		}
	}

	if len(matches) == 0 {
		return nil, fmt.Errorf("no entry found with ID prefix %s", prefix)
	}

	if len(matches) > 1 {
		return nil, fmt.Errorf("ambiguous ID prefix %s: matches %d entries", prefix, len(matches))
	}

	return matches[0], nil
}

func (s *JSONStorage) ListEntries() ([]*models.TimeEntry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.loadEntries()
}

func (s *JSONStorage) UpdateEntry(entry *models.TimeEntry) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if entry == nil {
		return fmt.Errorf("entry cannot be nil")
	}

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

func (s *JSONStorage) DeleteEntry(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if id == "" {
		return fmt.Errorf("entry ID cannot be empty")
	}

	entries, err := s.loadEntries()
	if err != nil {
		return err
	}

	// Find and remove the entry
	newEntries := make([]*models.TimeEntry, 0, len(entries))
	found := false
	for _, entry := range entries {
		if entry.ID == id {
			found = true
			continue
		}
		newEntries = append(newEntries, entry)
	}

	if !found {
		return fmt.Errorf("entry with ID %s not found", id)
	}

	return s.saveEntries(newEntries)
}

func (s *JSONStorage) DeleteEntryByPrefix(prefix string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if prefix == "" {
		return fmt.Errorf("prefix cannot be empty")
	}

	entries, err := s.loadEntries()
	if err != nil {
		return err
	}

	// Find matching entries
	var matches []*models.TimeEntry
	for _, entry := range entries {
		if len(entry.ID) >= len(prefix) && entry.ID[:len(prefix)] == prefix {
			matches = append(matches, entry)
		}
	}

	if len(matches) == 0 {
		return fmt.Errorf("no entry found with ID prefix %s", prefix)
	}

	if len(matches) > 1 {
		return fmt.Errorf("ambiguous ID prefix %s: matches %d entries", prefix, len(matches))
	}

	// Remove the matched entry
	newEntries := make([]*models.TimeEntry, 0, len(entries))
	for _, entry := range entries {
		if entry.ID != matches[0].ID {
			newEntries = append(newEntries, entry)
		}
	}

	return s.saveEntries(newEntries)
}

func (s *JSONStorage) ArchiveData(archivePath string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Read current data file
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("data file does not exist")
		}
		return fmt.Errorf("failed to read data file: %w", err)
	}

	// Create archive directory if it doesn't exist
	archiveDir := filepath.Dir(archivePath)
	if err := os.MkdirAll(archiveDir, 0755); err != nil {
		return fmt.Errorf("failed to create archive directory: %w", err)
	}

	// Write to archive file
	if err := os.WriteFile(archivePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write archive file: %w", err)
	}

	return nil
}

func (s *JSONStorage) RestoreData(archivePath string, merge bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Read archive file
	archiveData, err := os.ReadFile(archivePath)
	if err != nil {
		return fmt.Errorf("failed to read archive file: %w", err)
	}

	// Parse archived entries
	var archivedEntries []*models.TimeEntry
	if err := json.Unmarshal(archiveData, &archivedEntries); err != nil {
		return fmt.Errorf("failed to parse archive file: %w", err)
	}

	if merge {
		// Load current entries
		currentEntries, err := s.loadEntries()
		if err != nil {
			return err
		}

		// Merge entries, avoiding duplicates by ID
		entryMap := make(map[string]*models.TimeEntry)
		for _, entry := range currentEntries {
			entryMap[entry.ID] = entry
		}

		for _, entry := range archivedEntries {
			if _, exists := entryMap[entry.ID]; !exists {
				entryMap[entry.ID] = entry
			}
		}

		// Convert map back to slice
		mergedEntries := make([]*models.TimeEntry, 0, len(entryMap))
		for _, entry := range entryMap {
			mergedEntries = append(mergedEntries, entry)
		}

		return s.saveEntries(mergedEntries)
	} else {
		// Replace mode: backup current data first
		backupPath := s.filePath + ".backup"
		if _, err := os.Stat(s.filePath); err == nil {
			if err := os.Rename(s.filePath, backupPath); err != nil {
				return fmt.Errorf("failed to backup current data: %w", err)
			}
		}

		return s.saveEntries(archivedEntries)
	}
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
