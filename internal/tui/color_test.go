package tui

import (
	"testing"
)

func TestColorize(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		color    string
		expected string
	}{
		{
			name:     "red color",
			text:     "Hello",
			color:    "red",
			expected: "\033[31mHello\033[0m",
		},
		{
			name:     "green color",
			text:     "World",
			color:    "green",
			expected: "\033[32mWorld\033[0m",
		},
		{
			name:     "no color",
			text:     "Test",
			color:    "",
			expected: "Test",
		},
		{
			name:     "invalid color",
			text:     "Test",
			color:    "invalid",
			expected: "Test",
		},
		{
			name:     "bright color",
			text:     "Bright",
			color:    "bright-blue",
			expected: "\033[94mBright\033[0m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Colorize(tt.text, tt.color)
			if result != tt.expected {
				t.Errorf("Colorize(%q, %q) = %q, want %q", tt.text, tt.color, result, tt.expected)
			}
		})
	}
}

func TestGetColorCode(t *testing.T) {
	tests := []struct {
		name     string
		color    string
		expected string
	}{
		{
			name:     "red",
			color:    "red",
			expected: "\033[31m",
		},
		{
			name:     "green",
			color:    "green",
			expected: "\033[32m",
		},
		{
			name:     "invalid",
			color:    "invalid",
			expected: "",
		},
		{
			name:     "bright-magenta",
			color:    "bright-magenta",
			expected: "\033[95m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetColorCode(tt.color)
			if result != tt.expected {
				t.Errorf("GetColorCode(%q) = %q, want %q", tt.color, result, tt.expected)
			}
		})
	}
}

func TestSprint(t *testing.T) {
	result := Sprint("red", "Hello %s", "World")
	expected := "\033[31mHello World\033[0m"
	
	if result != expected {
		t.Errorf("Sprint(red, ...) = %q, want %q", result, expected)
	}
}
