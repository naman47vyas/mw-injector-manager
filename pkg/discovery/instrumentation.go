package discovery

import (
	"path/filepath"
	"regexp"
	"strings"
)

// detectInstrumentation detects Java agents and instrumentation in the process
func (d *discoverer) detectInstrumentation(javaProc *JavaProcess, cmdArgs []string) {
	// Reset instrumentation flags
	javaProc.HasJavaAgent = false
	javaProc.IsMiddlewareAgent = false
	javaProc.JavaAgentPath = ""
	javaProc.JavaAgentName = ""

	// Look for -javaagent arguments
	for _, arg := range cmdArgs {
		if strings.HasPrefix(arg, "-javaagent:") {
			agentPath := strings.TrimPrefix(arg, "-javaagent:")

			// Handle agent arguments (agent.jar=arg1,arg2)
			agentPath = d.extractAgentPath(agentPath)

			javaProc.HasJavaAgent = true
			javaProc.JavaAgentPath = agentPath
			javaProc.JavaAgentName = filepath.Base(agentPath)

			// Detect agent type
			agentType := d.detectAgentType(agentPath)
			if agentType == AgentMiddleware {
				javaProc.IsMiddlewareAgent = true
			}

			// For now, we'll use the first agent found
			// In future versions, we could support multiple agents
			break
		}
	}

	// Also check JAVA_TOOL_OPTIONS environment variable
	// This would require reading /proc/PID/environ which needs additional permissions
	// For now, we'll focus on command line detection
}

// extractAgentPath extracts the agent JAR path from javaagent argument
func (d *discoverer) extractAgentPath(agentArg string) string {
	// Handle cases like: agent.jar=option1,option2
	if idx := strings.Index(agentArg, "="); idx != -1 {
		return agentArg[:idx]
	}
	return agentArg
}

// detectAgentType determines what type of agent is being used
func (d *discoverer) detectAgentType(agentPath string) AgentType {
	agentPathLower := strings.ToLower(agentPath)
	agentName := strings.ToLower(filepath.Base(agentPath))

	// Middleware agent patterns
	middlewarePatterns := []*regexp.Regexp{
		regexp.MustCompile(`middleware`),
		regexp.MustCompile(`mw-`),
		regexp.MustCompile(`mw\.jar`),
		regexp.MustCompile(`middleware-javaagent`),
		regexp.MustCompile(`mw-javaagent`),
	}

	for _, pattern := range middlewarePatterns {
		if pattern.MatchString(agentPathLower) || pattern.MatchString(agentName) {
			return AgentMiddleware
		}
	}

	// OpenTelemetry agent patterns
	otelPatterns := []*regexp.Regexp{
		regexp.MustCompile(`opentelemetry`),
		regexp.MustCompile(`otel`),
		regexp.MustCompile(`opentelemetry-javaagent`),
		regexp.MustCompile(`otel-javaagent`),
	}

	for _, pattern := range otelPatterns {
		if pattern.MatchString(agentPathLower) || pattern.MatchString(agentName) {
			return AgentOpenTelemetry
		}
	}

	// Other well-known agents
	otherAgentPatterns := []*regexp.Regexp{
		regexp.MustCompile(`newrelic`),
		regexp.MustCompile(`datadog`),
		regexp.MustCompile(`appdynamics`),
		regexp.MustCompile(`dynatrace`),
		regexp.MustCompile(`elastic`),
		regexp.MustCompile(`jaeger`),
		regexp.MustCompile(`zipkin`),
		regexp.MustCompile(`skywalking`),
		regexp.MustCompile(`pinpoint`),
	}

	for _, pattern := range otherAgentPatterns {
		if pattern.MatchString(agentPathLower) || pattern.MatchString(agentName) {
			return AgentOther
		}
	}

	// If we found an agent but couldn't classify it
	return AgentOther
}

// isMiddlewareServerlessAgent checks if the agent is a Middleware serverless agent
func (d *discoverer) isMiddlewareServerlessAgent(agentPath string) bool {
	agentPathLower := strings.ToLower(agentPath)
	serverlessPatterns := []*regexp.Regexp{
		regexp.MustCompile(`serverless`),
		regexp.MustCompile(`lambda`),
		regexp.MustCompile(`aws`),
		regexp.MustCompile(`wo\.healthcheck`), // "without healthcheck" indicator
	}

	for _, pattern := range serverlessPatterns {
		if pattern.MatchString(agentPathLower) {
			return true
		}
	}

	return false
}

// extractMiddlewareAgentVersion attempts to extract version from Middleware agent path
func (d *discoverer) extractMiddlewareAgentVersion(agentPath string) string {
	// Common version patterns in agent paths
	versionPatterns := []*regexp.Regexp{
		regexp.MustCompile(`-(\d+\.\d+\.\d+)`),     // -1.7.0
		regexp.MustCompile(`_(\d+\.\d+\.\d+)`),     // _1.7.0
		regexp.MustCompile(`(\d+\.\d+\.\d+)\.jar`), // 1.7.0.jar
		regexp.MustCompile(`v(\d+\.\d+\.\d+)`),     // v1.7.0
	}

	for _, pattern := range versionPatterns {
		matches := pattern.FindStringSubmatch(agentPath)
		if len(matches) > 1 {
			return matches[1]
		}
	}

	return "unknown"
}

// detectEnvironmentVariables detects MW_* environment variables in the process
// Note: This requires reading /proc/PID/environ which may need elevated permissions
func (d *discoverer) detectEnvironmentVariables(javaProc *JavaProcess) {
	// This is a placeholder for environment variable detection
	// Implementation would require reading /proc/PID/environ file
	// which often requires root privileges or same user ownership

	// For now, we'll skip this to avoid permission issues
	// In a future version, we could:
	// 1. Try to read environ file with proper error handling
	// 2. Fall back to checking common locations for MW config files
	// 3. Use system calls if available
}

// GetAgentInfo returns detailed information about the detected agent
func (jp *JavaProcess) GetAgentInfo() *AgentInfo {
	if !jp.HasJavaAgent {
		return &AgentInfo{
			Type: AgentNone,
		}
	}

	agentType := AgentOther
	version := "unknown"
	isServerless := false

	// Determine agent type
	if jp.IsMiddlewareAgent {
		agentType = AgentMiddleware

		// Extract version and serverless detection for Middleware agents
		d := &discoverer{} // Create temporary discoverer for utility methods
		version = d.extractMiddlewareAgentVersion(jp.JavaAgentPath)
		isServerless = d.isMiddlewareServerlessAgent(jp.JavaAgentPath)
	} else {
		// Check for other agent types
		d := &discoverer{}
		detectedType := d.detectAgentType(jp.JavaAgentPath)
		if detectedType != AgentOther {
			agentType = detectedType
		}
	}

	return &AgentInfo{
		Type:         agentType,
		Path:         jp.JavaAgentPath,
		Name:         jp.JavaAgentName,
		Version:      version,
		IsServerless: isServerless,
	}
}

// FormatAgentStatus returns a human-readable agent status string
func (jp *JavaProcess) FormatAgentStatus() string {
	if !jp.HasJavaAgent {
		return "❌ None"
	}

	agentInfo := jp.GetAgentInfo()

	switch agentInfo.Type {
	case AgentMiddleware:
		if agentInfo.IsServerless {
			return "✅ MW (Serverless)"
		}
		return "✅ MW"
	case AgentOpenTelemetry:
		return "✅ OTel"
	case AgentOther:
		return "✅ Other"
	default:
		return "⚠️ Unknown"
	}
}

// HasInstrumentation checks if the process has any form of instrumentation
func (jp *JavaProcess) HasInstrumentation() bool {
	return jp.HasJavaAgent
}

// HasMiddlewareInstrumentation checks specifically for Middleware instrumentation
func (jp *JavaProcess) HasMiddlewareInstrumentation() bool {
	return jp.IsMiddlewareAgent
}

// IsServerless checks if this appears to be a serverless deployment
func (jp *JavaProcess) IsServerless() bool {
	if !jp.HasJavaAgent {
		return false
	}

	agentInfo := jp.GetAgentInfo()
	return agentInfo.IsServerless
}
