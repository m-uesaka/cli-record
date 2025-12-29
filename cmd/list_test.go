package cmd

import (
	"testing"
	"time"
)

func TestParseDateTime(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantErr   bool
		checkYear int
		checkMonth time.Month
		checkDay  int
	}{
		{
			name:      "date only format",
			input:     "2025-01-15",
			wantErr:   false,
			checkYear: 2025,
			checkMonth: time.January,
			checkDay:  15,
		},
		{
			name:      "date with time format",
			input:     "2025-01-15 14:30:00",
			wantErr:   false,
			checkYear: 2025,
			checkMonth: time.January,
			checkDay:  15,
		},
		{
			name:    "invalid format",
			input:   "2025/01/15",
			wantErr: true,
		},
		{
			name:    "invalid date",
			input:   "not-a-date",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseDateTime(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseDateTime(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if result.Year() != tt.checkYear || result.Month() != tt.checkMonth || result.Day() != tt.checkDay {
					t.Errorf("parseDateTime(%q) = %v, want year=%d, month=%v, day=%d",
						tt.input, result, tt.checkYear, tt.checkMonth, tt.checkDay)
				}
			}
		})
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		maxLen int
		want   string
	}{
		{
			name:   "no truncation needed",
			input:  "short",
			maxLen: 10,
			want:   "short",
		},
		{
			name:   "exact length",
			input:  "exactly10c",
			maxLen: 10,
			want:   "exactly10c",
		},
		{
			name:   "truncation needed",
			input:  "this is a very long string",
			maxLen: 10,
			want:   "this is...",
		},
		{
			name:   "truncation with small maxLen",
			input:  "hello world",
			maxLen: 5,
			want:   "he...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := truncate(tt.input, tt.maxLen)
			if result != tt.want {
				t.Errorf("truncate(%q, %d) = %q, want %q", tt.input, tt.maxLen, result, tt.want)
			}
		})
	}
}
