## 1. Project Setup

- [x] 1.1 Create `tui/` directory with Go module (`go mod init tact-tui`)
- [x] 1.2 Add dependencies: bubbletea, lipgloss, bubbles
- [x] 1.3 Create basic `main.go` with flag parsing for `--api` and env var support

## 2. Data Models

- [x] 2.1 Create `model/timecode.go` with TimeCode struct matching API response
- [x] 2.2 Create `model/worktype.go` with WorkType struct matching API response

## 3. API Client

- [x] 3.1 Create `api/client.go` with configurable base URL
- [x] 3.2 Implement `FetchTimeCodes() ([]TimeCode, error)`
- [x] 3.3 Implement `FetchWorkTypes() ([]WorkType, error)`

## 4. UI Components

- [x] 4.1 Create `ui/styles.go` with Lip Gloss styles (borders, colors, focus states)
- [x] 4.2 Create `ui/keys.go` with key bindings (j/k, h/l, enter, r, q)
- [x] 4.3 Create `ui/dashboard.go` with main dashboard model:
  - Model struct (timeCodes, workTypes, cursor, activePane, loading, err)
  - Init() - fetch initial data
  - Update() - handle key messages and data messages
  - View() - render split-pane layout with detail panel

## 5. Integration

- [x] 5.1 Wire up dashboard in main.go
- [x] 5.2 Add loading state display
- [x] 5.3 Add error display in status bar

## 6. Build Integration

- [x] 6.1 Add Makefile targets: `tui-build`, `tui-run`, `tui-dev`
- [x] 6.2 Update root README with TUI section

## Verification

Steps to independently verify the implementation:

1. **Start the backend:**
   ```bash
   make run
   ```

2. **Build and run TUI:**
   ```bash
   make tui-build
   ./tui/tact-tui
   ```

3. **Test navigation:**
   - Press `j`/`k` to move up/down in the list
   - Press `h`/`l` to switch between panes
   - Press `Enter` to toggle detail view
   - Press `r` to refresh data
   - Press `q` to quit

4. **Test configuration:**
   ```bash
   # With flag
   ./tui/tact-tui --api http://localhost:2100

   # With env var
   TACT_API_URL=http://localhost:2100 ./tui/tact-tui
   ```

5. **Test error handling:**
   ```bash
   # Stop backend, then run TUI
   ./tui/tact-tui
   # Should show error message, not crash
   # Press 'r' to retry
   ```
