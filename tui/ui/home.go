package ui

import (
	"fmt"
	"sort"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"tact-tui/api"
	"tact-tui/model"
)

const entriesLimit = 15

type Home struct {
	client  *api.Client
	entries []model.Entry
	cursor  int
	loading bool
	err     error
	width   int
	height  int
}

type entriesMsg struct{ entries []model.Entry }
type homeErrMsg struct{ err error }

func NewHome(client *api.Client) *Home {
	return &Home{
		client:  client,
		loading: true,
	}
}

func (h *Home) Init() tea.Cmd {
	return h.fetchEntries()
}

func (h *Home) Refresh() tea.Cmd {
	h.loading = true
	h.err = nil
	return h.fetchEntries()
}

func (h *Home) fetchEntries() tea.Cmd {
	return func() tea.Msg {
		entries, err := h.client.FetchEntries(entriesLimit)
		if err != nil {
			return homeErrMsg{err}
		}
		return entriesMsg{entries}
	}
}

func (h *Home) SetSize(width, height int) {
	h.width = width
	h.height = height
}

func (h *Home) SelectedEntry() *model.Entry {
	if len(h.entries) == 0 || h.cursor >= len(h.entries) {
		return nil
	}
	return &h.entries[h.cursor]
}

func (h *Home) Update(msg tea.Msg) (*Home, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		return h.handleKeyPress(msg)

	case entriesMsg:
		h.entries = msg.entries
		// Sort by created_at descending (newest first)
		sort.Slice(h.entries, func(i, j int) bool {
			return h.entries[i].CreatedAt.After(h.entries[j].CreatedAt.Time)
		})
		h.loading = false
		if h.cursor >= len(h.entries) {
			h.cursor = max(0, len(h.entries)-1)
		}
		return h, nil

	case homeErrMsg:
		h.err = msg.err
		h.loading = false
		return h, nil
	}

	return h, nil
}

func (h *Home) handleKeyPress(msg tea.KeyPressMsg) (*Home, tea.Cmd) {
	switch {
	case matchesKey(msg, keys.Up):
		if h.cursor > 0 {
			h.cursor--
		}
		return h, nil

	case matchesKey(msg, keys.Down):
		if h.cursor < len(h.entries)-1 {
			h.cursor++
		}
		return h, nil

	case matchesKey(msg, keys.Refresh):
		return h, h.Refresh()
	}

	return h, nil
}

func (h *Home) View() string {
	if h.width == 0 {
		return "Loading..."
	}

	var b strings.Builder

	// Title bar
	title := titleStyle.Render("Tact")
	hint := helpStyle.Render("[n] New Entry")
	titleBar := lipgloss.JoinHorizontal(
		lipgloss.Top,
		title,
		strings.Repeat(" ", max(0, h.width-lipgloss.Width(title)-lipgloss.Width(hint)-2)),
		hint,
	)
	b.WriteString(titleBar + "\n\n")

	// Entries section
	if h.loading {
		b.WriteString(statusStyle.Render("Loading...") + "\n")
	} else if h.err != nil {
		b.WriteString(errorStyle.Render(fmt.Sprintf("Error: %v", h.err)) + "\n")
	} else if len(h.entries) == 0 {
		b.WriteString(statusStyle.Render("No entries yet. Press [n] to create one.") + "\n")
	} else {
		var currentDate string
		for i, entry := range h.entries {
			entryDate := entry.EntryDate
			if len(entryDate) > 10 {
				entryDate = entryDate[:10]
			}

			// Add date header when date changes
			if entryDate != currentDate {
				if currentDate != "" {
					b.WriteString("\n") // Space between date groups
				}
				currentDate = entryDate
				dateHeader := h.formatDateHeader(entryDate)
				b.WriteString(headerStyle.Render(dateHeader) + "\n")
				b.WriteString(strings.Repeat("─", min(50, h.width-4)) + "\n")
			}

			line := h.renderEntryLine(i, entry)
			b.WriteString(line + "\n")
		}
	}

	// Help bar at bottom
	b.WriteString("\n")
	help := helpStyle.Render("[n] New  [Enter] Details  [m] Menu  [r] Refresh  [q] Quit")
	b.WriteString(help)

	return b.String()
}

func (h *Home) renderEntryLine(index int, entry model.Entry) string {
	cursor := "  "
	style := itemStyle
	if index == h.cursor {
		cursor = "> "
		style = selectedItemStyle
	}

	// Truncate raw text if too long
	rawText := entry.RawText
	maxLen := h.width - 20
	if maxLen < 20 {
		maxLen = 20
	}
	if len(rawText) > maxLen {
		rawText = rawText[:maxLen-3] + "..."
	}

	// Status with color
	status := h.renderStatus(entry.Status)

	line := fmt.Sprintf("%s%-*s  %s", cursor, maxLen, rawText, status)
	return style.Render(line)
}

func (h *Home) renderStatus(status string) string {
	switch status {
	case "parsed":
		return statusParsedStyle.Render("parsed ")
	case "pending":
		return statusPendingStyle.Render("pending")
	case "failed":
		return statusFailedStyle.Render("failed ")
	default:
		return statusStyle.Render(status)
	}
}

func (h *Home) formatDateHeader(dateStr string) string {
	// Parse the date string (YYYY-MM-DD format)
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return dateStr
	}

	// Format as "Monday - Jan 2, 2006"
	return t.Format("Monday - Jan 2, 2006")
}
