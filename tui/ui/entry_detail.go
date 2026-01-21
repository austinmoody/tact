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

// Focus field constants for edit mode
const (
	fieldUserInput = iota
	fieldDate
	fieldTimeCode
	fieldWorkType
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
	focusedField   int // 0=userInput, 1=date, 2=timeCode, 3=workType
	saving         bool
	validationErr  string

	// Code selection fields
	timeCodes        []model.TimeCode
	workTypes        []model.WorkType
	loadingCodes     bool
	selectedTimeCode int // index in timeCodes slice, -1 for none
	selectedWorkType int // index in workTypes slice, -1 for none

	// For learn flag comparison
	originalTimeCodeID *string
	originalWorkTypeID *string
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
				return m, m.enterEditMode()
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

	case entryDetailCodesFetchedMsg:
		m.timeCodes = msg.timeCodes
		m.workTypes = msg.workTypes
		m.loadingCodes = false
		m.initializeSelections()
		return m, nil

	case entryDetailCodesFetchErrMsg:
		m.err = msg.err
		m.loadingCodes = false
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
		if m.saving || m.loadingCodes {
			return m, nil
		}
		return m.saveChanges()

	case matchesKey(msg, keys.Tab):
		// Cycle focus forward through all 4 fields
		m.cycleFocusForward()
		return m, nil

	case matchesKey(msg, keys.ShiftTab):
		// Cycle focus backward through all 4 fields
		m.cycleFocusBackward()
		return m, nil

	case matchesKey(msg, keys.Down):
		// Navigate dropdown when on code fields
		if m.focusedField == fieldTimeCode && !m.loadingCodes {
			m.selectNextTimeCode()
			return m, nil
		} else if m.focusedField == fieldWorkType && !m.loadingCodes {
			m.selectNextWorkType()
			return m, nil
		}

	case matchesKey(msg, keys.Up):
		// Navigate dropdown when on code fields
		if m.focusedField == fieldTimeCode && !m.loadingCodes {
			m.selectPrevTimeCode()
			return m, nil
		} else if m.focusedField == fieldWorkType && !m.loadingCodes {
			m.selectPrevWorkType()
			return m, nil
		}
	}

	// Forward to focused text input (only for text fields)
	var cmd tea.Cmd
	if m.focusedField == fieldUserInput {
		m.userInputField, cmd = m.userInputField.Update(msg)
	} else if m.focusedField == fieldDate {
		m.dateField, cmd = m.dateField.Update(msg)
	}
	return m, cmd
}

func (m *EntryDetailModal) enterEditMode() tea.Cmd {
	m.editMode = true
	m.focusedField = fieldUserInput
	m.validationErr = ""
	m.err = nil
	m.loadingCodes = true

	// Reset field values from current entry
	m.userInputField.SetValue(m.entry.UserInput)
	date := m.entry.EntryDate
	if len(date) > 10 {
		date = date[:10]
	}
	m.dateField.SetValue(date)

	// Store original code values for learn flag comparison
	m.originalTimeCodeID = m.entry.TimeCodeID
	m.originalWorkTypeID = m.entry.WorkTypeID

	// Reset selections (will be initialized after fetch)
	m.selectedTimeCode = -1
	m.selectedWorkType = -1

	// Focus the user input field
	m.userInputField.Focus()
	m.dateField.Blur()

	// Fetch time codes and work types
	return m.fetchCodes()
}

func (m *EntryDetailModal) fetchCodes() tea.Cmd {
	return func() tea.Msg {
		timeCodes, err := m.client.FetchTimeCodes()
		if err != nil {
			return entryDetailCodesFetchErrMsg{err}
		}
		workTypes, err := m.client.FetchWorkTypes()
		if err != nil {
			return entryDetailCodesFetchErrMsg{err}
		}
		return entryDetailCodesFetchedMsg{timeCodes, workTypes}
	}
}

func (m *EntryDetailModal) initializeSelections() {
	// Find index of current time code
	m.selectedTimeCode = -1
	if m.entry.TimeCodeID != nil {
		for i, tc := range m.timeCodes {
			if tc.ID == *m.entry.TimeCodeID {
				m.selectedTimeCode = i
				break
			}
		}
	}

	// Find index of current work type
	m.selectedWorkType = -1
	if m.entry.WorkTypeID != nil {
		for i, wt := range m.workTypes {
			if wt.ID == *m.entry.WorkTypeID {
				m.selectedWorkType = i
				break
			}
		}
	}
}

func (m *EntryDetailModal) exitEditMode() {
	m.editMode = false
	m.validationErr = ""
	m.userInputField.Blur()
	m.dateField.Blur()
}

func (m *EntryDetailModal) cycleFocusForward() {
	m.blurCurrentField()
	m.focusedField = (m.focusedField + 1) % 4
	m.focusCurrentField()
}

func (m *EntryDetailModal) cycleFocusBackward() {
	m.blurCurrentField()
	m.focusedField = (m.focusedField + 3) % 4 // +3 is same as -1 mod 4
	m.focusCurrentField()
}

func (m *EntryDetailModal) blurCurrentField() {
	switch m.focusedField {
	case fieldUserInput:
		m.userInputField.Blur()
	case fieldDate:
		m.dateField.Blur()
	}
}

func (m *EntryDetailModal) focusCurrentField() {
	switch m.focusedField {
	case fieldUserInput:
		m.userInputField.Focus()
	case fieldDate:
		m.dateField.Focus()
	}
}

func (m *EntryDetailModal) selectNextTimeCode() {
	if len(m.timeCodes) == 0 {
		return
	}
	if m.selectedTimeCode < len(m.timeCodes)-1 {
		m.selectedTimeCode++
	}
}

func (m *EntryDetailModal) selectPrevTimeCode() {
	if len(m.timeCodes) == 0 {
		return
	}
	if m.selectedTimeCode > -1 {
		m.selectedTimeCode--
	}
}

func (m *EntryDetailModal) selectNextWorkType() {
	if len(m.workTypes) == 0 {
		return
	}
	if m.selectedWorkType < len(m.workTypes)-1 {
		m.selectedWorkType++
	}
}

func (m *EntryDetailModal) selectPrevWorkType() {
	if len(m.workTypes) == 0 {
		return
	}
	if m.selectedWorkType > -1 {
		m.selectedWorkType--
	}
}

func (m *EntryDetailModal) validateDate(date string) bool {
	// Validate YYYY-MM-DD format
	dateRegex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	return dateRegex.MatchString(date)
}

func (m *EntryDetailModal) getSelectedTimeCodeID() *string {
	if m.selectedTimeCode >= 0 && m.selectedTimeCode < len(m.timeCodes) {
		return &m.timeCodes[m.selectedTimeCode].ID
	}
	return nil
}

func (m *EntryDetailModal) getSelectedWorkTypeID() *string {
	if m.selectedWorkType >= 0 && m.selectedWorkType < len(m.workTypes) {
		return &m.workTypes[m.selectedWorkType].ID
	}
	return nil
}

func ptrStringEqual(a, b *string) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

func (m *EntryDetailModal) shouldLearn() bool {
	currentTimeCode := m.getSelectedTimeCodeID()
	currentWorkType := m.getSelectedWorkTypeID()

	timeCodeChanged := !ptrStringEqual(m.originalTimeCodeID, currentTimeCode)
	workTypeChanged := !ptrStringEqual(m.originalWorkTypeID, currentWorkType)

	return timeCodeChanged || workTypeChanged
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
	timeCodeID := m.getSelectedTimeCodeID()
	workTypeID := m.getSelectedWorkTypeID()
	learn := m.shouldLearn()

	return m, func() tea.Msg {
		update := api.EntryUpdate{
			UserInput:  &userInput,
			EntryDate:  &date,
			TimeCodeID: timeCodeID,
			WorkTypeID: workTypeID,
		}
		// learn=true when time_code or work_type changed (AI parsing corrections)
		// learn=false when only user_input or entry_date changed
		entry, err := m.client.UpdateEntry(m.entry.ID, update, learn)
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
type entryDetailCodesFetchedMsg struct {
	timeCodes []model.TimeCode
	workTypes []model.WorkType
}
type entryDetailCodesFetchErrMsg struct{ err error }

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
		if m.focusedField == fieldUserInput {
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
		if m.focusedField == fieldDate {
			label = "> Date:"
		}
		b.WriteString(labelStyle.Render(label))
		b.WriteString(" ")
		b.WriteString(m.dateField.View())
		b.WriteString("\n")

		// Time Code selection
		b.WriteString("\n")
		b.WriteString(m.renderTimeCodeField())

		// Work Type selection
		b.WriteString("\n")
		b.WriteString(m.renderWorkTypeField())
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
	} else if m.loadingCodes {
		b.WriteString(statusStyle.Render("Loading codes..."))
	} else if m.reparsing {
		b.WriteString(statusStyle.Render("Reparsing..."))
	} else if m.editMode {
		b.WriteString(helpStyle.Render("[Tab] Switch field  [j/k] Select code  [Enter] Save  [Esc] Cancel"))
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

func (m *EntryDetailModal) renderTimeCodeField() string {
	var b strings.Builder

	label := "Time Code:"
	if m.focusedField == fieldTimeCode {
		label = "> Time Code:"
	}
	b.WriteString(labelStyle.Render(label))

	// Show loading indicator or current selection
	if m.loadingCodes {
		b.WriteString(" Loading...")
		return b.String()
	}

	// Show current selection
	if m.selectedTimeCode >= 0 && m.selectedTimeCode < len(m.timeCodes) {
		tc := m.timeCodes[m.selectedTimeCode]
		b.WriteString(fmt.Sprintf(" [%s] %s", tc.ID, tc.Name))
	} else {
		b.WriteString(" (none)")
	}
	b.WriteString("\n")

	// Show dropdown list when focused
	if m.focusedField == fieldTimeCode && len(m.timeCodes) > 0 {
		b.WriteString(m.renderDropdownList(len(m.timeCodes), m.selectedTimeCode, func(i int) string {
			tc := m.timeCodes[i]
			return fmt.Sprintf("[%s] %s", tc.ID, tc.Name)
		}))
	}

	return b.String()
}

func (m *EntryDetailModal) renderWorkTypeField() string {
	var b strings.Builder

	label := "Work Type:"
	if m.focusedField == fieldWorkType {
		label = "> Work Type:"
	}
	b.WriteString(labelStyle.Render(label))

	// Show loading indicator or current selection
	if m.loadingCodes {
		b.WriteString(" Loading...")
		return b.String()
	}

	// Show current selection
	if m.selectedWorkType >= 0 && m.selectedWorkType < len(m.workTypes) {
		wt := m.workTypes[m.selectedWorkType]
		b.WriteString(fmt.Sprintf(" %s", wt.Name))
	} else {
		b.WriteString(" (none)")
	}
	b.WriteString("\n")

	// Show dropdown list when focused
	if m.focusedField == fieldWorkType && len(m.workTypes) > 0 {
		b.WriteString(m.renderDropdownList(len(m.workTypes), m.selectedWorkType, func(i int) string {
			return m.workTypes[i].Name
		}))
	}

	return b.String()
}

func (m *EntryDetailModal) renderDropdownList(total int, selected int, format func(int) string) string {
	var b strings.Builder

	// Show max 5 items, centered around selection
	maxVisible := 5

	// Calculate start and end indices
	start := 0
	end := total
	if total > maxVisible {
		// Center around selected item
		start = selected - maxVisible/2
		if start < 0 {
			start = 0
		}
		end = start + maxVisible
		if end > total {
			end = total
			start = end - maxVisible
		}
	}

	// Show scroll up indicator
	if start > 0 {
		b.WriteString(fmt.Sprintf("    ↑ %d more\n", start))
	}

	// Render visible items
	for i := start; i < end; i++ {
		prefix := "  "
		suffix := ""
		if i == selected {
			prefix = "→ "
			suffix = " ✓"
		}
		b.WriteString(fmt.Sprintf("  %s%s%s\n", prefix, format(i), suffix))
	}

	// Show scroll down indicator
	if end < total {
		b.WriteString(fmt.Sprintf("    ↓ %d more\n", total-end))
	}

	return b.String()
}
