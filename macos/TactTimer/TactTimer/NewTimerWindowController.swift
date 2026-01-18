import AppKit

class NewTimerWindowController: NSWindowController, NSTextFieldDelegate {

    var onTimerCreated: ((String) -> Void)?

    private var textField: NSTextField!
    private var startButton: NSButton!
    private var cancelButton: NSButton!

    convenience init() {
        let window = NSWindow(
            contentRect: NSRect(x: 0, y: 0, width: 320, height: 120),
            styleMask: [.titled, .closable],
            backing: .buffered,
            defer: false
        )
        window.title = "New Timer"
        window.center()
        window.isReleasedWhenClosed = false

        self.init(window: window)
        setupUI()
    }

    private func setupUI() {
        guard let contentView = window?.contentView else { return }

        // Description label
        let label = NSTextField(labelWithString: "What are you working on?")
        label.translatesAutoresizingMaskIntoConstraints = false
        contentView.addSubview(label)

        // Text field for description
        textField = NSTextField()
        textField.translatesAutoresizingMaskIntoConstraints = false
        textField.placeholderString = "Enter timer description..."
        textField.delegate = self
        contentView.addSubview(textField)

        // Cancel button
        cancelButton = NSButton(title: "Cancel", target: self, action: #selector(cancelClicked))
        cancelButton.translatesAutoresizingMaskIntoConstraints = false
        cancelButton.keyEquivalent = "\u{1b}" // Escape
        contentView.addSubview(cancelButton)

        // Start button
        startButton = NSButton(title: "Start", target: self, action: #selector(startClicked))
        startButton.translatesAutoresizingMaskIntoConstraints = false
        startButton.keyEquivalent = "\r" // Enter
        startButton.bezelStyle = .rounded
        contentView.addSubview(startButton)

        // Layout
        NSLayoutConstraint.activate([
            label.topAnchor.constraint(equalTo: contentView.topAnchor, constant: 20),
            label.leadingAnchor.constraint(equalTo: contentView.leadingAnchor, constant: 20),
            label.trailingAnchor.constraint(equalTo: contentView.trailingAnchor, constant: -20),

            textField.topAnchor.constraint(equalTo: label.bottomAnchor, constant: 8),
            textField.leadingAnchor.constraint(equalTo: contentView.leadingAnchor, constant: 20),
            textField.trailingAnchor.constraint(equalTo: contentView.trailingAnchor, constant: -20),

            cancelButton.topAnchor.constraint(equalTo: textField.bottomAnchor, constant: 16),
            cancelButton.trailingAnchor.constraint(equalTo: startButton.leadingAnchor, constant: -8),
            cancelButton.bottomAnchor.constraint(equalTo: contentView.bottomAnchor, constant: -20),

            startButton.topAnchor.constraint(equalTo: textField.bottomAnchor, constant: 16),
            startButton.trailingAnchor.constraint(equalTo: contentView.trailingAnchor, constant: -20),
            startButton.bottomAnchor.constraint(equalTo: contentView.bottomAnchor, constant: -20),
        ])
    }

    override func showWindow(_ sender: Any?) {
        super.showWindow(sender)
        textField.stringValue = ""
        window?.makeFirstResponder(textField)
    }

    @objc private func startClicked() {
        let description = textField.stringValue.trimmingCharacters(in: .whitespacesAndNewlines)

        if description.isEmpty {
            // Shake the text field to indicate validation error
            let animation = CAKeyframeAnimation(keyPath: "transform.translation.x")
            animation.timingFunction = CAMediaTimingFunction(name: .linear)
            animation.duration = 0.4
            animation.values = [-8, 8, -6, 6, -4, 4, -2, 2, 0]
            textField.layer?.add(animation, forKey: "shake")
            NSSound.beep()
            return
        }

        onTimerCreated?(description)
        close()
    }

    @objc private func cancelClicked() {
        close()
    }
}
