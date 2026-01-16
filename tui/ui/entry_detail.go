package ui

import (
	"fmt"
	"strings"

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
}

func NewEntryDetailModal(client *api.Client, entry *model.Entry, width int) *EntryDetailModal {
	return &EntryDetailModal{
		client: client,
		entry:  entry,
		width:  width,
	}
}

func (m *EntryDetailModal) Init() tea.Cmd {
	return nil
}

func (m *EntryDetailModal) Update(msg tea.Msg) (*EntryDetailModal, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case matchesKey(msg, keys.Escape):
			return m, func() tea.Msg { return ModalCloseMsg{} }

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
	}

	return m, nil
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

func (m *EntryDetailModal) View() string {
	if m.entry == nil {
		return modalStyle.Render("No entry selected")
	}

	var b strings.Builder

	b.WriteString(modalTitleStyle.Render("Entry Details"))
	b.WriteString("\n\n")

	// User input
	b.WriteString(labelStyle.Render("User Input:"))
	b.WriteString("\n")
	b.WriteString("  " + m.entry.UserInput)
	b.WriteString("\n\n")

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
	b.WriteString(labelStyle.Render("Date: "))
	date := m.entry.EntryDate
	if len(date) > 10 {
		date = date[:10]
	}
	b.WriteString(date)
	b.WriteString("\n")

	// Parse error if any
	if m.entry.ParseError != nil && *m.entry.ParseError != "" {
		b.WriteString("\n")
		b.WriteString(errorStyle.Render("Parse Error: " + *m.entry.ParseError))
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
	if m.reparsing {
		b.WriteString(statusStyle.Render("Reparsing..."))
	} else {
		b.WriteString(helpStyle.Render("[p] Reparse  [Esc] Close"))
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
