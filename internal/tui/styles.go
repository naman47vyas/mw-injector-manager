package tui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Color palette
	PrimaryColor   = lipgloss.Color("36")  // Cyan
	SecondaryColor = lipgloss.Color("135") // Purple
	AccentColor    = lipgloss.Color("46")  // Green
	WarningColor   = lipgloss.Color("226") // Yellow
	ErrorColor     = lipgloss.Color("196") // Red
	TextColor      = lipgloss.Color("255") // White
	DimColor       = lipgloss.Color("240") // Gray

	// Banner style
	BannerStyle = lipgloss.NewStyle().
			Foreground(PrimaryColor).
			Bold(true).
			Align(lipgloss.Center)

	// Title styles
	TitleStyle = lipgloss.NewStyle().
			Foreground(SecondaryColor).
			Bold(true).
			Align(lipgloss.Center).
			Padding(1, 0)

	// Status box style
	StatusBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(PrimaryColor).
			Padding(1, 2).
			Margin(1, 0)

	// Menu item styles
	MenuItemStyle = lipgloss.NewStyle().
			Padding(0, 2)

	SelectedMenuItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("0")).
				Background(AccentColor).
				Bold(true).
				Padding(0, 2)

		// Info styles
	InfoStyle = lipgloss.NewStyle().
			Foreground(DimColor).
			Italic(true).
			Align(lipgloss.Center)
	// Help style
	HelpStyle = lipgloss.NewStyle().
			Foreground(DimColor).
			Align(lipgloss.Center).
			Margin(1, 0)
)
