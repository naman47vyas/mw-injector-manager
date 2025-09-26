package tui

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	// types "github.com/naman47vyas/mw-injector-manager/internal/tui/types"

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
)

// MenuItem represents a menu item
type MenuItem struct {
	title       string
	description string
	action      AppState
}

func (m MenuItem) Title() string       { return m.title }
func (m MenuItem) Description() string { return m.description }
func (m MenuItem) FilterValue() string { return m.title }

// SystemStatus represents the current system status
//
//	type SystemStatus struct {
//		JavaServices    int
//		Instrumented    int
//		ActiveConfigs   int
//		PendingConfigs  int
//		SystemHealth    string
//		TelemetryStatus string
//		LastUpdate      time.Time
//	}
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
//
//	type Model struct {
//		state         AppState
//		width         int
//		height        int
//		list          list.Model
//		status        SystemStatus
//		showBanner    bool
//		showStatus    bool
//		animationStep int
//		services      []types.Service
//	}

type Model struct {
	state         AppState
	width         int
	height        int
	list          list.Model
	status        SystemStatus
	showBanner    bool
	showStatus    bool
	animationStep int
	discoverer    discovery.Discoverer
	ctx           context.Context
	cancel        context.CancelFunc
	isLoading     bool
	lastError     error
}

// NewModel creates a new model instance
// func NewModel() Model {
// 	// Create menu items
// 	items := []list.Item{
// 		MenuItem{
// 			title:       "📋 List Services",
// 			description: "View all Java services and their instrumentation status",
// 			action:      StateServiceList,
// 		},
// 		MenuItem{
// 			title:       "⚙️  Configure Service",
// 			description: "Configure MW environment variables for a specific service",
// 			action:      StateServiceConfig,
// 		},
// 		MenuItem{
// 			title:       "🔧 Enable MW Agent",
// 			description: "Enable Middleware.io instrumentation for selected services",
// 			action:      StateServiceConfig,
// 		},
// 		MenuItem{
// 			title:       "❌ Disable Instrumentation",
// 			description: "Remove instrumentation from selected services",
// 			action:      StateServiceConfig,
// 		},
// 		MenuItem{
// 			title:       "📊 View Telemetry",
// 			description: "Check telemetry data flow and collector connectivity",
// 			action:      StateHealthCheck,
// 		},
// 		MenuItem{
// 			title:       "🏥 Health Check",
// 			description: "Perform system health check and validation",
// 			action:      StateHealthCheck,
// 		},
// 		MenuItem{
// 			title:       "📤 Export Config",
// 			description: "Export current configurations to file",
// 			action:      StateServiceConfig,
// 		},
// 		MenuItem{
// 			title:       "❓ Help",
// 			description: "View documentation and troubleshooting guide",
// 			action:      StateHelp,
// 		},
// 	}

// 	// Create list
// 	delegate := list.NewDefaultDelegate()
// 	delegate.Styles.SelectedTitle = SelectedMenuItemStyle
// 	delegate.Styles.SelectedDesc = SelectedMenuItemStyle.Copy().Foreground(lipgloss.Color("0"))

// 	l := list.New(items, delegate, 0, 0)
// 	l.Title = "MW Injector Commands"
// 	l.SetShowStatusBar(false)
// 	l.SetFilteringEnabled(false)
// 	l.Styles.Title = TitleStyle

// 	// Initialize status
// 	status := SystemStatus{
// 		JavaServices:    3,
// 		Instrumented:    2,
// 		ActiveConfigs:   2,
// 		PendingConfigs:  1,
// 		SystemHealth:    "✅ Healthy",
// 		TelemetryStatus: "📈 Flowing",
// 		LastUpdate:      time.Now(),
// 	}

// 	services := []types.Service{
// 		{PID: 1234, Name: "microser-a.jar", Status: "MW", Config: "user-auth", LastSeen: "2 min ago"},
// 		{PID: 5678, Name: "microservice-b.jar", Status: "None", Config: "-", LastSeen: "5 min ago"},
// 		{PID: 9012, Name: "legacy-app.jar", Status: "OTel", Config: "default", LastSeen: "1 min ago"},
// 		{PID: 3456, Name: "data-processor.jar", Status: "Pending", Config: "new-config", LastSeen: "10 sec ago"},
// 	}

// 	return Model{
// 		state:      StateMain,
// 		list:       l,
// 		status:     status,
// 		showBanner: true,
// 		showStatus: true,
// 		services:   services,
// 	}
// }

func NewModel() Model {
	// Create context for discovery operations
	ctx, cancel := context.WithCancel(context.Background())

	// Create discoverer with optimized settings for TUI
	opts := discovery.DefaultDiscoveryOptions()
	opts.MaxConcurrency = 5         // Lower concurrency for UI responsiveness
	opts.Timeout = 10 * time.Second // Shorter timeout for better UX
	opts.IncludeMetrics = true
	discoverer := discovery.NewDiscovererWithOptions(ctx, opts)

	// Create menu items
	items := []list.Item{
		MenuItem{
			title:       "📋 List Services",
			description: "View all Java services and their instrumentation status",
			action:      StateServiceList,
		},
		MenuItem{
			title:       "⚙️ Configure Service",
			description: "Configure MW environment variables for a specific service",
			action:      StateServiceConfig,
		},
		MenuItem{
			title:       "🔧 Enable MW Agent",
			description: "Enable Middleware.io instrumentation for selected services",
			action:      StateServiceConfig,
		},
		MenuItem{
			title:       "❌ Disable Instrumentation",
			description: "Remove instrumentation from selected services",
			action:      StateServiceConfig,
		},
		MenuItem{
			title:       "📊 View Telemetry",
			description: "Check telemetry data flow and collector connectivity",
			action:      StateHealthCheck,
		},
		MenuItem{
			title:       "🏥 Health Check",
			description: "Perform system health check and validation",
			action:      StateHealthCheck,
		},
		MenuItem{
			title:       "📤 Export Config",
			description: "Export current configurations to file",
			action:      StateServiceConfig,
		},
		MenuItem{
			title:       "❓ Help",
			description: "View documentation and troubleshooting guide",
			action:      StateHelp,
		},
	}

	// Create list
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = SelectedMenuItemStyle
	delegate.Styles.SelectedDesc = SelectedMenuItemStyle.Copy().Foreground(lipgloss.Color("0"))

	l := list.New(items, delegate, 0, 0)
	l.Title = "MW Injector Commands"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = TitleStyle

	// Initialize with empty status
	// status := SystemStatus{
	// 	JavaServices:    0,
	// 	Instrumented:    0,
	// 	MiddlewareCount: 0,
	// 	OtherAgentCount: 0,
	// 	SystemHealth:    "🔄 Loading...",
	// 	TelemetryStatus: "🔄 Checking...",
	// 	LastUpdate:      time.Now(),
	// 	Processes:       []discovery.JavaProcess{},
	// }
	status := SystemStatus{
		JavaServices:    0,
		Instrumented:    0,
		MiddlewareCount: 0,
		OtherAgentCount: 0,
		SystemHealth:    "🔄 Loading...",
		TelemetryStatus: "🔄 Checking...",
		LastUpdate:      time.Now(),
		Processes:       []discovery.JavaProcess{},
	}

	return Model{
		state:      StateMain,
		list:       l,
		status:     status,
		showBanner: true,
		showStatus: true,
		discoverer: discoverer,
		ctx:        ctx,
		cancel:     cancel,
		isLoading:  false,
	}
}

// Init implements tea.Model
func (m Model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		verticalMargin := 0

		if m.showBanner {
			banner := BannerStyle.Render(GetBannerArt())
			subtitle := TitleStyle.Render(GetSubtitle())
			version := InfoStyle.Render(GetVersionInfo())
			verticalMargin += lipgloss.Height(banner) + lipgloss.Height(subtitle) + lipgloss.Height(version)
		}

		if m.showStatus {
			statusContent := m.renderStatus()
			statusBox := StatusBoxStyle.Render(statusContent)
			verticalMargin += lipgloss.Height(statusBox)
		}

		help := HelpStyle.Render("...")
		verticalMargin += lipgloss.Height(help)

		m.list.SetHeight(m.height - verticalMargin)
		m.list.SetWidth(m.width - 4)
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+c", "q"))):
			return m, tea.Quit

		case key.Matches(msg, key.NewBinding(key.WithKeys("b"))):
			m.showBanner = !m.showBanner
			return m, func() tea.Msg {
				return tea.WindowSizeMsg{Width: m.width, Height: m.height}
			}

		case key.Matches(msg, key.NewBinding(key.WithKeys("s"))):
			m.showStatus = !m.showStatus
			return m, func() tea.Msg {
				return tea.WindowSizeMsg{Width: m.width, Height: m.height}
			}

		case key.Matches(msg, key.NewBinding(key.WithKeys("r"))):
			m.status.LastUpdate = time.Now()
			return m, nil

		case key.Matches(msg, key.NewBinding(key.WithKeys("enter", " "))):
			if selectedItem, ok := m.list.SelectedItem().(MenuItem); ok {
				m.state = selectedItem.action
				return m, nil
			}
		}

		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}

	return m, nil
}

// View implements tea.Model
func (m Model) View() string {
	var sections []string

	// Banner section
	if m.showBanner {
		banner := BannerStyle.Render(GetBannerArt())
		subtitle := TitleStyle.Render(GetSubtitle())
		version := InfoStyle.Render(GetVersionInfo())
		sections = append(sections, banner, subtitle, version)
	}

	// Status section
	if m.showStatus {
		statusContent := m.renderStatus()
		sections = append(sections, StatusBoxStyle.Render(statusContent))
	}

	// Main content based on state
	switch m.state {
	case StateMain:
		sections = append(sections, m.list.View())

	case StateServiceList:
		sections = append(sections, m.renderServiceList())

	case StateHealthCheck:
		sections = append(sections, m.renderHealthCheck())

	case StateHelp:
		sections = append(sections, m.renderHelp())

	default:
		sections = append(sections, "Feature coming soon...")
	}

	// Help footer
	help := HelpStyle.Render("Press 'q' to quit • 'b' to toggle banner • 'r' to refresh • '?' for help")
	sections = append(sections, help)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// func (m Model) renderStatus() string {
// 	hostname, _ := os.Hostname()

// 	// Helper to create styled key-value pairs
// 	kv := func(key, value string) string {
// 		return fmt.Sprintf("%s %s", StatusKeyStyle.Render(key), StatusValueStyle.Render(value))
// 	}

// 	statusLines := []string{
// 		StatusHeaderStyle.Render("System Status"),
// 		kv("🖥️  Host:", hostname),
// 		kv("📅 Last Update:", m.status.LastUpdate.Format("15:04:05")),
// 		"",
// 		kv("☕ Java Services:", fmt.Sprintf("%d running, %d instrumented", m.status.JavaServices, m.status.Instrumented)),
// 		kv("⚙️  MW Configs:", fmt.Sprintf("%d active, %d pending", m.status.ActiveConfigs, m.status.PendingConfigs)),
// 		kv("🏥 Health:", m.status.SystemHealth),
// 		kv("📡 Telemetry:", m.status.TelemetryStatus),
// 	}

//		return strings.Join(statusLines, "\n")
//	}
func (m Model) renderStatus() string {
	hostname, _ := os.Hostname()

	statusLines := []string{
		fmt.Sprintf("🖥️ Host: %s", hostname),
		fmt.Sprintf("📅 Last Update: %s", m.status.LastUpdate.Format("15:04:05")),
		"",
		fmt.Sprintf("☕ Java Services: %d running, %d instrumented",
			m.status.JavaServices, m.status.Instrumented),
		fmt.Sprintf("⚙️ MW Agents: %d active, %d other agents",
			m.status.MiddlewareCount, m.status.OtherAgentCount),
		fmt.Sprintf("🏥 Health: %s", m.status.SystemHealth),
		fmt.Sprintf("📡 Telemetry: %s", m.status.TelemetryStatus),
		"",
		"🔧 System: LD_PRELOAD Injection Active",
		"🎯 Target: Middleware.io",
	}

	return strings.Join(statusLines, "\n")
}

// func (m Model) renderServiceList() string {
// 	columns := []tableColumn{
// 		{Header: "PID", Width: 7},
// 		{Header: "Service Name", Width: 35}, // Increased width for long service names
// 		{Header: "Status", Width: 12},       // Adjusted for emoji width
// 		{Header: "MW Config", Width: 16},
// 		{Header: "Last Seen", Width: 12},
// 	}

// 	// --- Build header ---
// 	headerCells := make([]string, len(columns))
// 	for i, col := range columns {
// 		headerText := padToWidth(col.Header, col.Width)
// 		headerCells[i] = lipgloss.NewStyle().
// 			Foreground(SecondaryColor).
// 			Bold(true).
// 			Render(headerText)
// 	}
// 	header := strings.Join(headerCells, TableSeparatorStyle.Render(" │ "))

// 	// --- Build separator line ---
// 	separatorCells := make([]string, len(columns))
// 	for i, col := range columns {
// 		separatorCells[i] = strings.Repeat("─", col.Width)
// 	}
// 	separator := TableSeparatorStyle.Render("─" + strings.Join(separatorCells, "─┼─") + "─")

// 	// --- Build rows ---
// 	var rows []string
// 	for _, s := range m.services {
// 		serviceData := []string{
// 			fmt.Sprint(s.PID),
// 			s.Name,
// 			s.Status,
// 			s.Config,
// 			s.LastSeen,
// 		}

// 		// Choose row style based on status
// 		var rowColor lipgloss.Color
// 		if strings.Contains(s.Status, "✅") {
// 			rowColor = AccentColor
// 		} else if strings.Contains(s.Status, "❌") {
// 			rowColor = ErrorColor
// 		} else if strings.Contains(s.Status, "⚠️") {
// 			rowColor = WarningColor
// 		} else {
// 			rowColor = TextColor
// 		}

// 		rowCells := make([]string, len(columns))
// 		for i, data := range serviceData {
// 			// Properly pad each cell to its column width
// 			cellText := padToWidth(data, columns[i].Width)
// 			rowCells[i] = lipgloss.NewStyle().
// 				Foreground(rowColor).
// 				Render(cellText)
// 		}

// 		row := strings.Join(rowCells, TableSeparatorStyle.Render(" │ "))
// 		rows = append(rows, row)
// 	}

// 	// Combine all parts
// 	var tableLines []string
// 	tableLines = append(tableLines, header)
// 	tableLines = append(tableLines, separator)
// 	tableLines = append(tableLines, strings.Join(rows, "\n"))

// 	tableContent := strings.Join(tableLines, "\n")

//		// Render the final table
//		return lipgloss.JoinVertical(lipgloss.Left,
//			TitleStyle.Render("📋 Java Services"),
//			"", // Add some spacing
//			TableBoxStyle.Render(tableContent),
//			"",
//			HelpStyle.Render("Press 'enter' to configure service • 'backspace' to go back"),
//		)
//	}
//
// renderServiceList renders the service list view with real process data
func (m Model) renderServiceList() string {
	content := []string{
		TitleStyle.Render("📋 Java Services"),
		"",
	}

	if m.isLoading {
		content = append(content, "🔄 Loading Java processes...")
		return strings.Join(content, "\n")
	}

	if len(m.status.Processes) == 0 {
		content = append(content, "No Java processes found.")
		if m.lastError != nil {
			content = append(content, fmt.Sprintf("Error: %v", m.lastError))
		}
		content = append(content, "", HelpStyle.Render("Press 'r' to refresh • 'backspace' to go back"))
		return strings.Join(content, "\n")
	}

	// Table header
	content = append(content,
		"┌─────────┬──────────────────────┬─────────┬──────────────────┬──────────────┐",
		"│ PID     │ Service Name         │ Status  │ JAR File         │ Memory       │",
		"├─────────┼──────────────────────┼─────────┼──────────────────┼──────────────┤",
	)

	// Table rows with real data
	for i, proc := range m.status.Processes {
		// Limit display to prevent overflow
		if i >= 10 {
			content = append(content, fmt.Sprintf("│ ... and %d more processes", len(m.status.Processes)-10))
			break
		}

		// Format each field with proper width
		pidStr := fmt.Sprintf("%-7d", proc.ProcessPID)
		serviceStr := fmt.Sprintf("%-20s", truncateString(proc.ServiceName, 20))
		statusStr := fmt.Sprintf("%-7s", proc.FormatAgentStatus())
		jarStr := fmt.Sprintf("%-16s", truncateString(proc.JarFile, 16))
		memoryStr := fmt.Sprintf("%-12s", fmt.Sprintf("%.1f%%", proc.MemoryPercent))

		content = append(content, fmt.Sprintf("│ %s │ %s │ %s │ %s │ %s │",
			pidStr, serviceStr, statusStr, jarStr, memoryStr))
	}

	// Table footer
	content = append(content,
		"└─────────┴──────────────────────┴─────────┴──────────────────┴──────────────┘",
		"",
		HelpStyle.Render("Press 'enter' to configure service • 'backspace' to go back • 'r' to refresh"),
	)

	return strings.Join(content, "\n")
}

// renderHealthCheck renders the health check view
func (m Model) renderHealthCheck() string {
	content := []string{
		TitleStyle.Render("🏥 System Health Check"),
		"",
		"✅ LD_PRELOAD injection: Active",
		"✅ Shared library: /usr/lib/middleware/libmwinjector.so loaded",
		"✅ Java agent: /usr/lib/middleware/middleware-javaagent-1.7.0.jar found",
		"✅ Configuration files: Valid",
		"✅ Network connectivity: Middleware.io reachable",
		"⚠️  Warning: 1 service without MW configuration",
		"",
		"📊 Telemetry Status:",
		"  • Traces: 1,234 spans/min",
		"  • Metrics: 56 metrics/min",
		"  • Logs: 789 logs/min",
		"  • Last export: 30 seconds ago",
		"",
		HelpStyle.Render("Press 'backspace' to go back"),
	}

	return strings.Join(content, "\n")
}

// renderHelp renders the help view
func (m Model) renderHelp() string {
	content := []string{
		TitleStyle.Render("❓ Help & Documentation"),
		"",
		"🔧 Key Features:",
		"  • Automatic Java service discovery via /proc filesystem",
		"  • Per-service Middleware.io configuration",
		"  • LD_PRELOAD shared library injection",
		"  • Real-time health monitoring",
		"",
		"⌨️  Keyboard Shortcuts:",
		"  • q, Ctrl+C: Quit application",
		"  • ↑/↓: Navigate menu items",
		"  • Enter/Space: Select menu item",
		"  • b: Toggle banner display",
		"  • r: Refresh system status",
		"  • Backspace: Go back to previous screen",
		"",
		"🔗 Resources:",
		"  • Documentation: https://docs.middleware.io/injector",
		"  • GitHub: https://github.com/middleware-labs/mw-injector",
		"  • Support: support@middleware.io",
		"",
		HelpStyle.Render("Press 'backspace' to go back"),
	}

	return strings.Join(content, "\n")
}
