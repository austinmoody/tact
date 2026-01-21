package ui

import (
	"fmt"
	"regexp"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"

	"tact-tui/api"
	"tact-tui/model"
)

type EntryDetailModal struct {
	client    *api.Client
	entry     *model.Entry
	width     int
	reparsing bool
	err       error

	// Edit mode fields
	editMode       bool
	userInputField textinput.Model
	dateField      textinput.Model
	focusedField   int // 0 = userInput, 1 = date
	saving         bool
	validationErr  string
}

func NewEntryDetailModal(client *api.Client, entry *model.Entry, width int) *EntryDetailModal {
	// Initialize text input fields
	userInput := textinput.New()
	userInput.Placeholder = "Enter description..."
	userInput.CharLimit = 500
	userInput.SetWidth(50)

	dateInput := textinput.New()
	dateInput.Placeholder = "YYYY-MM-DD"
	dateInput.CharLimit = 10
	dateInput.SetWidth(12)

	// Pre-populate with entry values
	if entry != nil {
		userInput.SetValue(entry.UserInput)
		date := entry.EntryDate
		if len(date) > 10 {
			date = date[:10]
		}
		dateInput.SetValue(date)
	}

	return &EntryDetailModal{
		client:         client,
		entry:          entry,
		width:          width,
		userInputField: userInput,
		dateField:      dateInput,
	}
}

func (m *EntryDetailModal) Init() tea.Cmd {
	return nil
}

func (m *EntryDetailModal) Update(msg tea.Msg) (*EntryDetailModal, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle edit mode
		if m.editMode {
			return m.handleEditModeKey(msg)
		}

		// View mode key handling
		switch {
		case matchesKey(msg, keys.Escape):
			return m, func() tea.Msg { return ModalCloseMsg{} }

		case matchesKey(msg, keys.Edit):
			if m.entry != nil && !m.reparsing {
				m.enterEditMode()
				return m, nil
			}

		case matchesKey(msg, keys.Reparse):
			if !m.reparsing && m.entry != nil {
				m.reparsing = true
				return m, m.reparse()
			}
			return m, nil
		}

	case entryDetailReparseOkMsg:
		m.entry = msg.entry
		m.reparsing = false
		return m, func() tea.Msg { return EntryReparseMsg{} }

	case entryDetailErrMsg:
		m.err = msg.err
		m.reparsing = false
		return m, nil

	case entryDetailUpdateOkMsg:
		m.entry = msg.entry
		m.saving = false
		m.editMode = false
		m.validationErr = ""
		return m, func() tea.Msg { return EntryUpdatedMsg{} }

	case entryDetailUpdateErrMsg:
		m.err = msg.err
		m.saving = false
		return m, nil
	}

	return m, nil
}

func (m *EntryDetailModal) handleEditModeKey(msg tea.KeyMsg) (*EntryDetailModal, tea.Cmd) {
	switch {
	case matchesKey(msg, keys.Escape):
		// Cancel edit mode
		m.exitEditMode()
		return m, nil

	case matchesKey(msg, keys.Enter):
		// Save changes
		if m.saving {
			return m, nil
		}
		return m.saveChanges()

	case matchesKey(msg, keys.Tab), matchesKey(msg, keys.ShiftTab):
		// Switch focus between fields
		m.toggleFocus()
		return m, nil

	default:
		// Forward to focused text input
		var cmd tea.Cmd
		if m.focusedField == 0 {
			m.userInputField, cmd = m.userInputField.Update(msg)
		} else {
			m.dateField, cmd = m.dateField.Update(msg)
		}
		return m, cmd
	}
}

func (m *EntryDetailModal) enterEditMode() {
	m.editMode = true
	m.focusedField = 0
	m.validationErr = ""
	m.err = nil

	// Reset field values from current entry
	m.userInputField.SetValue(m.entry.UserInput)
	date := m.entry.EntryDate
	if len(date) > 10 {
		date = date[:10]
	}
	m.dateField.SetValue(date)

	// Focus the user input field
	m.userInputField.Focus()
	m.dateField.Blur()
}

func (m *EntryDetailModal) exitEditMode() {
	m.editMode = false
	m.validationErr = ""
	m.userInputField.Blur()
	m.dateField.Blur()
}

func (m *EntryDetailModal) toggleFocus() {
	if m.focusedField == 0 {
		m.focusedField = 1
		m.userInputField.Blur()
		m.dateField.Focus()
	} else {
		m.focusedField = 0
		m.dateField.Blur()
		m.userInputField.Focus()
	}
}

func (m *EntryDetailModal) validateDate(date string) bool {
	// Validate YYYY-MM-DD format
	dateRegex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	return dateRegex.MatchString(date)
}

func (m *EntryDetailModal) saveChanges() (*EntryDetailModal, tea.Cmd) {
	// Validate date format
	date := m.dateField.Value()
	if !m.validateDate(date) {
		m.validationErr = "Invalid date format. Use YYYY-MM-DD"
		return m, nil
	}

	m.saving = true
	m.validationErr = ""

	userInput := m.userInputField.Value()
	return m, func() tea.Msg {
		update := api.EntryUpdate{
			UserInput: &userInput,
			EntryDate: &date,
		}
		// learn=false for user_input/entry_date edits - these aren't AI parsing corrections
		// learn=true should only be used when correcting time_code or work_type
		entry, err := m.client.UpdateEntry(m.entry.ID, update, false)
		if err != nil {
			return entryDetailUpdateErrMsg{err}
		}
		return entryDetailUpdateOkMsg{entry}
	}
}

func (m *EntryDetailModal) reparse() tea.Cmd {
	return func() tea.Msg {
		entry, err := m.client.ReparseEntry(m.entry.ID)
		if err != nil {
			return entryDetailErrMsg{err}
		}
		return entryDetailReparseOkMsg{entry}
	}
}

type entryDetailReparseOkMsg struct{ entry *model.Entry }
type entryDetailErrMsg struct{ err error }
type entryDetailUpdateOkMsg struct{ entry *model.Entry }
type entryDetailUpdateErrMsg struct{ err error }

func (m *EntryDetailModal) View() string {
	if m.entry == nil {
		return modalStyle.Render("No entry selected")
	}

	var b strings.Builder

	if m.editMode {
		b.WriteString(modalTitleStyle.Render("Edit Entry"))
	} else {
		b.WriteString(modalTitleStyle.Render("Entry Details"))
	}
	b.WriteString("\n\n")

	// User input
	if m.editMode {
		label := "User Input:"
		if m.focusedField == 0 {
			label = "> User Input:"
		}
		b.WriteString(labelStyle.Render(label))
		b.WriteString("\n")
		b.WriteString("  " + m.userInputField.View())
		b.WriteString("\n\n")
	} else {
		b.WriteString(labelStyle.Render("User Input:"))
		b.WriteString("\n")
		b.WriteString("  " + m.entry.UserInput)
		b.WriteString("\n\n")
	}

	// Status
	b.WriteString(labelStyle.Render("Status: "))
	b.WriteString(m.renderStatus(m.entry.Status))
	b.WriteString("\n\n")

	// Parsed fields
	b.WriteString(labelStyle.Render("Parsed Fields:"))
	b.WriteString("\n")

	// Duration
	if m.entry.DurationMinutes != nil {
		hours := *m.entry.DurationMinutes / 60
		mins := *m.entry.DurationMinutes % 60
		var duration string
		if hours > 0 && mins > 0 {
			duration = fmt.Sprintf("%dh %dm", hours, mins)
		} else if hours > 0 {
			duration = fmt.Sprintf("%dh", hours)
		} else {
			duration = fmt.Sprintf("%dm", mins)
		}
		confidence := ""
		if m.entry.ConfidenceDuration != nil {
			confidence = fmt.Sprintf(" (%.0f%%)", *m.entry.ConfidenceDuration*100)
		}
		b.WriteString(fmt.Sprintf("  Duration: %s%s\n", duration, confidence))
	} else {
		b.WriteString("  Duration: -\n")
	}

	// Time Code
	if m.entry.TimeCodeID != nil {
		confidence := ""
		if m.entry.ConfidenceTimeCode != nil {
			confidence = fmt.Sprintf(" (%.0f%%)", *m.entry.ConfidenceTimeCode*100)
		}
		b.WriteString(fmt.Sprintf("  Time Code: %s%s\n", *m.entry.TimeCodeID, confidence))
	} else {
		b.WriteString("  Time Code: -\n")
	}

	// Work Type
	if m.entry.WorkTypeID != nil {
		confidence := ""
		if m.entry.ConfidenceWorkType != nil {
			confidence = fmt.Sprintf(" (%.0f%%)", *m.entry.ConfidenceWorkType*100)
		}
		b.WriteString(fmt.Sprintf("  Work Type: %s%s\n", *m.entry.WorkTypeID, confidence))
	} else {
		b.WriteString("  Work Type: -\n")
	}

	// Parsed Description
	if m.entry.ParsedDescription != nil && *m.entry.ParsedDescription != "" {
		b.WriteString(fmt.Sprintf("  Description: %s\n", *m.entry.ParsedDescription))
	}

	// Overall confidence
	if m.entry.ConfidenceOverall != nil {
		b.WriteString(fmt.Sprintf("\n  Overall Confidence: %.0f%%\n", *m.entry.ConfidenceOverall*100))
	}

	// Parse notes (LLM reasoning and context info)
	if m.entry.ParseNotes != nil && *m.entry.ParseNotes != "" {
		b.WriteString("\n")
		b.WriteString(labelStyle.Render("Parse Notes:"))
		b.WriteString("\n")
		// Wrap notes to fit within modal width
		// Modal has padding (2 each side) + border (1 each side) + indent (2)
		wrapWidth := m.width - 10
		if wrapWidth < 30 {
			wrapWidth = 30
		}
		if wrapWidth > 70 {
			wrapWidth = 70
		}
		wrapped := wrapText(*m.entry.ParseNotes, wrapWidth)
		for _, line := range strings.Split(wrapped, "\n") {
			b.WriteString("  " + line + "\n")
		}
	}

	// Entry date
	b.WriteString("\n")
	if m.editMode {
		label := "Date:"
		if m.focusedField == 1 {
			label = "> Date:"
		}
		b.WriteString(labelStyle.Render(label))
		b.WriteString(" ")
		b.WriteString(m.dateField.View())
		b.WriteString("\n")
	} else {
		b.WriteString(labelStyle.Render("Date: "))
		date := m.entry.EntryDate
		if len(date) > 10 {
			date = date[:10]
		}
		b.WriteString(date)
		b.WriteString("\n")
	}

	// Parse error if any
	if m.entry.ParseError != nil && *m.entry.ParseError != "" {
		b.WriteString("\n")
		b.WriteString(errorStyle.Render("Parse Error: " + *m.entry.ParseError))
		b.WriteString("\n")
	}

	// Validation error
	if m.validationErr != "" {
		b.WriteString("\n")
		b.WriteString(errorStyle.Render(m.validationErr))
		b.WriteString("\n")
	}

	// Error from actions
	if m.err != nil {
		b.WriteString("\n")
		b.WriteString(errorStyle.Render("Error: " + m.err.Error()))
		b.WriteString("\n")
	}

	// Help
	b.WriteString("\n")
	if m.saving {
		b.WriteString(statusStyle.Render("Saving..."))
	} else if m.reparsing {
		b.WriteString(statusStyle.Render("Reparsing..."))
	} else if m.editMode {
		b.WriteString(helpStyle.Render("[Tab] Switch field  [Enter] Save  [Esc] Cancel"))
	} else {
		b.WriteString(helpStyle.Render("[e] Edit  [p] Reparse  [Esc] Close"))
	}

	return modalStyle.Render(b.String())
}

func (m *EntryDetailModal) renderStatus(status string) string {
	switch status {
	case "parsed":
		return statusParsedStyle.Render("Parsed")
	case "pending":
		return statusPendingStyle.Render("Pending")
	case "failed":
		return statusFailedStyle.Render("Failed")
	default:
		return statusStyle.Render(status)
	}
}
