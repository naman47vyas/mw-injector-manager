package tui

import (
	"context"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/naman47vyas/mw-injector-manager/pkg/discovery"
)

// AppState represents the current state of the application
type AppState int

const (
	StateMain AppState = iota
	StateServiceList
	StateServiceConfig
	StateHealthCheck
	StateHelp
	// --- New States Added ---
	StateProcessOptions    // State for showing options for a selected process
	StateProcessDetailView // State for showing the detailed view of a process
)

// MenuItem represents a main menu item
type MenuItem struct {
	title       string
	description string
	action      AppState
}

func (m MenuItem) Title() string       { return m.title }
func (m MenuItem) Description() string { return m.description }
func (m MenuItem) FilterValue() string { return m.title }

// --- New: Process Options Menu Item ---
type ProcessActionItem struct {
	title       string
	description string
}

func (p ProcessActionItem) Title() string       { return p.title }
func (p ProcessActionItem) Description() string { return p.description }
func (p ProcessActionItem) FilterValue() string { return p.title }

// SystemStatus represents the current system status
type SystemStatus struct {
	JavaServices    int
	Instrumented    int
	MiddlewareCount int
	OtherAgentCount int
	SystemHealth    string
	TelemetryStatus string
	LastUpdate      time.Time
	Processes       []discovery.JavaProcess
}

// Model represents the main application model
type Model struct {
	state                AppState
	width                int
	height               int
	list                 list.Model // Main menu list
	status               SystemStatus
	showBanner           bool
	showStatus           bool
	animationStep        int
	discoverer           discovery.Discoverer
	ctx                  context.Context
	cancel               context.CancelFunc
	selectedProcessIndex int // Track which process is selected in the list
	// --- New: List for process options ---
	processOptionsList list.Model
}

// NewModel creates a new model instance
func NewModel() Model {
	// Create context for discovery operations
	ctx, cancel := context.WithCancel(context.Background())

	// Create discoverer
	discoverer := discovery.NewDiscoverer(ctx)

	// --- Main Menu Items ---
	items := []list.Item{
		MenuItem{title: "📋 List Services", description: "View all Java services and their instrumentation status", action: StateServiceList},
		// ... (other menu items are the same)
		MenuItem{title: "⚙️ Configure Service", description: "Configure MW environment variables for a specific service", action: StateServiceConfig},
		MenuItem{title: "🔧 Enable MW Agent", description: "Enable Middleware.io instrumentation for selected services", action: StateServiceConfig},
		MenuItem{title: "❌ Disable Instrumentation", description: "Remove instrumentation from selected services", action: StateServiceConfig},
		MenuItem{title: "📊 View Telemetry", description: "Check telemetry data flow and collector connectivity", action: StateHealthCheck},
		MenuItem{title: "🏥 Health Check", description: "Perform system health check and validation", action: StateHealthCheck},
		MenuItem{title: "📤 Export Config", description: "Export current configurations to file", action: StateServiceConfig},
		MenuItem{title: "❓ Help", description: "View documentation and troubleshooting guide", action: StateHelp},
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = SelectedMenuItemStyle
	delegate.Styles.SelectedDesc = SelectedMenuItemStyle.Copy().Foreground(lipgloss.Color("0"))

	l := list.New(items, delegate, 0, 0)
	l.Title = "MW Injector Commands"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = TitleStyle

	// --- New: Process Options List ---
	processOptionsItems := []list.Item{
		ProcessActionItem{
			title:       "📄 View Details",
			description: "Show comprehensive details for this process",
		},
		ProcessActionItem{
			title:       "🛠️ Configure Instrumentation",
			description: "Set environment variables to enable the MW agent",
		},
		ProcessActionItem{
			title:       " Raw Data",
			description: "Display all discovered data for debugging",
		},
	}

	processDelegate := list.NewDefaultDelegate()
	processDelegate.Styles.SelectedTitle = SelectedMenuItemStyle
	processDelegate.Styles.SelectedDesc = SelectedMenuItemStyle.Copy().Foreground(lipgloss.Color("0"))
	processOptionsList := list.New(processOptionsItems, processDelegate, 0, 0)
	processOptionsList.SetShowStatusBar(false)
	processOptionsList.SetFilteringEnabled(false)
	processOptionsList.Styles.Title = TitleStyle
	processOptionsList.Title = "Process Options" // Will be updated dynamically

	status := SystemStatus{
		JavaServices:    0,
		Instrumented:    0,
		MiddlewareCount: 0,
		OtherAgentCount: 0,
		SystemHealth:    "🔄 Ready to scan",
		TelemetryStatus: "📡 Connected",
		LastUpdate:      time.Now(),
		Processes:       []discovery.JavaProcess{},
	}

	return Model{
		state:                StateMain,
		list:                 l,
		status:               status,
		showBanner:           true,
		showStatus:           true,
		discoverer:           discoverer,
		ctx:                  ctx,
		cancel:               cancel,
		selectedProcessIndex: -1,
		processOptionsList:   processOptionsList, // Add new list to model
	}
}

// Init implements tea.Model
func (m Model) Init() tea.Cmd {
	return tea.EnterAltScreen
}
