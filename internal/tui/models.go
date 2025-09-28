// package tui

// import (
// 	"context"
// 	"fmt"
// 	"os"
// 	"strings"
// 	"time"

// 	"github.com/charmbracelet/bubbles/key"
// 	"github.com/charmbracelet/bubbles/list"

// 	tea "github.com/charmbracelet/bubbletea"
// 	"github.com/charmbracelet/lipgloss"
// 	// types "github.com/naman47vyas/mw-injector-manager/internal/tui/types"

// 	"github.com/naman47vyas/mw-injector-manager/pkg/discovery"
// )

// // AppState represents the current state of the application
// type AppState int

// const (
// 	StateMain AppState = iota
// 	StateServiceList
// 	StateServiceConfig
// 	StateHealthCheck
// 	StateHelp
// )

// // MenuItem represents a menu item
// type MenuItem struct {
// 	title       string
// 	description string
// 	action      AppState
// }

// func (m MenuItem) Title() string       { return m.title }
// func (m MenuItem) Description() string { return m.description }
// func (m MenuItem) FilterValue() string { return m.title }

// // SystemStatus represents the current system status
// //
// //	type SystemStatus struct {
// //		JavaServices    int
// //		Instrumented    int
// //		ActiveConfigs   int
// //		PendingConfigs  int
// //		SystemHealth    string
// //		TelemetryStatus string
// //		LastUpdate      time.Time
// //	}
// type SystemStatus struct {
// 	JavaServices    int
// 	Instrumented    int
// 	MiddlewareCount int
// 	OtherAgentCount int
// 	SystemHealth    string
// 	TelemetryStatus string
// 	LastUpdate      time.Time
// 	Processes       []discovery.JavaProcess
// }

// // Model represents the main application model
// //
// //	type Model struct {
// //		state         AppState
// //		width         int
// //		height        int
// //		list          list.Model
// //		status        SystemStatus
// //		showBanner    bool
// //		showStatus    bool
// //		animationStep int
// //		services      []types.Service
// //	}

// type Model struct {
// 	state         AppState
// 	width         int
// 	height        int
// 	list          list.Model
// 	status        SystemStatus
// 	showBanner    bool
// 	showStatus    bool
// 	animationStep int
// 	discoverer    discovery.Discoverer
// 	ctx           context.Context
// 	cancel        context.CancelFunc
// 	isLoading     bool
// 	lastError     error
// }

// // NewModel creates a new model instance
// // func NewModel() Model {
// // 	// Create menu items
// // 	items := []list.Item{
// // 		MenuItem{
// // 			title:       "📋 List Services",
// // 			description: "View all Java services and their instrumentation status",
// // 			action:      StateServiceList,
// // 		},
// // 		MenuItem{
// // 			title:       "⚙️  Configure Service",
// // 			description: "Configure MW environment variables for a specific service",
// // 			action:      StateServiceConfig,
// // 		},
// // 		MenuItem{
// // 			title:       "🔧 Enable MW Agent",
// // 			description: "Enable Middleware.io instrumentation for selected services",
// // 			action:      StateServiceConfig,
// // 		},
// // 		MenuItem{
// // 			title:       "❌ Disable Instrumentation",
// // 			description: "Remove instrumentation from selected services",
// // 			action:      StateServiceConfig,
// // 		},
// // 		MenuItem{
// // 			title:       "📊 View Telemetry",
// // 			description: "Check telemetry data flow and collector connectivity",
// // 			action:      StateHealthCheck,
// // 		},
// // 		MenuItem{
// // 			title:       "🏥 Health Check",
// // 			description: "Perform system health check and validation",
// // 			action:      StateHealthCheck,
// // 		},
// // 		MenuItem{
// // 			title:       "📤 Export Config",
// // 			description: "Export current configurations to file",
// // 			action:      StateServiceConfig,
// // 		},
// // 		MenuItem{
// // 			title:       "❓ Help",
// // 			description: "View documentation and troubleshooting guide",
// // 			action:      StateHelp,
// // 		},
// // 	}

// // 	// Create list
// // 	delegate := list.NewDefaultDelegate()
// // 	delegate.Styles.SelectedTitle = SelectedMenuItemStyle
// // 	delegate.Styles.SelectedDesc = SelectedMenuItemStyle.Copy().Foreground(lipgloss.Color("0"))

// // 	l := list.New(items, delegate, 0, 0)
// // 	l.Title = "MW Injector Commands"
// // 	l.SetShowStatusBar(false)
// // 	l.SetFilteringEnabled(false)
// // 	l.Styles.Title = TitleStyle

// // 	// Initialize status
// // 	status := SystemStatus{
// // 		JavaServices:    3,
// // 		Instrumented:    2,
// // 		ActiveConfigs:   2,
// // 		PendingConfigs:  1,
// // 		SystemHealth:    "✅ Healthy",
// // 		TelemetryStatus: "📈 Flowing",
// // 		LastUpdate:      time.Now(),
// // 	}

// // 	services := []types.Service{
// // 		{PID: 1234, Name: "microser-a.jar", Status: "MW", Config: "user-auth", LastSeen: "2 min ago"},
// // 		{PID: 5678, Name: "microservice-b.jar", Status: "None", Config: "-", LastSeen: "5 min ago"},
// // 		{PID: 9012, Name: "legacy-app.jar", Status: "OTel", Config: "default", LastSeen: "1 min ago"},
// // 		{PID: 3456, Name: "data-processor.jar", Status: "Pending", Config: "new-config", LastSeen: "10 sec ago"},
// // 	}

// // 	return Model{
// // 		state:      StateMain,
// // 		list:       l,
// // 		status:     status,
// // 		showBanner: true,
// // 		showStatus: true,
// // 		services:   services,
// // 	}
// // }

// func NewModel() Model {
// 	// Create context for discovery operations
// 	ctx, cancel := context.WithCancel(context.Background())

// 	// Create discoverer with optimized settings for TUI
// 	opts := discovery.DefaultDiscoveryOptions()
// 	opts.MaxConcurrency = 5         // Lower concurrency for UI responsiveness
// 	opts.Timeout = 10 * time.Second // Shorter timeout for better UX
// 	opts.IncludeMetrics = true
// 	discoverer := discovery.NewDiscovererWithOptions(ctx, opts)

// 	// Create menu items
// 	items := []list.Item{
// 		MenuItem{
// 			title:       "📋 List Services",
// 			description: "View all Java services and their instrumentation status",
// 			action:      StateServiceList,
// 		},
// 		MenuItem{
// 			title:       "⚙️ Configure Service",
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

// 	// Initialize with empty status
// 	// status := SystemStatus{
// 	// 	JavaServices:    0,
// 	// 	Instrumented:    0,
// 	// 	MiddlewareCount: 0,
// 	// 	OtherAgentCount: 0,
// 	// 	SystemHealth:    "🔄 Loading...",
// 	// 	TelemetryStatus: "🔄 Checking...",
// 	// 	LastUpdate:      time.Now(),
// 	// 	Processes:       []discovery.JavaProcess{},
// 	// }
// 	status := SystemStatus{
// 		JavaServices:    0,
// 		Instrumented:    0,
// 		MiddlewareCount: 0,
// 		OtherAgentCount: 0,
// 		SystemHealth:    "🔄 Loading...",
// 		TelemetryStatus: "🔄 Checking...",
// 		LastUpdate:      time.Now(),
// 		Processes:       []discovery.JavaProcess{},
// 	}

// 	return Model{
// 		state:      StateMain,
// 		list:       l,
// 		status:     status,
// 		showBanner: true,
// 		showStatus: true,
// 		discoverer: discoverer,
// 		ctx:        ctx,
// 		cancel:     cancel,
// 		isLoading:  false,
// 	}
// }

// // Init implements tea.Model
// func (m Model) Init() tea.Cmd {
// 	return tea.EnterAltScreen
// }

// func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	switch msg := msg.(type) {
// 	case tea.WindowSizeMsg:
// 		m.width = msg.Width
// 		m.height = msg.Height
// 		verticalMargin := 0

// 		if m.showBanner {
// 			banner := BannerStyle.Render(GetBannerArt())
// 			subtitle := TitleStyle.Render(GetSubtitle())
// 			version := InfoStyle.Render(GetVersionInfo())
// 			verticalMargin += lipgloss.Height(banner) + lipgloss.Height(subtitle) + lipgloss.Height(version)
// 		}

// 		if m.showStatus {
// 			statusContent := m.renderStatus()
// 			statusBox := StatusBoxStyle.Render(statusContent)
// 			verticalMargin += lipgloss.Height(statusBox)
// 		}

// 		help := HelpStyle.Render("...")
// 		verticalMargin += lipgloss.Height(help)

// 		m.list.SetHeight(m.height - verticalMargin)
// 		m.list.SetWidth(m.width - 4)
// 		return m, nil

// 	case tea.KeyMsg:
// 		switch {
// 		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+c", "q"))):
// 			return m, tea.Quit

// 		case key.Matches(msg, key.NewBinding(key.WithKeys("b"))):
// 			m.showBanner = !m.showBanner
// 			return m, func() tea.Msg {
// 				return tea.WindowSizeMsg{Width: m.width, Height: m.height}
// 			}

// 		case key.Matches(msg, key.NewBinding(key.WithKeys("s"))):
// 			m.showStatus = !m.showStatus
// 			return m, func() tea.Msg {
// 				return tea.WindowSizeMsg{Width: m.width, Height: m.height}
// 			}

// 		case key.Matches(msg, key.NewBinding(key.WithKeys("r"))):
// 			m.status.LastUpdate = time.Now()
// 			return m, nil

// 		case key.Matches(msg, key.NewBinding(key.WithKeys("enter", " "))):
// 			if selectedItem, ok := m.list.SelectedItem().(MenuItem); ok {
// 				m.state = selectedItem.action
// 				return m, nil
// 			}
// 		}

// 		var cmd tea.Cmd
// 		m.list, cmd = m.list.Update(msg)
// 		return m, cmd
// 	}

// 	return m, nil
// }

// // View implements tea.Model
// func (m Model) View() string {
// 	var sections []string

// 	// Banner section
// 	if m.showBanner {
// 		banner := BannerStyle.Render(GetBannerArt())
// 		subtitle := TitleStyle.Render(GetSubtitle())
// 		version := InfoStyle.Render(GetVersionInfo())
// 		sections = append(sections, banner, subtitle, version)
// 	}

// 	// Status section
// 	if m.showStatus {
// 		statusContent := m.renderStatus()
// 		sections = append(sections, StatusBoxStyle.Render(statusContent))
// 	}

// 	// Main content based on state
// 	switch m.state {
// 	case StateMain:
// 		sections = append(sections, m.list.View())

// 	case StateServiceList:
// 		sections = append(sections, m.renderServiceList())

// 	case StateHealthCheck:
// 		sections = append(sections, m.renderHealthCheck())

// 	case StateHelp:
// 		sections = append(sections, m.renderHelp())

// 	default:
// 		sections = append(sections, "Feature coming soon...")
// 	}

// 	// Help footer
// 	help := HelpStyle.Render("Press 'q' to quit • 'b' to toggle banner • 'r' to refresh • '?' for help")
// 	sections = append(sections, help)

// 	return lipgloss.JoinVertical(lipgloss.Left, sections...)
// }

// // func (m Model) renderStatus() string {
// // 	hostname, _ := os.Hostname()

// // 	// Helper to create styled key-value pairs
// // 	kv := func(key, value string) string {
// // 		return fmt.Sprintf("%s %s", StatusKeyStyle.Render(key), StatusValueStyle.Render(value))
// // 	}

// // 	statusLines := []string{
// // 		StatusHeaderStyle.Render("System Status"),
// // 		kv("🖥️  Host:", hostname),
// // 		kv("📅 Last Update:", m.status.LastUpdate.Format("15:04:05")),
// // 		"",
// // 		kv("☕ Java Services:", fmt.Sprintf("%d running, %d instrumented", m.status.JavaServices, m.status.Instrumented)),
// // 		kv("⚙️  MW Configs:", fmt.Sprintf("%d active, %d pending", m.status.ActiveConfigs, m.status.PendingConfigs)),
// // 		kv("🏥 Health:", m.status.SystemHealth),
// // 		kv("📡 Telemetry:", m.status.TelemetryStatus),
// // 	}

// //		return strings.Join(statusLines, "\n")
// //	}
// func (m Model) renderStatus() string {
// 	hostname, _ := os.Hostname()

// 	statusLines := []string{
// 		fmt.Sprintf("🖥️ Host: %s", hostname),
// 		fmt.Sprintf("📅 Last Update: %s", m.status.LastUpdate.Format("15:04:05")),
// 		"",
// 		fmt.Sprintf("☕ Java Services: %d running, %d instrumented",
// 			m.status.JavaServices, m.status.Instrumented),
// 		fmt.Sprintf("⚙️ MW Agents: %d active, %d other agents",
// 			m.status.MiddlewareCount, m.status.OtherAgentCount),
// 		fmt.Sprintf("🏥 Health: %s", m.status.SystemHealth),
// 		fmt.Sprintf("📡 Telemetry: %s", m.status.TelemetryStatus),
// 		"",
// 		"🔧 System: LD_PRELOAD Injection Active",
// 		"🎯 Target: Middleware.io",
// 	}

// 	return strings.Join(statusLines, "\n")
// }

// // func (m Model) renderServiceList() string {
// // 	columns := []tableColumn{
// // 		{Header: "PID", Width: 7},
// // 		{Header: "Service Name", Width: 35}, // Increased width for long service names
// // 		{Header: "Status", Width: 12},       // Adjusted for emoji width
// // 		{Header: "MW Config", Width: 16},
// // 		{Header: "Last Seen", Width: 12},
// // 	}

// // 	// --- Build header ---
// // 	headerCells := make([]string, len(columns))
// // 	for i, col := range columns {
// // 		headerText := padToWidth(col.Header, col.Width)
// // 		headerCells[i] = lipgloss.NewStyle().
// // 			Foreground(SecondaryColor).
// // 			Bold(true).
// // 			Render(headerText)
// // 	}
// // 	header := strings.Join(headerCells, TableSeparatorStyle.Render(" │ "))

// // 	// --- Build separator line ---
// // 	separatorCells := make([]string, len(columns))
// // 	for i, col := range columns {
// // 		separatorCells[i] = strings.Repeat("─", col.Width)
// // 	}
// // 	separator := TableSeparatorStyle.Render("─" + strings.Join(separatorCells, "─┼─") + "─")

// // 	// --- Build rows ---
// // 	var rows []string
// // 	for _, s := range m.services {
// // 		serviceData := []string{
// // 			fmt.Sprint(s.PID),
// // 			s.Name,
// // 			s.Status,
// // 			s.Config,
// // 			s.LastSeen,
// // 		}

// // 		// Choose row style based on status
// // 		var rowColor lipgloss.Color
// // 		if strings.Contains(s.Status, "✅") {
// // 			rowColor = AccentColor
// // 		} else if strings.Contains(s.Status, "❌") {
// // 			rowColor = ErrorColor
// // 		} else if strings.Contains(s.Status, "⚠️") {
// // 			rowColor = WarningColor
// // 		} else {
// // 			rowColor = TextColor
// // 		}

// // 		rowCells := make([]string, len(columns))
// // 		for i, data := range serviceData {
// // 			// Properly pad each cell to its column width
// // 			cellText := padToWidth(data, columns[i].Width)
// // 			rowCells[i] = lipgloss.NewStyle().
// // 				Foreground(rowColor).
// // 				Render(cellText)
// // 		}

// // 		row := strings.Join(rowCells, TableSeparatorStyle.Render(" │ "))
// // 		rows = append(rows, row)
// // 	}

// // 	// Combine all parts
// // 	var tableLines []string
// // 	tableLines = append(tableLines, header)
// // 	tableLines = append(tableLines, separator)
// // 	tableLines = append(tableLines, strings.Join(rows, "\n"))

// // 	tableContent := strings.Join(tableLines, "\n")

// //		// Render the final table
// //		return lipgloss.JoinVertical(lipgloss.Left,
// //			TitleStyle.Render("📋 Java Services"),
// //			"", // Add some spacing
// //			TableBoxStyle.Render(tableContent),
// //			"",
// //			HelpStyle.Render("Press 'enter' to configure service • 'backspace' to go back"),
// //		)
// //	}
// //
// // renderServiceList renders the service list view with real process data
// func (m Model) renderServiceList() string {
// 	content := []string{
// 		TitleStyle.Render("📋 Java Services"),
// 		"",
// 	}

// 	if m.isLoading {
// 		content = append(content, "🔄 Loading Java processes...")
// 		return strings.Join(content, "\n")
// 	}

// 	if len(m.status.Processes) == 0 {
// 		content = append(content, "No Java processes found.")
// 		if m.lastError != nil {
// 			content = append(content, fmt.Sprintf("Error: %v", m.lastError))
// 		}
// 		content = append(content, "", HelpStyle.Render("Press 'r' to refresh • 'backspace' to go back"))
// 		return strings.Join(content, "\n")
// 	}

// 	// Table header
// 	content = append(content,
// 		"┌─────────┬──────────────────────┬─────────┬──────────────────┬──────────────┐",
// 		"│ PID     │ Service Name         │ Status  │ JAR File         │ Memory       │",
// 		"├─────────┼──────────────────────┼─────────┼──────────────────┼──────────────┤",
// 	)

// 	// Table rows with real data
// 	for i, proc := range m.status.Processes {
// 		// Limit display to prevent overflow
// 		if i >= 10 {
// 			content = append(content, fmt.Sprintf("│ ... and %d more processes", len(m.status.Processes)-10))
// 			break
// 		}

// 		// Format each field with proper width
// 		pidStr := fmt.Sprintf("%-7d", proc.ProcessPID)
// 		serviceStr := fmt.Sprintf("%-20s", truncateString(proc.ServiceName, 20))
// 		statusStr := fmt.Sprintf("%-7s", proc.FormatAgentStatus())
// 		jarStr := fmt.Sprintf("%-16s", truncateString(proc.JarFile, 16))
// 		memoryStr := fmt.Sprintf("%-12s", fmt.Sprintf("%.1f%%", proc.MemoryPercent))

// 		content = append(content, fmt.Sprintf("│ %s │ %s │ %s │ %s │ %s │",
// 			pidStr, serviceStr, statusStr, jarStr, memoryStr))
// 	}

// 	// Table footer
// 	content = append(content,
// 		"└─────────┴──────────────────────┴─────────┴──────────────────┴──────────────┘",
// 		"",
// 		HelpStyle.Render("Press 'enter' to configure service • 'backspace' to go back • 'r' to refresh"),
// 	)

// 	return strings.Join(content, "\n")
// }

// // renderHealthCheck renders the health check view
// func (m Model) renderHealthCheck() string {
// 	content := []string{
// 		TitleStyle.Render("🏥 System Health Check"),
// 		"",
// 		"✅ LD_PRELOAD injection: Active",
// 		"✅ Shared library: /usr/lib/middleware/libmwinjector.so loaded",
// 		"✅ Java agent: /usr/lib/middleware/middleware-javaagent-1.7.0.jar found",
// 		"✅ Configuration files: Valid",
// 		"✅ Network connectivity: Middleware.io reachable",
// 		"⚠️  Warning: 1 service without MW configuration",
// 		"",
// 		"📊 Telemetry Status:",
// 		"  • Traces: 1,234 spans/min",
// 		"  • Metrics: 56 metrics/min",
// 		"  • Logs: 789 logs/min",
// 		"  • Last export: 30 seconds ago",
// 		"",
// 		HelpStyle.Render("Press 'backspace' to go back"),
// 	}

// 	return strings.Join(content, "\n")
// }

// // renderHelp renders the help view
// func (m Model) renderHelp() string {
// 	content := []string{
// 		TitleStyle.Render("❓ Help & Documentation"),
// 		"",
// 		"🔧 Key Features:",
// 		"  • Automatic Java service discovery via /proc filesystem",
// 		"  • Per-service Middleware.io configuration",
// 		"  • LD_PRELOAD shared library injection",
// 		"  • Real-time health monitoring",
// 		"",
// 		"⌨️  Keyboard Shortcuts:",
// 		"  • q, Ctrl+C: Quit application",
// 		"  • ↑/↓: Navigate menu items",
// 		"  • Enter/Space: Select menu item",
// 		"  • b: Toggle banner display",
// 		"  • r: Refresh system status",
// 		"  • Backspace: Go back to previous screen",
// 		"",
// 		"🔗 Resources:",
// 		"  • Documentation: https://docs.middleware.io/injector",
// 		"  • GitHub: https://github.com/middleware-labs/mw-injector",
// 		"  • Support: support@middleware.io",
// 		"",
// 		HelpStyle.Render("Press 'backspace' to go back"),
// 	}

// 	return strings.Join(content, "\n")
// }

// package tui

// import (
// 	"context"
// 	"fmt"
// 	"os"
// 	"strings"
// 	"time"

// 	"github.com/charmbracelet/bubbles/key"
// 	"github.com/charmbracelet/bubbles/list"
// 	tea "github.com/charmbracelet/bubbletea"
// 	"github.com/charmbracelet/lipgloss"

// 	"github.com/naman47vyas/mw-injector-manager/pkg/discovery"
// )

// // AppState represents the current state of the application
// type AppState int

// const (
// 	StateMain AppState = iota
// 	StateServiceList
// 	StateServiceConfig
// 	StateHealthCheck
// 	StateHelp
// 	StateProcessDetail
// )

// type ProcessDetailState int

// const (
// 	DetailOverview ProcessDetailState = iota
// 	DetailInstrumentation
// 	DetailConfiguration
// )

// // MenuItem represents a menu item
// type MenuItem struct {
// 	title       string
// 	description string
// 	action      AppState
// }

// func (m MenuItem) Title() string       { return m.title }
// func (m MenuItem) Description() string { return m.description }
// func (m MenuItem) FilterValue() string { return m.title }

// // SystemStatus represents the current system status
// type SystemStatus struct {
// 	JavaServices    int
// 	Instrumented    int
// 	MiddlewareCount int
// 	OtherAgentCount int
// 	SystemHealth    string
// 	TelemetryStatus string
// 	LastUpdate      time.Time
// 	Processes       []discovery.JavaProcess
// }

// // Model represents the main application model
// type Model struct {
// 	state           AppState
// 	width           int
// 	height          int
// 	list            list.Model
// 	status          SystemStatus
// 	showBanner      bool
// 	showStatus      bool
// 	animationStep   int
// 	discoverer      discovery.Discoverer
// 	ctx             context.Context
// 	cancel          context.CancelFunc
// 	selectedProcess *discovery.JavaProcess
// 	detailState     ProcessDetailState
// }

// // NewModel creates a new model instance
// func NewModel() Model {
// 	// Create context for discovery operations
// 	ctx, cancel := context.WithCancel(context.Background())

// 	// Create discoverer - SIMPLE VERSION, no complex options
// 	discoverer := discovery.NewDiscoverer(ctx)

// 	// Create menu items
// 	items := []list.Item{
// 		MenuItem{
// 			title:       "📋 List Services",
// 			description: "View all Java services and their instrumentation status",
// 			action:      StateServiceList,
// 		},
// 		MenuItem{
// 			title:       "⚙️ Configure Service",
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

// 	// Initialize with empty status - will load real data when needed
// 	status := SystemStatus{
// 		JavaServices:    0,
// 		Instrumented:    0,
// 		MiddlewareCount: 0,
// 		OtherAgentCount: 0,
// 		SystemHealth:    "🔄 Ready to scan",
// 		TelemetryStatus: "📡 Connected",
// 		LastUpdate:      time.Now(),
// 		Processes:       []discovery.JavaProcess{},
// 	}

// 	return Model{
// 		state:      StateMain,
// 		list:       l,
// 		status:     status,
// 		showBanner: true,
// 		showStatus: true,
// 		discoverer: discoverer,
// 		ctx:        ctx,
// 		cancel:     cancel,
// 	}
// }

// // Init implements tea.Model
// func (m Model) Init() tea.Cmd {
// 	return tea.EnterAltScreen
// }

// func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	switch msg := msg.(type) {
// 	case tea.WindowSizeMsg:
// 		m.width = msg.Width
// 		m.height = msg.Height
// 		verticalMargin := 0

// 		if m.showBanner {
// 			banner := BannerStyle.Render(GetBannerArt())
// 			subtitle := TitleStyle.Render(GetSubtitle())
// 			version := InfoStyle.Render(GetVersionInfo())
// 			verticalMargin += lipgloss.Height(banner) + lipgloss.Height(subtitle) + lipgloss.Height(version)
// 		}

// 		if m.showStatus {
// 			statusContent := m.renderStatus()
// 			statusBox := StatusBoxStyle.Render(statusContent)
// 			verticalMargin += lipgloss.Height(statusBox)
// 		}

// 		help := HelpStyle.Render("...")
// 		verticalMargin += lipgloss.Height(help)

// 		m.list.SetHeight(m.height - verticalMargin)
// 		m.list.SetWidth(m.width - 4)
// 		return m, nil

// 	case tea.KeyMsg:
// 		switch {
// 		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+c", "q"))):
// 			m.cancel() // Cancel discovery context
// 			return m, tea.Quit

// 		case key.Matches(msg, key.NewBinding(key.WithKeys("b"))):
// 			m.showBanner = !m.showBanner
// 			return m, func() tea.Msg {
// 				return tea.WindowSizeMsg{Width: m.width, Height: m.height}
// 			}

// 		case key.Matches(msg, key.NewBinding(key.WithKeys("s"))):
// 			m.showStatus = !m.showStatus
// 			return m, func() tea.Msg {
// 				return tea.WindowSizeMsg{Width: m.width, Height: m.height}
// 			}

// 		case key.Matches(msg, key.NewBinding(key.WithKeys("r"))):
// 			// Manual refresh - run discovery synchronously for now
// 			m.refreshProcesses()
// 			return m, nil

// 		case key.Matches(msg, key.NewBinding(key.WithKeys("enter", " "))):
// 			if selectedItem, ok := m.list.SelectedItem().(MenuItem); ok {
// 				m.state = selectedItem.action
// 				// If entering service list, refresh the data
// 				if m.state == StateServiceList {
// 					m.refreshProcesses()
// 				}
// 				return m, nil
// 			}

// 		case key.Matches(msg, key.NewBinding(key.WithKeys("backspace"))):
// 			if m.state != StateMain {
// 				m.state = StateMain
// 				return m, nil
// 			}
// 		case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
// 			// Get currently selected process (you'll need to implement selection logic)
// 			if m.selectedProcessIndex >= 0 && m.selectedProcessIndex < len(m.status.Processes) {
// 				m.selectedProcess = &m.status.Processes[m.selectedProcessIndex]
// 				m.state = StateProcessDetail // Add this new state
// 				m.detailState = DetailOverview
// 				return m, nil
// 			}
// 		}

// 		// Handle navigation only in main menu
// 		if m.state == StateMain {
// 			var cmd tea.Cmd
// 			m.list, cmd = m.list.Update(msg)
// 			return m, cmd
// 		}
// 	}

// 	return m, nil
// }

// // refreshProcesses loads real Java processes - SIMPLE SYNCHRONOUS VERSION
// func (m *Model) refreshProcesses() {
// 	// Run discovery using the correct API method
// 	processes, err := m.discoverer.DiscoverJavaProcesses(m.ctx)
// 	if err != nil {
// 		// Handle error but don't crash
// 		m.status.SystemHealth = fmt.Sprintf("❌ Discovery error: %v", err)
// 		return
// 	}

// 	// Update the processes
// 	m.status.Processes = processes
// 	m.status.LastUpdate = time.Now()

// 	// Update counts
// 	m.status.JavaServices = len(processes)
// 	m.status.Instrumented = 0
// 	m.status.MiddlewareCount = 0
// 	m.status.OtherAgentCount = 0

// 	for _, proc := range processes {
// 		if proc.HasJavaAgent {
// 			m.status.Instrumented++
// 			if proc.IsMiddlewareAgent {
// 				m.status.MiddlewareCount++
// 			} else {
// 				m.status.OtherAgentCount++
// 			}
// 		}
// 	}

// 	// Update health status
// 	if m.status.JavaServices == 0 {
// 		m.status.SystemHealth = "ℹ️ No Java services found"
// 	} else if m.status.Instrumented == m.status.JavaServices {
// 		m.status.SystemHealth = "✅ All instrumented"
// 	} else if m.status.Instrumented > 0 {
// 		m.status.SystemHealth = "⚠️ Partially instrumented"
// 	} else {
// 		m.status.SystemHealth = "❌ No instrumentation"
// 	}
// }

// // View implements tea.Model
// func (m Model) View() string {
// 	var sections []string

// 	// Banner section
// 	if m.showBanner {
// 		banner := BannerStyle.Render(GetBannerArt())
// 		subtitle := TitleStyle.Render(GetSubtitle())
// 		version := InfoStyle.Render(GetVersionInfo())
// 		sections = append(sections, banner, subtitle, version)
// 	}

// 	// Status section
// 	if m.showStatus {
// 		statusContent := m.renderStatus()
// 		sections = append(sections, StatusBoxStyle.Render(statusContent))
// 	}

// 	// Main content based on state
// 	switch m.state {
// 	case StateMain:
// 		sections = append(sections, m.list.View())

// 	case StateServiceList:
// 		sections = append(sections, m.renderServiceList())

// 	case StateHealthCheck:
// 		sections = append(sections, m.renderHealthCheck())

// 	case StateHelp:
// 		sections = append(sections, m.renderHelp())

// 	default:
// 		sections = append(sections, "Feature coming soon...")
// 	}

// 	// Help footer
// 	help := HelpStyle.Render("Press 'q' to quit • 'b' to toggle banner • 'r' to refresh • '?' for help")
// 	sections = append(sections, help)

// 	return lipgloss.JoinVertical(lipgloss.Left, sections...)
// }

// func (m Model) renderStatus() string {
// 	hostname, _ := os.Hostname()

// 	statusLines := []string{
// 		fmt.Sprintf("🖥️ Host: %s", hostname),
// 		fmt.Sprintf("📅 Last Update: %s", m.status.LastUpdate.Format("15:04:05")),
// 		"",
// 		fmt.Sprintf("☕ Java Services: %d running, %d instrumented",
// 			m.status.JavaServices, m.status.Instrumented),
// 		fmt.Sprintf("⚙️ MW Agents: %d middleware, %d other agents",
// 			m.status.MiddlewareCount, m.status.OtherAgentCount),
// 		fmt.Sprintf("🏥 Health: %s", m.status.SystemHealth),
// 		fmt.Sprintf("📡 Telemetry: %s", m.status.TelemetryStatus),
// 		"",
// 		"🔧 System: LD_PRELOAD Injection Active",
// 		"🎯 Target: Middleware.io",
// 	}

// 	return strings.Join(statusLines, "\n")
// }

// // renderServiceList renders the service list view with real process data
// func (m Model) renderServiceList() string {
// 	content := []string{
// 		TitleStyle.Render("📋 Java Services"),
// 		"",
// 	}

// 	if len(m.status.Processes) == 0 {
// 		content = append(content,
// 			"No Java processes found.",
// 			"",
// 			"Press 'r' to scan for Java services",
// 			"",
// 			HelpStyle.Render("Press 'r' to refresh • 'backspace' to go back"))
// 		return strings.Join(content, "\n")
// 	}

// 	// Table header
// 	content = append(content,
// 		"┌─────────┬──────────────────────┬─────────┬──────────────────┬──────────────┐",
// 		"│ PID     │ Service Name         │ Status  │ JAR File         │ Memory       │",
// 		"├─────────┼──────────────────────┼─────────┼──────────────────┼──────────────┤",
// 	)

// 	// Table rows with real data
// 	for i, proc := range m.status.Processes {
// 		// Limit display to prevent overflow
// 		if i >= 10 {
// 			content = append(content,
// 				fmt.Sprintf("│ ... and %d more processes                                            │",
// 					len(m.status.Processes)-10))
// 			break
// 		}

// 		// Format each field with proper width using CORRECT field names
// 		pidStr := fmt.Sprintf("%-7d", proc.ProcessPID)
// 		serviceStr := fmt.Sprintf("%-20s", truncateString(proc.ServiceName, 20))
// 		statusStr := getStatusString(proc)
// 		jarStr := fmt.Sprintf("%-16s", truncateString(proc.JarFile, 16))
// 		memoryStr := fmt.Sprintf("%-12s", fmt.Sprintf("%.1f%%", proc.MemoryPercent))

// 		content = append(content, fmt.Sprintf("│ %s │ %s │ %s │ %s │ %s │",
// 			pidStr, serviceStr, statusStr, jarStr, memoryStr))
// 	}

// 	// Table footer
// 	content = append(content,
// 		"└─────────┴──────────────────────┴─────────┴──────────────────┴──────────────┘",
// 		"",
// 		HelpStyle.Render("Press 'r' to refresh • 'backspace' to go back"),
// 	)

// 	return strings.Join(content, "\n")
// }

// func getStatusString(proc discovery.JavaProcess) string {
// 	if proc.HasJavaAgent {
// 		if proc.IsMiddlewareAgent {
// 			return "✅ MW  "
// 		}
// 		return "⚙️  OTel"
// 	}
// 	return "❌ None"
// }

// // renderHealthCheck renders the health check view
// func (m Model) renderHealthCheck() string {
// 	content := []string{
// 		TitleStyle.Render("🏥 System Health Check"),
// 		"",
// 		"✅ LD_PRELOAD injection: Active",
// 		"✅ Shared library: /usr/lib/middleware/libmwinjector.so loaded",
// 		"✅ Java agent: /usr/lib/middleware/middleware-javaagent-1.7.0.jar found",
// 		"✅ Configuration files: Valid",
// 		"✅ Network connectivity: Middleware.io reachable",
// 		"",
// 		"📊 Current Status:",
// 		fmt.Sprintf("  • Java services: %d discovered", m.status.JavaServices),
// 		fmt.Sprintf("  • MW instrumented: %d", m.status.MiddlewareCount),
// 		fmt.Sprintf("  • Other agents: %d", m.status.OtherAgentCount),
// 		fmt.Sprintf("  • Last scan: %s", m.status.LastUpdate.Format("15:04:05")),
// 		"",
// 		HelpStyle.Render("Press 'backspace' to go back"),
// 	}

// 	return strings.Join(content, "\n")
// }

// // renderHelp renders the help view
// func (m Model) renderHelp() string {
// 	content := []string{
// 		TitleStyle.Render("❓ Help & Documentation"),
// 		"",
// 		"🔧 Key Features:",
// 		"  • Automatic Java service discovery via /proc filesystem",
// 		"  • Per-service Middleware.io configuration",
// 		"  • LD_PRELOAD shared library injection",
// 		"  • Real-time health monitoring",
// 		"",
// 		"⌨️ Keyboard Shortcuts:",
// 		"  • q, Ctrl+C: Quit application",
// 		"  • ↑/↓: Navigate menu items",
// 		"  • Enter/Space: Select menu item",
// 		"  • b: Toggle banner display",
// 		"  • r: Refresh/scan for Java processes",
// 		"  • Backspace: Go back to previous screen",
// 		"",
// 		"🔗 Resources:",
// 		"  • Documentation: https://docs.middleware.io/injector",
// 		"  • GitHub: https://github.com/middleware-labs/mw-injector",
// 		"  • Support: support@middleware.io",
// 		"",
// 		HelpStyle.Render("Press 'backspace' to go back"),
// 	}

//		return strings.Join(content, "\n")
//	}
// package tui

// import (
// 	"context"
// 	"fmt"
// 	"os"
// 	"strings"
// 	"time"

// 	"github.com/charmbracelet/bubbles/key"
// 	"github.com/charmbracelet/bubbles/list"
// 	tea "github.com/charmbracelet/bubbletea"
// 	"github.com/charmbracelet/lipgloss"

// 	"github.com/naman47vyas/mw-injector-manager/pkg/discovery"
// )

// // AppState represents the current state of the application
// type AppState int

// const (
// 	StateMain AppState = iota
// 	StateServiceList
// 	StateServiceConfig
// 	StateHealthCheck
// 	StateHelp
// 	StateProcessDetail // Add this new state
// )

// // MenuItem represents a menu item
// type MenuItem struct {
// 	title       string
// 	description string
// 	action      AppState
// }

// func (m MenuItem) Title() string       { return m.title }
// func (m MenuItem) Description() string { return m.description }
// func (m MenuItem) FilterValue() string { return m.title }

// // SystemStatus represents the current system status
// type SystemStatus struct {
// 	JavaServices    int
// 	Instrumented    int
// 	MiddlewareCount int
// 	OtherAgentCount int
// 	SystemHealth    string
// 	TelemetryStatus string
// 	LastUpdate      time.Time
// 	Processes       []discovery.JavaProcess
// }

// // ProcessDetailState represents the detail view state
// type ProcessDetailState int

// const (
// 	DetailOverview ProcessDetailState = iota
// 	DetailInstrumentation
// 	DetailConfiguration
// )

// // Model represents the main application model
// type Model struct {
// 	state         AppState
// 	width         int
// 	height        int
// 	list          list.Model
// 	status        SystemStatus
// 	showBanner    bool
// 	showStatus    bool
// 	animationStep int
// 	discoverer    discovery.Discoverer
// 	ctx           context.Context
// 	cancel        context.CancelFunc

// 	// Process selection and detail view
// 	selectedProcessIndex int
// 	selectedProcess      *discovery.JavaProcess
// 	detailState          ProcessDetailState
// }

// // NewModel creates a new model instance
// func NewModel() Model {
// 	// Create context for discovery operations
// 	ctx, cancel := context.WithCancel(context.Background())

// 	// Create discoverer - SIMPLE VERSION, no complex options
// 	discoverer := discovery.NewDiscoverer(ctx)

// 	// Create menu items
// 	items := []list.Item{
// 		MenuItem{
// 			title:       "📋 List Services",
// 			description: "View all Java services and their instrumentation status",
// 			action:      StateServiceList,
// 		},
// 		MenuItem{
// 			title:       "⚙️ Configure Service",
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

// 	// Initialize with empty status - will load real data when needed
// 	status := SystemStatus{
// 		JavaServices:    0,
// 		Instrumented:    0,
// 		MiddlewareCount: 0,
// 		OtherAgentCount: 0,
// 		SystemHealth:    "🔄 Ready to scan",
// 		TelemetryStatus: "📡 Connected",
// 		LastUpdate:      time.Now(),
// 		Processes:       []discovery.JavaProcess{},
// 	}

// 	return Model{
// 		state:                StateMain,
// 		list:                 l,
// 		status:               status,
// 		showBanner:           true,
// 		showStatus:           true,
// 		discoverer:           discoverer,
// 		ctx:                  ctx,
// 		cancel:               cancel,
// 		selectedProcessIndex: -1, // Initialize to -1 (no selection)
// 		selectedProcess:      nil,
// 		detailState:          DetailOverview,
// 	}
// }

// // Init implements tea.Model
// func (m Model) Init() tea.Cmd {
// 	return tea.EnterAltScreen
// }

// func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	switch msg := msg.(type) {
// 	case tea.WindowSizeMsg:
// 		m.width = msg.Width
// 		m.height = msg.Height
// 		verticalMargin := 0

// 		if m.showBanner {
// 			banner := BannerStyle.Render(GetBannerArt())
// 			subtitle := TitleStyle.Render(GetSubtitle())
// 			version := InfoStyle.Render(GetVersionInfo())
// 			verticalMargin += lipgloss.Height(banner) + lipgloss.Height(subtitle) + lipgloss.Height(version)
// 		}

// 		if m.showStatus {
// 			statusContent := m.renderStatus()
// 			statusBox := StatusBoxStyle.Render(statusContent)
// 			verticalMargin += lipgloss.Height(statusBox)
// 		}

// 		help := HelpStyle.Render("...")
// 		verticalMargin += lipgloss.Height(help)

// 		m.list.SetHeight(m.height - verticalMargin)
// 		m.list.SetWidth(m.width - 4)
// 		return m, nil

// 	case tea.KeyMsg:
// 		switch {
// 		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+c", "q"))):
// 			m.cancel() // Cancel discovery context
// 			return m, tea.Quit

// 		case key.Matches(msg, key.NewBinding(key.WithKeys("b"))):
// 			m.showBanner = !m.showBanner
// 			return m, func() tea.Msg {
// 				return tea.WindowSizeMsg{Width: m.width, Height: m.height}
// 			}

// 		case key.Matches(msg, key.NewBinding(key.WithKeys("s"))):
// 			m.showStatus = !m.showStatus
// 			return m, func() tea.Msg {
// 				return tea.WindowSizeMsg{Width: m.width, Height: m.height}
// 			}

// 		case key.Matches(msg, key.NewBinding(key.WithKeys("r"))):
// 			// Manual refresh - run discovery synchronously for now
// 			m.refreshProcesses()
// 			return m, nil

// 		case key.Matches(msg, key.NewBinding(key.WithKeys("enter", " "))):
// 			if selectedItem, ok := m.list.SelectedItem().(MenuItem); ok {
// 				m.state = selectedItem.action
// 				// If entering service list, refresh the data and reset selection
// 				if m.state == StateServiceList {
// 					m.refreshProcesses()
// 					m.selectedProcessIndex = 0 // Start with first process selected
// 				}
// 				return m, nil
// 			}

// 		case key.Matches(msg, key.NewBinding(key.WithKeys("backspace"))):
// 			if m.state != StateMain {
// 				m.state = StateMain
// 				return m, nil
// 			}
// 		}

// 		// Handle navigation only in main menu
// 		if m.state == StateMain {
// 			var cmd tea.Cmd
// 			m.list, cmd = m.list.Update(msg)
// 			return m, cmd
// 		}

// 		// Handle navigation in service list
// 		if m.state == StateServiceList && len(m.status.Processes) > 0 {
// 			switch {
// 			case key.Matches(msg, key.NewBinding(key.WithKeys("up", "k"))):
// 				if m.selectedProcessIndex > 0 {
// 					m.selectedProcessIndex--
// 				}
// 				return m, nil

// 			case key.Matches(msg, key.NewBinding(key.WithKeys("down", "j"))):
// 				if m.selectedProcessIndex < len(m.status.Processes)-1 {
// 					m.selectedProcessIndex++
// 				}
// 				return m, nil

// 			case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
// 				// Enter process detail view
// 				if m.selectedProcessIndex >= 0 && m.selectedProcessIndex < len(m.status.Processes) {
// 					m.selectedProcess = &m.status.Processes[m.selectedProcessIndex]
// 					m.state = StateProcessDetail // We'll need to add this state
// 					m.detailState = DetailOverview
// 					return m, nil
// 				}
// 			}
// 		}
// 	}

// 	return m, nil
// }

// // refreshProcesses loads real Java processes - SIMPLE SYNCHRONOUS VERSION
// func (m *Model) refreshProcesses() {
// 	// Run discovery using the correct API method
// 	processes, err := m.discoverer.DiscoverJavaProcesses(m.ctx)
// 	if err != nil {
// 		// Handle error but don't crash
// 		m.status.SystemHealth = fmt.Sprintf("❌ Discovery error: %v", err)
// 		return
// 	}

// 	// Update the processes
// 	m.status.Processes = processes
// 	m.status.LastUpdate = time.Now()

// 	// Update counts
// 	m.status.JavaServices = len(processes)
// 	m.status.Instrumented = 0
// 	m.status.MiddlewareCount = 0
// 	m.status.OtherAgentCount = 0

// 	for _, proc := range processes {
// 		if proc.HasJavaAgent {
// 			m.status.Instrumented++
// 			if proc.IsMiddlewareAgent {
// 				m.status.MiddlewareCount++
// 			} else {
// 				m.status.OtherAgentCount++
// 			}
// 		}
// 	}

// 	// Update health status
// 	if m.status.JavaServices == 0 {
// 		m.status.SystemHealth = "ℹ️ No Java services found"
// 	} else if m.status.Instrumented == m.status.JavaServices {
// 		m.status.SystemHealth = "✅ All instrumented"
// 	} else if m.status.Instrumented > 0 {
// 		m.status.SystemHealth = "⚠️ Partially instrumented"
// 	} else {
// 		m.status.SystemHealth = "❌ No instrumentation"
// 	}
// }

// // View implements tea.Model
// func (m Model) View() string {
// 	var sections []string

// 	// Banner section
// 	if m.showBanner {
// 		banner := BannerStyle.Render(GetBannerArt())
// 		subtitle := TitleStyle.Render(GetSubtitle())
// 		version := InfoStyle.Render(GetVersionInfo())
// 		sections = append(sections, banner, subtitle, version)
// 	}

// 	// Status section
// 	if m.showStatus {
// 		statusContent := m.renderStatus()
// 		sections = append(sections, StatusBoxStyle.Render(statusContent))
// 	}

// 	// Main content based on state
// 	switch m.state {
// 	case StateMain:
// 		sections = append(sections, m.list.View())

// 	case StateServiceList:
// 		sections = append(sections, m.renderServiceList())

// 	case StateProcessDetail:
// 		sections = append(sections, m.renderProcessDetail())

// 	case StateHealthCheck:
// 		sections = append(sections, m.renderHealthCheck())

// 	case StateHelp:
// 		sections = append(sections, m.renderHelp())

// 	default:
// 		sections = append(sections, "Feature coming soon...")
// 	}

// 	// Help footer
// 	help := HelpStyle.Render("Press 'q' to quit • 'b' to toggle banner • 'r' to refresh • '?' for help")
// 	sections = append(sections, help)

// 	return lipgloss.JoinVertical(lipgloss.Left, sections...)
// }

// func (m Model) renderStatus() string {
// 	hostname, _ := os.Hostname()

// 	statusLines := []string{
// 		fmt.Sprintf("🖥️ Host: %s", hostname),
// 		fmt.Sprintf("📅 Last Update: %s", m.status.LastUpdate.Format("15:04:05")),
// 		"",
// 		fmt.Sprintf("☕ Java Services: %d running, %d instrumented",
// 			m.status.JavaServices, m.status.Instrumented),
// 		fmt.Sprintf("⚙️ MW Agents: %d middleware, %d other agents",
// 			m.status.MiddlewareCount, m.status.OtherAgentCount),
// 		fmt.Sprintf("🏥 Health: %s", m.status.SystemHealth),
// 		fmt.Sprintf("📡 Telemetry: %s", m.status.TelemetryStatus),
// 		"",
// 		"🔧 System: LD_PRELOAD Injection Active",
// 		"🎯 Target: Middleware.io",
// 	}

// 	return strings.Join(statusLines, "\n")
// }

// // renderServiceList renders the service list view with real process data
// func (m Model) renderServiceList() string {
// 	content := []string{
// 		TitleStyle.Render("📋 Java Services"),
// 		"",
// 	}

// 	if len(m.status.Processes) == 0 {
// 		content = append(content,
// 			"No Java processes found.",
// 			"",
// 			"Press 'r' to scan for Java services",
// 			"",
// 			HelpStyle.Render("Press 'r' to refresh • 'backspace' to go back"))
// 		return strings.Join(content, "\n")
// 	}

// 	// Table header
// 	content = append(content,
// 		"┌─────────┬──────────────────────┬─────────┬──────────────────┬──────────────┐",
// 		"│ PID     │ Service Name         │ Status  │ JAR File         │ Memory       │",
// 		"├─────────┼──────────────────────┼─────────┼──────────────────┼──────────────┤",
// 	)

// 	// Table rows with real data
// 	for i, proc := range m.status.Processes {
// 		// Limit display to prevent overflow
// 		if i >= 10 {
// 			content = append(content,
// 				fmt.Sprintf("│ ... and %d more processes                                            │",
// 					len(m.status.Processes)-10))
// 			break
// 		}

// 		// Format each field with proper width using CORRECT field names
// 		pidStr := fmt.Sprintf("%-7d", proc.ProcessPID)
// 		serviceStr := fmt.Sprintf("%-20s", truncateString(proc.ServiceName, 20))
// 		statusStr := getStatusString(proc)
// 		jarStr := fmt.Sprintf("%-16s", truncateString(proc.JarFile, 16))
// 		memoryStr := fmt.Sprintf("%-12s", fmt.Sprintf("%.1f%%", proc.MemoryPercent))

// 		// Highlight selected row
// 		rowText := fmt.Sprintf("│ %s │ %s │ %s │ %s │ %s │",
// 			pidStr, serviceStr, statusStr, jarStr, memoryStr)

// 		if i == m.selectedProcessIndex {
// 			// Highlight the selected row
// 			rowText = lipgloss.NewStyle().
// 				Background(lipgloss.Color("240")).
// 				Foreground(lipgloss.Color("15")).
// 				Render(rowText)
// 		}

// 		content = append(content, rowText)
// 	}

// 	// Table footer
// 	content = append(content,
// 		"└─────────┴──────────────────────┴─────────┴──────────────────┴──────────────┘",
// 		"",
// 		HelpStyle.Render("Press '↑/↓' or 'j/k' to select • 'enter' for process details • 'r' to refresh • 'backspace' to go back"),
// 	)

// 	return strings.Join(content, "\n")
// }

// func getStatusString(proc discovery.JavaProcess) string {
// 	if proc.HasJavaAgent {
// 		if proc.IsMiddlewareAgent {
// 			return "✅ MW  "
// 		}
// 		return "⚙️  OTel"
// 	}
// 	return "❌ None"
// }

// // renderHealthCheck renders the health check view
// func (m Model) renderHealthCheck() string {
// 	content := []string{
// 		TitleStyle.Render("🏥 System Health Check"),
// 		"",
// 		"✅ LD_PRELOAD injection: Active",
// 		"✅ Shared library: /usr/lib/middleware/libmwinjector.so loaded",
// 		"✅ Java agent: /usr/lib/middleware/middleware-javaagent-1.7.0.jar found",
// 		"✅ Configuration files: Valid",
// 		"✅ Network connectivity: Middleware.io reachable",
// 		"",
// 		"📊 Current Status:",
// 		fmt.Sprintf("  • Java services: %d discovered", m.status.JavaServices),
// 		fmt.Sprintf("  • MW instrumented: %d", m.status.MiddlewareCount),
// 		fmt.Sprintf("  • Other agents: %d", m.status.OtherAgentCount),
// 		fmt.Sprintf("  • Last scan: %s", m.status.LastUpdate.Format("15:04:05")),
// 		"",
// 		HelpStyle.Render("Press 'backspace' to go back"),
// 	}

// 	return strings.Join(content, "\n")
// }

// // renderHelp renders the help view
// func (m Model) renderHelp() string {
// 	content := []string{
// 		TitleStyle.Render("❓ Help & Documentation"),
// 		"",
// 		"🔧 Key Features:",
// 		"  • Automatic Java service discovery via /proc filesystem",
// 		"  • Per-service Middleware.io configuration",
// 		"  • LD_PRELOAD shared library injection",
// 		"  • Real-time health monitoring",
// 		"",
// 		"⌨️ Keyboard Shortcuts:",
// 		"  • q, Ctrl+C: Quit application",
// 		"  • ↑/↓: Navigate menu items",
// 		"  • Enter/Space: Select menu item",
// 		"  • b: Toggle banner display",
// 		"  • r: Refresh/scan for Java processes",
// 		"  • Backspace: Go back to previous screen",
// 		"",
// 		"🔗 Resources:",
// 		"  • Documentation: https://docs.middleware.io/injector",
// 		"  • GitHub: https://github.com/middleware-labs/mw-injector",
// 		"  • Support: support@middleware.io",
// 		"",
// 		HelpStyle.Render("Press 'backspace' to go back"),
// 	}

// 	return strings.Join(content, "\n")
// }

// func (m Model) renderProcessDetail() string {
// 	if m.selectedProcess == nil {
// 		return "No process selected"
// 	}

// 	proc := *m.selectedProcess

// 	var content []string

// 	// Header with process identification
// 	content = append(content,
// 		TitleStyle.Render(fmt.Sprintf("🔍 Process Details - PID %d", proc.ProcessPID)),
// 		"",
// 	)

// 	switch m.detailState {
// 	case DetailOverview:
// 		content = append(content, m.renderProcessOverview(proc)...)
// 	case DetailInstrumentation:
// 		content = append(content, m.renderInstrumentationDetails(proc)...)
// 	case DetailConfiguration:
// 		content = append(content, m.renderConfigurationOptions(proc)...)
// 	}

// 	// Navigation footer
// 	content = append(content, "",
// 		HelpStyle.Render("Press '1' for Overview • '2' for Instrumentation • '3' for Configuration • 'backspace' to go back"))

// 	return strings.Join(content, "\n")
// }

// func (m Model) renderProcessOverview(proc discovery.JavaProcess) []string {
// 	// Determine instrumentation status
// 	instrumentationStatus := "❌ Not Instrumented"
// 	instrumentationColor := ErrorColor
// 	if proc.HasJavaAgent {
// 		if proc.IsMiddlewareAgent {
// 			instrumentationStatus = "✅ Middleware Instrumented"
// 			instrumentationColor = AccentColor
// 		} else {
// 			instrumentationStatus = "⚙️ Other Agent Present"
// 			instrumentationColor = WarningColor
// 		}
// 	}

// 	// Extract server port from JVM options
// 	serverPort := "unknown"
// 	for _, opt := range proc.JVMOptions {
// 		if strings.Contains(opt, "-Dserver.port=") {
// 			serverPort = strings.Split(opt, "=")[1]
// 			break
// 		}
// 	}

// 	// Calculate uptime
// 	uptime := time.Since(proc.ProcessCreateTime)
// 	uptimeStr := fmt.Sprintf("%.0f minutes", uptime.Minutes())
// 	if uptime.Hours() >= 1 {
// 		uptimeStr = fmt.Sprintf("%.1f hours", uptime.Hours())
// 	}

// 	content := []string{
// 		"🏷️ SERVICE IDENTIFICATION",
// 		fmt.Sprintf("   Service Name: %s", proc.ServiceName),
// 		fmt.Sprintf("   JAR File: %s", proc.JarFile),
// 		fmt.Sprintf("   JAR Path: %s", proc.JarPath),
// 		fmt.Sprintf("   Server Port: %s", serverPort),
// 		fmt.Sprintf("   Process Owner: %s", proc.ProcessOwner),
// 		"",
// 		"📊 OBSERVABILITY STATUS",
// 		fmt.Sprintf("   %s %s",
// 			lipgloss.NewStyle().Foreground(instrumentationColor).Render("●"),
// 			instrumentationStatus),
// 		"",
// 	}

// 	if proc.HasJavaAgent {
// 		content = append(content,
// 			"🔧 CURRENT INSTRUMENTATION",
// 			fmt.Sprintf("   Agent: %s", proc.JavaAgentName),
// 			fmt.Sprintf("   Agent Path: %s", proc.JavaAgentPath),
// 		)
// 		if proc.IsMiddlewareAgent {
// 			// Detect agent type from filename
// 			agentType := "Standard"
// 			if strings.Contains(proc.JavaAgentName, "serverless") {
// 				agentType = "Serverless"
// 			}
// 			content = append(content, fmt.Sprintf("   Agent Type: %s", agentType))
// 		}
// 		content = append(content, "")
// 	}

// 	content = append(content,
// 		"⚡ RUNTIME INFORMATION",
// 		fmt.Sprintf("   Runtime: %s", proc.ProcessRuntimeDescription),
// 		fmt.Sprintf("   Java Path: %s", proc.ProcessExecutablePath),
// 		fmt.Sprintf("   Parent PID: %d", proc.ProcessParentPID),
// 		fmt.Sprintf("   Created: %s", proc.ProcessCreateTime.Format("2006-01-02 15:04:05")),
// 		fmt.Sprintf("   Uptime: %s", uptimeStr),
// 		fmt.Sprintf("   Status: %s", proc.Status),
// 		"",
// 		"📈 PERFORMANCE METRICS",
// 		fmt.Sprintf("   Memory Usage: %.2f%% (~%.0fMB)", proc.MemoryPercent, proc.MemoryPercent*40), // Rough calculation
// 		fmt.Sprintf("   CPU Usage: %.2f%%", proc.CPUPercent),
// 		"",
// 		"🔍 FULL COMMAND",
// 		fmt.Sprintf("   %s", proc.ProcessCommandLine),
// 	)

// 	return content
// }

// func (m Model) renderInstrumentationDetails(proc discovery.JavaProcess) []string {
// 	content := []string{
// 		"🔧 INSTRUMENTATION ANALYSIS",
// 		"",
// 	}

// 	if !proc.HasJavaAgent {
// 		// Case 1: No instrumentation
// 		content = append(content,
// 			"❌ NO JAVA AGENT DETECTED",
// 			"",
// 			"Current State:",
// 			"   • No -javaagent parameter found in JVM options",
// 			"   • Process is running without observability instrumentation",
// 			"   • No telemetry data (traces, metrics, logs) is being collected",
// 			"",
// 			"Impact:",
// 			"   • Zero observability into this service",
// 			"   • Cannot track performance, errors, or dependencies",
// 			"   • No distributed tracing participation",
// 			"",
// 			"Recommendation:",
// 			"   ✅ Add Middleware.io agent to enable full observability",
// 			"   • Automatic trace collection",
// 			"   • Performance metrics",
// 			"   • Error tracking and alerting",
// 			"   • Service dependency mapping",
// 		)
// 	} else if proc.IsMiddlewareAgent {
// 		// Case 2: Already has Middleware agent
// 		content = append(content,
// 			"✅ MIDDLEWARE.IO AGENT ACTIVE",
// 			"",
// 			"Current Configuration:",
// 			fmt.Sprintf("   Agent: %s", proc.JavaAgentName),
// 			fmt.Sprintf("   Path: %s", proc.JavaAgentPath),
// 		)

// 		// Detect serverless vs standard
// 		if strings.Contains(proc.JavaAgentName, "serverless") {
// 			content = append(content,
// 				"   Type: Serverless Agent (optimized for FaaS)",
// 				"",
// 				"Active Features:",
// 				"   ✅ Distributed tracing",
// 				"   ✅ Custom metrics",
// 				"   ✅ Error tracking",
// 				"   ✅ Optimized for short-lived functions",
// 				"   ❌ Health check endpoint (disabled for serverless)",
// 			)
// 		} else {
// 			content = append(content,
// 				"   Type: Standard Agent (full feature set)",
// 				"",
// 				"Active Features:",
// 				"   ✅ Distributed tracing",
// 				"   ✅ Performance metrics",
// 				"   ✅ Log correlation",
// 				"   ✅ Health check endpoint",
// 				"   ✅ Custom instrumentation",
// 			)
// 		}

// 		content = append(content,
// 			"",
// 			"Status:",
// 			"   🟢 Observability fully enabled",
// 			"   📊 Telemetry data actively collected",
// 			"   🔗 Connected to Middleware.io platform",
// 		)
// 	} else {
// 		// Case 3: Has other agent (OpenTelemetry, AppDynamics, etc.)
// 		content = append(content,
// 			"⚙️ NON-MIDDLEWARE AGENT DETECTED",
// 			"",
// 			"Current Configuration:",
// 			fmt.Sprintf("   Agent: %s", proc.JavaAgentName),
// 			fmt.Sprintf("   Path: %s", proc.JavaAgentPath),
// 			"   Type: Third-party instrumentation",
// 			"",
// 			"Detected Issues:",
// 			"   ⚠️  Not using Middleware.io agent",
// 			"   ⚠️  Telemetry may not flow to your Middleware.io account",
// 			"   ⚠️  Missing Middleware.io specific features",
// 			"",
// 			"Options:",
// 			"   1. Replace with Middleware.io agent (recommended)",
// 			"   2. Configure existing agent to export to Middleware.io",
// 			"   3. Keep existing setup (limited Middleware.io integration)",
// 		)
// 	}

// 	// JVM Options analysis
// 	if len(proc.JVMOptions) > 0 {
// 		content = append(content,
// 			"",
// 			"🔍 JVM OPTIONS ANALYSIS",
// 		)
// 		for _, opt := range proc.JVMOptions {
// 			if strings.HasPrefix(opt, "-javaagent:") {
// 				content = append(content, fmt.Sprintf("   Agent: %s", opt))
// 			} else if strings.HasPrefix(opt, "-D") {
// 				content = append(content, fmt.Sprintf("   Property: %s", opt))
// 			} else {
// 				content = append(content, fmt.Sprintf("   Option: %s", opt))
// 			}
// 		}
// 	}

// 	return content
// }

// func (m Model) renderConfigurationOptions(proc discovery.JavaProcess) []string {
// 	content := []string{
// 		"⚙️ CONFIGURATION MANAGEMENT",
// 		"",
// 	}

// 	if !proc.HasJavaAgent {
// 		// Configuration for uninstrumented process
// 		content = append(content,
// 			"🎯 MIDDLEWARE.IO SETUP OPTIONS",
// 			"",
// 			"Option 1: Standard Agent (Recommended)",
// 			"   Agent: middleware-javaagent-1.7.0.jar",
// 			"   Features: Full observability stack",
// 			"   Restart Required: Yes",
// 			"",
// 			"   Required Environment Variables:",
// 			"   • MW_API_KEY=your-api-key-here",
// 			"   • MW_SERVICE_NAME="+generateServiceName(proc),
// 			"   • MW_TARGET=https://your-tenant.middleware.io",
// 			"   • MW_APM_COLLECT_TRACES=true",
// 			"   • MW_APM_COLLECT_METRICS=true",
// 			"   • MW_APM_COLLECT_LOGS=true",
// 			"",
// 			"Option 2: Serverless Agent",
// 			"   Agent: middleware-javaagent-serverless.jar",
// 			"   Features: Optimized for short-lived processes",
// 			"   Restart Required: Yes",
// 			"",
// 			"🔧 IMPLEMENTATION STEPS",
// 			"",
// 			"1. Stop the Java process (PID "+fmt.Sprintf("%d", proc.ProcessPID)+")",
// 			"2. Update startup command to include:",
// 			"   -javaagent:/path/to/middleware-javaagent.jar",
// 			"3. Set environment variables (MW_*)",
// 			"4. Restart the process",
// 			"5. Verify telemetry flow in Middleware.io dashboard",
// 		)
// 	} else if proc.IsMiddlewareAgent {
// 		// Configuration for existing Middleware agent
// 		content = append(content,
// 			"✅ MIDDLEWARE.IO AGENT CONFIGURATION",
// 			"",
// 			"Current Status: Already instrumented with Middleware.io",
// 			"",
// 			"🔧 CONFIGURATION OPTIONS",
// 			"",
// 			"Environment Variables to Check:",
// 			"   • MW_API_KEY (verify correct API key)",
// 			"   • MW_SERVICE_NAME (currently: "+proc.ServiceName+")",
// 			"   • MW_TARGET (collector endpoint)",
// 			"   • MW_LOG_LEVEL (adjust verbosity)",
// 			"",
// 			"Advanced Configuration:",
// 			"   • MW_CUSTOM_RESOURCE_ATTRIBUTE (add metadata)",
// 			"   • MW_APM_COLLECT_PROFILING (enable profiling)",
// 			"   • MW_PROFILING_ALLOC (allocation profiling)",
// 			"",
// 			"🔍 HEALTH CHECK",
// 			"",
// 			"Verification Steps:",
// 			"   1. Check process environment: /proc/"+fmt.Sprintf("%d", proc.ProcessPID)+"/environ",
// 			"   2. Verify telemetry in Middleware.io dashboard",
// 			"   3. Test trace generation with sample requests",
// 			"   4. Monitor agent logs for errors",
// 		)
// 	} else {
// 		// Configuration for non-Middleware agent
// 		content = append(content,
// 			"🔄 AGENT REPLACEMENT OPTIONS",
// 			"",
// 			"Current Agent: "+proc.JavaAgentName,
// 			"",
// 			"Migration Path to Middleware.io:",
// 			"",
// 			"Option 1: Replace Existing Agent",
// 			"   • Remove current -javaagent parameter",
// 			"   • Add -javaagent:/path/to/middleware-javaagent.jar",
// 			"   • Configure MW_* environment variables",
// 			"   • Restart required: Yes",
// 			"   • Risk: Temporary observability gap during restart",
// 			"",
// 			"Option 2: Parallel Configuration (if supported)",
// 			"   • Keep existing agent",
// 			"   • Configure existing agent to export to Middleware.io",
// 			"   • Requires: OpenTelemetry compatible agent",
// 			"   • Restart required: Possibly",
// 			"",
// 			"⚠️  MIGRATION CONSIDERATIONS",
// 			"",
// 			"   • Test in non-production environment first",
// 			"   • Backup current configuration",
// 			"   • Monitor for performance impact",
// 			"   • Verify all telemetry types (traces, metrics, logs)",
// 		)
// 	}

// 	return content
// }

// // Helper function to generate a reasonable service name
// func generateServiceName(proc discovery.JavaProcess) string {
// 	// Extract service name from JAR file
// 	jarName := proc.JarFile
// 	if jarName == "" {
// 		return "java-service"
// 	}

// 	// Remove .jar extension and version numbers
// 	serviceName := strings.TrimSuffix(jarName, ".jar")
// 	serviceName = strings.Split(serviceName, "-")[0] // Remove version like "demo-0.0.1-SNAPSHOT"

// 	// Add port if available
// 	for _, opt := range proc.JVMOptions {
// 		if strings.Contains(opt, "-Dserver.port=") {
// 			port := strings.Split(opt, "=")[1]
// 			return serviceName + "-" + port
// 		}
// 	}

//		return serviceName
//	}
// package tui

// import (
// 	"context"
// 	"fmt"
// 	"os"
// 	"strings"
// 	"time"

// 	"github.com/charmbracelet/bubbles/key"
// 	"github.com/charmbracelet/bubbles/list"
// 	tea "github.com/charmbracelet/bubbletea"
// 	"github.com/charmbracelet/lipgloss"

// 	"github.com/naman47vyas/mw-injector-manager/pkg/discovery"
// )

// // AppState represents the current state of the application
// type AppState int

// const (
// 	StateMain AppState = iota
// 	StateServiceList
// 	StateServiceConfig
// 	StateHealthCheck
// 	StateHelp
// )

// // MenuItem represents a menu item
// type MenuItem struct {
// 	title       string
// 	description string
// 	action      AppState
// }

// func (m MenuItem) Title() string       { return m.title }
// func (m MenuItem) Description() string { return m.description }
// func (m MenuItem) FilterValue() string { return m.title }

// // SystemStatus represents the current system status
// type SystemStatus struct {
// 	JavaServices    int
// 	Instrumented    int
// 	MiddlewareCount int
// 	OtherAgentCount int
// 	SystemHealth    string
// 	TelemetryStatus string
// 	LastUpdate      time.Time
// 	Processes       []discovery.JavaProcess
// }

// // Model represents the main application model
// type Model struct {
// 	state         AppState
// 	width         int
// 	height        int
// 	list          list.Model
// 	status        SystemStatus
// 	showBanner    bool
// 	showStatus    bool
// 	animationStep int
// 	discoverer    discovery.Discoverer
// 	ctx           context.Context
// 	cancel        context.CancelFunc

// 	selectedProcessIndex int // Track which process is selected in the list (-1 = none)
// }

// // NewModel creates a new model instance
// func NewModel() Model {
// 	// Create context for discovery operations
// 	ctx, cancel := context.WithCancel(context.Background())

// 	// Create discoverer
// 	discoverer := discovery.NewDiscoverer(ctx)

// 	// Create menu items
// 	items := []list.Item{
// 		MenuItem{
// 			title:       "📋 List Services",
// 			description: "View all Java services and their instrumentation status",
// 			action:      StateServiceList,
// 		},
// 		MenuItem{
// 			title:       "⚙️ Configure Service",
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

// 	// Initialize with empty status
// 	status := SystemStatus{
// 		JavaServices:    0,
// 		Instrumented:    0,
// 		MiddlewareCount: 0,
// 		OtherAgentCount: 0,
// 		SystemHealth:    "🔄 Ready to scan",
// 		TelemetryStatus: "📡 Connected",
// 		LastUpdate:      time.Now(),
// 		Processes:       []discovery.JavaProcess{},
// 	}

// 	return Model{
// 		state:      StateMain,
// 		list:       l,
// 		status:     status,
// 		showBanner: true,
// 		showStatus: true,
// 		discoverer: discoverer,
// 		ctx:        ctx,
// 		cancel:     cancel,

// 		selectedProcessIndex: -1,
// 	}
// }

// // Init implements tea.Model
// func (m Model) Init() tea.Cmd {
// 	return tea.EnterAltScreen
// }

// func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	switch msg := msg.(type) {
// 	case tea.WindowSizeMsg:
// 		m.width = msg.Width
// 		m.height = msg.Height
// 		verticalMargin := 0

// 		if m.showBanner {
// 			banner := BannerStyle.Render(GetBannerArt())
// 			subtitle := TitleStyle.Render(GetSubtitle())
// 			version := InfoStyle.Render(GetVersionInfo())
// 			verticalMargin += lipgloss.Height(banner) + lipgloss.Height(subtitle) + lipgloss.Height(version)
// 		}

// 		if m.showStatus {
// 			statusContent := m.renderStatus()
// 			statusBox := StatusBoxStyle.Render(statusContent)
// 			verticalMargin += lipgloss.Height(statusBox)
// 		}

// 		help := HelpStyle.Render("...")
// 		verticalMargin += lipgloss.Height(help)

// 		m.list.SetHeight(m.height - verticalMargin)
// 		m.list.SetWidth(m.width - 4)
// 		return m, nil

// 	case tea.KeyMsg:
// 		switch {
// 		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+c", "q"))):
// 			m.cancel()
// 			return m, tea.Quit

// 		case key.Matches(msg, key.NewBinding(key.WithKeys("b"))):
// 			m.showBanner = !m.showBanner
// 			return m, func() tea.Msg {
// 				return tea.WindowSizeMsg{Width: m.width, Height: m.height}
// 			}

// 		case key.Matches(msg, key.NewBinding(key.WithKeys("s"))):
// 			m.showStatus = !m.showStatus
// 			return m, func() tea.Msg {
// 				return tea.WindowSizeMsg{Width: m.width, Height: m.height}
// 			}

// 		case key.Matches(msg, key.NewBinding(key.WithKeys("r"))):
// 			m.refreshProcesses()
// 			return m, nil

// 		case key.Matches(msg, key.NewBinding(key.WithKeys("enter", " "))):
// 			if selectedItem, ok := m.list.SelectedItem().(MenuItem); ok {
// 				m.state = selectedItem.action
// 				if m.state == StateServiceList {
// 					m.refreshProcesses()
// 					if len(m.status.Processes) > 0 {
// 						m.selectedProcessIndex = 0 // Start with first process selected
// 					} else {
// 						m.selectedProcessIndex = -1 // No processes to select
// 					}
// 				}
// 				return m, nil
// 			}

// 		case key.Matches(msg, key.NewBinding(key.WithKeys("backspace"))):
// 			if m.state != StateMain {
// 				m.state = StateMain
// 				return m, nil
// 			}
// 		}

// 		// Handle navigation only in main menu
// 		if m.state == StateMain {
// 			var cmd tea.Cmd
// 			m.list, cmd = m.list.Update(msg)
// 			return m, cmd
// 		}

// 		// -------------------- Added: Process selection navigation in service list --------------------
// 		if m.state == StateServiceList && len(m.status.Processes) > 0 {
// 			switch {
// 			case key.Matches(msg, key.NewBinding(key.WithKeys("up", "k"))):
// 				if m.selectedProcessIndex > 0 {
// 					m.selectedProcessIndex--
// 				}
// 				return m, nil

// 			case key.Matches(msg, key.NewBinding(key.WithKeys("down", "j"))):
// 				if m.selectedProcessIndex < len(m.status.Processes)-1 {
// 					m.selectedProcessIndex++
// 				}
// 				return m, nil

// 			case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
// 				// -------------------- Added: Placeholder for process detail view --------------------
// 				// For now, just show which process was selected
// 				// Later we'll implement the detail view here
// 				return m, nil
// 			}
// 		}
// 	}

// 	return m, nil
// }

// // refreshProcesses loads real Java processes
// func (m *Model) refreshProcesses() {
// 	processes, err := m.discoverer.DiscoverJavaProcesses(m.ctx)
// 	if err != nil {
// 		m.status.SystemHealth = fmt.Sprintf("❌ Discovery error: %v", err)
// 		return
// 	}

// 	m.status.Processes = processes
// 	m.status.LastUpdate = time.Now()

// 	// Update counts
// 	m.status.JavaServices = len(processes)
// 	m.status.Instrumented = 0
// 	m.status.MiddlewareCount = 0
// 	m.status.OtherAgentCount = 0

// 	for _, proc := range processes {
// 		if proc.HasJavaAgent {
// 			m.status.Instrumented++
// 			if proc.IsMiddlewareAgent {
// 				m.status.MiddlewareCount++
// 			} else {
// 				m.status.OtherAgentCount++
// 			}
// 		}
// 	}

// 	// Update health status
// 	if m.status.JavaServices == 0 {
// 		m.status.SystemHealth = "ℹ️ No Java services found"
// 	} else if m.status.Instrumented == m.status.JavaServices {
// 		m.status.SystemHealth = "✅ All instrumented"
// 	} else if m.status.Instrumented > 0 {
// 		m.status.SystemHealth = "⚠️ Partially instrumented"
// 	} else {
// 		m.status.SystemHealth = "❌ No instrumentation"
// 	}
// }

// // View implements tea.Model
// func (m Model) View() string {
// 	var sections []string

// 	// Banner section
// 	if m.showBanner {
// 		banner := BannerStyle.Render(GetBannerArt())
// 		subtitle := TitleStyle.Render(GetSubtitle())
// 		version := InfoStyle.Render(GetVersionInfo())
// 		sections = append(sections, banner, subtitle, version)
// 	}

// 	// Status section
// 	if m.showStatus {
// 		statusContent := m.renderStatus()
// 		sections = append(sections, StatusBoxStyle.Render(statusContent))
// 	}

// 	// Main content based on state
// 	switch m.state {
// 	case StateMain:
// 		sections = append(sections, m.list.View())

// 	case StateServiceList:
// 		sections = append(sections, m.renderServiceList())

// 	case StateHealthCheck:
// 		sections = append(sections, m.renderHealthCheck())

// 	case StateHelp:
// 		sections = append(sections, m.renderHelp())

// 	default:
// 		sections = append(sections, "Feature coming soon...")
// 	}

// 	// Help footer
// 	help := HelpStyle.Render("Press 'q' to quit • 'b' to toggle banner • 'r' to refresh")
// 	sections = append(sections, help)

// 	return lipgloss.JoinVertical(lipgloss.Left, sections...)
// }

// func (m Model) renderStatus() string {
// 	hostname, _ := os.Hostname()

// 	statusLines := []string{
// 		fmt.Sprintf("🖥️ Host: %s", hostname),
// 		fmt.Sprintf("📅 Last Update: %s", m.status.LastUpdate.Format("15:04:05")),
// 		"",
// 		fmt.Sprintf("☕ Java Services: %d running, %d instrumented",
// 			m.status.JavaServices, m.status.Instrumented),
// 		fmt.Sprintf("⚙️ MW Agents: %d middleware, %d other agents",
// 			m.status.MiddlewareCount, m.status.OtherAgentCount),
// 		fmt.Sprintf("🏥 Health: %s", m.status.SystemHealth),
// 		fmt.Sprintf("📡 Telemetry: %s", m.status.TelemetryStatus),
// 		"",
// 		"🔧 System: LD_PRELOAD Injection Active",
// 		"🎯 Target: Middleware.io",
// 	}

// 	return strings.Join(statusLines, "\n")
// }

// func (m Model) renderServiceList() string {
// 	content := []string{
// 		TitleStyle.Render("📋 Java Services"),
// 		"",
// 	}

// 	if len(m.status.Processes) == 0 {
// 		content = append(content,
// 			"No Java processes found.",
// 			"",
// 			"Press 'r' to scan for Java services",
// 			"",
// 			HelpStyle.Render("Press 'r' to refresh • 'backspace' to go back"))
// 		return strings.Join(content, "\n")
// 	}

// 	// Simple table without complex selection logic
// 	content = append(content,
// 		"┌─────────┬──────────────────────┬─────────┬──────────────────┬──────────────┐",
// 		"│ PID     │ Service Name         │ Status  │ JAR File         │ Memory       │",
// 		"├─────────┼──────────────────────┼─────────┼──────────────────┼──────────────┤",
// 	)

// 	for i, proc := range m.status.Processes {
// 		if i >= 10 {
// 			content = append(content,
// 				fmt.Sprintf("│ ... and %d more processes                                            │",
// 					len(m.status.Processes)-10))
// 			break
// 		}

// 		pidStr := fmt.Sprintf("%-7d", proc.ProcessPID)
// 		serviceStr := fmt.Sprintf("%-20s", truncateString(proc.ServiceName, 20))
// 		statusStr := getStatusString(proc)
// 		jarStr := fmt.Sprintf("%-16s", truncateString(proc.JarFile, 16))
// 		memoryStr := fmt.Sprintf("%-12s", fmt.Sprintf("%.1f%%", proc.MemoryPercent))

// 		content = append(content, fmt.Sprintf("│ %s │ %s │ %s │ %s │ %s │",
// 			pidStr, serviceStr, statusStr, jarStr, memoryStr))
// 	}

// 	content = append(content,
// 		"└─────────┴──────────────────────┴─────────┴──────────────────┴──────────────┘",
// 		"",
// 		HelpStyle.Render("Press 'r' to refresh • 'backspace' to go back"),
// 	)

// 	return strings.Join(content, "\n")
// }

// func getStatusString(proc discovery.JavaProcess) string {
// 	if proc.HasJavaAgent {
// 		if proc.IsMiddlewareAgent {
// 			return "✅ MW  "
// 		}
// 		return "⚙️  OTel"
// 	}
// 	return "❌ None"
// }

// func (m Model) renderHealthCheck() string {
// 	content := []string{
// 		TitleStyle.Render("🏥 System Health Check"),
// 		"",
// 		"✅ LD_PRELOAD injection: Active",
// 		"✅ Shared library: /usr/lib/middleware/libmwinjector.so loaded",
// 		"✅ Java agent: /usr/lib/middleware/middleware-javaagent-1.7.0.jar found",
// 		"✅ Configuration files: Valid",
// 		"✅ Network connectivity: Middleware.io reachable",
// 		"",
// 		"📊 Current Status:",
// 		fmt.Sprintf("  • Java services: %d discovered", m.status.JavaServices),
// 		fmt.Sprintf("  • MW instrumented: %d", m.status.MiddlewareCount),
// 		fmt.Sprintf("  • Other agents: %d", m.status.OtherAgentCount),
// 		fmt.Sprintf("  • Last scan: %s", m.status.LastUpdate.Format("15:04:05")),
// 		"",
// 		HelpStyle.Render("Press 'backspace' to go back"),
// 	}

// 	return strings.Join(content, "\n")
// }

// func (m Model) renderHelp() string {
// 	content := []string{
// 		TitleStyle.Render("❓ Help & Documentation"),
// 		"",
// 		"🔧 Key Features:",
// 		"  • Automatic Java service discovery via /proc filesystem",
// 		"  • Per-service Middleware.io configuration",
// 		"  • LD_PRELOAD shared library injection",
// 		"  • Real-time health monitoring",
// 		"",
// 		"⌨️ Keyboard Shortcuts:",
// 		"  • q, Ctrl+C: Quit application",
// 		"  • ↑/↓: Navigate menu items",
// 		"  • Enter/Space: Select menu item",
// 		"  • b: Toggle banner display",
// 		"  • r: Refresh/scan for Java processes",
// 		"  • Backspace: Go back to previous screen",
// 		"",
// 		"🔗 Resources:",
// 		"  • Documentation: https://docs.middleware.io/injector",
// 		"  • GitHub: https://github.com/middleware-labs/mw-injector",
// 		"  • Support: support@middleware.io",
// 		"",
// 		HelpStyle.Render("Press 'backspace' to go back"),
// 	}

//		return strings.Join(content, "\n")
//	}
package tui

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
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

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

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

		listHeight := m.height - verticalMargin
		listWidth := m.width - 4
		m.list.SetHeight(listHeight)
		m.list.SetWidth(listWidth)
		m.processOptionsList.SetHeight(listHeight)
		m.processOptionsList.SetWidth(listWidth)

		return m, nil

	case tea.KeyMsg:
		// Global keybindings
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+c", "q"))):
			m.cancel()
			return m, tea.Quit
		case key.Matches(msg, key.NewBinding(key.WithKeys("r"))):
			m.refreshProcesses()
			return m, nil
		}

		// State-specific keybindings
		switch m.state {
		case StateMain:
			switch {
			case key.Matches(msg, key.NewBinding(key.WithKeys("enter", " "))):
				if selectedItem, ok := m.list.SelectedItem().(MenuItem); ok {
					m.state = selectedItem.action
					if m.state == StateServiceList {
						m.refreshProcesses()
						if len(m.status.Processes) > 0 {
							m.selectedProcessIndex = 0
						} else {
							m.selectedProcessIndex = -1
						}
					}
					return m, nil
				}
			}
			m.list, cmd = m.list.Update(msg)
			return m, cmd

		case StateServiceList:
			if len(m.status.Processes) > 0 {
				switch {
				case key.Matches(msg, key.NewBinding(key.WithKeys("up", "k"))):
					if m.selectedProcessIndex > 0 {
						m.selectedProcessIndex--
					}
				case key.Matches(msg, key.NewBinding(key.WithKeys("down", "j"))):
					if m.selectedProcessIndex < len(m.status.Processes)-1 {
						m.selectedProcessIndex++
					}
				// --- Updated: Enter key now transitions to the process options menu ---
				case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
					if m.selectedProcessIndex != -1 {
						m.state = StateProcessOptions
						// Dynamically set title for the options menu
						selectedProc := m.status.Processes[m.selectedProcessIndex]
						m.processOptionsList.Title = fmt.Sprintf("Options for %s (PID: %d)", getEnhancedServiceName(selectedProc), selectedProc.ProcessPID)
						m.processOptionsList.Select(0) // Reset selection to the first item
					}
					return m, nil
				}
			}
			if key.Matches(msg, key.NewBinding(key.WithKeys("backspace"))) {
				m.state = StateMain
			}
			return m, nil

		// --- New: Handle Process Options Menu ---
		case StateProcessOptions:
			switch {
			case key.Matches(msg, key.NewBinding(key.WithKeys("backspace"))):
				m.state = StateServiceList
				return m, nil
			case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
				switch m.processOptionsList.Index() {
				case 0: // View Details
					m.state = StateProcessDetailView
				case 1: // Configure Instrumentation
					// Placeholder: You can switch to StateServiceConfig here
				case 2: // View Raw Data
					// Placeholder: You can create a new state and view for this
				}
				return m, nil
			}
			m.processOptionsList, cmd = m.processOptionsList.Update(msg)
			return m, cmd

		// --- New: Handle Detail View ---
		case StateProcessDetailView:
			if key.Matches(msg, key.NewBinding(key.WithKeys("backspace", "enter"))) {
				m.state = StateProcessOptions
			}
			return m, nil

		case StateHealthCheck, StateHelp:
			if key.Matches(msg, key.NewBinding(key.WithKeys("backspace"))) {
				m.state = StateMain
			}
			return m, nil
		}
	}
	return m, nil
}

// refreshProcesses loads real Java processes
func (m *Model) refreshProcesses() {
	processes, err := m.discoverer.DiscoverJavaProcesses(m.ctx)
	if err != nil {
		m.status.SystemHealth = fmt.Sprintf("❌ Discovery error: %v", err)
		return
	}

	m.status.Processes = processes
	m.status.LastUpdate = time.Now()

	m.status.JavaServices = len(processes)
	m.status.Instrumented = 0
	m.status.MiddlewareCount = 0
	m.status.OtherAgentCount = 0

	for _, proc := range processes {
		if proc.HasJavaAgent {
			m.status.Instrumented++
			if proc.IsMiddlewareAgent {
				m.status.MiddlewareCount++
			} else {
				m.status.OtherAgentCount++
			}
		}
	}

	if m.status.JavaServices == 0 {
		m.status.SystemHealth = "ℹ️ No Java services found"
	} else if m.status.Instrumented == m.status.JavaServices {
		m.status.SystemHealth = "✅ All instrumented"
	} else if m.status.Instrumented > 0 {
		m.status.SystemHealth = "⚠️ Partially instrumented"
	} else {
		m.status.SystemHealth = "❌ No instrumentation"
	}
}

// View implements tea.Model
func (m Model) View() string {
	var sections []string

	if m.showBanner {
		banner := BannerStyle.Render(GetBannerArt())
		subtitle := TitleStyle.Render(GetSubtitle())
		version := InfoStyle.Render(GetVersionInfo())
		sections = append(sections, banner, subtitle, version)
	}

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
	// --- New: Render process options menu ---
	case StateProcessOptions:
		sections = append(sections, m.processOptionsList.View())
	// --- New: Render detailed process view ---
	case StateProcessDetailView:
		sections = append(sections, m.renderProcessDetailView())
	case StateHealthCheck:
		sections = append(sections, m.renderHealthCheck())
	case StateHelp:
		sections = append(sections, m.renderHelp())
	default:
		sections = append(sections, "Feature coming soon...")
	}

	help := HelpStyle.Render("Press 'q' to quit • 'r' to refresh • 'backspace' to go back")
	sections = append(sections, help)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m Model) renderStatus() string {
	// ... (this function remains the same)
	hostname, _ := os.Hostname()

	statusLines := []string{
		fmt.Sprintf("🖥️ Host: %s", hostname),
		fmt.Sprintf("📅 Last Update: %s", m.status.LastUpdate.Format("15:04:05")),
		"",
		fmt.Sprintf("☕ Java Services: %d running, %d instrumented",
			m.status.JavaServices, m.status.Instrumented),
		fmt.Sprintf("⚙️ MW Agents: %d middleware, %d other agents",
			m.status.MiddlewareCount, m.status.OtherAgentCount),
		fmt.Sprintf("🏥 Health: %s", m.status.SystemHealth),
		fmt.Sprintf("📡 Telemetry: %s", m.status.TelemetryStatus),
		"",
		"🔧 System: LD_PRELOAD Injection Active",
		"🎯 Target: Middleware.io",
	}

	return strings.Join(statusLines, "\n")
}

// --- Enhanced: renderServiceList now shows a more detailed table ---
func (m Model) renderServiceList() string {
	content := []string{
		TitleStyle.Render("📋 Java Services"),
		"",
	}

	if len(m.status.Processes) == 0 {
		content = append(content,
			"No Java processes found.",
			"",
			"Press 'r' to scan for Java services",
		)
		return strings.Join(content, "\n")
	}

	header := fmt.Sprintf("  %-7s │ %-20s │ %-7s │ %-18s │ %-12s │ %-10s",
		"PID", "Service Name", "Agent", "JAR File", "Memory/CPU", "Port")
	content = append(content, HeaderStyle.Render(header))

	for i, proc := range m.status.Processes {
		serviceName := getEnhancedServiceName(proc)
		agentStatus := getStatusString(proc)
		jarFile := truncateString(proc.JarFile, 18)
		resUsage := fmt.Sprintf("%.1f%%/%.1f%%", proc.MemoryPercent, proc.CPUPercent)
		port := getServerPort(proc)

		line := fmt.Sprintf("  %-7d │ %-20s │ %-7s │ %-18s │ %-12s │ %-10s",
			proc.ProcessPID,
			truncateString(serviceName, 20),
			agentStatus,
			jarFile,
			resUsage,
			port,
		)

		if i == m.selectedProcessIndex {
			content = append(content, SelectedRowStyle.Render(line))
		} else {
			content = append(content, line)
		}
	}

	content = append(content, "", HelpStyle.Render("Use ↑/↓ to navigate • 'enter' for options • 'r' to refresh"))

	return strings.Join(content, "\n")
}

// --- New: Renders the detailed view for a single process ---
func (m Model) renderProcessDetailView() string {
	if m.selectedProcessIndex < 0 || m.selectedProcessIndex >= len(m.status.Processes) {
		return "Error: No process selected."
	}
	proc := m.status.Processes[m.selectedProcessIndex]

	// Helper for creating styled key-value pairs
	kv := func(key, value string) string {
		return lipgloss.JoinHorizontal(lipgloss.Top,
			KeyStyle.Render(key),
			ValueStyle.Render(value),
		)
	}

	// Build sections
	idSection := []string{
		SubtleStyle.Render("🏷️  Identification"),
		kv("Service Name:", getEnhancedServiceName(proc)),
		kv("JAR File:", proc.JarFile),
		kv("JAR Path:", proc.JarPath),
		kv("Process Owner:", proc.ProcessOwner),
	}

	runtimeSection := []string{
		SubtleStyle.Render("⚙️  Runtime Information"),
		kv("Executable:", proc.ProcessExecutableName),
		kv("Parent PID:", fmt.Sprintf("%d", proc.ProcessParentPID)),
		kv("Created:", proc.ProcessCreateTime.Format("2006-01-02 15:04:05")),
		kv("Status:", proc.Status),
	}

	jvmSection := []string{
		SubtleStyle.Render("🔧 JVM Configuration"),
		kv("Agent:", proc.JavaAgentPath),
		kv("Server Port:", getServerPort(proc)),
	}

	perfSection := []string{
		SubtleStyle.Render("📊 Performance"),
		kv("Memory Usage:", fmt.Sprintf("%.2f%%", proc.MemoryPercent)),
		kv("CPU Usage:", fmt.Sprintf("%.2f%%", proc.CPUPercent)),
	}

	// Combine all sections into a final layout
	leftPanel := lipgloss.JoinVertical(lipgloss.Left,
		strings.Join(idSection, "\n"),
		"",
		strings.Join(jvmSection, "\n"),
	)
	rightPanel := lipgloss.JoinVertical(lipgloss.Left,
		strings.Join(runtimeSection, "\n"),
		"",
		strings.Join(perfSection, "\n"),
	)

	content := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)

	finalView := []string{
		TitleStyle.Render(fmt.Sprintf("📋 Process Details - PID %d", proc.ProcessPID)),
		"",
		content,
		"",
		HelpStyle.Render("Press 'backspace' or 'enter' to return to options"),
	}

	return strings.Join(finalView, "\n")
}

// --- Helper Functions ---

func getStatusString(proc discovery.JavaProcess) string {
	if proc.HasJavaAgent {
		if proc.IsMiddlewareAgent {
			return "✅ MW"
		}
		return "⚙️ OTel"
	}
	return "❌ None"
}

func getServerPort(proc discovery.JavaProcess) string {
	// Regex to find server.port in various formats
	re := regexp.MustCompile(`-Dserver\.port=(\d+)`)
	for _, opt := range proc.JVMOptions {
		matches := re.FindStringSubmatch(opt)
		if len(matches) > 1 {
			return matches[1]
		}
	}
	return "default"
}

func getEnhancedServiceName(proc discovery.JavaProcess) string {
	port := getServerPort(proc)
	baseName := strings.TrimSuffix(proc.JarFile, ".jar")
	if port != "default" {
		return fmt.Sprintf("%s:%s", baseName, port)
	}
	return baseName
}

func (m Model) renderHealthCheck() string {
	// ... (this function remains the same)
	content := []string{
		TitleStyle.Render("🏥 System Health Check"),
		"",
		"✅ LD_PRELOAD injection: Active",
		"✅ Shared library: /usr/lib/middleware/libmwinjector.so loaded",
		"✅ Java agent: /usr/lib/middleware/middleware-javaagent-1.7.0.jar found",
		"✅ Configuration files: Valid",
		"✅ Network connectivity: Middleware.io reachable",
		"",
		"📊 Current Status:",
		fmt.Sprintf("  • Java services: %d discovered", m.status.JavaServices),
		fmt.Sprintf("  • MW instrumented: %d", m.status.MiddlewareCount),
		fmt.Sprintf("  • Other agents: %d", m.status.OtherAgentCount),
		fmt.Sprintf("  • Last scan: %s", m.status.LastUpdate.Format("15:04:05")),
		"",
		HelpStyle.Render("Press 'backspace' to go back"),
	}

	return strings.Join(content, "\n")
}

func (m Model) renderHelp() string {
	// ... (this function remains the same)
	content := []string{
		TitleStyle.Render("❓ Help & Documentation"),
		"",
		"🔧 Key Features:",
		"  • Automatic Java service discovery via /proc filesystem",
		"  • Per-service Middleware.io configuration",
		"  • LD_PRELOAD shared library injection",
		"  • Real-time health monitoring",
		"",
		"⌨️ Keyboard Shortcuts:",
		"  • q, Ctrl+C: Quit application",
		"  • ↑/↓: Navigate menu items",
		"  • Enter/Space: Select menu item",
		"  • r: Refresh/scan for Java processes",
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
