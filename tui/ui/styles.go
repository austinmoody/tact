package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	// Colors
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}

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
			Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"})

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
			Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#626262"})

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
			Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#AAAAAA"})

	// Active/Inactive status styles
	activeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#73F59F"))

	inactiveStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#777777"))
)

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
