# Tasks

## 1. Add Input Sanitization Utility

- [ ] 1.1 Add `sanitizeInput()` function to `tui/ui/styles.go` or new `tui/ui/utils.go`
- [ ] 1.2 Function filters control characters (0x00-0x1F, 0x7F) except standard whitespace

## 2. Apply Sanitization to Entry Input Modal

- [ ] 2.1 Update `entry_input.go` to sanitize input value after Update
- [ ] 2.2 Test paste with control characters is filtered

## 3. Fix Arrow Key Handling in Text Inputs

- [ ] 3.1 Update `entry_input.go` to explicitly handle KeyLeft/KeyRight
- [ ] 3.2 Update `timecode_edit.go` to filter raw escape sequences
- [ ] 3.3 Update `worktype_edit.go` to filter raw escape sequences
- [ ] 3.4 Test left/right arrows work for cursor movement without inserting characters

## 4. Expand Time Code Add Form

- [ ] 4.1 Change `TimeCodeEditModal` to set `inputCount = 5` in add mode
- [ ] 4.2 Update `View()` to render all fields in add mode
- [ ] 4.3 Update `updateFocus()` to cycle through all 5 fields in add mode
- [ ] 4.4 Update `save()` to include all fields when creating time code
- [ ] 4.5 Update API client `CreateTimeCode` to accept all fields

## 5. Implement Responsive Form Widths

- [ ] 5.1 Add width parameter to modal constructors (EntryInputModal, TimeCodeEditModal, WorkTypeEditModal)
- [ ] 5.2 Pass terminal width from App when creating modals
- [ ] 5.3 Calculate input width based on available space (min 30, max 80)
- [ ] 5.4 Apply calculated width to all textinput.Model instances

## 6. Testing

- [ ] 6.1 Manual test: paste text with control characters
- [ ] 6.2 Manual test: arrow keys in all text fields
- [ ] 6.3 Manual test: add time code with all fields populated
- [ ] 6.4 Manual test: resize terminal and verify form adapts

## Dependencies

- Task 2 depends on Task 1
- Task 3 depends on Task 1
- Tasks 4, 5 are independent
- Task 6 depends on all others
