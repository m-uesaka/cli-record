package tui

import (
	"fmt"
)

// ANSI color codes
const (
	ColorReset = "\033[0m"
	
	// Regular colors
	ColorBlack   = "\033[30m"
	ColorRed     = "\033[31m"
	ColorGreen   = "\033[32m"
	ColorYellow  = "\033[33m"
	ColorBlue    = "\033[34m"
	ColorMagenta = "\033[35m"
	ColorCyan    = "\033[36m"
	ColorWhite   = "\033[37m"
	ColorGray    = "\033[90m"
	
	// Bright colors
	ColorBrightRed     = "\033[91m"
	ColorBrightGreen   = "\033[92m"
	ColorBrightYellow  = "\033[93m"
	ColorBrightBlue    = "\033[94m"
	ColorBrightMagenta = "\033[95m"
	ColorBrightCyan    = "\033[96m"
	ColorBrightWhite   = "\033[97m"
)

var colorCodes = map[string]string{
	"black":          ColorBlack,
	"red":            ColorRed,
	"green":          ColorGreen,
	"yellow":         ColorYellow,
	"blue":           ColorBlue,
	"magenta":        ColorMagenta,
	"cyan":           ColorCyan,
	"white":          ColorWhite,
	"gray":           ColorGray,
	"bright-red":     ColorBrightRed,
	"bright-green":   ColorBrightGreen,
	"bright-yellow":  ColorBrightYellow,
	"bright-blue":    ColorBrightBlue,
	"bright-magenta": ColorBrightMagenta,
	"bright-cyan":    ColorBrightCyan,
	"bright-white":   ColorBrightWhite,
}

// Colorize applies a color to text
func Colorize(text, color string) string {
	if color == "" {
		return text
	}
	
	if code, ok := colorCodes[color]; ok {
		return code + text + ColorReset
	}
	
	return text
}

// GetColorCode returns the ANSI code for a color name
func GetColorCode(color string) string {
	if code, ok := colorCodes[color]; ok {
		return code
	}
	return ""
}

// Sprint formats and colorizes text
func Sprint(color, format string, args ...interface{}) string {
	text := fmt.Sprintf(format, args...)
	return Colorize(text, color)
}
