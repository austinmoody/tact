package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"tact-tui/api"
)

type EntryInputModal struct {
	client    *api.Client
	textInput textinput.Model
	err       error
	saving    bool
}

func NewEntryInputModal(client *api.Client) *EntryInputModal {
	ti := textinput.New()
	ti.Placeholder = "2h meeting with client ABC123"
	ti.Focus()
	ti.CharLimit = 500
	ti.Width = 50

	return &EntryInputModal{
		client:    client,
		textInput: ti,
	}
}

func (m *EntryInputModal) Init() tea.Cmd {
	return textinput.Blink
}

func (m *EntryInputModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle keys directly by type to prevent control characters
		switch msg.Type {
		case tea.KeyUp, tea.KeyDown, tea.KeyShiftUp, tea.KeyShiftDown, tea.KeyCtrlUp, tea.KeyCtrlDown:
			return m, nil
		case tea.KeyEsc:
			return m, func() tea.Msg { return ModalCloseMsg{} }
		case tea.KeyEnter:
			if m.textInput.Value() != "" && !m.saving {
				m.saving = true
				return m, m.createEntry()
			}
			return m, nil
		}

		// Consume any escape sequences that slipped through
		keyStr := msg.String()
		if keyStr == "up" || keyStr == "down" {
			return m, nil
		}
		if len(keyStr) > 1 && keyStr[0] == '\x1b' {
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m *EntryInputModal) createEntry() tea.Cmd {
	return func() tea.Msg {
		_, err := m.client.CreateEntry(m.textInput.Value())
		if err != nil {
			return entryInputErrMsg{err}
		}
		return EntryCreatedMsg{}
	}
}

type entryInputErrMsg struct{ err error }

func (m *EntryInputModal) View() string {
	var b strings.Builder

	b.WriteString(modalTitleStyle.Render("New Time Entry"))
	b.WriteString("\n\n")

	b.WriteString(labelStyle.Render("Enter time entry (e.g. '2h meeting ABC123'):"))
	b.WriteString("\n")

	inputStyle := focusedInputStyle
	b.WriteString(inputStyle.Render(m.textInput.View()))
	b.WriteString("\n\n")

	if m.err != nil {
		b.WriteString(errorStyle.Render("Error: " + m.err.Error()))
		b.WriteString("\n\n")
	}

	if m.saving {
		b.WriteString(statusStyle.Render("Creating entry..."))
	} else {
		b.WriteString(helpStyle.Render("[Enter] Create  [Esc] Cancel"))
	}

	return modalStyle.Render(b.String())
}
