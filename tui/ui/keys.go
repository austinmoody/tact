package ui

import "github.com/charmbracelet/bubbletea"

type keyMap struct {
	Up      []string
	Down    []string
	Left    []string
	Right   []string
	Enter   []string
	Refresh []string
	Quit    []string
}

var keys = keyMap{
	Up:      []string{"k", "up"},
	Down:    []string{"j", "down"},
	Left:    []string{"h", "left"},
	Right:   []string{"l", "right"},
	Enter:   []string{"enter"},
	Refresh: []string{"r"},
	Quit:    []string{"q", "ctrl+c"},
}

func matchesKey(msg tea.KeyMsg, bindings []string) bool {
	for _, binding := range bindings {
		if msg.String() == binding {
			return true
		}
	}
	return false
}
