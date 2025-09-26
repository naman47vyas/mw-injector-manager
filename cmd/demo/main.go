// Package main provides a simple demo of the process discovery engine
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/naman47vyas/mw-injector-manager/pkg/discovery"
)

func main() {
	var (
		currentUserOnly  = flag.Bool("user", false, "Show only current user processes")
		instrumentedOnly = flag.Bool("instrumented", false, "Show only instrumented processes")
		middlewareOnly   = flag.Bool("mw", false, "Show only Middleware processes")
		verbose          = flag.Bool("v", true, "Verbose output with detailed information")
		timeout          = flag.Duration("timeout", 10*time.Second, "Discovery timeout")
		jsonOutput       = flag.Bool("json", false, "Output in JSON format")
	)
	flag.Parse()

	fmt.Println("🔍 MW Injector - Java Process Discovery Demo")
	fmt.Println("=" + strings.Repeat("=", 47))
	fmt.Println()

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	var processes []discovery.JavaProcess
	var err error

	// Choose discovery method based on flags
	switch {
	case *middlewareOnly:
		fmt.Println("🔍 Searching for Middleware.io instrumented processes...")
		processes, err = discovery.FindMiddlewareProcesses(ctx)
	case *instrumentedOnly:
		fmt.Println("🔍 Searching for instrumented Java processes...")
		processes, err = discovery.FindInstrumentedProcesses(ctx)
	case *currentUserOnly:
		fmt.Println("🔍 Searching for current user Java processes...")
		processes, err = discovery.FindCurrentUserJavaProcesses(ctx)
	default:
		fmt.Println("🔍 Searching for all Java processes...")
		processes, err = discovery.FindAllJavaProcesses(ctx)
	}

	if err != nil {
		log.Fatalf("❌ Discovery failed: %v", err)
	}

	fmt.Printf("✅ Found %d Java process(es)\n\n", len(processes))

	if len(processes) == 0 {
		fmt.Println("No Java processes found.")
		fmt.Println()
		fmt.Println("💡 Tips:")
		fmt.Println("   • Make sure Java processes are running")
		fmt.Println("   • Try running with different filters")
		fmt.Println("   • Check permissions (some processes may be hidden)")
		return
	}

	if *jsonOutput {
		outputJSON(processes)
		return
	}

	// Sort processes by PID for consistent output
	sort.Slice(processes, func(i, j int) bool {
		return processes[i].ProcessPID < processes[j].ProcessPID
	})

	if *verbose {
		outputVerbose(processes)
	} else {
		outputTable(processes)
	}

	// Summary
	fmt.Println()
	printSummary(processes)
}

// outputTable prints processes in a table format
func outputTable(processes []discovery.JavaProcess) {
	fmt.Printf("%-8s %-20s %-15s %-25s %-10s\n", "PID", "Service", "Status", "JAR/Class", "Memory")
	fmt.Println(strings.Repeat("-", 80))

	for _, proc := range processes {
		serviceName := proc.ServiceName
		if serviceName == "" {
			serviceName = "-"
		}

		jarOrClass := proc.JarFile
		if jarOrClass == "" {
			jarOrClass = proc.MainClass
		}
		if jarOrClass == "" {
			jarOrClass = "-"
		}

		memoryStr := fmt.Sprintf("%.1f%%", proc.MemoryPercent)

		fmt.Printf("%-8d %-20s %-15s %-25s %-10s\n",
			proc.ProcessPID,
			truncate(serviceName, 20),
			proc.FormatAgentStatus(),
			truncate(jarOrClass, 25),
			memoryStr,
		)
	}
}

// outputVerbose prints detailed information about each process
func outputVerbose(processes []discovery.JavaProcess) {
	for i, proc := range processes {
		if i > 0 {
			fmt.Println(strings.Repeat("-", 60))
		}

		fmt.Printf("🔍 Process #%d\n", i+1)
		fmt.Printf("   PID: %d\n", proc.ProcessPID)
		fmt.Printf("   Owner: %s\n", proc.ProcessOwner)
		fmt.Printf("   Service Name: %s\n", getValueOrDefault(proc.ServiceName, "auto-detected"))
		fmt.Printf("   Status: %s\n", proc.Status)
		fmt.Printf("   Created: %s\n", proc.ProcessCreateTime.Format("2006-01-02 15:04:05"))

		fmt.Println()
		fmt.Printf("   ☕ Java Information:\n")
		fmt.Printf("      Executable: %s\n", proc.ProcessExecutablePath)
		fmt.Printf("      JAR File: %s\n", getValueOrDefault(proc.JarFile, "none"))
		fmt.Printf("      JAR Path: %s\n", getValueOrDefault(proc.JarPath, "none"))
		fmt.Printf("      Main Class: %s\n", getValueOrDefault(proc.MainClass, "none"))
		fmt.Printf("      Runtime Version: %s\n", getValueOrDefault(proc.ProcessRuntimeVersion, "unknown"))

		if len(proc.JVMOptions) > 0 {
			fmt.Printf("      JVM Options: %d options\n", len(proc.JVMOptions))
			for _, opt := range proc.JVMOptions {
				fmt.Printf("         %s\n", opt)
			}
		}

		fmt.Println()
		fmt.Printf("   🔧 Instrumentation:\n")
		fmt.Printf("      Has Agent: %t\n", proc.HasJavaAgent)
		if proc.HasJavaAgent {
			agentInfo := proc.GetAgentInfo()
			fmt.Printf("      Agent Type: %s\n", agentInfo.Type.String())
			fmt.Printf("      Agent Path: %s\n", proc.JavaAgentPath)
			fmt.Printf("      Agent Name: %s\n", proc.JavaAgentName)
			fmt.Printf("      Is Middleware: %t\n", proc.IsMiddlewareAgent)
			if agentInfo.IsServerless {
				fmt.Printf("      Serverless Mode: %t\n", agentInfo.IsServerless)
			}
			if agentInfo.Version != "unknown" {
				fmt.Printf("      Agent Version: %s\n", agentInfo.Version)
			}
		}

		fmt.Println()
		fmt.Printf("   📊 Metrics:\n")
		fmt.Printf("      Memory: %.1f%%\n", proc.MemoryPercent)
		fmt.Printf("      CPU: %.1f%%\n", proc.CPUPercent)

		fmt.Println()
		fmt.Printf("   💻 Command Line:\n")
		fmt.Printf("      %s\n", proc.ProcessCommandLine)
	}
}

// outputJSON prints processes in JSON format
func outputJSON(processes []discovery.JavaProcess) {
	// Simple JSON output for demonstration
	fmt.Println("[")
	for i, proc := range processes {
		if i > 0 {
			fmt.Println(",")
		}
		fmt.Printf(`  {
    "pid": %d,
    "service_name": "%s",
    "owner": "%s",
    "jar_file": "%s",
    "has_agent": %t,
    "is_middleware": %t,
    "agent_status": "%s",
    "memory_percent": %.2f,
    "cpu_percent": %.2f
  }`, proc.ProcessPID, proc.ServiceName, proc.ProcessOwner, proc.JarFile,
			proc.HasJavaAgent, proc.IsMiddlewareAgent, proc.FormatAgentStatus(),
			proc.MemoryPercent, proc.CPUPercent)
	}
	fmt.Println("\n]")
}

// printSummary prints a summary of discovered processes
func printSummary(processes []discovery.JavaProcess) {
	var (
		totalProcesses    = len(processes)
		instrumentedCount = 0
		middlewareCount   = 0
		otelCount         = 0
		otherAgentCount   = 0
		users             = make(map[string]int)
	)

	for _, proc := range processes {
		users[proc.ProcessOwner]++

		if proc.HasInstrumentation() {
			instrumentedCount++

			agentInfo := proc.GetAgentInfo()
			switch agentInfo.Type {
			case discovery.AgentMiddleware:
				middlewareCount++
			case discovery.AgentOpenTelemetry:
				otelCount++
			default:
				otherAgentCount++
			}
		}
	}

	fmt.Println("📊 Discovery Summary:")
	fmt.Printf("   Total Java Processes: %d\n", totalProcesses)
	fmt.Printf("   Instrumented: %d (%.1f%%)\n", instrumentedCount,
		float64(instrumentedCount)/float64(totalProcesses)*100)
	fmt.Printf("   Middleware Agents: %d\n", middlewareCount)
	fmt.Printf("   OpenTelemetry Agents: %d\n", otelCount)
	fmt.Printf("   Other Agents: %d\n", otherAgentCount)
	fmt.Printf("   Unique Users: %d\n", len(users))

	if len(users) > 0 {
		fmt.Println()
		fmt.Println("👥 Processes by User:")
		for user, count := range users {
			fmt.Printf("   %s: %d process(es)\n", user, count)
		}
	}

	// Recommendations
	fmt.Println()
	fmt.Println("💡 Recommendations:")
	if instrumentedCount == 0 {
		fmt.Println("   • No instrumentation detected - consider adding Middleware.io agents")
	} else if middlewareCount == 0 {
		fmt.Println("   • Non-Middleware agents detected - consider switching to MW agents")
	} else if middlewareCount < totalProcesses {
		fmt.Printf("   • %d processes lack MW instrumentation\n", totalProcesses-middlewareCount)
	} else {
		fmt.Println("   • ✅ All processes have Middleware.io instrumentation!")
	}
}

// Helper functions
func truncate(s string, maxLen int) string {
	if s == "" {
		return "-"
	}
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

func getValueOrDefault(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}
