package tui

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/charmbracelet/bubbles/textinput"
	config "github.com/naman47vyas/mw-injector-manager/pkg/config"
	"github.com/naman47vyas/mw-injector-manager/pkg/discovery"
)

// ClearScreen clears the terminal screen
func ClearScreen() {
	fmt.Print("\033[2J\033[H")
}

// PrintSeparator prints a visual separator line
func PrintSeparator(char string, length int) {
	for i := 0; i < length; i++ {
		fmt.Print(char)
	}
	fmt.Println()
}

// PrintHeader prints a styled header
func PrintHeader(text string) {
	fmt.Printf("%s%s%s%s\n", Bold, Cyan, text, Reset)
	PrintSeparator("=", len(text))
}

// Confirm prompts user for confirmation
func Confirm(message string) bool {
	fmt.Printf("%s%s (y/N): %s", Yellow, message, Reset)
	var response string
	fmt.Scanln(&response)
	return response == "y" || response == "Y" || response == "yes"
}

func truncateString(s string, maxWidth int) string {
	if maxWidth <= 0 {
		return ""
	}

	// Handle empty string
	if len(s) == 0 {
		return s
	}

	// If string fits, return as is
	if utf8.RuneCountInString(s) <= maxWidth {
		return s
	}

	// Truncate and add ellipsis
	runes := []rune(s)
	if maxWidth <= 3 {
		return string(runes[:maxWidth])
	}

	return string(runes[:maxWidth-3]) + "..."
}

// Helper function to pad string to exact width
func padToWidth(s string, width int) string {
	runeCount := utf8.RuneCountInString(s)
	if runeCount >= width {
		return truncateString(s, width)
	}

	padding := width - runeCount
	return s + strings.Repeat(" ", padding)
}

// Add this struct definition near your other types
type tableColumn struct {
	Header string
	Width  int
}

// Helper function to calculate visual width of string (handles emojis/unicode)
func visualWidth(s string) int {
	// This is a simplified approach - for production you might want
	// to use a proper Unicode width calculation library
	width := 0
	for _, r := range s {
		if r < 127 {
			width++
		} else {
			// Most emojis and Unicode chars take 2 spaces visually
			width += 2
		}
	}
	return width
}

// Helper function to truncate string considering visual width
func truncateToVisualWidth(s string, maxWidth int) string {
	if maxWidth <= 0 {
		return ""
	}

	if visualWidth(s) <= maxWidth {
		return s
	}

	// Truncate character by character until we fit
	runes := []rune(s)
	result := ""
	currentWidth := 0

	for _, r := range runes {
		charWidth := 1
		if r >= 127 {
			charWidth = 2
		}

		if currentWidth+charWidth+3 > maxWidth { // +3 for "..."
			result += "..."
			break
		}

		result += string(r)
		currentWidth += charWidth
	}

	return result
}

// Helper function to pad string to exact visual width
func padToVisualWidth(s string, width int) string {
	currentWidth := visualWidth(s)
	if currentWidth >= width {
		return truncateToVisualWidth(s, width)
	}

	padding := width - currentWidth
	return s + strings.Repeat(" ", padding)
}

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

// func (m *Model) saveConfiguration() error {
// 	proc := m.status.Processes[m.selectedProcessIndex]

// 	config := config.ProcessConfiguration{
// 		PID:                 proc.ProcessPID,
// 		ServiceName:         proc.ServiceName,
// 		MWAPIKey:            m.configFormInputs[0].Value(),
// 		MWTarget:            m.configFormInputs[1].Value(),
// 		MWServiceName:       m.configFormInputs[2].Value(),
// 		MWLogLevel:          m.configFormInputs[3].Value(),
// 		MWAPMCollectTraces:  m.configToEdit.MWAPMCollectTraces,
// 		MWAPMCollectMetrics: m.configToEdit.MWAPMCollectMetrics,
// 		MWAPMCollectLogs:    m.configToEdit.MWAPMCollectLogs,
// 		CreatedAt:           time.Now(),
// 		UpdatedAt:           time.Now(),
// 	}

// 	// Save to file - you'll implement this next
// 	return saveConfigToFile(config)
// }

// Sync input values back to the config struct
func (m *Model) syncInputsToConfig() {
	m.configToEdit.MWAPIKey = m.configFormInputs[0].Value()
	m.configToEdit.MWTarget = m.configFormInputs[1].Value()
	m.configToEdit.MWServiceName = m.configFormInputs[2].Value()
	m.configToEdit.MWLogLevel = m.configFormInputs[3].Value()
	m.configToEdit.MWCustomResourceAttribute = m.configFormInputs[4].Value()
	m.configToEdit.OtelTracesSampler = m.configFormInputs[5].Value()
	m.configToEdit.OtelTracesSamplerArg = m.configFormInputs[6].Value()
}

//------------------

var configPersistence = config.NewConfigPersistence()

// Update saveConfiguration function
func (m *Model) saveConfiguration() error {
	// Sync current input values to config
	m.syncInputsToConfig()

	// Validate configuration
	if err := m.configToEdit.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Update timestamps
	m.configToEdit.UpdatedAt = time.Now()
	if m.configToEdit.CreatedAt.IsZero() {
		m.configToEdit.CreatedAt = time.Now()
	}

	// Save to file using persistence layer
	if err := configPersistence.Save(m.configToEdit); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	return nil
}

// Add function to load existing configuration
func (m *Model) loadExistingConfiguration() (*config.ProcessConfiguration, error) {
	proc := m.status.Processes[m.selectedProcessIndex]

	// Try to load existing config
	existingConfig, err := configPersistence.Load(proc.ServiceName)
	if err != nil {
		// No existing config, return default
		return nil, nil
	}

	return existingConfig, nil
}

// Update initializeConfigForm to load existing config if available
func (m *Model) initializeConfigForm() {
	proc := m.status.Processes[m.selectedProcessIndex]

	// Try to load existing config first
	existingConfig, err := m.loadExistingConfiguration()
	if err == nil && existingConfig != nil {
		// Use existing config
		m.configToEdit = existingConfig
	} else {
		// Create default config
		defaultConfig := config.DefaultConfiguration()
		defaultConfig.PID = proc.ProcessPID
		defaultConfig.ServiceName = proc.ServiceName
		defaultConfig.MWServiceName = proc.ServiceName
		defaultConfig.JavaAgentPath = proc.JavaAgentPath
		m.configToEdit = &defaultConfig
	}

	// Initialize text inputs (7 inputs total)
	m.configFormInputs = make([]textinput.Model, 7)

	inputs := []struct {
		placeholder string
		value       string
	}{
		{"MW_API_KEY", m.configToEdit.MWAPIKey},
		{"MW_TARGET", m.configToEdit.MWTarget},
		{"MW_SERVICE_NAME", m.configToEdit.MWServiceName},
		{"MW_LOG_LEVEL", m.configToEdit.MWLogLevel},
		{"MW_CUSTOM_RESOURCE_ATTRIBUTE", m.configToEdit.MWCustomResourceAttribute},
		{"OTEL_TRACES_SAMPLER", m.configToEdit.OtelTracesSampler},
		{"OTEL_TRACES_SAMPLER_ARG", m.configToEdit.OtelTracesSamplerArg},
	}

	for i, input := range inputs {
		ti := textinput.New()
		ti.Placeholder = input.placeholder
		ti.SetValue(input.value)
		ti.CharLimit = 256

		if i == 0 {
			ti.Focus()
		}

		m.configFormInputs[i] = ti
	}

	m.configFormFocusIndex = 0
	m.configFormSection = "required"
}

func (m *Model) focusConfigInput() {
	for i := range m.configFormInputs {
		if i == m.configFormFocusIndex {
			m.configFormInputs[i].Focus()
		} else {
			m.configFormInputs[i].Blur()
		}
	}
}
