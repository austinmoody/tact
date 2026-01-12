package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"tact-tui/api"
	"tact-tui/model"
)

type WorkTypeEditModal struct {
	client   *api.Client
	workType *model.WorkType // nil for add mode
	isEdit   bool

	nameInput textinput.Model

	saving bool
	err    error
}

func NewWorkTypeEditModal(client *api.Client, wt *model.WorkType) *WorkTypeEditModal {
	isEdit := wt != nil

	nameInput := textinput.New()
	nameInput.Placeholder = "Work Type Name"
	nameInput.CharLimit = 100
	nameInput.Width = 40
	nameInput.Focus()

	if isEdit {
		nameInput.SetValue(wt.Name)
	}

	return &WorkTypeEditModal{
		client:    client,
		workType:  wt,
		isEdit:    isEdit,
		nameInput: nameInput,
	}
}

func (m *WorkTypeEditModal) Init() tea.Cmd {
	return textinput.Blink
}

func (m *WorkTypeEditModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle keys directly by type to prevent control characters
		switch msg.Type {
		case tea.KeyUp, tea.KeyDown, tea.KeyShiftUp, tea.KeyShiftDown, tea.KeyCtrlUp, tea.KeyCtrlDown:
			return m, nil
		case tea.KeyEsc:
			return m, func() tea.Msg { return ModalCloseMsg{} }
		case tea.KeyEnter:
			if !m.saving && m.nameInput.Value() != "" {
				m.saving = true
				return m, m.save()
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
			// Update existing work type
			updates := api.WorkTypeUpdate{Name: &name}
			_, err := m.client.UpdateWorkType(m.workType.ID, updates)
			if err != nil {
				return workTypeEditErrMsg{err}
			}
			return WorkTypeUpdatedMsg{}
		}

		// Create new work type
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

	// Name field
	b.WriteString(labelStyle.Render("Name:"))
	b.WriteString("\n")
	b.WriteString(focusedInputStyle.Render(m.nameInput.View()))
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
		b.WriteString(helpStyle.Render("[Enter] Save  [Esc] Cancel"))
	}

	return modalStyle.Render(b.String())
}
