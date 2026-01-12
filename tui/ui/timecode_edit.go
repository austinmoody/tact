package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"tact-tui/api"
	"tact-tui/model"
)

type TimeCodeEditModal struct {
	client   *api.Client
	timeCode *model.TimeCode // nil for add mode
	isEdit   bool

	// Input fields
	idInput          textinput.Model
	nameInput        textinput.Model
	descriptionInput textinput.Model
	keywordsInput    textinput.Model
	examplesInput    textinput.Model

	focusIndex int
	inputCount int

	saving bool
	err    error
}

func NewTimeCodeEditModal(client *api.Client, tc *model.TimeCode) *TimeCodeEditModal {
	isEdit := tc != nil

	idInput := textinput.New()
	idInput.Placeholder = "ABC123"
	idInput.CharLimit = 50
	idInput.Width = 40

	nameInput := textinput.New()
	nameInput.Placeholder = "Time Code Name"
	nameInput.CharLimit = 100
	nameInput.Width = 40

	descriptionInput := textinput.New()
	descriptionInput.Placeholder = "Optional description"
	descriptionInput.CharLimit = 500
	descriptionInput.Width = 40

	keywordsInput := textinput.New()
	keywordsInput.Placeholder = "keyword1, keyword2, keyword3"
	keywordsInput.CharLimit = 500
	keywordsInput.Width = 40

	examplesInput := textinput.New()
	examplesInput.Placeholder = "2h on project, 30m meeting"
	examplesInput.CharLimit = 500
	examplesInput.Width = 40

	inputCount := 2 // Add mode: just ID and Name
	if isEdit {
		// In edit mode, ID is readonly so we have: name, description, keywords, examples
		inputCount = 4
		idInput.SetValue(tc.ID)
		nameInput.SetValue(tc.Name)
		if tc.Description != "" {
			descriptionInput.SetValue(tc.Description)
		}
		if len(tc.Keywords) > 0 {
			keywordsInput.SetValue(strings.Join(tc.Keywords, ", "))
		}
		if len(tc.Examples) > 0 {
			examplesInput.SetValue(strings.Join(tc.Examples, ", "))
		}
		nameInput.Focus()
	} else {
		idInput.Focus()
	}

	return &TimeCodeEditModal{
		client:           client,
		timeCode:         tc,
		isEdit:           isEdit,
		idInput:          idInput,
		nameInput:        nameInput,
		descriptionInput: descriptionInput,
		keywordsInput:    keywordsInput,
		examplesInput:    examplesInput,
		focusIndex:       0,
		inputCount:       inputCount,
	}
}

func (m *TimeCodeEditModal) Init() tea.Cmd {
	return textinput.Blink
}

func (m *TimeCodeEditModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle navigation and consume arrow keys to prevent control chars in input
		switch msg.Type {
		case tea.KeyUp, tea.KeyShiftUp, tea.KeyCtrlUp:
			m.prevInput()
			return m, nil
		case tea.KeyDown, tea.KeyShiftDown, tea.KeyCtrlDown:
			m.nextInput()
			return m, nil
		case tea.KeyLeft, tea.KeyRight, tea.KeyShiftLeft, tea.KeyShiftRight:
			// Allow left/right for cursor movement within text - pass to focused input only
		case tea.KeyEsc:
			return m, func() tea.Msg { return ModalCloseMsg{} }
		case tea.KeyTab:
			m.nextInput()
			return m, nil
		case tea.KeyShiftTab:
			m.prevInput()
			return m, nil
		case tea.KeyEnter:
			if !m.saving {
				m.saving = true
				return m, m.save()
			}
			return m, nil
		}

		// Also check string representation for terminals that send raw sequences
		keyStr := msg.String()
		if keyStr == "up" {
			m.prevInput()
			return m, nil
		}
		if keyStr == "down" {
			m.nextInput()
			return m, nil
		}
		// Consume any escape sequences that slipped through
		if len(keyStr) > 1 && keyStr[0] == '\x1b' {
			return m, nil
		}
	}

	// Update only the focused input to avoid passing keys to all inputs
	var cmd tea.Cmd
	if m.isEdit {
		switch m.focusIndex {
		case 0:
			m.nameInput, cmd = m.nameInput.Update(msg)
		case 1:
			m.descriptionInput, cmd = m.descriptionInput.Update(msg)
		case 2:
			m.keywordsInput, cmd = m.keywordsInput.Update(msg)
		case 3:
			m.examplesInput, cmd = m.examplesInput.Update(msg)
		}
	} else {
		switch m.focusIndex {
		case 0:
			m.idInput, cmd = m.idInput.Update(msg)
		case 1:
			m.nameInput, cmd = m.nameInput.Update(msg)
		}
	}

	return m, cmd
}

func (m *TimeCodeEditModal) nextInput() {
	m.focusIndex = (m.focusIndex + 1) % m.inputCount
	m.updateFocus()
}

func (m *TimeCodeEditModal) prevInput() {
	m.focusIndex--
	if m.focusIndex < 0 {
		m.focusIndex = m.inputCount - 1
	}
	m.updateFocus()
}

func (m *TimeCodeEditModal) updateFocus() {
	m.idInput.Blur()
	m.nameInput.Blur()
	m.descriptionInput.Blur()
	m.keywordsInput.Blur()
	m.examplesInput.Blur()

	if m.isEdit {
		// In edit mode: name(0), description(1), keywords(2), examples(3)
		switch m.focusIndex {
		case 0:
			m.nameInput.Focus()
		case 1:
			m.descriptionInput.Focus()
		case 2:
			m.keywordsInput.Focus()
		case 3:
			m.examplesInput.Focus()
		}
	} else {
		// In add mode: id(0), name(1)
		switch m.focusIndex {
		case 0:
			m.idInput.Focus()
		case 1:
			m.nameInput.Focus()
		}
	}
}

func (m *TimeCodeEditModal) save() tea.Cmd {
	return func() tea.Msg {
		if m.isEdit {
			// Update existing time code
			name := m.nameInput.Value()
			desc := m.descriptionInput.Value()
			keywords := parseKeywords(m.keywordsInput.Value())
			examples := parseKeywords(m.examplesInput.Value()) // Same parsing as keywords

			updates := api.TimeCodeUpdate{
				Name:        &name,
				Description: &desc,
				Keywords:    keywords,
				Examples:    examples,
			}
			_, err := m.client.UpdateTimeCode(m.timeCode.ID, updates)
			if err != nil {
				return timeCodeEditErrMsg{err}
			}
			return TimeCodeUpdatedMsg{}
		}

		// Create new time code
		id := m.idInput.Value()
		name := m.nameInput.Value()
		if id == "" || name == "" {
			return timeCodeEditErrMsg{err: nil} // Validation error
		}

		_, err := m.client.CreateTimeCode(id, name, "")
		if err != nil {
			return timeCodeEditErrMsg{err}
		}
		return TimeCodeCreatedMsg{}
	}
}

func parseKeywords(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	var result []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}

type timeCodeEditErrMsg struct{ err error }

func (m *TimeCodeEditModal) View() string {
	var b strings.Builder

	if m.isEdit {
		b.WriteString(modalTitleStyle.Render("Edit Time Code"))
	} else {
		b.WriteString(modalTitleStyle.Render("Add Time Code"))
	}
	b.WriteString("\n\n")

	// ID field (readonly in edit mode)
	b.WriteString(labelStyle.Render("ID:"))
	b.WriteString("\n")
	if m.isEdit {
		b.WriteString(statusStyle.Render("  " + m.timeCode.ID + " (readonly)"))
	} else {
		style := inputStyle
		if m.focusIndex == 0 {
			style = focusedInputStyle
		}
		b.WriteString(style.Render(m.idInput.View()))
	}
	b.WriteString("\n\n")

	// Name field
	b.WriteString(labelStyle.Render("Name:"))
	b.WriteString("\n")
	nameIdx := 1
	if m.isEdit {
		nameIdx = 0
	}
	style := inputStyle
	if m.focusIndex == nameIdx {
		style = focusedInputStyle
	}
	b.WriteString(style.Render(m.nameInput.View()))
	b.WriteString("\n\n")

	// Additional fields (edit mode only)
	if m.isEdit {
		b.WriteString(labelStyle.Render("Description:"))
		b.WriteString("\n")
		style := inputStyle
		if m.focusIndex == 1 {
			style = focusedInputStyle
		}
		b.WriteString(style.Render(m.descriptionInput.View()))
		b.WriteString("\n\n")

		b.WriteString(labelStyle.Render("Keywords (comma separated):"))
		b.WriteString("\n")
		style = inputStyle
		if m.focusIndex == 2 {
			style = focusedInputStyle
		}
		b.WriteString(style.Render(m.keywordsInput.View()))
		b.WriteString("\n\n")

		b.WriteString(labelStyle.Render("Examples (comma separated):"))
		b.WriteString("\n")
		style = inputStyle
		if m.focusIndex == 3 {
			style = focusedInputStyle
		}
		b.WriteString(style.Render(m.examplesInput.View()))
		b.WriteString("\n\n")
	}

	// Error
	if m.err != nil {
		b.WriteString(errorStyle.Render("Error: " + m.err.Error()))
		b.WriteString("\n\n")
	}

	// Help
	if m.saving {
		b.WriteString(statusStyle.Render("Saving..."))
	} else {
		b.WriteString(helpStyle.Render("[Tab] Next Field  [Enter] Save  [Esc] Cancel"))
	}

	return modalStyle.Render(b.String())
}
