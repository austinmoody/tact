## 1. Timer Model & Core Logic

- [x] 1.1 Create `tui/timer/` package directory
- [x] 1.2 Implement `Timer` struct and `TimerState` constants in `model.go`
- [x] 1.3 Add timer methods: `Pause()`, `Resume()`, `Stop()`, `TotalElapsedSeconds()`
- [x] 1.4 Implement duration formatting utilities in `formatter.go` (display format + API format)

## 2. Timer Manager & Persistence

- [x] 2.1 Implement `TimerManager` struct in `manager.go`
- [x] 2.2 Add timer lifecycle methods: `StartTimer`, `PauseTimer`, `ResumeTimer`, `StopTimer`, `DeleteTimer`
- [x] 2.3 Add query methods: `RunningTimer`, `ActiveTimers`, `CompletedToday`
- [x] 2.4 Implement JSON persistence to `~/.tact/timers.json` (load/save)
- [x] 2.5 Add auto-cleanup of completed timers from previous days on load
- [x] 2.6 Ensure single-running-timer enforcement (auto-pause on start/resume)

## 3. Timer Panel UI

- [x] 3.1 Create `TimerPanel` struct in `ui/timer_panel.go`
- [x] 3.2 Implement panel rendering with two sections: Active Timers and Completed Today
- [x] 3.3 Add cursor navigation between timers (up/down, vim keys)
- [x] 3.4 Implement input mode for new timer description
- [x] 3.5 Add empty state display when no timers exist
- [x] 3.6 Style timer items with state indicators (Running/Paused) and elapsed time

## 4. Timer Panel Actions

- [x] 4.1 Implement `[n]` - New timer (enter input mode, create on submit)
- [x] 4.2 Implement `[p]` - Pause selected running timer
- [x] 4.3 Implement `[r]` - Resume selected paused timer
- [x] 4.4 Implement `[s]` - Stop timer and create entry via API
- [x] 4.5 Implement `[d]` - Delete selected timer
- [x] 4.6 Add error handling and user feedback for API calls

## 5. App Integration

- [x] 5.1 Add `ModalTimerPanel` constant to modal enum in `app.go`
- [x] 5.2 Add `timerManager` and `timerPanel` fields to `App` struct
- [x] 5.3 Initialize `TimerManager` in `NewApp` (load persisted timers)
- [x] 5.4 Add `[t]` key binding to toggle timer panel from home screen
- [x] 5.5 Route timer panel messages in `App.Update()`
- [x] 5.6 Handle panel close (Esc or `[t]` toggle)

## 6. Real-time Updates

- [x] 6.1 Define `timerTickMsg` message type
- [x] 6.2 Implement tick command that fires every second
- [x] 6.3 Start tick when timer panel opens or a timer is running
- [x] 6.4 Stop tick when panel closes and no timer is running
- [x] 6.5 Update elapsed time display on each tick

## 7. Home Screen Status Indicator

- [x] 7.1 Add timer status rendering to home screen `View()`
- [x] 7.2 Display format: `⏱ Working on: {description} [{elapsed}]`
- [x] 7.3 Only show indicator when a timer is in running state
- [x] 7.4 Position indicator at bottom of home screen (above help text)

## 8. Key Bindings & Help Text

- [x] 8.1 Add `[t]` to home screen help text
- [x] 8.2 Add panel key bindings to `keys.go` if using centralized bindings
- [x] 8.3 Display help text in timer panel footer showing available actions
