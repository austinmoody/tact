package ui

import (
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"

	"tact-tui/api"
)

type EntryInputModal struct {
	client    *api.Client
	textInput textinput.Model
	err       error
	saving    bool
	width     int
}

func NewEntryInputModal(client *api.Client, width int) *EntryInputModal {
	inputWidth := calculateInputWidth(width)

	ti := textinput.New()
	ti.Placeholder = "2h meeting with client ABC123"
	ti.Focus()
	ti.CharLimit = 500
	ti.SetWidth(inputWidth)

	return &EntryInputModal{
		client:    client,
		textInput: ti,
		width:     width,
	}
}

func (m *EntryInputModal) Init() tea.Cmd {
	return textinput.Blink
}

func (m *EntryInputModal) Update(msg tea.Msg) (*EntryInputModal, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.Key().Code {
		case tea.KeyEscape:
			return m, func() tea.Msg { return ModalCloseMsg{} }
		case tea.KeyEnter:
			if m.textInput.Value() != "" && !m.saving {
				m.saving = true
				return m, m.createEntry()
			}
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

	style := focusedInputStyle
	b.WriteString(style.Render(m.textInput.View()))
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
