import Foundation

enum TimerState: String, Codable {
    case running
    case paused
}

struct TactTimer: Codable, Identifiable {
    let id: UUID
    var description: String
    var state: TimerState
    var startedAt: Date?           // When the current running period started
    var accumulatedSeconds: Int    // Total seconds accumulated before current period
    
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
}
