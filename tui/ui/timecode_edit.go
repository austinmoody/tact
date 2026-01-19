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
	idInput         textinput.Model
	projectSelector *ProjectSelector
	nameInput       textinput.Model

	focusIndex int
	inputCount int

	saving bool
	err    error
}

func NewTimeCodeEditModal(client *api.Client, tc *model.TimeCode, projects []model.Project, width int) *TimeCodeEditModal {
	isEdit := tc != nil
	inputWidth := calculateInputWidth(width)

	idInput := textinput.New()
	idInput.Placeholder = "ABC123"
	idInput.CharLimit = 50
	idInput.SetWidth(inputWidth)

	// Initialize project selector with current project if editing
	selectedProjectID := ""
	if isEdit && tc.ProjectID != "" {
		selectedProjectID = tc.ProjectID
	} else if len(projects) > 0 {
		selectedProjectID = projects[0].ID
	}
	projectSelector := NewProjectSelector(projects, selectedProjectID, inputWidth)

	nameInput := textinput.New()
	nameInput.Placeholder = "Time Code Name"
	nameInput.CharLimit = 100
	nameInput.SetWidth(inputWidth)

	// Add mode: ID, Project, Name (3 fields)
	// Edit mode: Project, Name (2 fields, ID readonly)
	inputCount := 3
	if isEdit {
		inputCount = 2
		idInput.SetValue(tc.ID)
		nameInput.SetValue(tc.Name)
		projectSelector.Focus()
	} else {
		idInput.Focus()
	}

	return &TimeCodeEditModal{
		client:          client,
		timeCode:        tc,
		isEdit:          isEdit,
		width:           width,
		idInput:         idInput,
		projectSelector: projectSelector,
		nameInput:       nameInput,
		focusIndex:      0,
		inputCount:      inputCount,
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
			// If project selector is focused, let it handle up/down for selection
			if m.projectSelector.Focused() {
				m.projectSelector, _ = m.projectSelector.Update(msg)
				return m, nil
			}
			m.prevInput()
			return m, nil
		case tea.KeyDown:
			// If project selector is focused, let it handle up/down for selection
			if m.projectSelector.Focused() {
				m.projectSelector, _ = m.projectSelector.Update(msg)
				return m, nil
			}
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

		// Let project selector handle j/k for navigation when focused
		if m.projectSelector.Focused() {
			if msg.String() == "j" || msg.String() == "k" {
				m.projectSelector, _ = m.projectSelector.Update(msg)
				return m, nil
			}
		}
	}

	// Update only the focused input
	var cmd tea.Cmd
	if m.isEdit {
		// Edit mode: Project, Name
		switch m.focusIndex {
		case 0:
			m.projectSelector, cmd = m.projectSelector.Update(msg)
		case 1:
			m.nameInput, cmd = m.nameInput.Update(msg)
		}
	} else {
		// Add mode: ID, Project, Name
		switch m.focusIndex {
		case 0:
			m.idInput, cmd = m.idInput.Update(msg)
		case 1:
			m.projectSelector, cmd = m.projectSelector.Update(msg)
		case 2:
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
	m.projectSelector.Blur()
	m.nameInput.Blur()

	if m.isEdit {
		// Edit mode: Project, Name
		switch m.focusIndex {
		case 0:
			m.projectSelector.Focus()
		case 1:
			m.nameInput.Focus()
		}
	} else {
		// Add mode: ID, Project, Name
		switch m.focusIndex {
		case 0:
			m.idInput.Focus()
		case 1:
			m.projectSelector.Focus()
		case 2:
			m.nameInput.Focus()
		}
	}
}

func (m *TimeCodeEditModal) save() tea.Cmd {
	return func() tea.Msg {
		if m.isEdit {
			projectID := m.projectSelector.SelectedProjectID()
			name := m.nameInput.Value()

			updates := api.TimeCodeUpdate{
				ProjectID: &projectID,
				Name:      &name,
			}
			_, err := m.client.UpdateTimeCode(m.timeCode.ID, updates)
			if err != nil {
				return timeCodeEditErrMsg{err}
			}
			return TimeCodeUpdatedMsg{}
		}

		id := m.idInput.Value()
		projectID := m.projectSelector.SelectedProjectID()
		name := m.nameInput.Value()
		if id == "" || name == "" || projectID == "" {
			return timeCodeEditErrMsg{err: nil}
		}

		_, err := m.client.CreateTimeCode(id, projectID, name)
		if err != nil {
			return timeCodeEditErrMsg{err}
		}
		return TimeCodeCreatedMsg{}
	}
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

	// Project field
	b.WriteString(labelStyle.Render("Project:"))
	if m.projectSelector.Focused() {
		b.WriteString(helpStyle.Render(" (↑/↓ to select)"))
	}
	b.WriteString("\n")
	b.WriteString(m.projectSelector.View())
	b.WriteString("\n\n")

	// Name field
	b.WriteString(labelStyle.Render("Name:"))
	b.WriteString("\n")
	nameIdx := 2 // Add mode: ID(0), Project(1), Name(2)
	if m.isEdit {
		nameIdx = 1 // Edit mode: Project(0), Name(1)
	}
	style := inputStyle
	if m.focusIndex == nameIdx {
		style = focusedInputStyle
	}
	b.WriteString(style.Render(m.nameInput.View()))
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
