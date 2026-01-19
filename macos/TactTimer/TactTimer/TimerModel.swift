import Foundation

enum TimerState: String, Codable {
    case running
    case paused
    case stopped
}

struct TactTimer: Codable, Identifiable {
    let id: UUID
    var description: String
    var state: TimerState
    var startedAt: Date?           // When the current running period started
    var accumulatedSeconds: Int    // Total seconds accumulated before current period
    var stoppedAt: Date?           // When the timer was stopped (for cleanup)
    
    var totalElapsedSeconds: Int {
        var total = accumulatedSeconds
        if state == .running, let started = startedAt {
            total += Int(Date().timeIntervalSince(started))
        }
        return total
    }
    
    var formattedElapsedTime: String {
        let total = totalElapsedSeconds
        let hours = total / 3600
        let minutes = (total % 3600) / 60
        let seconds = total % 60

        if hours > 0 {
            return String(format: "%d:%02d:%02d", hours, minutes, seconds)
        } else {
            return String(format: "%02d:%02d", minutes, seconds)
        }
    }

    /// Formatted duration for completed timers (e.g., "45m", "1h30m")
    var formattedFinalDuration: String {
        let totalMinutes = (accumulatedSeconds + 30) / 60  // Round to nearest minute
        if totalMinutes < 1 {
            return "1m"
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
    
    init(id: UUID = UUID(), description: String) {
        self.id = id
        self.description = description
        self.state = .running
        self.startedAt = Date()
        self.accumulatedSeconds = 0
    }
    
    mutating func pause() {
        guard state == .running, let started = startedAt else { return }
        accumulatedSeconds += Int(Date().timeIntervalSince(started))
        startedAt = nil
        state = .paused
    }
    
    mutating func resume() {
        guard state == .paused else { return }
        startedAt = Date()
        state = .running
    }

    mutating func stop() {
        // Accumulate any remaining running time
        if state == .running, let started = startedAt {
            accumulatedSeconds += Int(Date().timeIntervalSince(started))
        }
        startedAt = nil
        state = .stopped
        stoppedAt = Date()
    }
}
