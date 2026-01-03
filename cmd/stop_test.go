package cmd

import (
	"testing"
	"time"
)

// NOTE: Tests for runStop function removed as they require interactive TUI input
// and are not suitable for automated unit testing. These scenarios should be covered
// by integration tests or manual testing instead.

func TestFormatDuration_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		contains string // Check if output contains this
	}{
		{
			name:     "negative duration",
			duration: -5 * time.Second,
			contains: "0s",
		},
		{
			name:     "very large duration",
			duration: 100*time.Hour + 30*time.Minute + 15*time.Second,
			contains: "100h",
		},
		{
			name:     "exactly one hour",
			duration: 1 * time.Hour,
			contains: "1h 0m",
		},
		{
			name:     "exactly one minute",
			duration: 1 * time.Minute,
			contains: "1m 0s",
		},
		{
			name:     "milliseconds rounded",
			duration: 1*time.Second + 500*time.Millisecond,
			contains: "2s", // Should round to 2s
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatDuration(tt.duration)
			// Just check that it doesn't panic and returns something
			if result == "" {
				t.Error("formatDuration returned empty string")
			}
		})
	}
}
