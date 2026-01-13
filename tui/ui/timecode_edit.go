package ui

import (
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"

	"tact-tui/api"
	"tact-tui/model"
)

type TimeCodeEditModal struct {
	client   *api.Client
	timeCode *model.TimeCode // nil for add mode
	isEdit   bool
	width    int

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

func NewTimeCodeEditModal(client *api.Client, tc *model.TimeCode, width int) *TimeCodeEditModal {
	isEdit := tc != nil
	inputWidth := calculateInputWidth(width)

	idInput := textinput.New()
	idInput.Placeholder = "ABC123"
	idInput.CharLimit = 50
	idInput.SetWidth(inputWidth)

	nameInput := textinput.New()
	nameInput.Placeholder = "Time Code Name"
	nameInput.CharLimit = 100
	nameInput.SetWidth(inputWidth)

	descriptionInput := textinput.New()
	descriptionInput.Placeholder = "Optional description"
	descriptionInput.CharLimit = 500
	descriptionInput.SetWidth(inputWidth)

	keywordsInput := textinput.New()
	keywordsInput.Placeholder = "keyword1, keyword2, keyword3"
	keywordsInput.CharLimit = 500
	keywordsInput.SetWidth(inputWidth)

	examplesInput := textinput.New()
	examplesInput.Placeholder = "2h on project, 30m meeting"
	examplesInput.CharLimit = 500
	examplesInput.SetWidth(inputWidth)

	// Both add and edit modes now have all fields
	inputCount := 5 // Add mode: ID, Name, Description, Keywords, Examples
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
		width:            width,
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

func (m *TimeCodeEditModal) Update(msg tea.Msg) (*TimeCodeEditModal, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.Key().Code {
		case tea.KeyUp:
			m.prevInput()
			return m, nil
		case tea.KeyDown:
			m.nextInput()
			return m, nil
		case tea.KeyEscape:
			return m, func() tea.Msg { return ModalCloseMsg{} }
		case tea.KeyTab:
			if msg.Key().Mod&tea.ModShift != 0 {
				m.prevInput()
			} else {
				m.nextInput()
			}
			return m, nil
		case tea.KeyEnter:
			if !m.saving {
				m.saving = true
				return m, m.save()
			}
			return m, nil
		}
	}

	// Update only the focused input
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
		case 2:
			m.descriptionInput, cmd = m.descriptionInput.Update(msg)
		case 3:
			m.keywordsInput, cmd = m.keywordsInput.Update(msg)
		case 4:
			m.examplesInput, cmd = m.examplesInput.Update(msg)
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
		switch m.focusIndex {
		case 0:
			m.idInput.Focus()
		case 1:
			m.nameInput.Focus()
		case 2:
			m.descriptionInput.Focus()
		case 3:
			m.keywordsInput.Focus()
		case 4:
			m.examplesInput.Focus()
		}
	}
}

func (m *TimeCodeEditModal) save() tea.Cmd {
	return func() tea.Msg {
		if m.isEdit {
			name := m.nameInput.Value()
			desc := m.descriptionInput.Value()
			keywords := parseKeywords(m.keywordsInput.Value())
			examples := parseKeywords(m.examplesInput.Value())

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

		id := m.idInput.Value()
		name := m.nameInput.Value()
		if id == "" || name == "" {
			return timeCodeEditErrMsg{err: nil}
		}

		desc := m.descriptionInput.Value()
		keywords := parseKeywords(m.keywordsInput.Value())
		examples := parseKeywords(m.examplesInput.Value())

		_, err := m.client.CreateTimeCode(id, name, desc, keywords, examples)
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

	// Description field
	b.WriteString(labelStyle.Render("Description:"))
	b.WriteString("\n")
	descIdx := 2
	if m.isEdit {
		descIdx = 1
	}
	style = inputStyle
	if m.focusIndex == descIdx {
		style = focusedInputStyle
	}
	b.WriteString(style.Render(m.descriptionInput.View()))
	b.WriteString("\n\n")

	// Keywords field
	b.WriteString(labelStyle.Render("Keywords (comma separated):"))
	b.WriteString("\n")
	keywordsIdx := 3
	if m.isEdit {
		keywordsIdx = 2
	}
	style = inputStyle
	if m.focusIndex == keywordsIdx {
		style = focusedInputStyle
	}
	b.WriteString(style.Render(m.keywordsInput.View()))
	b.WriteString("\n\n")

	// Examples field
	b.WriteString(labelStyle.Render("Examples (comma separated):"))
	b.WriteString("\n")
	examplesIdx := 4
	if m.isEdit {
		examplesIdx = 3
	}
	style = inputStyle
	if m.focusIndex == examplesIdx {
		style = focusedInputStyle
	}
	b.WriteString(style.Render(m.examplesInput.View()))
	b.WriteString("\n\n")

	// Error
	if m.err != nil {
		b.WriteString(errorStyle.Render("Error: " + m.err.Error()))
		b.WriteString("\n\n")
	}

	// Help
	if m.saving {
		b.WriteString(statusStyle.Render("Saving..."))
	} else {
		b.WriteString(helpStyle.Render("[Tab] Next  [Enter] Save  [Esc] Cancel"))
	}

	return modalStyle.Render(b.String())
}
