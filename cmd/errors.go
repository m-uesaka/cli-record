package cmd

import (
	"fmt"
	"strings"
)

// ErrorWithSuggestion wraps an error with a helpful suggestion
type ErrorWithSuggestion struct {
	Err        error
	Suggestion string
}

func (e *ErrorWithSuggestion) Error() string {
	if e.Suggestion != "" {
		return fmt.Sprintf("%v\n\nSuggestion: %s", e.Err, e.Suggestion)
	}
	return e.Err.Error()
}

// NewErrorWithSuggestion creates a new error with a suggestion
func NewErrorWithSuggestion(err error, suggestion string) *ErrorWithSuggestion {
	return &ErrorWithSuggestion{
		Err:        err,
		Suggestion: suggestion,
	}
}

// ValidateDateFormat validates date format and provides helpful error
func ValidateDateFormat(dateStr string) error {
	if dateStr == "" {
		return nil
	}

	// Try parsing to see if it's valid
	_, err := parseDateTime(dateStr)
	if err != nil {
		return NewErrorWithSuggestion(
			fmt.Errorf("invalid date format: %s", dateStr),
			"Use format: YYYY-MM-DD or YYYY-MM-DD HH:MM:SS\nExample: 2025-01-15 or 2025-01-15 14:30:00",
		)
	}
	return nil
}

// SuggestCommand suggests a similar command when user makes a typo
func SuggestCommand(input string, available []string) string {
	input = strings.ToLower(input)
	
	// Simple Levenshtein distance check
	minDistance := 1000
	suggestion := ""
	
	for _, cmd := range available {
		distance := levenshteinDistance(input, strings.ToLower(cmd))
		if distance < minDistance && distance <= 2 {
			minDistance = distance
			suggestion = cmd
		}
	}
	
	if suggestion != "" {
		return fmt.Sprintf("Did you mean '%s'?", suggestion)
	}
	return ""
}

// levenshteinDistance calculates the Levenshtein distance between two strings
func levenshteinDistance(s1, s2 string) int {
	if len(s1) == 0 {
		return len(s2)
	}
	if len(s2) == 0 {
		return len(s1)
	}

	matrix := make([][]int, len(s1)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(s2)+1)
		matrix[i][0] = i
	}
	for j := range matrix[0] {
		matrix[0][j] = j
	}

	for i := 1; i <= len(s1); i++ {
		for j := 1; j <= len(s2); j++ {
			cost := 0
			if s1[i-1] != s2[j-1] {
				cost = 1
			}
			matrix[i][j] = min(
				matrix[i-1][j]+1,      // deletion
				matrix[i][j-1]+1,      // insertion
				matrix[i-1][j-1]+cost, // substitution
			)
		}
	}

	return matrix[len(s1)][len(s2)]
}

func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

// FormatErrorMessage formats error messages in a user-friendly way
func FormatErrorMessage(err error) string {
	if err == nil {
		return ""
	}

	// Check if it's an error with suggestion
	if errWithSugg, ok := err.(*ErrorWithSuggestion); ok {
		return errWithSugg.Error()
	}

	return err.Error()
}
