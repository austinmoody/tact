package ui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"

	"tact-tui/api"
	"tact-tui/model"
)

// ContextOwner represents what owns the context (project or time code)
type ContextOwner struct {
	ProjectID  *string
	TimeCodeID *string
	Name       string // Display name
}

type ContextListModal struct {
	client   *api.Client
	owner    ContextOwner
	contexts []model.ContextDocument
	cursor   int
	loading  bool
	err      error
	width    int
}

type contextListMsg struct{ contexts []model.ContextDocument }
type contextListErrMsg struct{ err error }

func NewContextListModal(client *api.Client, owner ContextOwner, width int) *ContextListModal {
	return &ContextListModal{
		client:  client,
		owner:   owner,
		loading: true,
		width:   width,
	}
}

func (m *ContextListModal) Init() tea.Cmd {
	return m.fetchContexts()
}

func (m *ContextListModal) Refresh() tea.Cmd {
	m.loading = true
	m.err = nil
	return m.fetchContexts()
}

func (m *ContextListModal) fetchContexts() tea.Cmd {
	return func() tea.Msg {
		var contexts []model.ContextDocument
		var err error

		if m.owner.ProjectID != nil {
			contexts, err = m.client.FetchProjectContext(*m.owner.ProjectID)
		} else if m.owner.TimeCodeID != nil {
			contexts, err = m.client.FetchTimeCodeContext(*m.owner.TimeCodeID)
		}

		if err != nil {
			return contextListErrMsg{err}
		}
		return contextListMsg{contexts}
	}
}

func (m *ContextListModal) SelectedContext() *model.ContextDocument {
	if len(m.contexts) == 0 || m.cursor >= len(m.contexts) {
		return nil
	}
	return &m.contexts[m.cursor]
}

func (m *ContextListModal) Update(msg tea.Msg) (*ContextListModal, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		return m.handleKeyPress(msg)

	case contextListMsg:
		m.contexts = msg.contexts
		m.loading = false
		if m.cursor >= len(m.contexts) {
			m.cursor = max(0, len(m.contexts)-1)
		}
		return m, nil

	case contextListErrMsg:
		m.err = msg.err
		m.loading = false
		return m, nil
	}

	return m, nil
}

func (m *ContextListModal) handleKeyPress(msg tea.KeyPressMsg) (*ContextListModal, tea.Cmd) {
	switch {
	case matchesKey(msg, keys.Escape):
		return m, func() tea.Msg { return ModalCloseMsg{} }

	case matchesKey(msg, keys.Up):
		if m.cursor > 0 {
			m.cursor--
		}
		return m, nil

	case matchesKey(msg, keys.Down):
		if m.cursor < len(m.contexts)-1 {
			m.cursor++
		}
		return m, nil

	case matchesKey(msg, keys.Add):
		return m, func() tea.Msg { return OpenContextAddMsg{Owner: m.owner} }

	case matchesKey(msg, keys.Edit):
		if ctx := m.SelectedContext(); ctx != nil {
			return m, func() tea.Msg { return OpenContextEditMsg{Context: ctx, Owner: m.owner} }
		}
		return m, nil

	case matchesKey(msg, keys.Delete):
		if ctx := m.SelectedContext(); ctx != nil {
			return m, m.deleteContext(ctx.ID)
		}
		return m, nil

	case matchesKey(msg, keys.Refresh):
		return m, m.Refresh()
	}

	return m, nil
}

func (m *ContextListModal) deleteContext(id string) tea.Cmd {
	return func() tea.Msg {
		if err := m.client.DeleteContext(id); err != nil {
			return contextListErrMsg{err}
		}
		return ContextDeletedMsg{}
	}
}

func (m *ContextListModal) View() string {
	var b strings.Builder

	title := fmt.Sprintf("Context for: %s", m.owner.Name)
	b.WriteString(modalTitleStyle.Render(title))
	b.WriteString("\n\n")

	if m.loading {
		b.WriteString(statusStyle.Render("Loading..."))
		b.WriteString("\n")
	} else if m.err != nil {
		b.WriteString(errorStyle.Render(fmt.Sprintf("Error: %v", m.err)))
		b.WriteString("\n")
	} else if len(m.contexts) == 0 {
		b.WriteString(statusStyle.Render("No context documents. Press [a] to add one."))
		b.WriteString("\n")
	} else {
		maxWidth := min(60, m.width-10)
		for i, ctx := range m.contexts {
			line := m.renderContextLine(i, ctx, maxWidth)
			b.WriteString(line + "\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("[a] Add  [e] Edit  [d] Delete  [r] Refresh  [Esc] Close"))

	return modalStyle.Render(b.String())
}

func (m *ContextListModal) renderContextLine(index int, ctx model.ContextDocument, maxWidth int) string {
	cursor := "  "
	style := itemStyle
	if index == m.cursor {
		cursor = "> "
		style = selectedItemStyle
	}

	// Truncate content preview
	content := strings.ReplaceAll(ctx.Content, "\n", " ")
	if len(content) > maxWidth-4 {
		content = content[:maxWidth-7] + "..."
	}

	line := fmt.Sprintf("%s%s", cursor, content)
	return style.Render(line)
}
