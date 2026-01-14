package ui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"tact-tui/model"
)

// ProjectSelector is a component for selecting a project from a list.
type ProjectSelector struct {
	projects     []model.Project
	selectedIdx  int
	focused      bool
	width        int
	maxVisible   int
	scrollOffset int
}

func NewProjectSelector(projects []model.Project, selectedProjectID string, width int) *ProjectSelector {
	selectedIdx := 0
	for i, p := range projects {
		if p.ID == selectedProjectID {
			selectedIdx = i
			break
		}
	}

	return &ProjectSelector{
		projects:    projects,
		selectedIdx: selectedIdx,
		width:       width,
		maxVisible:  5,
	}
}

func (s *ProjectSelector) Focus() {
	s.focused = true
}

func (s *ProjectSelector) Blur() {
	s.focused = false
}

func (s *ProjectSelector) Focused() bool {
	return s.focused
}

func (s *ProjectSelector) SelectedProject() *model.Project {
	if len(s.projects) == 0 {
		return nil
	}
	return &s.projects[s.selectedIdx]
}

func (s *ProjectSelector) SelectedProjectID() string {
	if p := s.SelectedProject(); p != nil {
		return p.ID
	}
	return ""
}

func (s *ProjectSelector) Update(msg tea.Msg) (*ProjectSelector, tea.Cmd) {
	if !s.focused {
		return s, nil
	}

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.Key().Code {
		case tea.KeyUp:
			if s.selectedIdx > 0 {
				s.selectedIdx--
				s.adjustScroll()
			}
			return s, nil
		case tea.KeyDown:
			if s.selectedIdx < len(s.projects)-1 {
				s.selectedIdx++
				s.adjustScroll()
			}
			return s, nil
		}

		// Also handle j/k for vim-style navigation
		switch msg.String() {
		case "k":
			if s.selectedIdx > 0 {
				s.selectedIdx--
				s.adjustScroll()
			}
			return s, nil
		case "j":
			if s.selectedIdx < len(s.projects)-1 {
				s.selectedIdx++
				s.adjustScroll()
			}
			return s, nil
		}
	}

	return s, nil
}

func (s *ProjectSelector) adjustScroll() {
	// Ensure selected item is visible
	if s.selectedIdx < s.scrollOffset {
		s.scrollOffset = s.selectedIdx
	} else if s.selectedIdx >= s.scrollOffset+s.maxVisible {
		s.scrollOffset = s.selectedIdx - s.maxVisible + 1
	}
}

func (s *ProjectSelector) View() string {
	if len(s.projects) == 0 {
		return statusStyle.Render("  (no projects available)")
	}

	var b strings.Builder

	// Determine visible range
	start := s.scrollOffset
	end := start + s.maxVisible
	if end > len(s.projects) {
		end = len(s.projects)
	}

	// Show scroll indicator if needed
	if s.scrollOffset > 0 {
		b.WriteString(helpStyle.Render("  ↑ more"))
		b.WriteString("\n")
	}

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("170")).
		Bold(true)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	for i := start; i < end; i++ {
		p := s.projects[i]
		prefix := "  "
		style := normalStyle

		if i == s.selectedIdx {
			prefix = "▸ "
			if s.focused {
				style = selectedStyle
			}
		}

		line := prefix + p.ID + " - " + p.Name
		// Truncate if too long
		maxLen := s.width - 4
		if len(line) > maxLen && maxLen > 3 {
			line = line[:maxLen-3] + "..."
		}

		b.WriteString(style.Render(line))
		if i < end-1 {
			b.WriteString("\n")
		}
	}

	// Show scroll indicator if needed
	if end < len(s.projects) {
		b.WriteString("\n")
		b.WriteString(helpStyle.Render("  ↓ more"))
	}

	return b.String()
}
