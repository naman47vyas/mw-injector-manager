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
