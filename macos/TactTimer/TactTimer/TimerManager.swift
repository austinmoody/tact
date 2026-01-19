import Foundation

class TimerManager: ObservableObject {

    private static let storageKey = "TactTimers"

    @Published private(set) var timers: [TactTimer] = []

    var hasActiveTimers: Bool {
        !activeTimers.isEmpty
    }

    var runningTimer: TactTimer? {
        timers.first { $0.state == .running }
    }

    var timerCount: Int {
        activeTimers.count
    }

    var runningCount: Int {
        timers.filter { $0.state == .running }.count
    }

    /// Active timers (running or paused)
    var activeTimers: [TactTimer] {
        timers.filter { $0.state == .running || $0.state == .paused }
    }

    /// Completed timers from today only
    var completedTodayTimers: [TactTimer] {
        let startOfToday = Calendar.current.startOfDay(for: Date())
        return timers.filter { timer in
            guard timer.state == .stopped, let stoppedAt = timer.stoppedAt else { return false }
            return stoppedAt >= startOfToday
        }
    }

    init() {
        load()
    }

    // MARK: - Timer Operations

    func startNewTimer(description: String) {
        // Pause any currently running timer
        pauseRunningTimer()

        // Create and add new timer
        let timer = TactTimer(description: description)
        timers.insert(timer, at: 0)
        save()
    }

    func pauseTimer(id: UUID) {
        guard let index = timers.firstIndex(where: { $0.id == id }) else { return }
        timers[index].pause()
        save()
    }

    func resumeTimer(id: UUID) {
        // Pause any currently running timer first
        pauseRunningTimer()

        guard let index = timers.firstIndex(where: { $0.id == id }) else { return }
        timers[index].resume()
        save()
    }

    func pauseRunningTimer() {
        for i in timers.indices {
            if timers[i].state == .running {
                timers[i].pause()
            }
        }
        save()
    }

    func stopTimer(id: UUID, completion: @escaping (Result<Void, Error>) -> Void) {
        guard let index = timers.firstIndex(where: { $0.id == id }) else {
            completion(.failure(TimerError.notFound))
            return
        }

        // Format the time entry
        let formattedEntry = TimeFormatter.formatEntry(
            seconds: timers[index].totalElapsedSeconds,
            description: timers[index].description
        )

        // Call API to create entry
        APIClient.shared.createEntry(userInput: formattedEntry) { [weak self] result in
            DispatchQueue.main.async {
                switch result {
                case .success:
                    // Mark as stopped instead of removing
                    self?.timers[index].stop()
                    self?.save()
                    completion(.success(()))
                case .failure(let error):
                    completion(.failure(error))
                }
            }
        }
    }

    func removeTimer(id: UUID) {
        timers.removeAll { $0.id == id }
        save()
    }

    // MARK: - Persistence

    func save() {
        do {
            let data = try JSONEncoder().encode(timers)
            UserDefaults.standard.set(data, forKey: Self.storageKey)
        } catch {
            print("Failed to save timers: \(error)")
        }
    }

    private func load() {
        guard let data = UserDefaults.standard.data(forKey: Self.storageKey) else { return }

        do {
            var loadedTimers = try JSONDecoder().decode([TactTimer].self, from: data)

            // Clean up completed timers from previous days
            let startOfToday = Calendar.current.startOfDay(for: Date())
            loadedTimers.removeAll { timer in
                guard timer.state == .stopped, let stoppedAt = timer.stoppedAt else { return false }
                return stoppedAt < startOfToday
            }

            timers = loadedTimers
            save() // Persist the cleanup
        } catch {
            print("Failed to load timers: \(error)")
            timers = []
        }
    }
}

enum TimerError: LocalizedError {
    case notFound

    var errorDescription: String? {
        switch self {
        case .notFound:
            return "Timer not found"
        }
    }
}
