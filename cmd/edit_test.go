package cmd

import (
	"testing"
	"time"

	"github.com/m-uesaka/cli-record/internal/models"
)

func TestParseEditDateTime_Valid(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "valid datetime",
			input:    "2025-12-29 14:30:00",
			expected: "2025-12-29 14:30:00",
		},
		{
			name:     "midnight",
			input:    "2025-01-01 00:00:00",
			expected: "2025-01-01 00:00:00",
		},
		{
			name:     "end of day",
			input:    "2025-12-31 23:59:59",
			expected: "2025-12-31 23:59:59",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseEditDateTime(tt.input)
			if err != nil {
				t.Errorf("parseEditDateTime(%s) returned error: %v", tt.input, err)
			}
			
			formatted := result.Format("2006-01-02 15:04:05")
			if formatted != tt.expected {
				t.Errorf("parseEditDateTime(%s) = %s, want %s", tt.input, formatted, tt.expected)
			}
		})
	}
}

func TestParseEditDateTime_Invalid(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "empty string",
			input: "",
		},
		{
			name:  "invalid format",
			input: "2025-12-29",
		},
		{
			name:  "wrong format",
			input: "12/29/2025 14:30:00",
		},
		{
			name:  "missing seconds",
			input: "2025-12-29 14:30",
		},
		{
			name:  "invalid date",
			input: "2025-13-40 14:30:00",
		},
		{
			name:  "invalid time",
			input: "2025-12-29 25:70:00",
		},
		{
			name:  "random text",
			input: "not a date",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseEditDateTime(tt.input)
			if err == nil {
				t.Errorf("parseEditDateTime(%s) expected error, got nil", tt.input)
			}
		})
	}
}

func TestValidateEndTimeAfterStart(t *testing.T) {
	startTime := time.Date(2025, 12, 29, 10, 0, 0, 0, time.UTC)
	
	tests := []struct {
		name      string
		endTime   *time.Time
		shouldErr bool
	}{
		{
			name: "valid end time",
			endTime: func() *time.Time {
				t := startTime.Add(2 * time.Hour)
				return &t
			}(),
			shouldErr: false,
		},
		{
			name: "end time before start",
			endTime: func() *time.Time {
				t := startTime.Add(-2 * time.Hour)
				return &t
			}(),
			shouldErr: true,
		},
		{
			name: "end time same as start",
			endTime: func() *time.Time {
				t := startTime
				return &t
			}(),
			shouldErr: true,
		},
		{
			name:      "nil end time",
			endTime:   nil,
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entry := &models.TimeEntry{
				StartTime: startTime,
				EndTime:   tt.endTime,
			}
			
			var hasError bool
			if entry.EndTime != nil && entry.EndTime.Before(entry.StartTime) {
				hasError = true
			}
			
			if hasError != tt.shouldErr {
				t.Errorf("validation hasError = %v, want %v", hasError, tt.shouldErr)
			}
		})
	}
}
