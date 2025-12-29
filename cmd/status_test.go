package cmd

import (
	"testing"
	"time"
)

func TestFormatElapsedTime(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		expected string
	}{
		{
			name:     "just now",
			duration: 30 * time.Second,
			expected: "just now",
		},
		{
			name:     "one minute",
			duration: 1 * time.Minute,
			expected: "1 minute",
		},
		{
			name:     "multiple minutes",
			duration: 5 * time.Minute,
			expected: "5 minutes",
		},
		{
			name:     "one hour exactly",
			duration: 1 * time.Hour,
			expected: "1 hour",
		},
		{
			name:     "one hour with minutes",
			duration: 1*time.Hour + 30*time.Minute,
			expected: "1 hour 30 minutes",
		},
		{
			name:     "multiple hours",
			duration: 3 * time.Hour,
			expected: "3 hours",
		},
		{
			name:     "multiple hours with minutes",
			duration: 2*time.Hour + 15*time.Minute,
			expected: "2 hours 15 minutes",
		},
		{
			name:     "one day exactly",
			duration: 24 * time.Hour,
			expected: "1 day",
		},
		{
			name:     "one day with hours",
			duration: 24*time.Hour + 5*time.Hour,
			expected: "1 day 5 hours",
		},
		{
			name:     "multiple days",
			duration: 3*24*time.Hour + 12*time.Hour,
			expected: "3 days 12 hours",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatElapsedTime(tt.duration)
			if result != tt.expected {
				t.Errorf("formatElapsedTime(%v) = %q, want %q", tt.duration, result, tt.expected)
			}
		})
	}
}

func TestFormatTags(t *testing.T) {
	tests := []struct {
		name     string
		tags     []string
		expected string
	}{
		{
			name:     "empty tags",
			tags:     []string{},
			expected: "",
		},
		{
			name:     "single tag",
			tags:     []string{"development"},
			expected: "development",
		},
		{
			name:     "multiple tags",
			tags:     []string{"development", "review", "testing"},
			expected: "development, review, testing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatTags(tt.tags)
			if result != tt.expected {
				t.Errorf("formatTags(%v) = %q, want %q", tt.tags, result, tt.expected)
			}
		})
	}
}
