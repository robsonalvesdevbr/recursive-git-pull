package colors

import (
	"fmt"
	"os"
)

// ANSI color codes
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	Gray   = "\033[37m"
	White  = "\033[97m"

	// Bold variants
	BoldRed    = "\033[1;31m"
	BoldGreen  = "\033[1;32m"
	BoldYellow = "\033[1;33m"
	BoldBlue   = "\033[1;34m"
	BoldWhite  = "\033[1;97m"
)

var forceNoColor bool

// SetForceNoColor allows disabling colors programmatically
func SetForceNoColor(noColor bool) {
	forceNoColor = noColor
}

// ColorsEnabled checks if colors should be displayed
func ColorsEnabled() bool {
	// Check if colors are force-disabled
	if forceNoColor {
		return false
	}
	
	// Check if output is being piped or redirected
	if fileInfo, _ := os.Stdout.Stat(); (fileInfo.Mode() & os.ModeCharDevice) == 0 {
		return false
	}
	
	// Check NO_COLOR environment variable
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	
	// Check TERM environment variable
	term := os.Getenv("TERM")
	if term == "" || term == "dumb" {
		return false
	}
	
	return true
}

// Colorize applies color to text if colors are enabled
func Colorize(color, text string) string {
	if !ColorsEnabled() {
		return text
	}
	return color + text + Reset
}

// Success returns green colored text for success messages
func Success(text string) string {
	return Colorize(BoldGreen, text)
}

// Error returns red colored text for error messages
func Error(text string) string {
	return Colorize(BoldRed, text)
}

// Warning returns yellow colored text for warning messages
func Warning(text string) string {
	return Colorize(BoldYellow, text)
}

// Info returns blue colored text for info messages
func Info(text string) string {
	return Colorize(BoldBlue, text)
}

// Dim returns gray colored text for less important messages
func Dim(text string) string {
	return Colorize(Gray, text)
}

// Bold returns bold white text
func Bold(text string) string {
	return Colorize(BoldWhite, text)
}

// SuccessIcon returns a green checkmark symbol
func SuccessIcon() string {
	return Success("✓")
}

// ErrorIcon returns a red X symbol
func ErrorIcon() string {
	return Error("✗")
}

// WarningIcon returns a yellow warning symbol
func WarningIcon() string {
	return Warning("⚠")
}

// InfoIcon returns a blue info symbol
func InfoIcon() string {
	return Info("ℹ")
}

// Printf with color support
func Printf(color, format string, args ...interface{}) {
	fmt.Printf(Colorize(color, format), args...)
}

// Println with color support
func Println(color, text string) {
	fmt.Println(Colorize(color, text))
}