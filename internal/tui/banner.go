// package tui

// import (
// 	"fmt"
// 	"os"
// 	"runtime"
// 	"time"
// )

// const (
// 	Version = "1.0.0"
// 	Build   = "20250923"
// )

// // Banner displays the welcome banner for MW Injector Manager
// func Banner() {
// 	// Clear screen (optional)
// 	// fmt.Print("\033[2J\033[H")

// 	fmt.Printf("%s%s", Cyan, Bold)
// 	fmt.Sprintf(`
//                                             ███                  ███   ████████           ███████████    █████
//                                             ░░░                  ░░░   ███░░░░███         ░█░░░███░░░█  ███░░░███
//  █████████████   █████ ███ █████            ████  ████████       █████░░░    ░███  ██████ ░   ░███  ░  ███   ░░███ ████████
// ░░███░░███░░███ ░░███ ░███░░███  ██████████░░███ ░░███░░███     ░░███    ██████░  ███░░███    ░███    ░███    ░███░░███░░███
//  ░███ ░███ ░███  ░███ ░███ ░███ ░░░░░░░░░░  ░███  ░███ ░███      ░███   ░░░░░░███░███ ░░░     ░███    ░███    ░███ ░███ ░░░
//  ░███ ░███ ░███  ░░███████████              ░███  ░███ ░███      ░███  ███   ░███░███  ███    ░███    ░░███   ███  ░███
//  █████░███ █████  ░░████░████               █████ ████ █████     ░███ ░░████████ ░░██████     █████    ░░░█████░   █████
// ░░░░░ ░░░ ░░░░░    ░░░░ ░░░░               ░░░░░ ░░░░ ░░░░░      ░███  ░░░░░░░░   ░░░░░░     ░░░░░       ░░░░░░   ░░░░░
//                                                              ███ ░███
//                                                             ░░██████
//                                                              ░░░░░░                                                         `)
// fmt.Print(Reset)

// 	// Subtitle with gradient effect
// 	fmt.Printf("%s%s", Purple, Bold)
// 	fmt.Println("                    🚀 OpenTelemetry Fleet Management System 🚀")
// 	fmt.Print(Reset)

// 	// Version and system info
// 	fmt.Printf("%s%s", Yellow, Dim)
// 	fmt.Printf("                          Version %s (Build %s)\n", Version, Build)
// 	fmt.Printf("                       Running on %s/%s • Go %s\n",
// 		runtime.GOOS, runtime.GOARCH, runtime.Version())
// 	fmt.Print(Reset)

// 	// Separator line with special chars
// 	fmt.Printf("%s", Blue)
// 	fmt.Println("  ╔══════════════════════════════════════════════════════════════════════════════╗")
// 	fmt.Print(Reset)

// 	// Status indicators
// 	hostname, _ := os.Hostname()
// 	currentTime := time.Now().Format("2006-01-02 15:04:05 MST")

// 	fmt.Printf("  %s║%s %s🖥️  Host:%s %s%-20s%s %s📅 %s%s%s %s║%s\n",
// 		Blue, Reset, Green, Reset, White, hostname, Reset,
// 		Green, Reset, currentTime, Green, Blue, Reset)

// 	fmt.Printf("  %s║%s %s🔧 System:%s %sLD_PRELOAD Injection Active%s %s🎯 Target:%s %sMiddleware.io%s %s║%s\n",
// 		Blue, Reset, Green, Reset, Yellow, Reset,
// 		Green, Reset, Cyan, Reset, Blue, Reset)

// 	fmt.Printf("  %s║%s %s📊 Status:%s %sScanning Java Services...%s %s🔍 Mode:%s %sInteractive%s %s  ║%s\n",
// 		Blue, Reset, Green, Reset, Yellow, Reset,
// 		Green, Reset, Purple, Reset, Blue, Reset)

// 	fmt.Printf("%s", Blue)
// 	fmt.Println("  ╚══════════════════════════════════════════════════════════════════════════════╝")
// 	fmt.Print(Reset)

// 	fmt.Printf("%s%s", Dim, White)
// 	fmt.Println("  💡 Tip: Use 'mw-injector-manager --help' for available commands")
// 	fmt.Println("  🔗 Docs: https://docs.middleware.io/injector-manager")
// 	fmt.Print(Reset)

// 	fmt.Println()
// }

// // AnimatedWelcome displays an animated welcome sequence
// func AnimatedWelcome() {
// 	steps := []string{
// 		"🔍 Discovering Java services...",
// 		"📋 Loading configurations...",
// 		"🔧 Initializing LD_PRELOAD hooks...",
// 		"🚀 MW Injector Manager ready!",
// 	}

// 	for i, step := range steps {
// 		fmt.Printf("\r%s%s%s", Yellow, step, Reset)
// 		if i < len(steps)-1 {
// 			time.Sleep(800 * time.Millisecond)
// 		}
// 	}
// 	fmt.Println()
// 	fmt.Println()
// }

//	func WelcomeScreen() {
//		Banner()
//		AnimatedWelcome()
//	}
package tui

import (
	"fmt"
	"runtime"
)

const (
	Version = "1.0.0"
	Build   = "20250923"
)

func GetBannerArt() string {
	return ``
}

func GetSubtitle() string {
	return "🚀 OpenTelemetry Injection System 🚀"
}

func GetVersionInfo() string {
	return fmt.Sprintf("Version %s (Build %s) • Running on %s/%s • Go %s",
		Version, Build, runtime.GOOS, runtime.GOARCH, runtime.Version())
}
