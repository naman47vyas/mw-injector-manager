package main

import (
	"fmt"
	"log"
	"os"

	"github.com/naman47vyas/mw-injector-manager/internal/tui"
)

func main() {
	// Check for compact mode or help flags
	if len(os.Args) > 1 {
		switch os.Args[1] {
		// case "--compact":
		// 	tui.CompactBanner()
		// 	return
		case "--help", "-h":
			printHelp()
			return
		case "--version", "-v":
			fmt.Printf("MW Injector v%s (Build %s)\n", tui.Version, tui.Build)
			return
		}
	}

	// Run the Bubble Tea TUI
	if err := tui.RunApp(); err != nil {
		log.Fatalf("Error running application: %v", err)
	}
}

func printHelp() {
	fmt.Printf(`MW Injector - OpenTelemetry Injection System

Usage:
  mw-injector                 Start interactive TUI
  mw-injector --compact       Show compact banner only
  mw-injector --help          Show this help message
  mw-injector --version       Show version information

Interactive Mode Commands:
  - List and manage Java services
  - Configure Middleware.io instrumentation  
  - Monitor system health and telemetry
  - Export and import configurations

For more information, visit: https://docs.middleware.io/injector
`)
}

// =============================================================================
// Alternative: Advanced features example
// =============================================================================

// internal/tui/advanced.go - Additional advanced components
/*
package tui

import (
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// ServiceConfigModel for service configuration screen
type ServiceConfigModel struct {
	inputs      []textinput.Model
	focusIndex  int
	spinner     spinner.Model
	progress    progress.Model
	loading     bool
}

// Add loading states, progress bars, text inputs for configuration
// Add real-time updates using tea.Tick commands
// Add confirmation dialogs using bubbles/dialog
*/

