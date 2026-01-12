package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type menuItem struct {
	key   string
	label string
}

type MenuModal struct {
	items  []menuItem
	cursor int
}

func NewMenuModal() *MenuModal {
	return &MenuModal{
		items: []menuItem{
			{key: "timecodes", label: "Time Codes"},
			{key: "worktypes", label: "Work Types"},
		},
		cursor: 0,
	}
}

func (m *MenuModal) Init() tea.Cmd {
	return nil
}

func (m *MenuModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case matchesKey(msg, keys.Escape):
			return m, func() tea.Msg { return ModalCloseMsg{} }

		case matchesKey(msg, keys.Up):
			if m.cursor > 0 {
				m.cursor--
			}
			return m, nil

		case matchesKey(msg, keys.Down):
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
			return m, nil

		case matchesKey(msg, keys.Enter):
			if m.cursor < len(m.items) {
				selection := m.items[m.cursor].key
				return m, func() tea.Msg { return MenuSelectMsg{Selection: selection} }
			}
			return m, nil
		}
	}

	return m, nil
}

func (m *MenuModal) View() string {
	var b strings.Builder

	b.WriteString(modalTitleStyle.Render("Menu"))
	b.WriteString("\n\n")

	for i, item := range m.items {
		cursor := "  "
		style := itemStyle
		if i == m.cursor {
			cursor = "> "
			style = selectedItemStyle
		}
		b.WriteString(style.Render(cursor + item.label))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("[Enter] Select  [Esc] Cancel"))

	return modalStyle.Render(b.String())
}
