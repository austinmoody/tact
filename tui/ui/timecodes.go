package ui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"tact-tui/api"
	"tact-tui/model"
)

type TimeCodesScreen struct {
	client    *api.Client
	timeCodes []model.TimeCode
	cursor    int
	loading   bool
	err       error
	width     int
	height    int
}

type timeCodesMsg struct{ timeCodes []model.TimeCode }
type timeCodesErrMsg struct{ err error }

func NewTimeCodesScreen(client *api.Client) *TimeCodesScreen {
	return &TimeCodesScreen{
		client:  client,
		loading: true,
	}
}

func (s *TimeCodesScreen) Init() tea.Cmd {
	return s.fetchTimeCodes()
}

func (s *TimeCodesScreen) Refresh() tea.Cmd {
	s.loading = true
	s.err = nil
	return s.fetchTimeCodes()
}

func (s *TimeCodesScreen) fetchTimeCodes() tea.Cmd {
	return func() tea.Msg {
		timeCodes, err := s.client.FetchTimeCodes()
		if err != nil {
			return timeCodesErrMsg{err}
		}
		return timeCodesMsg{timeCodes}
	}
}

func (s *TimeCodesScreen) SetSize(width, height int) {
	s.width = width
	s.height = height
}

func (s *TimeCodesScreen) SelectedTimeCode() *model.TimeCode {
	if len(s.timeCodes) == 0 || s.cursor >= len(s.timeCodes) {
		return nil
	}
	return &s.timeCodes[s.cursor]
}

func (s *TimeCodesScreen) Update(msg tea.Msg) (*TimeCodesScreen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		return s.handleKeyPress(msg)

	case timeCodesMsg:
		s.timeCodes = msg.timeCodes
		s.loading = false
		if s.cursor >= len(s.timeCodes) {
			s.cursor = max(0, len(s.timeCodes)-1)
		}
		return s, nil

	case timeCodesErrMsg:
		s.err = msg.err
		s.loading = false
		return s, nil
	}

	return s, nil
}

func (s *TimeCodesScreen) handleKeyPress(msg tea.KeyPressMsg) (*TimeCodesScreen, tea.Cmd) {
	switch {
	case matchesKey(msg, keys.Up):
		if s.cursor > 0 {
			s.cursor--
		}
		return s, nil

	case matchesKey(msg, keys.Down):
		if s.cursor < len(s.timeCodes)-1 {
			s.cursor++
		}
		return s, nil

	case matchesKey(msg, keys.Add):
		return s, func() tea.Msg { return OpenTimeCodeAddMsg{} }

	case matchesKey(msg, keys.Edit):
		if tc := s.SelectedTimeCode(); tc != nil {
			return s, func() tea.Msg { return OpenTimeCodeEditMsg{TimeCode: tc} }
		}
		return s, nil

	case matchesKey(msg, keys.Delete):
		if tc := s.SelectedTimeCode(); tc != nil {
			return s, s.deleteTimeCode(tc.ID)
		}
		return s, nil

	case matchesKey(msg, keys.Refresh):
		return s, s.Refresh()
	}

	return s, nil
}

func (s *TimeCodesScreen) deleteTimeCode(id string) tea.Cmd {
	return func() tea.Msg {
		if err := s.client.DeleteTimeCode(id); err != nil {
			return timeCodesErrMsg{err}
		}
		return TimeCodeDeletedMsg{}
	}
}

func (s *TimeCodesScreen) View() string {
	if s.width == 0 {
		return "Loading..."
	}

	var b strings.Builder

	// Title bar
	title := titleStyle.Render("Time Codes")
	hint := helpStyle.Render("[Esc] Back")
	titleBar := lipgloss.JoinHorizontal(
		lipgloss.Top,
		title,
		strings.Repeat(" ", max(0, s.width-lipgloss.Width(title)-lipgloss.Width(hint)-2)),
		hint,
	)
	b.WriteString(titleBar + "\n\n")

	// Header
	b.WriteString(headerStyle.Render("Manage Time Codes") + "\n")
	b.WriteString(strings.Repeat("─", min(60, s.width-4)) + "\n")

	if s.loading {
		b.WriteString(statusStyle.Render("Loading...") + "\n")
	} else if s.err != nil {
		b.WriteString(errorStyle.Render(fmt.Sprintf("Error: %v", s.err)) + "\n")
	} else if len(s.timeCodes) == 0 {
		b.WriteString(statusStyle.Render("No time codes. Press [a] to add one.") + "\n")
	} else {
		for i, tc := range s.timeCodes {
			line := s.renderTimeCodeLine(i, tc)
			b.WriteString(line + "\n")
		}
	}

	// Help bar at bottom
	b.WriteString("\n")
	help := helpStyle.Render("[a] Add  [e] Edit  [d] Delete  [r] Refresh  [Esc] Back")
	b.WriteString(help)

	return b.String()
}

func (s *TimeCodesScreen) renderTimeCodeLine(index int, tc model.TimeCode) string {
	cursor := "  "
	style := itemStyle
	if index == s.cursor {
		cursor = "> "
		style = selectedItemStyle
	}

	// Show ID, Name, and active status
	status := ""
	if !tc.Active {
		status = inactiveStyle.Render(" [inactive]")
	}

	name := tc.Name
	if len(name) > 30 {
		name = name[:27] + "..."
	}

	line := fmt.Sprintf("%s%-10s  %-30s%s", cursor, tc.ID, name, status)
	return style.Render(line)
}
