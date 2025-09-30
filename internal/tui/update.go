package tui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

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
			if m.state != StateConfigureInstrumentation {
				m.refreshProcesses()
				return m, nil
			}
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
				return m, tea.ClearScreen
			case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
				switch m.processOptionsList.Index() {
				case 0: // View Details
					m.state = StateProcessDetailView
					return m, tea.ClearScreen
				case 1: // Configure Instrumentation
					// Placeholder: You can switch to StateServiceConfig here
					m.initializeConfigForm()
					m.state = StateConfigureInstrumentation
					return m, tea.ClearScreen
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

		case StateConfigureInstrumentation:
			switch {
			case key.Matches(msg, key.NewBinding(key.WithKeys("esc"))):
				m.state = StateProcessOptions
				return m, tea.ClearScreen

			case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
				// Validate and save
				if err := m.configToEdit.Validate(); err != nil {
					// Show error (you could add an error field to Model)
					return m, nil
				}
				if err := m.saveConfiguration(); err != nil {
					// Show error
					return m, nil
				}
				m.state = StateProcessOptions
				return m, nil

			case key.Matches(msg, key.NewBinding(key.WithKeys("right"))):
				// Switch to next section
				sections := []string{"required", "features", "advanced", "preview"}
				for i, s := range sections {
					if s == m.configFormSection && i < len(sections)-1 {
						m.configFormSection = sections[i+1]
						break
					}
				}
				return m, tea.ClearScreen

			case key.Matches(msg, key.NewBinding(key.WithKeys("left"))):
				// Switch to previous section
				sections := []string{"required", "features", "advanced", "preview"}
				for i, s := range sections {
					if s == m.configFormSection && i > 0 {
						m.configFormSection = sections[i-1]
						break
					}
				}
				return m, tea.ClearScreen

			case key.Matches(msg, key.NewBinding(key.WithKeys("tab"))):
				// Section-aware tab navigation
				switch m.configFormSection {
				case "required":
					// Only cycle through 3 text inputs (0-2)
					m.configFormFocusIndex = (m.configFormFocusIndex + 1) % 3
				case "features":
					// Cycle through 5 toggles (3-7)
					if m.configFormFocusIndex < 3 || m.configFormFocusIndex >= 8 {
						m.configFormFocusIndex = 3
					} else {
						m.configFormFocusIndex++
						if m.configFormFocusIndex >= 8 {
							m.configFormFocusIndex = 3
						}
					}
				case "advanced":
					// Cycle through 4 text inputs (3-6)
					if m.configFormFocusIndex < 3 || m.configFormFocusIndex >= 7 {
						m.configFormFocusIndex = 3
					} else {
						m.configFormFocusIndex++
						if m.configFormFocusIndex >= 7 {
							m.configFormFocusIndex = 3
						}
					}
				case "preview":
					// No inputs in preview, do nothing
					return m, nil
				}
				m.focusConfigInput()
				return m, nil

			case key.Matches(msg, key.NewBinding(key.WithKeys("shift+tab"))):
				// Section-aware reverse tab navigation
				switch m.configFormSection {
				case "required":
					m.configFormFocusIndex--
					if m.configFormFocusIndex < 0 {
						m.configFormFocusIndex = 2
					}
				case "features":
					if m.configFormFocusIndex <= 3 || m.configFormFocusIndex > 7 {
						m.configFormFocusIndex = 7
					} else {
						m.configFormFocusIndex--
					}
				case "advanced":
					if m.configFormFocusIndex <= 3 || m.configFormFocusIndex > 6 {
						m.configFormFocusIndex = 6
					} else {
						m.configFormFocusIndex--
					}
				case "preview":
					return m, nil
				}
				m.focusConfigInput()
				return m, nil
			case key.Matches(msg, key.NewBinding(key.WithKeys(" "))):
				// Toggle boolean fields - only in features section
				if m.configFormSection == "features" && m.configFormFocusIndex >= 3 && m.configFormFocusIndex <= 7 {
					toggleIndex := m.configFormFocusIndex - 3
					switch toggleIndex {
					case 0:
						m.configToEdit.MWAPMCollectTraces = !m.configToEdit.MWAPMCollectTraces
					case 1:
						m.configToEdit.MWAPMCollectMetrics = !m.configToEdit.MWAPMCollectMetrics
					case 2:
						m.configToEdit.MWAPMCollectLogs = !m.configToEdit.MWAPMCollectLogs
					case 3:
						m.configToEdit.MWAPMCollectProfiling = !m.configToEdit.MWAPMCollectProfiling
					case 4:
						m.configToEdit.MWEnableGzip = !m.configToEdit.MWEnableGzip
					}
				}
				return m, nil
			}

			// Update text inputs
			if m.configFormFocusIndex < len(m.configFormInputs) {
				m.configFormInputs[m.configFormFocusIndex], cmd = m.configFormInputs[m.configFormFocusIndex].Update(msg)

				// Sync values to config
				m.syncInputsToConfig()
			}
			return m, cmd
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
