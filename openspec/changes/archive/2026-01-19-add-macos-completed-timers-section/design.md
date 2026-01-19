## Context

The macOS Tact Timer app currently removes timers from the list immediately after they are successfully saved to the API. Users want to see completed timers from today so they can quickly start new timers with the same description without retyping.

### Stakeholders
- End users who frequently switch between tasks throughout the day

### Constraints
- Must not affect API integration (completed timers are already saved)
- Must not clutter the UI with old history (only show today's completions)
- Starting from a completed timer creates a NEW timer (not restarting the old one)

## Goals / Non-Goals

### Goals
- Show completed timers from today in the Timer List window
- Allow starting a new timer from a completed timer's description
- Automatically clean up completed timers from previous days

### Non-Goals
- Full timer history across multiple days
- Editing completed timer descriptions
- Re-submitting completed timers to the API

## Decisions

### Decision: Add `stopped` State and `stoppedAt` Timestamp

Extend `TimerState` enum to include `stopped` and add `stoppedAt: Date?` to `TactTimer`.

```swift
enum TimerState: String, Codable {
    case running
    case paused
    case stopped
}

struct TactTimer: Codable, Identifiable {
    // ... existing fields ...
    var stoppedAt: Date?  // When the timer was stopped (for cleanup)
}
```

**Rationale**: Using the existing model keeps persistence simple - all timers (active and completed) are stored in the same UserDefaults key.

### Decision: Separate Sections in Timer List

The Timer List window shows two sections:
1. **Active Timers** (top) - running and paused timers with Pause/Resume/Stop controls
2. **Completed Today** (bottom) - stopped timers with "Start New" button only

A visual separator (section header) distinguishes the two areas.

**Rationale**: Clear visual separation prevents confusion between active and completed timers. Users can quickly scan active timers at the top.

### Decision: "Start New" Creates Fresh Timer

Clicking "Start New" on a completed timer:
1. Copies the description
2. Calls `startNewTimer(description:)` to create a brand new timer
3. Does NOT modify or remove the completed timer entry

**Rationale**: The completed timer represents work already saved to the API. Starting new preserves that record while allowing the user to continue similar work.

### Decision: Auto-Cleanup on Launch

On app launch, remove any completed timers where `stoppedAt` is before the start of today (midnight local time).

**Rationale**: Keeps the completed section relevant without manual cleanup. Users only see today's work.

## Architecture

```
TimerListWindow
├── Active Timers Section
│   └── TimerCellView (Pause/Resume, Stop)
├── Separator + "Completed Today" Header
└── Completed Timers Section
    └── CompletedTimerCellView (description, duration, "Start New")
```

### Data Flow

1. **Stop Timer**: Instead of `removeTimer()`, set `state = .stopped` and `stoppedAt = Date()`
2. **Display**: Filter timers by state for each section
3. **Start New**: Call existing `startNewTimer(description:)` with copied description
4. **Cleanup**: On `load()`, filter out completed timers from previous days

## Risks / Trade-offs

### Risk: UserDefaults Size Growth
Completed timers accumulate in storage until cleanup.

**Mitigation**: Auto-cleanup on launch keeps only today's completions. Typical usage of ~10-20 timers/day is negligible.

### Trade-off: No Multi-Day History
Users cannot see yesterday's completed timers.

**Rationale**: Keeps scope minimal. The TUI and API provide historical data if needed. This feature is specifically for same-day quick-restart.
