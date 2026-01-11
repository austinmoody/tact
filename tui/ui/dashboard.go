package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"tact-tui/api"
	"tact-tui/model"
)

type pane int

const (
	timeCodesPane pane = iota
	workTypesPane
)

type Dashboard struct {
	client      *api.Client
	timeCodes   []model.TimeCode
	workTypes   []model.WorkType
	activePane  pane
	codeCursor  int
	typeCursor  int
	showDetail  bool
	loading     bool
	err         error
	width       int
	height      int
}

// Messages
type timeCodesMsg struct{ codes []model.TimeCode }
type workTypesMsg struct{ types []model.WorkType }
type errMsg struct{ err error }

func NewDashboard(client *api.Client) *Dashboard {
	return &Dashboard{
		client:     client,
		activePane: timeCodesPane,
		loading:    true,
	}
}

func (d *Dashboard) Init() tea.Cmd {
	return tea.Batch(d.fetchTimeCodes(), d.fetchWorkTypes())
}

func (d *Dashboard) fetchTimeCodes() tea.Cmd {
	return func() tea.Msg {
		codes, err := d.client.FetchTimeCodes()
		if err != nil {
			return errMsg{err}
		}
		return timeCodesMsg{codes}
	}
}

func (d *Dashboard) fetchWorkTypes() tea.Cmd {
	return func() tea.Msg {
		types, err := d.client.FetchWorkTypes()
		if err != nil {
			return errMsg{err}
		}
		return workTypesMsg{types}
	}
}

func (d *Dashboard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return d.handleKeyPress(msg)

	case tea.WindowSizeMsg:
		d.width = msg.Width
		d.height = msg.Height
		return d, nil

	case timeCodesMsg:
		d.timeCodes = msg.codes
		d.loading = d.workTypes == nil
		return d, nil

	case workTypesMsg:
		d.workTypes = msg.types
		d.loading = d.timeCodes == nil
		return d, nil

	case errMsg:
		d.err = msg.err
		d.loading = false
		return d, nil
	}

	return d, nil
}

func (d *Dashboard) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case matchesKey(msg, keys.Quit):
		return d, tea.Quit

	case matchesKey(msg, keys.Refresh):
		d.loading = true
		d.err = nil
		return d, tea.Batch(d.fetchTimeCodes(), d.fetchWorkTypes())

	case matchesKey(msg, keys.Up):
		d.moveCursor(-1)
		return d, nil

	case matchesKey(msg, keys.Down):
		d.moveCursor(1)
		return d, nil

	case matchesKey(msg, keys.Left):
		d.activePane = timeCodesPane
		return d, nil

	case matchesKey(msg, keys.Right):
		d.activePane = workTypesPane
		return d, nil

	case matchesKey(msg, keys.Enter):
		d.showDetail = !d.showDetail
		return d, nil
	}

	return d, nil
}

func (d *Dashboard) moveCursor(delta int) {
	if d.activePane == timeCodesPane {
		d.codeCursor += delta
		if d.codeCursor < 0 {
			d.codeCursor = 0
		}
		if d.codeCursor >= len(d.timeCodes) {
			d.codeCursor = max(0, len(d.timeCodes)-1)
		}
	} else {
		d.typeCursor += delta
		if d.typeCursor < 0 {
			d.typeCursor = 0
		}
		if d.typeCursor >= len(d.workTypes) {
			d.typeCursor = max(0, len(d.workTypes)-1)
		}
	}
}

func (d *Dashboard) View() string {
	if d.width == 0 {
		return "Loading..."
	}

	var b strings.Builder

	// Title bar
	title := titleStyle.Render("Tact Dashboard")
	quitHint := helpStyle.Render("[q] quit")
	titleBar := lipgloss.JoinHorizontal(
		lipgloss.Top,
		title,
		strings.Repeat(" ", max(0, d.width-lipgloss.Width(title)-lipgloss.Width(quitHint)-2)),
		quitHint,
	)
	b.WriteString(titleBar + "\n\n")

	// Calculate pane dimensions
	paneWidth := (d.width - 4) / 2
	paneHeight := d.height - 10 // Leave room for title, detail, and help

	// Render panes
	leftPane := d.renderTimeCodesPane(paneWidth, paneHeight)
	rightPane := d.renderWorkTypesPane(paneWidth, paneHeight)

	panes := lipgloss.JoinHorizontal(lipgloss.Top, leftPane, " ", rightPane)
	b.WriteString(panes + "\n")

	// Detail panel
	if d.showDetail {
		detail := d.renderDetailPanel(d.width - 2)
		b.WriteString(detail + "\n")
	}

	// Status bar
	if d.err != nil {
		b.WriteString(errorStyle.Render(fmt.Sprintf("Error: %v", d.err)) + "\n")
	} else if d.loading {
		b.WriteString(statusStyle.Render("Loading...") + "\n")
	}

	// Help
	help := helpStyle.Render("[j/k] navigate  [h/l] switch pane  [enter] details  [r] refresh  [q] quit")
	b.WriteString("\n" + help)

	return b.String()
}

func (d *Dashboard) renderTimeCodesPane(width, height int) string {
	style := paneStyle
	if d.activePane == timeCodesPane {
		style = focusedPaneStyle
	}

	var content strings.Builder
	content.WriteString(headerStyle.Render("Time Codes") + "\n")

	if len(d.timeCodes) == 0 {
		content.WriteString(statusStyle.Render("No time codes"))
	} else {
		for i, tc := range d.timeCodes {
			cursor := "  "
			s := itemStyle
			if i == d.codeCursor && d.activePane == timeCodesPane {
				cursor = "> "
				s = selectedItemStyle
			}
			line := fmt.Sprintf("%s%s %s", cursor, tc.ID, tc.Name)
			content.WriteString(s.Render(line) + "\n")
		}
	}

	return style.Width(width).Height(height).Render(content.String())
}

func (d *Dashboard) renderWorkTypesPane(width, height int) string {
	style := paneStyle
	if d.activePane == workTypesPane {
		style = focusedPaneStyle
	}

	var content strings.Builder
	content.WriteString(headerStyle.Render("Work Types") + "\n")

	if len(d.workTypes) == 0 {
		content.WriteString(statusStyle.Render("No work types"))
	} else {
		for i, wt := range d.workTypes {
			cursor := "  "
			s := itemStyle
			if i == d.typeCursor && d.activePane == workTypesPane {
				cursor = "> "
				s = selectedItemStyle
			}
			line := fmt.Sprintf("%s%s", cursor, wt.Name)
			content.WriteString(s.Render(line) + "\n")
		}
	}

	return style.Width(width).Height(height).Render(content.String())
}

func (d *Dashboard) renderDetailPanel(width int) string {
	var content string

	if d.activePane == timeCodesPane && len(d.timeCodes) > 0 {
		tc := d.timeCodes[d.codeCursor]
		content = fmt.Sprintf(
			"%s: %s\nDescription: %s\nKeywords: %s\nExamples: %s",
			tc.ID,
			tc.Name,
			tc.Description,
			strings.Join(tc.Keywords, ", "),
			strings.Join(tc.Examples, ", "),
		)
	} else if d.activePane == workTypesPane && len(d.workTypes) > 0 {
		wt := d.workTypes[d.typeCursor]
		content = fmt.Sprintf("%s: %s", wt.ID, wt.Name)
	} else {
		content = "No item selected"
	}

	return detailStyle.Width(width).Render(content)
}
