## Context

Tact Timer is a native macOS app for tracking time. It integrates with the existing Tact backend API to save time entries. The app is dock-centric: users interact primarily via the dock icon's right-click menu, with minimal windows for input and timer management.

### Stakeholders
- End user wanting quick time tracking without terminal/browser
- Existing Tact backend (API consumer)

### Constraints
- macOS only (no cross-platform requirement)
- Must work with existing Tact API (no backend changes)
- No JavaScript/TypeScript (project constraint)

## Goals / Non-Goals

### Goals
- Quick timer start/stop from Dock
- Multiple concurrent timers (one active, others paused)
- Persist timers across app restarts
- Save completed timers to Tact backend API
- Simple, focused UX - minimal windows and clicks

### Non-Goals
- Editing time codes, projects, or work types (use TUI for that)
- Viewing historical entries (use TUI for that)
- Menu bar presence (explicitly Dock-only per requirements)
- iOS/iPad support
- Offline queue for API failures (keep it simple for v1)

## Decisions

### Decision: Dock Menu as Primary Interface
The right-click dock menu is the main interaction point. It shows:
- Current timer status (count of timers, which is running)
- Quick actions (Start New, View All, Preferences)

**Alternatives considered:**
- Menu bar app: Rejected per requirements
- Main window always visible: Too heavy for quick tracking

### Decision: Timer State Machine
Each timer has a state: `running`, `paused`, or `stopped`.
- Only one timer can be `running` at a time
- Starting a new timer automatically pauses the current running one
- Stopping a timer saves it to the API and removes it from the list

**State transitions:**
```
[New] --start--> [Running] --pause--> [Paused]
                     |                    |
                     +----<--resume-------+
                     |
                 --stop--> [API Call] --> [Removed]
```

### Decision: Persistence via UserDefaults + JSON
Timers are serialized to JSON and stored in UserDefaults. On launch, timers are restored. Elapsed time is calculated based on stored start timestamps.

**Alternatives considered:**
- SQLite: Overkill for a few timers
- Core Data: Too complex for this use case
- File-based JSON: UserDefaults is simpler and standard for app state

### Decision: Time Entry Format
When stopping a timer, the app formats the entry as:
```
{duration} {description}
```
Example: `"45m standup meeting"` or `"1h30m Project Alpha work"`

Duration formatting:
- < 60 minutes: `Xm` (e.g., "45m")
- >= 60 minutes: `Xh` or `XhYm` (e.g., "1h", "1h30m")

This leverages the existing LLM parsing in the backend.

### Decision: Swift + AppKit
Pure AppKit for dock integration and windows. No SwiftUI to keep dependencies minimal and dock behavior predictable.

**Alternatives considered:**
- SwiftUI: Would still need AppKit for dock; adds complexity
- Objective-C: Swift is more maintainable

## Architecture

```
TactTimer.app/
├── AppDelegate.swift       # App lifecycle, dock menu setup
├── TimerManager.swift      # Timer state, persistence, business logic
├── Timer.swift             # Timer model (id, description, state, elapsed)
├── DockMenuController.swift # Builds and updates dock menu
├── TimerListWindow.swift   # Window showing all timers
├── NewTimerWindow.swift    # Popup for entering description
├── PreferencesWindow.swift # API URL configuration
├── APIClient.swift         # HTTP client for Tact API
└── TimeFormatter.swift     # Duration formatting utilities
```

### Data Flow

1. **Start Timer**: User right-clicks dock → "Start New Timer" → popup appears → user enters description → TimerManager creates timer in `running` state
2. **Pause/Resume**: Dock menu or timer list → TimerManager updates state
3. **Stop Timer**: Dock menu or timer list → TimerManager formats entry → APIClient POSTs to `/entries` → timer removed from list
4. **Persistence**: On any state change, TimerManager serializes all timers to UserDefaults

## Risks / Trade-offs

### Risk: API Unavailable
If the backend is down when stopping a timer, the entry would be lost.

**Mitigation (v1)**: Show error alert, keep timer in list so user can retry. Don't auto-queue (keep it simple).

**Future**: Could add offline queue if this becomes a pain point.

### Risk: App Crash Loses Running Time
If app crashes, elapsed time since last persist could be lost.

**Mitigation**: Persist timer state on every state change. Store `startedAt` timestamp, not just elapsed seconds, so time can be reconstructed on restore.

### Trade-off: No Historical View
Users can't see past entries in this app.

**Rationale**: Keep scope minimal. TUI already provides this. Adding it here would bloat the app.

## Open Questions

None - all clarified during proposal discussion.
