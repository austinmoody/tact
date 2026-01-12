# Proposal: Enhance TUI Input Handling

## Summary

Fix control character issues in TUI text inputs and improve form usability.

## Motivation

Several usability issues exist with the current TUI text input handling:

1. **Ctrl+V paste inserts control characters** - When pasting text, escape sequences or control characters can appear in the input field
2. **Arrow keys insert control codes** - Left/right arrow keys in text fields sometimes insert visible control characters instead of moving the cursor
3. **Time code add form is incomplete** - Users must add a time code with just ID/Name, then edit it to access keywords, examples, and description fields
4. **Fixed-width forms** - Input fields use hardcoded widths (40-50 chars) regardless of terminal size, making the time codes screen feel cramped

## Proposed Changes

### 1. Filter Control Characters in Text Inputs

Add input sanitization to strip non-printable characters (except newlines where appropriate) from text input values. This handles both paste operations and any escape sequences that slip through key filtering.

### 2. Improve Arrow Key Handling

Enhance the key filtering in all text input modals to properly consume arrow key escape sequences before they reach the underlying textinput component.

### 3. Expand Time Code Add Form

Change the "Add Time Code" modal to include all fields (ID, Name, Description, Keywords, Examples) rather than just ID and Name. This provides a consistent experience with the edit modal and eliminates the need for a two-step add-then-edit workflow.

### 4. Responsive Form Widths

Make form input widths responsive to the terminal width. Pass the current window dimensions to modals and calculate appropriate input widths based on available space.

## Scope

- **In Scope**:
  - Input sanitization for control characters
  - Arrow key escape sequence filtering
  - Expanded add forms for time codes
  - Responsive input widths based on terminal size

- **Out of Scope**:
  - Multi-line text areas (future enhancement)
  - Custom clipboard handling
  - Form validation beyond current behavior

## Affected Components

- `tui/ui/entry_input.go` - New entry modal
- `tui/ui/timecode_edit.go` - Time code add/edit modal
- `tui/ui/worktype_edit.go` - Work type add/edit modal
- `tui/ui/styles.go` - Modal rendering utilities
- `tui/ui/app.go` - Pass dimensions to modals

## Dependencies

None - these are internal TUI improvements.
