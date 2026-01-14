package ui

import (
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"

	"tact-tui/api"
	"tact-tui/model"
)

type ProjectEditModal struct {
	client  *api.Client
	project *model.Project // nil for add mode
	isEdit  bool
	width   int

	idInput          textinput.Model
	nameInput        textinput.Model
	descriptionInput textinput.Model
	focusIndex       int // 0=id, 1=name, 2=description

	saving bool
	err    error
}

func NewProjectEditModal(client *api.Client, p *model.Project, width int) *ProjectEditModal {
	isEdit := p != nil
	inputWidth := calculateInputWidth(width)

	idInput := textinput.New()
	idInput.Placeholder = "project-id"
	idInput.CharLimit = 50
	idInput.SetWidth(inputWidth)

	nameInput := textinput.New()
	nameInput.Placeholder = "Project Name"
	nameInput.CharLimit = 100
	nameInput.SetWidth(inputWidth)

	descriptionInput := textinput.New()
	descriptionInput.Placeholder = "Description (optional)"
	descriptionInput.CharLimit = 500
	descriptionInput.SetWidth(inputWidth)

	if isEdit {
		idInput.SetValue(p.ID)
		nameInput.SetValue(p.Name)
		descriptionInput.SetValue(p.Description)
		// Focus name input for edit mode (ID is not editable)
		nameInput.Focus()
	} else {
		// Focus ID input for add mode
		idInput.Focus()
	}

	return &ProjectEditModal{
		client:           client,
		project:          p,
		isEdit:           isEdit,
		width:            width,
		idInput:          idInput,
		nameInput:        nameInput,
		descriptionInput: descriptionInput,
		focusIndex:       0,
	}
}

func (m *ProjectEditModal) Init() tea.Cmd {
	return textinput.Blink
}

func (m *ProjectEditModal) Update(msg tea.Msg) (*ProjectEditModal, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.Key().Code {
		case tea.KeyEscape:
			return m, func() tea.Msg { return ModalCloseMsg{} }
		case tea.KeyEnter:
			if !m.saving && m.isValid() {
				m.saving = true
				return m, m.save()
			}
			return m, nil
		case tea.KeyTab:
			if msg.Key().Mod&tea.ModShift != 0 {
				return m.prevField()
			}
			return m.nextField()
		case tea.KeyUp:
			return m.prevField()
		case tea.KeyDown:
			return m.nextField()
		}
	}

	return m.updateInputs(msg)
}

func (m *ProjectEditModal) isValid() bool {
	if m.isEdit {
		return m.nameInput.Value() != ""
	}
	return m.idInput.Value() != "" && m.nameInput.Value() != ""
}

func (m *ProjectEditModal) nextField() (*ProjectEditModal, tea.Cmd) {
	// In edit mode, skip ID field (index 0)
	if m.isEdit {
		if m.focusIndex == 1 {
			m.focusIndex = 2
		} else {
			m.focusIndex = 1
		}
	} else {
		m.focusIndex = (m.focusIndex + 1) % 3
	}
	return m, m.updateFocus()
}

func (m *ProjectEditModal) prevField() (*ProjectEditModal, tea.Cmd) {
	// In edit mode, skip ID field (index 0)
	if m.isEdit {
		if m.focusIndex == 2 {
			m.focusIndex = 1
		} else {
			m.focusIndex = 2
		}
	} else {
		m.focusIndex--
		if m.focusIndex < 0 {
			m.focusIndex = 2
		}
	}
	return m, m.updateFocus()
}

func (m *ProjectEditModal) updateFocus() tea.Cmd {
	m.idInput.Blur()
	m.nameInput.Blur()
	m.descriptionInput.Blur()

	switch m.focusIndex {
	case 0:
		m.idInput.Focus()
	case 1:
		m.nameInput.Focus()
	case 2:
		m.descriptionInput.Focus()
	}

	return textinput.Blink
}

func (m *ProjectEditModal) updateInputs(msg tea.Msg) (*ProjectEditModal, tea.Cmd) {
	var cmds []tea.Cmd

	var cmd tea.Cmd
	m.idInput, cmd = m.idInput.Update(msg)
	cmds = append(cmds, cmd)

	m.nameInput, cmd = m.nameInput.Update(msg)
	cmds = append(cmds, cmd)

	m.descriptionInput, cmd = m.descriptionInput.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *ProjectEditModal) save() tea.Cmd {
	return func() tea.Msg {
		id := m.idInput.Value()
		name := m.nameInput.Value()
		description := m.descriptionInput.Value()

		if m.isEdit {
			updates := api.ProjectUpdate{
				Name:        &name,
				Description: &description,
			}
			_, err := m.client.UpdateProject(m.project.ID, updates)
			if err != nil {
				return projectEditErrMsg{err}
			}
			return ProjectUpdatedMsg{}
		}

		_, err := m.client.CreateProject(id, name, description)
		if err != nil {
			return projectEditErrMsg{err}
		}
		return ProjectCreatedMsg{}
	}
}

type projectEditErrMsg struct{ err error }

func (m *ProjectEditModal) View() string {
	var b strings.Builder

	if m.isEdit {
		b.WriteString(modalTitleStyle.Render("Edit Project"))
	} else {
		b.WriteString(modalTitleStyle.Render("Add Project"))
	}
	b.WriteString("\n\n")

	// ID field (disabled in edit mode)
	b.WriteString(labelStyle.Render("ID:"))
	b.WriteString("\n")
	if m.isEdit {
		b.WriteString(disabledInputStyle.Render(m.idInput.Value()))
	} else if m.focusIndex == 0 {
		b.WriteString(focusedInputStyle.Render(m.idInput.View()))
	} else {
		b.WriteString(inputStyle.Render(m.idInput.View()))
	}
	b.WriteString("\n\n")

	// Name field
	b.WriteString(labelStyle.Render("Name:"))
	b.WriteString("\n")
	if m.focusIndex == 1 {
		b.WriteString(focusedInputStyle.Render(m.nameInput.View()))
	} else {
		b.WriteString(inputStyle.Render(m.nameInput.View()))
	}
	b.WriteString("\n\n")

	// Description field
	b.WriteString(labelStyle.Render("Description:"))
	b.WriteString("\n")
	if m.focusIndex == 2 {
		b.WriteString(focusedInputStyle.Render(m.descriptionInput.View()))
	} else {
		b.WriteString(inputStyle.Render(m.descriptionInput.View()))
	}
	b.WriteString("\n\n")

	if m.err != nil {
		b.WriteString(errorStyle.Render("Error: " + m.err.Error()))
		b.WriteString("\n\n")
	}

	if m.saving {
		b.WriteString(statusStyle.Render("Saving..."))
	} else {
		b.WriteString(helpStyle.Render("[Tab] Next Field  [Enter] Save  [Esc] Cancel"))
	}

	return modalStyle.Render(b.String())
}
