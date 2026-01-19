package timer

import "fmt"

// FormatElapsed formats seconds for display (MM:SS or H:MM:SS)
func FormatElapsed(seconds int) string {
	if seconds < 0 {
		seconds = 0
	}

	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	secs := seconds % 60

	if hours > 0 {
		return fmt.Sprintf("%d:%02d:%02d", hours, minutes, secs)
	}
	return fmt.Sprintf("%02d:%02d", minutes, secs)
}

// FormatDuration formats seconds for API submission (e.g., "1h30m", "45m")
// Rounds to the nearest minute, minimum 1m
func FormatDuration(seconds int) string {
	// Round to nearest minute
	totalMinutes := (seconds + 30) / 60
	if totalMinutes < 1 {
		totalMinutes = 1
	}

	hours := totalMinutes / 60
	minutes := totalMinutes % 60

	if hours == 0 {
		return fmt.Sprintf("%dm", minutes)
	}
	if minutes == 0 {
		return fmt.Sprintf("%dh", hours)
	}
	return fmt.Sprintf("%dh%dm", hours, minutes)
}

// FormatEntry creates an entry string for the API (e.g., "1h30m Fix auth bug")
func FormatEntry(seconds int, description string) string {
	return fmt.Sprintf("%s %s", FormatDuration(seconds), description)
}
