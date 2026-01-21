# Proposal: TUI Entry Editing

## Problem Statement

Currently, users can view entry details in a read-only modal (press Enter on an entry) and trigger a reparse, but cannot edit the entry's core fields directly. When a user notices their entry has incorrect text or is assigned to the wrong date, they have no way to fix it from the TUI.

## Proposed Solution

Add editing capability to the entry detail modal, allowing users to:

1. **Edit User Input** - Modify the original text of the time entry
2. **Edit Entry Date** - Change which date the entry belongs to
3. **Save with Learn Flag** - When saving, optionally trigger the backend's "learn" feature that creates context documents from corrections

### User Flow

1. User selects an entry on home screen and presses Enter
2. Entry detail modal opens (existing behavior)
3. User presses `e` to enter edit mode
4. Modal transforms to show editable text input for user_input
5. User can tab to date field and modify it (format: YYYY-MM-DD)
6. User presses Enter to save (with learn=true) or Escape to cancel
7. Entry refreshes in the modal and on home screen

### API Integration

The backend already supports entry updates via `PATCH /entries/{entry_id}` with optional `?learn=true` query parameter:

```python
@router.patch("/{entry_id}", response_model=EntryResponse)
def update_entry(
    entry_id: str,
    data: EntryUpdate,
    learn: bool = Query(True, description="Create context document from correction"),
    session: Session = Depends(get_session),
)
```

The `learn` flag (default true) creates a context document from the correction, helping the AI learn from user corrections.

### Implementation Components

1. **API Client** (`tui/api/client.go`)
   - Add `UpdateEntry(id string, userInput *string, entryDate *string, learn bool)` method
   - Uses PATCH with `?learn=` query parameter

2. **Entry Detail Modal** (`tui/ui/entry_detail.go`)
   - Add edit mode state
   - Add text input fields for user_input and entry_date
   - Handle mode switching (view/edit)
   - Add save/cancel commands

3. **Key Bindings** (`tui/ui/keys.go`)
   - Add `Edit` key binding (e key)

## Alternatives Considered

1. **Separate edit modal** - More complexity, worse UX
2. **Inline editing on home screen** - Too cramped, limited space for multi-line input
3. **Delete and re-create** - Loses entry history and is error-prone

## Success Criteria

- User can edit user_input text in the entry detail modal
- User can change entry_date
- Changes are saved via API with learn flag enabled
- Entry list on home screen refreshes to show updated entry
- Entry appears under correct date header after date change
