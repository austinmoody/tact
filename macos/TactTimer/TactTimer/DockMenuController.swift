import AppKit

class DockMenuController {

    private let timerManager: TimerManager

    var onStartNewTimer: (() -> Void)?
    var onViewAllTimers: (() -> Void)?
    var onOpenPreferences: (() -> Void)?

    init(timerManager: TimerManager) {
        self.timerManager = timerManager
    }

    func buildMenu() -> NSMenu {
        let menu = NSMenu()

        // Status header
        let statusText = buildStatusText()
        let statusItem = NSMenuItem(title: statusText, action: nil, keyEquivalent: "")
        statusItem.isEnabled = false
        menu.addItem(statusItem)

        menu.addItem(NSMenuItem.separator())

        // Start New Timer
        let startItem = NSMenuItem(
            title: "Start New Timer...",
            action: #selector(startNewTimerClicked),
            keyEquivalent: "n"
        )
        startItem.target = self
        menu.addItem(startItem)

        // View All Timers (if any exist)
        if timerManager.timerCount > 0 {
            let viewItem = NSMenuItem(
                title: "View All Timers",
                action: #selector(viewAllTimersClicked),
                keyEquivalent: "l"
            )
            viewItem.target = self
            menu.addItem(viewItem)
        }

        menu.addItem(NSMenuItem.separator())

        // Preferences
        let prefsItem = NSMenuItem(
            title: "Preferences...",
            action: #selector(preferencesClicked),
            keyEquivalent: ","
        )
        prefsItem.target = self
        menu.addItem(prefsItem)

        return menu
    }

    private func buildStatusText() -> String {
        let count = timerManager.timerCount
        let running = timerManager.runningCount

        if count == 0 {
            return "No active timers"
        } else if count == 1 {
            return running == 1 ? "1 timer (running)" : "1 timer (paused)"
        } else {
            let runningText = running == 1 ? "1 running" : (running > 0 ? "\(running) running" : "all paused")
            return "\(count) timers (\(runningText))"
        }
    }

    @objc private func startNewTimerClicked() {
        onStartNewTimer?()
    }

    @objc private func viewAllTimersClicked() {
        onViewAllTimers?()
    }

    @objc private func preferencesClicked() {
        onOpenPreferences?()
    }
}
