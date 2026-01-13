package ui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"tact-tui/api"
	"tact-tui/model"
)

type WorkTypesScreen struct {
	client    *api.Client
	workTypes []model.WorkType
	cursor    int
	loading   bool
	err       error
	width     int
	height    int
}

type workTypesMsg struct{ workTypes []model.WorkType }
type workTypesErrMsg struct{ err error }

func NewWorkTypesScreen(client *api.Client) *WorkTypesScreen {
	return &WorkTypesScreen{
		client:  client,
		loading: true,
	}
}

func (s *WorkTypesScreen) Init() tea.Cmd {
	return s.fetchWorkTypes()
}

func (s *WorkTypesScreen) Refresh() tea.Cmd {
	s.loading = true
	s.err = nil
	return s.fetchWorkTypes()
}

func (s *WorkTypesScreen) fetchWorkTypes() tea.Cmd {
	return func() tea.Msg {
		workTypes, err := s.client.FetchWorkTypes()
		if err != nil {
			return workTypesErrMsg{err}
		}
		return workTypesMsg{workTypes}
	}
}

func (s *WorkTypesScreen) SetSize(width, height int) {
	s.width = width
	s.height = height
}

func (s *WorkTypesScreen) SelectedWorkType() *model.WorkType {
	if len(s.workTypes) == 0 || s.cursor >= len(s.workTypes) {
		return nil
	}
	return &s.workTypes[s.cursor]
}

func (s *WorkTypesScreen) Update(msg tea.Msg) (*WorkTypesScreen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		return s.handleKeyPress(msg)

	case workTypesMsg:
		s.workTypes = msg.workTypes
		s.loading = false
		if s.cursor >= len(s.workTypes) {
			s.cursor = max(0, len(s.workTypes)-1)
		}
		return s, nil

	case workTypesErrMsg:
		s.err = msg.err
		s.loading = false
		return s, nil
	}

	return s, nil
}

func (s *WorkTypesScreen) handleKeyPress(msg tea.KeyPressMsg) (*WorkTypesScreen, tea.Cmd) {
	switch {
	case matchesKey(msg, keys.Up):
		if s.cursor > 0 {
			s.cursor--
		}
		return s, nil

	case matchesKey(msg, keys.Down):
		if s.cursor < len(s.workTypes)-1 {
			s.cursor++
		}
		return s, nil

	case matchesKey(msg, keys.Add):
		return s, func() tea.Msg { return OpenWorkTypeAddMsg{} }

	case matchesKey(msg, keys.Edit):
		if wt := s.SelectedWorkType(); wt != nil {
			return s, func() tea.Msg { return OpenWorkTypeEditMsg{WorkType: wt} }
		}
		return s, nil

	case matchesKey(msg, keys.Delete):
		if wt := s.SelectedWorkType(); wt != nil {
			return s, s.deleteWorkType(wt.ID)
		}
		return s, nil

	case matchesKey(msg, keys.Refresh):
		return s, s.Refresh()
	}

	return s, nil
}

func (s *WorkTypesScreen) deleteWorkType(id string) tea.Cmd {
	return func() tea.Msg {
		if err := s.client.DeleteWorkType(id); err != nil {
			return workTypesErrMsg{err}
		}
		return WorkTypeDeletedMsg{}
	}
}

func (s *WorkTypesScreen) View() string {
	if s.width == 0 {
		return "Loading..."
	}

	var b strings.Builder

	// Title bar
	title := titleStyle.Render("Work Types")
	hint := helpStyle.Render("[Esc] Back")
	titleBar := lipgloss.JoinHorizontal(
		lipgloss.Top,
		title,
		strings.Repeat(" ", max(0, s.width-lipgloss.Width(title)-lipgloss.Width(hint)-2)),
		hint,
	)
	b.WriteString(titleBar + "\n\n")

	// Header
	b.WriteString(headerStyle.Render("Manage Work Types") + "\n")
	b.WriteString(strings.Repeat("─", min(60, s.width-4)) + "\n")

	if s.loading {
		b.WriteString(statusStyle.Render("Loading...") + "\n")
	} else if s.err != nil {
		b.WriteString(errorStyle.Render(fmt.Sprintf("Error: %v", s.err)) + "\n")
	} else if len(s.workTypes) == 0 {
		b.WriteString(statusStyle.Render("No work types. Press [a] to add one.") + "\n")
	} else {
		for i, wt := range s.workTypes {
			line := s.renderWorkTypeLine(i, wt)
			b.WriteString(line + "\n")
		}
	}

	// Help bar at bottom
	b.WriteString("\n")
	help := helpStyle.Render("[a] Add  [e] Edit  [d] Delete  [r] Refresh  [Esc] Back")
	b.WriteString(help)

	return b.String()
}

func (s *WorkTypesScreen) renderWorkTypeLine(index int, wt model.WorkType) string {
	cursor := "  "
	style := itemStyle
	if index == s.cursor {
		cursor = "> "
		style = selectedItemStyle
	}

	// Show ID and Name
	status := ""
	if !wt.Active {
		status = inactiveStyle.Render(" [inactive]")
	}

	name := wt.Name
	if len(name) > 40 {
		name = name[:37] + "..."
	}

	line := fmt.Sprintf("%s%-40s%s", cursor, name, status)
	return style.Render(line)
}
