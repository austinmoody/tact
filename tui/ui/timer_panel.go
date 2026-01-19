package ui

import (
	"strings"
	"time"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"

	"tact-tui/api"
	"tact-tui/timer"
)

// TimerPanel is a floating panel for timer management
type TimerPanel struct {
	manager   *timer.Manager
	client    *api.Client
	cursor    int
	inputMode bool
	input     textinput.Model
	width     int
	height    int
	err       error
	saving    bool
}

// NewTimerPanel creates a new timer panel
func NewTimerPanel(manager *timer.Manager, client *api.Client, width, height int) *TimerPanel {
	inputWidth := calculateInputWidth(width)

	ti := textinput.New()
	ti.Placeholder = "What are you working on?"
	ti.CharLimit = 200
	ti.SetWidth(inputWidth)

	return &TimerPanel{
		manager: manager,
		client:  client,
		input:   ti,
		width:   width,
		height:  height,
	}
}

// Init initializes the timer panel
func (p *TimerPanel) Init() tea.Cmd {
	return nil
}

// SetSize updates the panel dimensions
func (p *TimerPanel) SetSize(width, height int) {
	p.width = width
	p.height = height
	p.input.SetWidth(calculateInputWidth(width))
}

// Update handles messages for the timer panel
func (p *TimerPanel) Update(msg tea.Msg) (*TimerPanel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if p.inputMode {
			return p.updateInputMode(msg)
		}
		return p.updateNormalMode(msg)
	}

	if p.inputMode {
		var cmd tea.Cmd
		p.input, cmd = p.input.Update(msg)
		return p, cmd
	}

	return p, nil
}

func (p *TimerPanel) updateInputMode(msg tea.KeyPressMsg) (*TimerPanel, tea.Cmd) {
	switch msg.Key().Code {
	case tea.KeyEscape:
		p.inputMode = false
		p.input.Reset()
		return p, nil
	case tea.KeyEnter:
		if p.input.Value() != "" {
			p.manager.StartTimer(p.input.Value())
			p.inputMode = false
			p.input.Reset()
			p.cursor = 0
			return p, timerTick()
		}
		return p, nil
	}

	var cmd tea.Cmd
	p.input, cmd = p.input.Update(msg)
	return p, cmd
}

func (p *TimerPanel) updateNormalMode(msg tea.KeyPressMsg) (*TimerPanel, tea.Cmd) {
	allTimers := p.allDisplayTimers()

	switch msg.Key().Code {
	case tea.KeyEscape:
		return p, func() tea.Msg { return ModalCloseMsg{} }
	}

	switch msg.String() {
	case "t":
		// Toggle - close panel
		return p, func() tea.Msg { return ModalCloseMsg{} }

	case "n":
		// New timer
		p.inputMode = true
		p.input.Focus()
		return p, textinput.Blink

	case "j", "down":
		if len(allTimers) > 0 {
			p.cursor = (p.cursor + 1) % len(allTimers)
		}
		return p, nil

	case "k", "up":
		if len(allTimers) > 0 {
			p.cursor = (p.cursor - 1 + len(allTimers)) % len(allTimers)
		}
		return p, nil

	case "p":
		// Pause selected timer
		if t := p.selectedTimer(); t != nil && t.IsRunning() {
			p.manager.PauseTimer(t.ID)
		}
		return p, nil

	case "r":
		// Resume selected timer
		if t := p.selectedTimer(); t != nil && t.IsPaused() {
			p.manager.ResumeTimer(t.ID)
			return p, timerTick()
		}
		return p, nil

	case "s":
		// Stop selected timer and create entry
		if t := p.selectedTimer(); t != nil && !t.IsStopped() {
			return p, p.stopTimerAndCreateEntry(t.ID)
		}
		return p, nil

	case "d":
		// Delete selected timer
		if t := p.selectedTimer(); t != nil {
			p.manager.DeleteTimer(t.ID)
			allTimers = p.allDisplayTimers()
			if p.cursor >= len(allTimers) && p.cursor > 0 {
				p.cursor--
			}
		}
		return p, nil
	}

	return p, nil
}

func (p *TimerPanel) stopTimerAndCreateEntry(id string) tea.Cmd {
	return func() tea.Msg {
		t := p.manager.StopTimer(id)
		if t == nil {
			return timerEntryErrorMsg{err: nil}
		}

		// Create entry via API
		entry := timer.FormatEntry(t.TotalElapsedSeconds(), t.Description)
		_, err := p.client.CreateEntry(entry)
		if err != nil {
			return timerEntryErrorMsg{err: err}
		}
		return timerEntryCreatedMsg{}
	}
}

// timerEntryCreatedMsg is sent when a timer entry is successfully created
type timerEntryCreatedMsg struct{}

// timerEntryErrorMsg is sent when creating a timer entry fails
type timerEntryErrorMsg struct{ err error }

// allDisplayTimers returns active timers followed by completed today
func (p *TimerPanel) allDisplayTimers() []*timer.Timer {
	var all []*timer.Timer
	all = append(all, p.manager.ActiveTimers()...)
	all = append(all, p.manager.CompletedToday()...)
	return all
}

// selectedTimer returns the currently selected timer
func (p *TimerPanel) selectedTimer() *timer.Timer {
	all := p.allDisplayTimers()
	if p.cursor >= 0 && p.cursor < len(all) {
		return all[p.cursor]
	}
	return nil
}

// View renders the timer panel
func (p *TimerPanel) View() string {
	var b strings.Builder

	b.WriteString(modalTitleStyle.Render("Timers"))
	b.WriteString("\n\n")

	if p.inputMode {
		b.WriteString(labelStyle.Render("What are you working on?"))
		b.WriteString("\n")
		b.WriteString(focusedInputStyle.Render(p.input.View()))
		b.WriteString("\n\n")
		b.WriteString(helpStyle.Render("[Enter] Start  [Esc] Cancel"))
		return modalStyle.Render(b.String())
	}

	active := p.manager.ActiveTimers()
	completed := p.manager.CompletedToday()
	cursorIndex := 0

	// Active Timers section
	if len(active) > 0 {
		b.WriteString(headerStyle.Render("Active Timers"))
		b.WriteString("\n")
		for _, t := range active {
			prefix := "  "
			style := itemStyle
			if cursorIndex == p.cursor {
				prefix = "> "
				style = selectedItemStyle
			}

			stateStr := ""
			if t.IsRunning() {
				stateStr = statusParsedStyle.Render(" [Running]")
			} else if t.IsPaused() {
				stateStr = statusPendingStyle.Render(" [Paused]")
			}

			elapsed := timer.FormatElapsed(t.TotalElapsedSeconds())
			line := prefix + truncateString(t.Description, 30) + " " + elapsed + stateStr
			b.WriteString(style.Render(line))
			b.WriteString("\n")
			cursorIndex++
		}
		b.WriteString("\n")
	}

	// Completed Today section
	if len(completed) > 0 {
		b.WriteString(headerStyle.Render("Completed Today"))
		b.WriteString("\n")
		for _, t := range completed {
			prefix := "  "
			style := itemStyle
			if cursorIndex == p.cursor {
				prefix = "> "
				style = selectedItemStyle
			}

			duration := timer.FormatDuration(t.AccumulatedSeconds)
			line := prefix + truncateString(t.Description, 30) + " " + duration
			b.WriteString(style.Render(line))
			b.WriteString("\n")
			cursorIndex++
		}
		b.WriteString("\n")
	}

	// Empty state
	if len(active) == 0 && len(completed) == 0 {
		b.WriteString(statusStyle.Render("No timers yet. Press [n] to start one."))
		b.WriteString("\n\n")
	}

	// Error display
	if p.err != nil {
		b.WriteString(errorStyle.Render("Error: " + p.err.Error()))
		b.WriteString("\n\n")
	}

	// Help text
	help := "[n] New  [p] Pause  [r] Resume  [s] Stop  [d] Delete  [Esc] Close"
	b.WriteString(helpStyle.Render(help))

	return modalStyle.Render(b.String())
}

// truncateString truncates a string to max length with ellipsis
func truncateString(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

// timerTickMsg is sent every second to update elapsed time
type timerTickMsg struct{}

// timerTick returns a command that sends a tick message every second
func timerTick() tea.Cmd {
	return tea.Tick(1_000_000_000, func(_ time.Time) tea.Msg {
		return timerTickMsg{}
	})
}
