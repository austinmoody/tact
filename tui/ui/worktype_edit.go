package ui

import (
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"

	"tact-tui/api"
	"tact-tui/model"
)

type WorkTypeEditModal struct {
	client   *api.Client
	workType *model.WorkType // nil for add mode
	isEdit   bool
	width    int

	nameInput textinput.Model

	saving bool
	err    error
}

func NewWorkTypeEditModal(client *api.Client, wt *model.WorkType, width int) *WorkTypeEditModal {
	isEdit := wt != nil
	inputWidth := calculateInputWidth(width)

	nameInput := textinput.New()
	nameInput.Placeholder = "Work Type Name"
	nameInput.CharLimit = 100
	nameInput.SetWidth(inputWidth)
	nameInput.Focus()

	if isEdit {
		nameInput.SetValue(wt.Name)
	}

	return &WorkTypeEditModal{
		client:    client,
		workType:  wt,
		isEdit:    isEdit,
		width:     width,
		nameInput: nameInput,
	}
}

func (m *WorkTypeEditModal) Init() tea.Cmd {
	return textinput.Blink
}

func (m *WorkTypeEditModal) Update(msg tea.Msg) (*WorkTypeEditModal, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.Key().Code {
		case tea.KeyEscape:
			return m, func() tea.Msg { return ModalCloseMsg{} }
		case tea.KeyEnter:
			if !m.saving && m.nameInput.Value() != "" {
				m.saving = true
				return m, m.save()
			}
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.nameInput, cmd = m.nameInput.Update(msg)
	return m, cmd
}

func (m *WorkTypeEditModal) save() tea.Cmd {
	return func() tea.Msg {
		name := m.nameInput.Value()
		if name == "" {
			return workTypeEditErrMsg{err: nil}
		}

		if m.isEdit {
			updates := api.WorkTypeUpdate{Name: &name}
			_, err := m.client.UpdateWorkType(m.workType.ID, updates)
			if err != nil {
				return workTypeEditErrMsg{err}
			}
			return WorkTypeUpdatedMsg{}
		}

		_, err := m.client.CreateWorkType(name)
		if err != nil {
			return workTypeEditErrMsg{err}
		}
		return WorkTypeCreatedMsg{}
	}
}

type workTypeEditErrMsg struct{ err error }

func (m *WorkTypeEditModal) View() string {
	var b strings.Builder

	if m.isEdit {
		b.WriteString(modalTitleStyle.Render("Edit Work Type"))
	} else {
		b.WriteString(modalTitleStyle.Render("Add Work Type"))
	}
	b.WriteString("\n\n")

	b.WriteString(labelStyle.Render("Name:"))
	b.WriteString("\n")
	b.WriteString(focusedInputStyle.Render(m.nameInput.View()))
	b.WriteString("\n\n")

	if m.err != nil {
		b.WriteString(errorStyle.Render("Error: " + m.err.Error()))
		b.WriteString("\n\n")
	}

	if m.saving {
		b.WriteString(statusStyle.Render("Saving..."))
	} else {
		b.WriteString(helpStyle.Render("[Enter] Save  [Esc] Cancel"))
	}

	return modalStyle.Render(b.String())
}
