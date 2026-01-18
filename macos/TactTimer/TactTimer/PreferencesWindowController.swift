import AppKit

class PreferencesWindowController: NSWindowController, NSTextFieldDelegate {

    private var apiURLField: NSTextField!
    private var saveButton: NSButton!

    convenience init() {
        let window = NSWindow(
            contentRect: NSRect(x: 0, y: 0, width: 400, height: 140),
            styleMask: [.titled, .closable],
            backing: .buffered,
            defer: false
        )
        window.title = "Preferences"
        window.center()
        window.isReleasedWhenClosed = false

        self.init(window: window)
        setupUI()
    }

    private func setupUI() {
        guard let contentView = window?.contentView else { return }

        // API URL label
        let label = NSTextField(labelWithString: "API URL:")
        label.translatesAutoresizingMaskIntoConstraints = false
        label.alignment = .right
        contentView.addSubview(label)

        // API URL field
        apiURLField = NSTextField()
        apiURLField.translatesAutoresizingMaskIntoConstraints = false
        apiURLField.placeholderString = "http://localhost:2100"
        apiURLField.delegate = self
        contentView.addSubview(apiURLField)

        // Help text
        let helpLabel = NSTextField(labelWithString: "The base URL of your Tact backend API")
        helpLabel.translatesAutoresizingMaskIntoConstraints = false
        helpLabel.font = .systemFont(ofSize: 11)
        helpLabel.textColor = .secondaryLabelColor
        contentView.addSubview(helpLabel)

        // Save button
        saveButton = NSButton(title: "Save", target: self, action: #selector(saveClicked))
        saveButton.translatesAutoresizingMaskIntoConstraints = false
        saveButton.bezelStyle = .rounded
        saveButton.keyEquivalent = "\r"
        contentView.addSubview(saveButton)

        // Cancel button
        let cancelButton = NSButton(title: "Cancel", target: self, action: #selector(cancelClicked))
        cancelButton.translatesAutoresizingMaskIntoConstraints = false
        cancelButton.keyEquivalent = "\u{1b}"
        contentView.addSubview(cancelButton)

        // Layout
        NSLayoutConstraint.activate([
            label.topAnchor.constraint(equalTo: contentView.topAnchor, constant: 20),
            label.leadingAnchor.constraint(equalTo: contentView.leadingAnchor, constant: 20),
            label.widthAnchor.constraint(equalToConstant: 60),

            apiURLField.centerYAnchor.constraint(equalTo: label.centerYAnchor),
            apiURLField.leadingAnchor.constraint(equalTo: label.trailingAnchor, constant: 8),
            apiURLField.trailingAnchor.constraint(equalTo: contentView.trailingAnchor, constant: -20),

            helpLabel.topAnchor.constraint(equalTo: apiURLField.bottomAnchor, constant: 4),
            helpLabel.leadingAnchor.constraint(equalTo: apiURLField.leadingAnchor),

            cancelButton.topAnchor.constraint(equalTo: helpLabel.bottomAnchor, constant: 20),
            cancelButton.trailingAnchor.constraint(equalTo: saveButton.leadingAnchor, constant: -8),
            cancelButton.bottomAnchor.constraint(equalTo: contentView.bottomAnchor, constant: -20),

            saveButton.topAnchor.constraint(equalTo: helpLabel.bottomAnchor, constant: 20),
            saveButton.trailingAnchor.constraint(equalTo: contentView.trailingAnchor, constant: -20),
            saveButton.bottomAnchor.constraint(equalTo: contentView.bottomAnchor, constant: -20),
        ])
    }

    override func showWindow(_ sender: Any?) {
        super.showWindow(sender)
        apiURLField.stringValue = APIClient.shared.baseURL
        window?.makeFirstResponder(apiURLField)
    }

    @objc private func saveClicked() {
        let url = apiURLField.stringValue.trimmingCharacters(in: .whitespacesAndNewlines)

        // Basic URL validation
        if url.isEmpty {
            apiURLField.stringValue = "http://localhost:2100"
        }

        // Remove trailing slash if present
        var normalizedURL = url
        if normalizedURL.hasSuffix("/") {
            normalizedURL = String(normalizedURL.dropLast())
        }

        APIClient.shared.baseURL = normalizedURL
        close()
    }

    @objc private func cancelClicked() {
        close()
    }
}
