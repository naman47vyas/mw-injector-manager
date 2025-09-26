package tui

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// ClearScreen clears the terminal screen
func ClearScreen() {
	fmt.Print("\033[2J\033[H")
}

// PrintSeparator prints a visual separator line
func PrintSeparator(char string, length int) {
	for i := 0; i < length; i++ {
		fmt.Print(char)
	}
	fmt.Println()
}

// PrintHeader prints a styled header
func PrintHeader(text string) {
	fmt.Printf("%s%s%s%s\n", Bold, Cyan, text, Reset)
	PrintSeparator("=", len(text))
}

// Confirm prompts user for confirmation
func Confirm(message string) bool {
	fmt.Printf("%s%s (y/N): %s", Yellow, message, Reset)
	var response string
	fmt.Scanln(&response)
	return response == "y" || response == "Y" || response == "yes"
}

func truncateString(s string, maxWidth int) string {
	if maxWidth <= 0 {
		return ""
	}

	// Handle empty string
	if len(s) == 0 {
		return s
	}

	// If string fits, return as is
	if utf8.RuneCountInString(s) <= maxWidth {
		return s
	}

	// Truncate and add ellipsis
	runes := []rune(s)
	if maxWidth <= 3 {
		return string(runes[:maxWidth])
	}

	return string(runes[:maxWidth-3]) + "..."
}

// Helper function to pad string to exact width
func padToWidth(s string, width int) string {
	runeCount := utf8.RuneCountInString(s)
	if runeCount >= width {
		return truncateString(s, width)
	}

	padding := width - runeCount
	return s + strings.Repeat(" ", padding)
}

// Add this struct definition near your other types
type tableColumn struct {
	Header string
	Width  int
}

// Helper function to calculate visual width of string (handles emojis/unicode)
func visualWidth(s string) int {
	// This is a simplified approach - for production you might want
	// to use a proper Unicode width calculation library
	width := 0
	for _, r := range s {
		if r < 127 {
			width++
		} else {
			// Most emojis and Unicode chars take 2 spaces visually
			width += 2
		}
	}
	return width
}

// Helper function to truncate string considering visual width
func truncateToVisualWidth(s string, maxWidth int) string {
	if maxWidth <= 0 {
		return ""
	}

	if visualWidth(s) <= maxWidth {
		return s
	}

	// Truncate character by character until we fit
	runes := []rune(s)
	result := ""
	currentWidth := 0

	for _, r := range runes {
		charWidth := 1
		if r >= 127 {
			charWidth = 2
		}

		if currentWidth+charWidth+3 > maxWidth { // +3 for "..."
			result += "..."
			break
		}

		result += string(r)
		currentWidth += charWidth
	}

	return result
}

// Helper function to pad string to exact visual width
func padToVisualWidth(s string, width int) string {
	currentWidth := visualWidth(s)
	if currentWidth >= width {
		return truncateToVisualWidth(s, width)
	}

	padding := width - currentWidth
	return s + strings.Repeat(" ", padding)
}
