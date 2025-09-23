package tui

const (
	// ANSI color codes
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	White  = "\033[37m"
	Bold   = "\033[1m"
	Dim    = "\033[2m"
)

// ColorText applies color and formatting to text
func ColorText(text, color, format string) string {
	return color + format + text + Reset
}
