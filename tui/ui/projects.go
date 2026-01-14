package ui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"tact-tui/api"
	"tact-tui/model"
)

type ProjectsScreen struct {
	client   *api.Client
	projects []model.Project
	cursor   int
	loading  bool
	err      error
	width    int
	height   int
}

type projectsMsg struct{ projects []model.Project }
type projectsErrMsg struct{ err error }

func NewProjectsScreen(client *api.Client) *ProjectsScreen {
	return &ProjectsScreen{
		client:  client,
		loading: true,
	}
}

func (s *ProjectsScreen) Init() tea.Cmd {
	return s.fetchProjects()
}

func (s *ProjectsScreen) Refresh() tea.Cmd {
	s.loading = true
	s.err = nil
	return s.fetchProjects()
}

func (s *ProjectsScreen) fetchProjects() tea.Cmd {
	return func() tea.Msg {
		projects, err := s.client.FetchProjects()
		if err != nil {
			return projectsErrMsg{err}
		}
		return projectsMsg{projects}
	}
}

func (s *ProjectsScreen) SetSize(width, height int) {
	s.width = width
	s.height = height
}

func (s *ProjectsScreen) SelectedProject() *model.Project {
	if len(s.projects) == 0 || s.cursor >= len(s.projects) {
		return nil
	}
	return &s.projects[s.cursor]
}

func (s *ProjectsScreen) Update(msg tea.Msg) (*ProjectsScreen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		return s.handleKeyPress(msg)

	case projectsMsg:
		s.projects = msg.projects
		s.loading = false
		if s.cursor >= len(s.projects) {
			s.cursor = max(0, len(s.projects)-1)
		}
		return s, nil

	case projectsErrMsg:
		s.err = msg.err
		s.loading = false
		return s, nil
	}

	return s, nil
}

func (s *ProjectsScreen) handleKeyPress(msg tea.KeyPressMsg) (*ProjectsScreen, tea.Cmd) {
	switch {
	case matchesKey(msg, keys.Up):
		if s.cursor > 0 {
			s.cursor--
		}
		return s, nil

	case matchesKey(msg, keys.Down):
		if s.cursor < len(s.projects)-1 {
			s.cursor++
		}
		return s, nil

	case matchesKey(msg, keys.Add):
		return s, func() tea.Msg { return OpenProjectAddMsg{} }

	case matchesKey(msg, keys.Edit):
		if p := s.SelectedProject(); p != nil {
			return s, func() tea.Msg { return OpenProjectEditMsg{Project: p} }
		}
		return s, nil

	case matchesKey(msg, keys.Delete):
		if p := s.SelectedProject(); p != nil {
			return s, s.deleteProject(p.ID)
		}
		return s, nil

	case matchesKey(msg, keys.Context):
		if p := s.SelectedProject(); p != nil {
			return s, func() tea.Msg { return OpenProjectContextMsg{Project: p} }
		}
		return s, nil

	case matchesKey(msg, keys.Refresh):
		return s, s.Refresh()
	}

	return s, nil
}

func (s *ProjectsScreen) deleteProject(id string) tea.Cmd {
	return func() tea.Msg {
		if err := s.client.DeleteProject(id); err != nil {
			return projectsErrMsg{err}
		}
		return ProjectDeletedMsg{}
	}
}

func (s *ProjectsScreen) View() string {
	if s.width == 0 {
		return "Loading..."
	}

	var b strings.Builder

	// Title bar
	title := titleStyle.Render("Projects")
	hint := helpStyle.Render("[Esc] Back")
	titleBar := lipgloss.JoinHorizontal(
		lipgloss.Top,
		title,
		strings.Repeat(" ", max(0, s.width-lipgloss.Width(title)-lipgloss.Width(hint)-2)),
		hint,
	)
	b.WriteString(titleBar + "\n\n")

	// Header
	b.WriteString(headerStyle.Render("Manage Projects") + "\n")
	b.WriteString(strings.Repeat("─", min(60, s.width-4)) + "\n")

	if s.loading {
		b.WriteString(statusStyle.Render("Loading...") + "\n")
	} else if s.err != nil {
		b.WriteString(errorStyle.Render(fmt.Sprintf("Error: %v", s.err)) + "\n")
	} else if len(s.projects) == 0 {
		b.WriteString(statusStyle.Render("No projects. Press [a] to add one.") + "\n")
	} else {
		for i, p := range s.projects {
			line := s.renderProjectLine(i, p)
			b.WriteString(line + "\n")
		}
	}

	// Help bar at bottom
	b.WriteString("\n")
	help := helpStyle.Render("[a] Add  [e] Edit  [d] Delete  [c] Context  [r] Refresh  [Esc] Back")
	b.WriteString(help)

	return b.String()
}

func (s *ProjectsScreen) renderProjectLine(index int, p model.Project) string {
	cursor := "  "
	style := itemStyle
	if index == s.cursor {
		cursor = "> "
		style = selectedItemStyle
	}

	// Show ID and Name
	status := ""
	if !p.Active {
		status = inactiveStyle.Render(" [inactive]")
	}

	id := p.ID
	if len(id) > 12 {
		id = id[:12]
	}

	name := p.Name
	if len(name) > 30 {
		name = name[:27] + "..."
	}

	line := fmt.Sprintf("%s%-12s  %-30s%s", cursor, id, name, status)
	return style.Render(line)
}
