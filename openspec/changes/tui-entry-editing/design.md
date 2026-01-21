## Context

The TUI currently has an `EntryDetailModal` (`tui/ui/entry_detail.go`) that displays read-only entry information and allows reparsing. The backend already supports entry updates via `PATCH /entries/{entry_id}` with a `learn` query parameter that triggers context document creation from corrections.

Current entry detail modal:
- Shows user input, status, parsed fields, date, and parse notes
- Supports `[p]` to reparse and `[Esc]` to close
- Uses Bubble Tea framework patterns

The TUI API client (`tui/api/client.go`) has no `UpdateEntry` method - this must be added.

## Goals / Non-Goals

**Goals:**
- Allow editing `user_input` and `entry_date` fields from the entry detail modal
- Integrate with backend's learn flag to improve AI context from corrections
- Maintain consistent UX patterns with existing modals
- Refresh entry list after edits to show changes immediately

**Non-Goals:**
- Editing other fields (duration, time_code, work_type) - these are AI-parsed
- Bulk editing multiple entries
- Undo/redo functionality
- Offline edit queueing

## Decisions

### 1. Edit mode toggle vs separate modal

**Decision**: Add edit mode toggle to existing `EntryDetailModal`

**Rationale**:
- Keeps user in context (same modal, same entry)
- Less code than creating a new modal
- Consistent with "detail + actions" pattern already used

**Alternative considered**: Separate `EntryEditModal` - rejected due to added complexity and context switching

### 2. Text input component

**Decision**: Use Bubble Tea's `textinput` component from `charm.land/bubbletea/v2`

**Rationale**:
- Already available in the Bubble Tea ecosystem
- Handles cursor movement, text editing natively
- Used elsewhere in the TUI for consistency

### 3. Date input format

**Decision**: Use simple text input with YYYY-MM-DD format validation

**Rationale**:
- Matches the API's expected format
- Simple to implement
- User can easily type or edit dates

**Alternative considered**: Date picker widget - rejected as overly complex for this use case

### 4. Learn flag behavior

**Decision**: Always enable learn=true when saving edits

**Rationale**:
- User is explicitly correcting an entry, so the AI should learn
- Matches the backend's default behavior
- Simplifies UX (no extra checkbox)

### 5. Entry refresh strategy

**Decision**: Return `EntryUpdatedMsg` from modal, triggering home screen to refetch entries

**Rationale**:
- Ensures list is consistent with backend state
- Handles date changes (entry moves to different date group)
- Simple message-passing pattern already used for reparsing

## Risks / Trade-offs

**Risk**: User accidentally enters invalid date format
→ **Mitigation**: Validate format client-side before API call, show error in modal

**Risk**: API update fails silently
→ **Mitigation**: Display error message in modal, keep edit mode active so user can retry

**Risk**: Race condition if entry changes while editing
→ **Mitigation**: Accept this for MVP; backend handles concurrent updates via ORM

**Trade-off**: Single save button for both fields vs field-by-field saves
→ **Accepted**: Single save is simpler UX; user edits all needed fields then saves once

## Implementation Components

```
tui/
├── api/
│   └── client.go          # Add UpdateEntry method
└── ui/
    ├── entry_detail.go    # Add edit mode, text inputs, save/cancel
    ├── keys.go            # Add Edit key binding
    └── messages.go        # Add EntryUpdatedMsg (if not exists)
```

### API Client Changes

```go
type EntryUpdate struct {
    UserInput *string `json:"user_input,omitempty"`
    EntryDate *string `json:"entry_date,omitempty"`
}

func (c *Client) UpdateEntry(id string, update EntryUpdate, learn bool) (*model.Entry, error)
```

### Modal State Changes

```go
type EntryDetailModal struct {
    // existing fields...
    editMode       bool
    userInputField textinput.Model
    dateField      textinput.Model
    focusedField   int  // 0 = userInput, 1 = date
    saving         bool
}
```

### Key Bindings

- `e` - Enter edit mode
- `Tab` - Switch between fields in edit mode
- `Enter` - Save changes (in edit mode)
- `Esc` - Cancel edit mode (first press) or close modal (second press)
