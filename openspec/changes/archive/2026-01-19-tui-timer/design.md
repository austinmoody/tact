## Context

The TUI is built with Bubble Tea v2 (Go) using a screen/modal architecture. The App struct manages screen state (Home, TimeCodes, WorkTypes, Projects) and modal overlays. Each screen and modal is a self-contained component that handles its own rendering and input.

The Mac app already has a fully working timer implementation using `TactTimer` struct with states (running/paused/stopped), accumulated seconds tracking, and persistence via UserDefaults. We'll mirror this design in Go.

**Current TUI architecture:**
- `App` routes messages to screens/modals
- Modals overlay screens and capture input when active
- Shared styles in `styles.go`
- API client for backend communication
- Message-driven async operations

## Goals / Non-Goals

**Goals:**
- Add timer functionality to TUI matching Mac app behavior
- Floating timer panel accessible via `[t]` from home screen
- Persistent status indicator on home screen for running timer
- Local file-based persistence independent of Mac app
- Multiple timers with single-running enforcement
- Create entries via existing API on timer stop

**Non-Goals:**
- Syncing timer state with Mac app
- Backend API changes for timer storage
- Timer functionality on screens other than home (can navigate back)
- Concurrent running timers

## Decisions

### 1. Timer Model Structure

Mirror the Mac app's `TactTimer` struct in Go:

```go
type TimerState string
const (
    TimerRunning TimerState = "running"
    TimerPaused  TimerState = "paused"
    TimerStopped TimerState = "stopped"
)

type Timer struct {
    ID                 string     `json:"id"`
    Description        string     `json:"description"`
    State              TimerState `json:"state"`
    StartedAt          *time.Time `json:"started_at,omitempty"`
    AccumulatedSeconds int        `json:"accumulated_seconds"`
    StoppedAt          *time.Time `json:"stopped_at,omitempty"`
}
```

**Rationale:** Proven design from Mac app. `AccumulatedSeconds` + `StartedAt` pattern handles pause/resume correctly without complex interval tracking.

### 2. Timer Manager

Create a `TimerManager` that handles lifecycle operations:

```go
type TimerManager struct {
    timers       []Timer
    persistPath  string
}

func (tm *TimerManager) StartTimer(description string) *Timer
func (tm *TimerManager) PauseTimer(id string)
func (tm *TimerManager) ResumeTimer(id string)
func (tm *TimerManager) StopTimer(id string) *Timer
func (tm *TimerManager) DeleteTimer(id string)
func (tm *TimerManager) RunningTimer() *Timer
func (tm *TimerManager) ActiveTimers() []Timer    // running + paused
func (tm *TimerManager) CompletedToday() []Timer  // stopped today
```

**Rationale:** Encapsulates timer logic and persistence, can be shared across UI components. Mirrors Mac app's `TimerManager` class.

### 3. UI Component: TimerPanel

New modal-style component for timer management:

```go
type TimerPanel struct {
    manager      *TimerManager
    client       *api.Client
    cursor       int
    inputMode    bool
    input        textinput.Model
    width, height int
    // ... tick tracking for elapsed time updates
}
```

**Key features:**
- Two sections: Active Timers (running/paused) and Completed Today
- Cursor navigation between timers
- Input mode for new timer description
- Real-time elapsed time display via tick messages

**Rationale:** Panel as overlay (like other modals) keeps home screen visible behind it. Consistent with existing modal patterns.

### 4. Integration Point: Home Screen

Add to `App` struct:
- `timerManager *TimerManager` - shared timer state
- `timerPanel *TimerPanel` - panel component
- `ModalTimerPanel` constant - new modal type

Add to home screen rendering:
- Status bar at bottom showing running timer (if any)
- Format: `⏱ Working on: {description} [{elapsed}]`

**Rationale:** Status indicator provides always-visible feedback without cluttering the main UI. Using existing modal system for panel.

### 5. Persistence Location

Use `~/.tact/timers.json`:

```go
func defaultPersistPath() string {
    home, _ := os.UserHomeDir()
    return filepath.Join(home, ".tact", "timers.json")
}
```

**Rationale:**
- `~/.tact/` aligns with potential future TUI config
- JSON format for human-readability and debugging
- Separate from Mac app's UserDefaults

**Alternatives considered:**
- XDG config dir (`~/.config/tact/`) - more standard on Linux but `.tact` simpler
- SQLite - overkill for simple timer state

### 6. Real-time Updates

Use Bubble Tea's tick mechanism for elapsed time updates:

```go
type timerTickMsg time.Time

func timerTick() tea.Cmd {
    return tea.Tick(time.Second, func(t time.Time) tea.Msg {
        return timerTickMsg(t)
    })
}
```

Only tick when:
1. Timer panel is open, OR
2. A timer is running (for home screen indicator)

**Rationale:** Bubble Tea's tick is the standard pattern for time-based updates. Conditional ticking avoids unnecessary CPU usage.

### 7. Key Bindings

**Home screen:**
- `[t]` - Toggle timer panel

**Timer panel:**
- `[n]` - New timer (enter input mode)
- `[p]` - Pause selected timer (if running)
- `[r]` - Resume selected timer (if paused)
- `[s]` - Stop selected timer (creates entry)
- `[d]` - Delete selected timer
- `[j/k]` or `[↑/↓]` - Navigate timer list
- `[Esc]` or `[t]` - Close panel

**Rationale:** Consistent with existing TUI key patterns. `[p]`/`[r]`/`[s]` are mnemonic for Pause/Resume/Stop.

### 8. Entry Creation on Stop

When stopping a timer:
1. Calculate total elapsed seconds
2. Format duration as `{hours}h{minutes}m` (e.g., "1h30m", "45m")
3. Create user_input as `"{duration} {description}"`
4. POST to `/entries` endpoint
5. Handle success/error with user feedback

**Rationale:** Matches Mac app behavior exactly. Reuses existing entry parsing pipeline.

## Risks / Trade-offs

**[Risk] Timer state lost if file write fails** → Write after every state change with error logging. Consider backup file on startup.

**[Risk] Clock skew during suspend/resume** → Use `time.Now()` at calculation time (not stored timestamps). Mac app handles this the same way.

**[Risk] Stale elapsed time display** → 1-second tick ensures max 1-second drift. Acceptable for time tracking.

**[Trade-off] No cross-device sync** → Intentional simplicity. Users who need sync can use Mac app or future backend integration.

**[Trade-off] Panel vs dedicated screen** → Panel keeps context visible but has limited space. Can add dedicated screen later if needed.

## File Structure

New files to create:
```
tui/
├── timer/
│   ├── model.go      # Timer struct, state constants
│   ├── manager.go    # TimerManager with persistence
│   └── formatter.go  # Duration formatting utilities
└── ui/
    └── timer_panel.go # TimerPanel component
```

Existing files to modify:
```
tui/ui/
├── app.go     # Add timerManager, timerPanel, ModalTimerPanel
├── home.go    # Add status indicator rendering
└── keys.go    # Add timer-related key bindings
```
