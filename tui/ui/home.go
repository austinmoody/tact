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
	"tact-tui/timer"
)

const entriesLimit = 15

type Home struct {
	client        *api.Client
	timerManager  *timer.Manager
	entries       []model.Entry
	timeCodes     []model.TimeCode
	timeCodeNames map[string]string // Lookup map: ID -> Name
	cursor        int
	loading       bool
	err           error
	width         int
	height        int
}

type entriesMsg struct{ entries []model.Entry }
type homeTimeCodesMsg struct{ timeCodes []model.TimeCode }
type homeErrMsg struct{ err error }

func NewHome(client *api.Client, timerManager *timer.Manager) *Home {
	return &Home{
		client:       client,
		timerManager: timerManager,
		loading:      true,
	}
}

func (h *Home) Init() tea.Cmd {
	return tea.Batch(h.fetchEntries(), h.fetchTimeCodes())
}

func (h *Home) Refresh() tea.Cmd {
	h.loading = true
	h.err = nil
	return tea.Batch(h.fetchEntries(), h.fetchTimeCodes())
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

func (h *Home) fetchTimeCodes() tea.Cmd {
	return func() tea.Msg {
		timeCodes, err := h.client.FetchTimeCodes()
		if err != nil {
			return homeErrMsg{err}
		}
		return homeTimeCodesMsg{timeCodes}
	}
}

func (h *Home) SetSize(width, height int) {
	h.width = width
	h.height = height
	h.clampCursor()
}

// calculateAvailableEntryLines computes how many lines are available for entries
// based on terminal height minus fixed UI elements.
func (h *Home) calculateAvailableEntryLines() int {
	// Fixed UI overhead:
	// - Title bar + blank line: 2 lines
	// - Help bar + blank line before: 2 lines
	// - Scroll indicator (reserved): 1 line
	// - Safety margin: 5 lines
	// - Timer status (if running): +2 lines
	fixedLines := 10
	if h.timerManager.RunningTimer() != nil {
		fixedLines += 2
	}

	available := h.height - fixedLines
	// Minimum of 1 entry line regardless of calculated space
	if available < 1 {
		available = 1
	}
	return available
}

func (h *Home) SelectedEntry() *model.Entry {
	if len(h.entries) == 0 || h.cursor >= len(h.entries) {
		return nil
	}
	return &h.entries[h.cursor]
}

func (h *Home) TimeCodeNames() map[string]string {
	return h.timeCodeNames
}

func (h *Home) Update(msg tea.Msg) (*Home, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		return h.handleKeyPress(msg)

	case entriesMsg:
		h.entries = msg.entries
		// Sort by entry_date descending, then created_at descending within each date
		sort.Slice(h.entries, func(i, j int) bool {
			if h.entries[i].EntryDate != h.entries[j].EntryDate {
				return h.entries[i].EntryDate > h.entries[j].EntryDate
			}
			return h.entries[i].CreatedAt.After(h.entries[j].CreatedAt.Time)
		})
		h.loading = false
		h.clampCursor()
		return h, nil

	case homeErrMsg:
		h.err = msg.err
		h.loading = false
		return h, nil

	case homeTimeCodesMsg:
		h.timeCodes = msg.timeCodes
		// Build lookup map from ID to Name
		h.timeCodeNames = make(map[string]string)
		for _, tc := range h.timeCodes {
			h.timeCodeNames[tc.ID] = tc.Name
		}
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
		maxVisible := h.calculateMaxVisibleEntries()
		maxCursor := min(len(h.entries)-1, maxVisible-1)
		if h.cursor < maxCursor {
			h.cursor++
		}
		return h, nil

	case matchesKey(msg, keys.Refresh):
		return h, h.Refresh()
	}

	return h, nil
}

func (h *Home) View() string {
	if h.width == 0 || h.height == 0 {
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
		maxVisible := h.calculateMaxVisibleEntries()
		var currentDate string
		entriesRendered := 0

		for i, entry := range h.entries {
			if entriesRendered >= maxVisible {
				break
			}

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
			entriesRendered++
		}

		// Show scroll indicator if entries are hidden
		hiddenCount := len(h.entries) - entriesRendered
		if hiddenCount > 0 {
			indicator := fmt.Sprintf("↓ %d more entries", hiddenCount)
			b.WriteString(helpStyle.Render(indicator) + "\n")
		}
	}

	// Timer status indicator (if a timer is running)
	if running := h.timerManager.RunningTimer(); running != nil {
		b.WriteString("\n")
		elapsed := timer.FormatElapsed(running.TotalElapsedSeconds())
		desc := running.Description
		if len(desc) > 40 {
			desc = desc[:37] + "..."
		}
		timerStatus := fmt.Sprintf("⏱ Working on: %s [%s]", desc, elapsed)
		b.WriteString(statusParsedStyle.Render(timerStatus))
		b.WriteString("\n")
	}

	// Help bar at bottom
	b.WriteString("\n")
	help := helpStyle.Render("[n] New  [t] Timer  [Enter] Details  [m] Menu  [r] Refresh  [q] Quit")
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

	// Time code display (ID + Name, truncated)
	const timeCodeColWidth = 25
	timeCodeDisplay := h.getTimeCodeDisplay(entry.TimeCodeID)

	// Truncate user input to fit: width - cursor(2) - timeCodeCol - gap(2) - status(7) - gap(2)
	userInput := entry.UserInput
	maxLen := h.width - 2 - timeCodeColWidth - 2 - 7 - 2
	if maxLen < 20 {
		maxLen = 20
	}
	if len(userInput) > maxLen {
		userInput = userInput[:maxLen-3] + "..."
	}

	// Status with color
	status := h.renderStatus(entry.Status)

	line := fmt.Sprintf("%s%-*s  %-*s  %s", cursor, maxLen, userInput, timeCodeColWidth, timeCodeDisplay, status)
	return style.Render(line)
}

func (h *Home) getTimeCodeDisplay(timeCodeID *string) string {
	if timeCodeID == nil || *timeCodeID == "" {
		return ""
	}
	id := *timeCodeID
	name, ok := h.timeCodeNames[id]
	if !ok || name == "" {
		return id
	}

	// Combine ID and name, truncate if needed
	display := fmt.Sprintf("%s %s", id, name)
	const maxDisplay = 25
	if len(display) > maxDisplay {
		// Keep ID visible, truncate name
		availableForName := maxDisplay - len(id) - 4 // space + "..."
		if availableForName > 0 && len(name) > availableForName {
			display = fmt.Sprintf("%s %s...", id, name[:availableForName])
		} else {
			display = display[:maxDisplay-3] + "..."
		}
	}
	return display
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

// clampCursor ensures the cursor stays within the visible range of entries.
func (h *Home) clampCursor() {
	if len(h.entries) == 0 {
		h.cursor = 0
		return
	}
	maxVisible := h.calculateMaxVisibleEntries()
	if h.cursor >= maxVisible {
		h.cursor = maxVisible - 1
	}
	if h.cursor < 0 {
		h.cursor = 0
	}
}

// calculateMaxVisibleEntries determines how many entries can be displayed
// given the current terminal height and accounting for date headers.
func (h *Home) calculateMaxVisibleEntries() int {
	if len(h.entries) == 0 {
		return 0
	}

	// If height isn't set yet, show limited entries
	if h.height == 0 {
		return min(5, len(h.entries))
	}

	availableLines := h.calculateAvailableEntryLines()
	linesUsed := 0
	visibleCount := 0
	var currentDate string

	for _, entry := range h.entries {
		entryDate := entry.EntryDate
		if len(entryDate) > 10 {
			entryDate = entryDate[:10]
		}

		linesNeeded := 1 // Entry line

		// Date header costs
		if entryDate != currentDate {
			if currentDate != "" {
				linesNeeded++ // Space between date groups
			}
			linesNeeded += 2 // Date header + separator
			currentDate = entryDate
		}

		if linesUsed+linesNeeded > availableLines {
			break
		}

		linesUsed += linesNeeded
		visibleCount++
	}

	// Always show at least 1 entry if entries exist
	if visibleCount == 0 && len(h.entries) > 0 {
		visibleCount = 1
	}

	return visibleCount
}
