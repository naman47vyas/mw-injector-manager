package tui

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
type SystemStatus struct {
	JavaServices    int
	Instrumented    int
	ActiveConfigs   int
	PendingConfigs  int
	SystemHealth    string
	TelemetryStatus string
	LastUpdate      time.Time
}

// Model represents the main application model
type Model struct {
	state         AppState
	width         int
	height        int
	list          list.Model
	status        SystemStatus
	showBanner    bool
	showStatus    bool
	animationStep int
}

// NewModel creates a new model instance
func NewModel() Model {
	// Create menu items
	items := []list.Item{
		MenuItem{
			title:       "📋 List Services",
			description: "View all Java services and their instrumentation status",
			action:      StateServiceList,
		},
		MenuItem{
			title:       "⚙️  Configure Service",
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

	// Initialize status
	status := SystemStatus{
		JavaServices:    3,
		Instrumented:    2,
		ActiveConfigs:   2,
		PendingConfigs:  1,
		SystemHealth:    "✅ Healthy",
		TelemetryStatus: "📈 Flowing",
		LastUpdate:      time.Now(),
	}

	return Model{
		state:      StateMain,
		list:       l,
		status:     status,
		showBanner: true,
		showStatus: true,
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

// renderStatus renders the system status section
func (m Model) renderStatus() string {
	hostname, _ := os.Hostname()

	statusLines := []string{
		fmt.Sprintf("🖥️  Host: %s", hostname),
		fmt.Sprintf("📅 Last Update: %s", m.status.LastUpdate.Format("15:04:05")),
		"",
		fmt.Sprintf("☕ Java Services: %d running, %d instrumented",
			m.status.JavaServices, m.status.Instrumented),
		fmt.Sprintf("⚙️  MW Configs: %d active, %d pending",
			m.status.ActiveConfigs, m.status.PendingConfigs),
		fmt.Sprintf("🏥 Health: %s", m.status.SystemHealth),
		fmt.Sprintf("📡 Telemetry: %s", m.status.TelemetryStatus),
		"",
		"🔧 System: LD_PRELOAD Injection Active",
		"🎯 Target: Middleware.io",
	}

	return strings.Join(statusLines, "\n")
}

// renderServiceList renders the service list view
func (m Model) renderServiceList() string {
	content := []string{
		TitleStyle.Render("📋 Java Services"),
		"",
		"┌─────────┬──────────────────────┬─────────┬──────────────────┬──────────────┐",
		"│ PID     │ Service Name         │ Status  │ MW Config        │ Last Seen    │",
		"├─────────┼──────────────────────┼─────────┼──────────────────┼──────────────┤",
		"│ 1234    │ microservice-a.jar   │ ✅ MW   │ user-auth        │ 2 min ago    │",
		"│ 5678    │ microservice-b.jar   │ ❌ None │ -                │ 5 min ago    │",
		"│ 9012    │ legacy-app.jar       │ ✅ OTel │ default          │ 1 min ago    │",
		"└─────────┴──────────────────────┴─────────┴──────────────────┴──────────────┘",
		"",
		HelpStyle.Render("Press 'enter' to configure service • 'backspace' to go back"),
	}

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
