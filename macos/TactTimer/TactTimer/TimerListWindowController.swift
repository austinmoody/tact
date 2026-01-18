import AppKit
import Combine

class TimerListWindowController: NSWindowController, NSTableViewDataSource, NSTableViewDelegate {

    private let timerManager: TimerManager
    private var tableView: NSTableView!
    private var emptyLabel: NSTextField!
    private var updateTimer: Timer?
    private var cancellables = Set<AnyCancellable>()

    var onNewTimer: (() -> Void)?

    init(timerManager: TimerManager) {
        self.timerManager = timerManager

        let window = NSWindow(
            contentRect: NSRect(x: 0, y: 0, width: 450, height: 300),
            styleMask: [.titled, .closable, .resizable, .miniaturizable],
            backing: .buffered,
            defer: false
        )
        window.title = "Active Timers"
        window.center()
        window.isReleasedWhenClosed = false
        window.minSize = NSSize(width: 350, height: 200)

        super.init(window: window)
        setupUI()
        bindToManager()
    }

    required init?(coder: NSCoder) {
        fatalError("init(coder:) has not been implemented")
    }

    private func setupUI() {
        guard let contentView = window?.contentView else { return }

        // Bottom bar with New Timer button
        let bottomBar = NSView()
        bottomBar.translatesAutoresizingMaskIntoConstraints = false
        bottomBar.wantsLayer = true
        bottomBar.layer?.backgroundColor = NSColor.windowBackgroundColor.cgColor
        contentView.addSubview(bottomBar)

        // Separator line above bottom bar
        let separator = NSBox()
        separator.translatesAutoresizingMaskIntoConstraints = false
        separator.boxType = .separator
        contentView.addSubview(separator)

        // New Timer button
        let newTimerButton = NSButton(title: "+ New Timer", target: self, action: #selector(newTimerClicked))
        newTimerButton.translatesAutoresizingMaskIntoConstraints = false
        newTimerButton.bezelStyle = .rounded
        newTimerButton.keyEquivalent = "n"
        newTimerButton.keyEquivalentModifierMask = .command
        bottomBar.addSubview(newTimerButton)

        // Empty state label
        emptyLabel = NSTextField(labelWithString: "No active timers\nClick \"+ New Timer\" to start tracking")
        emptyLabel.translatesAutoresizingMaskIntoConstraints = false
        emptyLabel.alignment = .center
        emptyLabel.textColor = .secondaryLabelColor
        emptyLabel.maximumNumberOfLines = 2
        emptyLabel.isHidden = true
        contentView.addSubview(emptyLabel)

        // Table view
        tableView = NSTableView()
        tableView.dataSource = self
        tableView.delegate = self
        tableView.rowHeight = 60
        tableView.usesAlternatingRowBackgroundColors = true

        // Single column for custom cell
        let column = NSTableColumn(identifier: NSUserInterfaceItemIdentifier("TimerCell"))
        column.title = "Timers"
        column.width = 400
        tableView.addTableColumn(column)
        tableView.headerView = nil

        // Scroll view
        let scrollView = NSScrollView()
        scrollView.translatesAutoresizingMaskIntoConstraints = false
        scrollView.documentView = tableView
        scrollView.hasVerticalScroller = true
        scrollView.autohidesScrollers = true
        contentView.addSubview(scrollView)

        // Layout
        NSLayoutConstraint.activate([
            // Bottom bar
            bottomBar.leadingAnchor.constraint(equalTo: contentView.leadingAnchor),
            bottomBar.trailingAnchor.constraint(equalTo: contentView.trailingAnchor),
            bottomBar.bottomAnchor.constraint(equalTo: contentView.bottomAnchor),
            bottomBar.heightAnchor.constraint(equalToConstant: 44),

            // Separator
            separator.leadingAnchor.constraint(equalTo: contentView.leadingAnchor),
            separator.trailingAnchor.constraint(equalTo: contentView.trailingAnchor),
            separator.bottomAnchor.constraint(equalTo: bottomBar.topAnchor),

            // New Timer button in bottom bar
            newTimerButton.centerYAnchor.constraint(equalTo: bottomBar.centerYAnchor),
            newTimerButton.leadingAnchor.constraint(equalTo: bottomBar.leadingAnchor, constant: 12),

            // Scroll view above bottom bar
            scrollView.topAnchor.constraint(equalTo: contentView.topAnchor),
            scrollView.leadingAnchor.constraint(equalTo: contentView.leadingAnchor),
            scrollView.trailingAnchor.constraint(equalTo: contentView.trailingAnchor),
            scrollView.bottomAnchor.constraint(equalTo: separator.topAnchor),

            // Empty label centered in scroll area
            emptyLabel.centerXAnchor.constraint(equalTo: scrollView.centerXAnchor),
            emptyLabel.centerYAnchor.constraint(equalTo: scrollView.centerYAnchor),
        ])
    }

    @objc private func newTimerClicked() {
        onNewTimer?()
    }

    private func bindToManager() {
        timerManager.$timers
            .receive(on: DispatchQueue.main)
            .sink { [weak self] _ in
                self?.refresh()
            }
            .store(in: &cancellables)
    }

    override func showWindow(_ sender: Any?) {
        super.showWindow(sender)
        refresh()
        startUpdateTimer()
    }

    override func close() {
        stopUpdateTimer()
        super.close()
    }

    func refresh() {
        tableView.reloadData()
        updateEmptyState()
    }

    private func updateEmptyState() {
        let isEmpty = timerManager.timers.isEmpty
        emptyLabel.isHidden = !isEmpty
        tableView.isHidden = isEmpty
    }

    private func startUpdateTimer() {
        stopUpdateTimer()
        updateTimer = Timer.scheduledTimer(withTimeInterval: 1.0, repeats: true) { [weak self] _ in
            self?.updateElapsedTimes()
        }
    }

    private func stopUpdateTimer() {
        updateTimer?.invalidate()
        updateTimer = nil
    }

    private func updateElapsedTimes() {
        // Update visible rows only
        for row in 0..<timerManager.timers.count {
            if let cellView = tableView.view(atColumn: 0, row: row, makeIfNecessary: false) as? TimerCellView {
                cellView.updateElapsedTime(timerManager.timers[row])
            }
        }
    }

    // MARK: - NSTableViewDataSource

    func numberOfRows(in tableView: NSTableView) -> Int {
        return timerManager.timers.count
    }

    // MARK: - NSTableViewDelegate

    func tableView(_ tableView: NSTableView, viewFor tableColumn: NSTableColumn?, row: Int) -> NSView? {
        let timer = timerManager.timers[row]

        let cellView = TimerCellView()
        cellView.configure(with: timer)

        cellView.onPauseResume = { [weak self] in
            guard let self = self else { return }
            let timer = self.timerManager.timers[row]
            if timer.state == .running {
                self.timerManager.pauseTimer(id: timer.id)
            } else {
                self.timerManager.resumeTimer(id: timer.id)
            }
        }

        cellView.onStop = { [weak self] in
            guard let self = self else { return }
            let timer = self.timerManager.timers[row]
            self.stopTimer(timer)
        }

        return cellView
    }

    func tableView(_ tableView: NSTableView, heightOfRow row: Int) -> CGFloat {
        return 60
    }

    private func stopTimer(_ timer: TactTimer) {
        timerManager.stopTimer(id: timer.id) { [weak self] result in
            switch result {
            case .success:
                break // Timer removed automatically via binding
            case .failure(let error):
                self?.showErrorAlert(error: error, timer: timer)
            }
        }
    }

    private func showErrorAlert(error: Error, timer: TactTimer) {
        guard let window = self.window else { return }

        let alert = NSAlert()
        alert.messageText = "Failed to Save Entry"
        alert.informativeText = "\(error.localizedDescription)\n\nThe timer has been kept in the list. You can try again later."
        alert.alertStyle = .warning
        alert.addButton(withTitle: "OK")
        alert.beginSheetModal(for: window, completionHandler: nil)
    }
}

// MARK: - Timer Cell View

class TimerCellView: NSTableCellView {

    var onPauseResume: (() -> Void)?
    var onStop: (() -> Void)?

    private let descriptionLabel = NSTextField(labelWithString: "")
    private let elapsedLabel = NSTextField(labelWithString: "00:00")
    private let stateLabel = NSTextField(labelWithString: "")
    private let pauseResumeButton = NSButton()
    private let stopButton = NSButton()

    override init(frame frameRect: NSRect) {
        super.init(frame: frameRect)
        setupUI()
    }

    required init?(coder: NSCoder) {
        fatalError("init(coder:) has not been implemented")
    }

    private func setupUI() {
        descriptionLabel.translatesAutoresizingMaskIntoConstraints = false
        descriptionLabel.font = .systemFont(ofSize: 13, weight: .medium)
        descriptionLabel.lineBreakMode = .byTruncatingTail
        addSubview(descriptionLabel)

        elapsedLabel.translatesAutoresizingMaskIntoConstraints = false
        elapsedLabel.font = .monospacedDigitSystemFont(ofSize: 14, weight: .regular)
        elapsedLabel.textColor = .secondaryLabelColor
        addSubview(elapsedLabel)

        stateLabel.translatesAutoresizingMaskIntoConstraints = false
        stateLabel.font = .systemFont(ofSize: 11)
        stateLabel.textColor = .tertiaryLabelColor
        addSubview(stateLabel)

        pauseResumeButton.translatesAutoresizingMaskIntoConstraints = false
        pauseResumeButton.bezelStyle = .rounded
        pauseResumeButton.target = self
        pauseResumeButton.action = #selector(pauseResumeClicked)
        addSubview(pauseResumeButton)

        stopButton.translatesAutoresizingMaskIntoConstraints = false
        stopButton.title = "Stop"
        stopButton.bezelStyle = .rounded
        stopButton.target = self
        stopButton.action = #selector(stopClicked)
        addSubview(stopButton)

        NSLayoutConstraint.activate([
            descriptionLabel.topAnchor.constraint(equalTo: topAnchor, constant: 8),
            descriptionLabel.leadingAnchor.constraint(equalTo: leadingAnchor, constant: 12),
            descriptionLabel.trailingAnchor.constraint(lessThanOrEqualTo: pauseResumeButton.leadingAnchor, constant: -8),

            elapsedLabel.topAnchor.constraint(equalTo: descriptionLabel.bottomAnchor, constant: 4),
            elapsedLabel.leadingAnchor.constraint(equalTo: leadingAnchor, constant: 12),

            stateLabel.centerYAnchor.constraint(equalTo: elapsedLabel.centerYAnchor),
            stateLabel.leadingAnchor.constraint(equalTo: elapsedLabel.trailingAnchor, constant: 8),

            pauseResumeButton.centerYAnchor.constraint(equalTo: centerYAnchor),
            pauseResumeButton.trailingAnchor.constraint(equalTo: stopButton.leadingAnchor, constant: -8),
            pauseResumeButton.widthAnchor.constraint(equalToConstant: 70),

            stopButton.centerYAnchor.constraint(equalTo: centerYAnchor),
            stopButton.trailingAnchor.constraint(equalTo: trailingAnchor, constant: -12),
            stopButton.widthAnchor.constraint(equalToConstant: 50),
        ])
    }

    func configure(with timer: TactTimer) {
        descriptionLabel.stringValue = timer.description
        updateElapsedTime(timer)
        updateState(timer)
    }

    func updateElapsedTime(_ timer: TactTimer) {
        elapsedLabel.stringValue = timer.formattedElapsedTime
    }

    private func updateState(_ timer: TactTimer) {
        if timer.state == .running {
            stateLabel.stringValue = "Running"
            stateLabel.textColor = NSColor.systemGreen
            pauseResumeButton.title = "Pause"
        } else {
            stateLabel.stringValue = "Paused"
            stateLabel.textColor = .tertiaryLabelColor
            pauseResumeButton.title = "Resume"
        }
    }

    @objc private func pauseResumeClicked() {
        onPauseResume?()
    }

    @objc private func stopClicked() {
        onStop?()
    }
}
