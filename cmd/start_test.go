package cmd

import (
	"reflect"
	"testing"
	"time"
)

func TestParseTags(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: []string{},
		},
		{
			name:     "single tag",
			input:    "development",
			expected: []string{"development"},
		},
		{
			name:     "multiple tags",
			input:    "development,review,testing",
			expected: []string{"development", "review", "testing"},
		},
		{
			name:     "tags with spaces",
			input:    "development, review, testing",
			expected: []string{"development", "review", "testing"},
		},
		{
			name:     "tags with extra spaces",
			input:    "  development  ,  review  ,  testing  ",
			expected: []string{"development", "review", "testing"},
		},
		{
			name:     "tags with empty values",
			input:    "development,,review",
			expected: []string{"development", "review"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseTags(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("parseTags(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		expected string
	}{
		{
			name:     "only seconds",
			duration: 45 * time.Second,
			expected: "45s",
		},
		{
			name:     "minutes and seconds",
			duration: 2*time.Minute + 30*time.Second,
			expected: "2m 30s",
		},
		{
			name:     "hours, minutes and seconds",
			duration: 2*time.Hour + 30*time.Minute + 45*time.Second,
			expected: "2h 30m 45s",
		},
		{
			name:     "only hours",
			duration: 3 * time.Hour,
			expected: "3h 0m 0s",
		},
		{
			name:     "zero duration",
			duration: 0,
			expected: "0s",
		},
		{
			name:     "rounding seconds",
			duration: 2*time.Minute + 30*time.Second + 500*time.Millisecond,
			expected: "2m 31s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatDuration(tt.duration)
			if result != tt.expected {
				t.Errorf("formatDuration(%v) = %q, want %q", tt.duration, result, tt.expected)
			}
		})
	}
}
