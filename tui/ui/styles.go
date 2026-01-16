package ui

import (
	"strings"

	"charm.land/lipgloss/v2"
)

var (
	// Colors (dark theme)
	subtle    = lipgloss.Color("#383838")
	highlight = lipgloss.Color("#7D56F4")
	special   = lipgloss.Color("#73F59F")

	// Title style
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(highlight).
			Padding(0, 1)

	// Pane styles
	paneStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(subtle).
			Padding(0, 1)

	focusedPaneStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(highlight).
				Padding(0, 1)

	// List item styles
	itemStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	selectedItemStyle = lipgloss.NewStyle().
				Foreground(special).
				Bold(true).
				PaddingLeft(0)

	// Detail panel style
	detailStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(subtle).
			Padding(0, 1)

	// Status bar style
	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#777777"))

	// Entry status styles
	statusParsedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#73F59F"))

	statusPendingStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFD93D"))

	statusFailedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FF6B6B"))

	// Error style
	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B6B")).
			Bold(true)

	// Help style
	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262"))

	// Header style for panes
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(highlight).
			MarginBottom(1)

	// Modal styles
	modalStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(highlight).
			Padding(1, 2)

	modalTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(highlight).
			MarginBottom(1)

	// Input field style
	inputStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(subtle).
			Padding(0, 1)

	focusedInputStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(highlight).
				Padding(0, 1)

	// Label style
	labelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#AAAAAA"))

	// Active/Inactive status styles
	activeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#73F59F"))

	inactiveStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#777777"))

	// Disabled input style (for read-only fields)
	disabledInputStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("#444444")).
				Foreground(lipgloss.Color("#666666")).
				Padding(0, 1)
)

// calculateInputWidth determines input field width based on terminal width.
// Returns a width between 30 and 80 characters.
func calculateInputWidth(termWidth int) int {
	// Modal padding: 2 chars each side, border: 1 each side, internal padding: ~4
	margins := 12
	available := termWidth - margins

	// Clamp between 30 and 80
	if available < 30 {
		return 30
	}
	if available > 80 {
		return 80
	}
	return available
}

// wrapText wraps text to fit within the specified width.
// It breaks on word boundaries when possible.
func wrapText(text string, width int) string {
	if width <= 0 {
		return text
	}

	var result strings.Builder
	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}

	lineLen := 0
	for i, word := range words {
		wordLen := len(word)

		if i == 0 {
			result.WriteString(word)
			lineLen = wordLen
		} else if lineLen+1+wordLen <= width {
			result.WriteString(" ")
			result.WriteString(word)
			lineLen += 1 + wordLen
		} else {
			result.WriteString("\n")
			result.WriteString(word)
			lineLen = wordLen
		}
	}

	return result.String()
}

// Helper function to render a modal overlay
func renderModalOverlay(background, modal string, width, height int) string {
	// Simple overlay - just center the modal
	modalWidth := lipgloss.Width(modal)
	modalHeight := lipgloss.Height(modal)

	// Calculate padding
	paddingLeft := max(0, (width-modalWidth)/2)
	paddingTop := max(0, (height-modalHeight)/3)

	// Build the overlay
	var b strings.Builder
	for i := 0; i < paddingTop; i++ {
		b.WriteString("\n")
	}
	lines := strings.Split(modal, "\n")
	for _, line := range lines {
		b.WriteString(strings.Repeat(" ", paddingLeft))
		b.WriteString(line)
		b.WriteString("\n")
	}

	return b.String()
}
