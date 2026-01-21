## 1. API Client

- [ ] 1.1 Add `EntryUpdate` struct to `tui/api/client.go` with `UserInput` and `EntryDate` pointer fields
- [ ] 1.2 Implement `UpdateEntry(id string, update EntryUpdate, learn bool)` method using PATCH request

## 2. Key Bindings

- [ ] 2.1 Add `Edit` key binding (e key) to `tui/ui/keys.go`

## 3. Entry Detail Modal - State

- [ ] 3.1 Add edit mode fields to `EntryDetailModal` struct: `editMode`, `userInputField`, `dateField`, `focusedField`, `saving`
- [ ] 3.2 Import `textinput` package from Bubble Tea
- [ ] 3.3 Initialize text input fields in `NewEntryDetailModal` (set values from entry)

## 4. Entry Detail Modal - Update Logic

- [ ] 4.1 Handle `e` key press to enter edit mode (populate fields, focus userInput)
- [ ] 4.2 Handle `Tab` key to switch focus between userInput and dateField in edit mode
- [ ] 4.3 Handle `Esc` key to cancel edit mode (return to view mode without saving)
- [ ] 4.4 Handle `Enter` key to trigger save in edit mode
- [ ] 4.5 Forward key messages to focused text input when in edit mode
- [ ] 4.6 Add date format validation (YYYY-MM-DD) before save
- [ ] 4.7 Implement save command that calls `UpdateEntry` API with learn=true
- [ ] 4.8 Handle save success: update entry, exit edit mode, return `EntryUpdatedMsg`
- [ ] 4.9 Handle save error: display error message, remain in edit mode

## 5. Entry Detail Modal - View

- [ ] 5.1 Render text input fields when in edit mode (replacing static text)
- [ ] 5.2 Show visual indicator for focused field
- [ ] 5.3 Update help text for edit mode: `[Tab] Switch field  [Enter] Save  [Esc] Cancel`
- [ ] 5.4 Show "Saving..." status while save is in progress
- [ ] 5.5 Display validation errors inline

## 6. Home Screen Integration

- [ ] 6.1 Add `EntryUpdatedMsg` type to messages (if not exists)
- [ ] 6.2 Handle `EntryUpdatedMsg` in home screen to refresh entry list
