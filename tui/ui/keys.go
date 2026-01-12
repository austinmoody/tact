package ui

import "github.com/charmbracelet/bubbletea"

type keyMap struct {
	Up         []string
	Down       []string
	Left       []string
	Right      []string
	Enter      []string
	Escape     []string
	Tab        []string
	ShiftTab   []string
	Refresh    []string
	Quit       []string
	NewEntry   []string
	Menu       []string
	Add        []string
	Edit       []string
	Delete     []string
	Reparse    []string
}

var keys = keyMap{
	Up:       []string{"k", "up"},
	Down:     []string{"j", "down"},
	Left:     []string{"h", "left"},
	Right:    []string{"l", "right"},
	Enter:    []string{"enter"},
	Escape:   []string{"esc"},
	Tab:      []string{"tab"},
	ShiftTab: []string{"shift+tab"},
	Refresh:  []string{"r"},
	Quit:     []string{"q", "ctrl+c"},
	NewEntry: []string{"n"},
	Menu:     []string{"m"},
	Add:      []string{"a"},
	Edit:     []string{"e"},
	Delete:   []string{"d"},
	Reparse:  []string{"p"},
}

func matchesKey(msg tea.KeyMsg, bindings []string) bool {
	for _, binding := range bindings {
		if msg.String() == binding {
			return true
		}
	}
	return false
}
