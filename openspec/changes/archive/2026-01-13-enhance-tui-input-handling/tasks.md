# Tasks

## 1. Upgrade to Bubble Tea v2

The original approach of sanitizing escape sequences was replaced with upgrading to Bubble Tea v2, which has proper input handling via the x/input package.

- [x] 1.1 Update go.mod dependencies to charm.land/bubbletea/v2, bubbles/v2, lipgloss/v2
- [x] 1.2 Update imports across all UI files
- [x] 1.3 Migrate `tea.KeyMsg` to `tea.KeyPressMsg` with `msg.Key().Code` API
- [x] 1.4 Migrate `View()` to return `tea.View` with `AltScreen` field
- [x] 1.5 Update `textinput.Width` to `textinput.SetWidth()`
- [x] 1.6 Fix Shift+Tab detection using `msg.Key().Mod&tea.ModShift`
- [x] 1.7 Update color definitions for lipgloss v2 (simple hex format)
- [x] 1.8 Route `tea.PasteMsg` to modals in app.go for paste support

## 2. Expand Time Code Add Form

- [x] 2.1 Change `TimeCodeEditModal` to set `inputCount = 5` in add mode
- [x] 2.2 Update `View()` to render all fields in add mode
- [x] 2.3 Update `updateFocus()` to cycle through all 5 fields in add mode
- [x] 2.4 Update `save()` to include all fields when creating time code
- [x] 2.5 Update API client `CreateTimeCode` to accept all fields

## 3. Implement Responsive Form Widths

- [x] 3.1 Add width parameter to modal constructors (EntryInputModal, TimeCodeEditModal, WorkTypeEditModal)
- [x] 3.2 Pass terminal width from App when creating modals
- [x] 3.3 Calculate input width based on available space (min 30, max 80)
- [x] 3.4 Apply calculated width to all textinput.Model instances

## 4. Testing

- [x] 4.1 Manual test: paste text into input fields
- [x] 4.2 Manual test: arrow keys in all text fields
- [x] 4.3 Manual test: add time code with all fields populated
- [x] 4.4 Manual test: colors display correctly

## Notes

The v2 migration resolved issues 1 and 2 from the proposal (control characters and arrow keys) at the framework level rather than through application-level sanitization. This is a cleaner solution that leverages proper terminal input parsing in the underlying x/input package.

Removed code:
- `sanitizeInput()` function (no longer needed)
- `looksLikeEscapeFragment()` function (no longer needed)
- Arrow key filtering workarounds (handled by v2)
