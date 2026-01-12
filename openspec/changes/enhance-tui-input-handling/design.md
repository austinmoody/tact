## Context

The TUI uses Bubble Tea with the `textinput` component from Bubbles. Control characters can appear in inputs through:
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

### Control Character Filtering

**Decision:** Add a sanitization function that filters the input value after each update, stripping non-printable ASCII characters (0x00-0x1F, 0x7F) except for standard whitespace.

**Rationale:**
- Filtering at the value level catches all sources (paste, key sequences, etc.)
- Simple and reliable approach
- Minimal performance impact on short input strings

**Implementation:**
```go
func sanitizeInput(s string) string {
    var result strings.Builder
    for _, r := range s {
        // Allow printable ASCII and extended unicode
        // Filter control chars except tab and common whitespace
        if r >= 32 || r == '\t' {
            result.WriteRune(r)
        }
    }
    return result.String()
}
```

### Arrow Key Handling

**Decision:** Explicitly handle `tea.KeyLeft` and `tea.KeyRight` in Update methods and only pass them to the focused input, consuming any raw escape sequences.

**Rationale:**
- Bubble Tea's key types should already handle arrow keys
- The issue occurs when terminals send raw escape sequences that aren't recognized
- Checking for `\x1b[` prefix catches unrecognized sequences

**Implementation:**
```go
case tea.KeyMsg:
    // Handle recognized keys
    switch msg.Type {
    case tea.KeyLeft, tea.KeyRight:
        // Pass to focused input for cursor movement
    }

    // Filter any unrecognized escape sequences
    if strings.HasPrefix(msg.String(), "\x1b") {
        return m, nil
    }
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
    // Modal padding: 2 chars each side, border: 1 each side, label margin: 2
    margins := 10
    available := termWidth - margins

    // Min 30, max 80
    return max(30, min(80, available))
}
```
