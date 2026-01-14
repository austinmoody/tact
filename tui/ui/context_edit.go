package ui

import (
	"strings"

	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"

	"tact-tui/api"
	"tact-tui/model"
)

type ContextEditModal struct {
	client  *api.Client
	owner   ContextOwner
	context *model.ContextDocument // nil for add mode
	isEdit  bool
	width   int
	height  int

	textarea textarea.Model

	saving bool
	err    error
}

func NewContextEditModal(client *api.Client, owner ContextOwner, ctx *model.ContextDocument, width, height int) *ContextEditModal {
	isEdit := ctx != nil

	ta := textarea.New()
	ta.Placeholder = "Enter context content here...\n\nThis content will be used to help parse time entries."
	ta.CharLimit = 5000
	ta.SetWidth(min(70, width-10))
	ta.SetHeight(min(10, height-15))
	ta.Focus()

	if isEdit {
		ta.SetValue(ctx.Content)
	}

	return &ContextEditModal{
		client:   client,
		owner:    owner,
		context:  ctx,
		isEdit:   isEdit,
		width:    width,
		height:   height,
		textarea: ta,
	}
}

func (m *ContextEditModal) Init() tea.Cmd {
	return textarea.Blink
}

func (m *ContextEditModal) Update(msg tea.Msg) (*ContextEditModal, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		// Check for Ctrl+S to save
		if msg.Key().Code == 's' && msg.Key().Mod&tea.ModCtrl != 0 {
			if !m.saving && m.textarea.Value() != "" {
				m.saving = true
				return m, m.save()
			}
			return m, nil
		}

		switch msg.Key().Code {
		case tea.KeyEscape:
			return m, func() tea.Msg { return ModalCloseMsg{} }
		}
	}

	// Pass all other input to textarea
	var cmd tea.Cmd
	m.textarea, cmd = m.textarea.Update(msg)
	return m, cmd
}

func (m *ContextEditModal) save() tea.Cmd {
	return func() tea.Msg {
		content := strings.TrimSpace(m.textarea.Value())
		if content == "" {
			return contextEditErrMsg{err: nil}
		}

		if m.isEdit {
			_, err := m.client.UpdateContext(m.context.ID, content)
			if err != nil {
				return contextEditErrMsg{err}
			}
			return ContextUpdatedMsg{}
		}

		// Create new context
		var err error
		if m.owner.ProjectID != nil {
			_, err = m.client.CreateProjectContext(*m.owner.ProjectID, content)
		} else if m.owner.TimeCodeID != nil {
			_, err = m.client.CreateTimeCodeContext(*m.owner.TimeCodeID, content)
		}

		if err != nil {
			return contextEditErrMsg{err}
		}
		return ContextCreatedMsg{}
	}
}

type contextEditErrMsg struct{ err error }

func (m *ContextEditModal) View() string {
	var b strings.Builder

	if m.isEdit {
		b.WriteString(modalTitleStyle.Render("Edit Context"))
	} else {
		b.WriteString(modalTitleStyle.Render("Add Context"))
	}
	b.WriteString("\n\n")

	b.WriteString(labelStyle.Render("Content:"))
	b.WriteString("\n")
	b.WriteString(focusedInputStyle.Render(m.textarea.View()))
	b.WriteString("\n\n")

	if m.err != nil {
		b.WriteString(errorStyle.Render("Error: " + m.err.Error()))
		b.WriteString("\n\n")
	}

	if m.saving {
		b.WriteString(statusStyle.Render("Saving..."))
	} else {
		b.WriteString(helpStyle.Render("[Ctrl+S] Save  [Esc] Cancel"))
	}

	return modalStyle.Render(b.String())
}
