package tui

import "fmt"

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
