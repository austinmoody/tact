## Context

The TUI uses Bubble Tea with the `textinput` component from Bubbles. Control characters were appearing in inputs through:
1. Paste operations (Ctrl+V) that include escape sequences
2. Arrow key escape sequences not being consumed before reaching the input
3. Terminal-specific key encoding differences

### Constraints
- Must work with standard Bubble Tea message flow
- Cannot intercept OS-level clipboard - must filter at input level
- Should not break legitimate cursor movement within text fields

## Goals / Non-Goals

**Goals:**
- Eliminate control character insertion in all text input fields
- Provide full form fields when adding time codes (not just ID/Name)
- Make forms responsive to terminal width

**Non-Goals:**
- Custom clipboard implementation
- Multi-line text editing
- Rich text or formatting support

## Decisions

### Bubble Tea v2 Migration

**Decision:** Upgrade to Bubble Tea v2 instead of implementing application-level sanitization.

**Rationale:**
- Bubble Tea v2 uses the `x/input` package for proper terminal input parsing
- Handles escape sequences at the framework level
- Cleaner solution than filtering at the application layer
- Future-proof as the framework handles terminal differences

**Implementation:**
- Updated dependencies to `charm.land/bubbletea/v2`, `bubbles/v2`, `lipgloss/v2`
- Migrated `tea.KeyMsg` to `tea.KeyPressMsg` with `msg.Key().Code` API
- Changed `View()` to return `tea.View` struct with `AltScreen` field
- Updated key matching to use new v2 patterns
- Routed `tea.PasteMsg` to modals for paste support (v2 separates paste from key events)

**Key API Changes:**
```go
// Key handling (v1 -> v2)
case tea.KeyMsg:           ->  case tea.KeyPressMsg:
msg.Type == tea.KeyEsc     ->  msg.Key().Code == tea.KeyEscape
msg.Type == tea.KeyEnter   ->  msg.Key().Code == tea.KeyEnter

// Shift detection
msg.Key().Mod&tea.ModShift != 0

// View rendering
func (m *Model) View() string           ->  func (m *Model) View() tea.View
return someString                       ->  v := tea.NewView(someString)
                                            v.AltScreen = true
                                            return v

// Text input width
ti.Width = 50                           ->  ti.SetWidth(50)

// Paste handling (separate message type in v2)
case tea.PasteMsg:
    content := msg.Content
```

### Expanded Add Form

**Decision:** Change `TimeCodeEditModal` to show all 5 fields in add mode (ID, Name, Description, Keywords, Examples) instead of just 2.

**Rationale:**
- Consistent with edit mode
- Eliminates friction of add-then-edit workflow
- All fields are optional except ID and Name (existing behavior)

**Changes:**
- Set `inputCount = 5` in add mode
- Show all fields in View() for add mode
- Update focus cycling to include all fields

### Responsive Widths

**Decision:** Accept `width` parameter when creating modals, calculate input width as `min(maxWidth, availableWidth - margins)`.

**Rationale:**
- Simple calculation based on modal margins and padding
- Falls back to current fixed width on narrow terminals
- No complex resize handling needed - modals are recreated on open

**Implementation:**
```go
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
```
