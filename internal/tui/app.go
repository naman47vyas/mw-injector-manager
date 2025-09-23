package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// RunApp starts the Bubble Tea application
func RunApp() error {
	p := tea.NewProgram(
		NewModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	_, err := p.Run()
	return err
}
