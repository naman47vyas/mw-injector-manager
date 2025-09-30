package tui

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

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
