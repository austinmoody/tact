import Foundation

struct TimeFormatter {

    /// Formats seconds into a duration string for the API (e.g., "45m", "1h30m")
    static func formatDuration(seconds: Int) -> String {
        let totalMinutes = (seconds + 30) / 60  // Round to nearest minute

        if totalMinutes < 1 {
            return "1m"  // Minimum 1 minute
        }

        let hours = totalMinutes / 60
        let minutes = totalMinutes % 60

        if hours == 0 {
            return "\(minutes)m"
        } else if minutes == 0 {
            return "\(hours)h"
        } else {
            return "\(hours)h\(minutes)m"
        }
    }

    /// Formats a complete entry for the API: "{duration} {description}"
    static func formatEntry(seconds: Int, description: String) -> String {
        let duration = formatDuration(seconds: seconds)
        return "\(duration) \(description)"
    }

    /// Formats seconds into display time (e.g., "00:45:23" or "12:34")
    static func formatDisplay(seconds: Int) -> String {
        let hours = seconds / 3600
        let minutes = (seconds % 3600) / 60
        let secs = seconds % 60

        if hours > 0 {
            return String(format: "%d:%02d:%02d", hours, minutes, secs)
        } else {
            return String(format: "%02d:%02d", minutes, secs)
        }
    }
}
