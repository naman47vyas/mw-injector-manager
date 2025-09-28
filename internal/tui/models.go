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
