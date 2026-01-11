## 1. Entry Model and API Client

- [ ] 1.1 Create `tui/model/entry.go` with Entry struct matching API response
- [ ] 1.2 Add `FetchEntries(limit int)` to API client
- [ ] 1.3 Add `CreateEntry(rawText string)` to API client
- [ ] 1.4 Add `ReparseEntry(id string)` to API client

## 2. App Router

- [ ] 2.1 Create `tui/ui/app.go` as root model managing screen state
- [ ] 2.2 Implement screen routing (home, timecodes, worktypes)
- [ ] 2.3 Implement modal overlay system
- [ ] 2.4 Update `main.go` to use new App model

## 3. Home Screen

- [ ] 3.1 Create `tui/ui/home.go` with entries list view
- [ ] 3.2 Display last 5 entries with status and date
- [ ] 3.3 Implement cursor navigation (j/k)
- [ ] 3.4 Add keybindings for new entry (n), menu (m), refresh (r)
- [ ] 3.5 Handle Enter to open entry detail modal

## 4. Entry Input Modal

- [ ] 4.1 Create `tui/ui/modal/entry_input.go`
- [ ] 4.2 Implement text input field with Bubble Tea textinput
- [ ] 4.3 Handle Enter to submit and create entry via API
- [ ] 4.4 Handle Esc to cancel
- [ ] 4.5 Show loading state during API call

## 5. Entry Detail Modal

- [ ] 5.1 Create `tui/ui/modal/entry_detail.go`
- [ ] 5.2 Display all entry fields (raw_text, status, duration, time_code, etc.)
- [ ] 5.3 Implement reparse action (p key)
- [ ] 5.4 Handle Esc to close
- [ ] 5.5 Show loading/success feedback for reparse

## 6. Menu Modal

- [ ] 6.1 Create `tui/ui/modal/menu.go`
- [ ] 6.2 Display menu options (Time Codes, Work Types)
- [ ] 6.3 Implement cursor navigation
- [ ] 6.4 Handle Enter to navigate to selected screen
- [ ] 6.5 Handle Esc to close

## 7. Time Codes Management

- [ ] 7.1 Add `CreateTimeCode`, `UpdateTimeCode`, `DeleteTimeCode` to API client
- [ ] 7.2 Create `tui/ui/timecodes.go` management screen
- [ ] 7.3 Display time codes list with active/inactive status
- [ ] 7.4 Implement cursor navigation and selection
- [ ] 7.5 Create `tui/ui/modal/timecode_edit.go` for edit modal
- [ ] 7.6 Implement quick-add modal (ID + Name only)
- [ ] 7.7 Implement edit modal (Name, Description, Keywords)
- [ ] 7.8 Implement deactivate action (d key)
- [ ] 7.9 Handle Esc to return to home

## 8. Work Types Management

- [ ] 8.1 Add `CreateWorkType`, `UpdateWorkType`, `DeleteWorkType` to API client
- [ ] 8.2 Create `tui/ui/worktypes.go` management screen
- [ ] 8.3 Display work types list with active/inactive status
- [ ] 8.4 Implement cursor navigation and selection
- [ ] 8.5 Create `tui/ui/modal/worktype_edit.go` for edit modal
- [ ] 8.6 Implement quick-add modal (Name only, ID auto-generated)
- [ ] 8.7 Implement edit modal (Name)
- [ ] 8.8 Implement deactivate action (d key)
- [ ] 8.9 Handle Esc to return to home

## 9. Styling and Polish

- [ ] 9.1 Update `tui/ui/styles.go` with modal styles
- [ ] 9.2 Add status color coding (parsed=green, pending=yellow, failed=red)
- [ ] 9.3 Add loading spinners/indicators
- [ ] 9.4 Ensure consistent styling across all screens

## 10. Cleanup

- [ ] 10.1 Remove old `dashboard.go` (replaced by home.go)
- [ ] 10.2 Update keybindings documentation in keys.go
- [ ] 10.3 Test all navigation flows
- [ ] 10.4 Verify `make build` succeeds for TUI

## Verification

1. **Start the backend:**
   ```bash
   make docker-up
   docker compose exec ollama ollama pull llama3.2:3b
   ```

2. **Run the TUI:**
   ```bash
   cd tui && go run .
   ```

3. **Test entry workflow:**
   - Press `n` to open new entry modal
   - Type "2h working on alpha" and press Enter
   - Entry should appear in list with status "pending"
   - Wait for parsing, then press `r` to refresh
   - Entry should show "parsed" status
   - Press Enter on entry to view details
   - Press `p` to reparse, verify it resets to pending

4. **Test menu navigation:**
   - Press `m` to open menu
   - Select "Time Codes" and press Enter
   - Verify time codes list displays
   - Press `a` to add new time code
   - Press `Esc` to return home

5. **Test time code management:**
   - Navigate to Time Codes via menu
   - Press `a`, enter ID "TEST-01" and name "Test Code"
   - Verify code appears in list
   - Select it and press `e` to edit
   - Change name and save
   - Press `d` to deactivate

6. **Test work type management:**
   - Navigate to Work Types via menu
   - Press `a`, enter name "Testing"
   - Verify work type appears with auto-generated ID "testing"
   - Select and press `e` to edit name
   - Press `d` to deactivate
