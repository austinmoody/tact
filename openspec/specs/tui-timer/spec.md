# tui-timer Specification

## Purpose
TBD - created by archiving change tui-timer. Update Purpose after archive.
## Requirements
### Requirement: Timer panel toggle
The system SHALL display a floating timer panel overlay when the user presses `[t]` from the home screen, and close it when pressed again or `[Esc]` is pressed.

#### Scenario: Open timer panel
- **WHEN** user presses `[t]` on the home screen
- **THEN** a floating panel overlay appears showing timer controls and timer list

#### Scenario: Close timer panel with toggle
- **WHEN** user presses `[t]` while the timer panel is open
- **THEN** the panel closes and returns to the home screen

#### Scenario: Close timer panel with escape
- **WHEN** user presses `[Esc]` while the timer panel is open
- **THEN** the panel closes and returns to the home screen

### Requirement: Start new timer
The system SHALL allow users to start a new timer by entering a description.

#### Scenario: Start timer with description
- **WHEN** user presses `[n]` in the timer panel and enters a description
- **THEN** a new timer is created with state "running" and the description is stored

#### Scenario: Start timer pauses existing running timer
- **WHEN** user starts a new timer while another timer is running
- **THEN** the existing running timer is paused and the new timer starts running

### Requirement: Only one timer running at a time
The system SHALL enforce that at most one timer can be in the "running" state at any time.

#### Scenario: Single running timer enforcement
- **WHEN** a timer is started or resumed
- **THEN** any other running timer is automatically paused

### Requirement: Pause timer
The system SHALL allow users to pause a running timer, preserving accumulated time.

#### Scenario: Pause running timer
- **WHEN** user selects a running timer and presses `[p]` or the pause action
- **THEN** the timer state changes to "paused" and accumulated seconds are preserved

### Requirement: Resume timer
The system SHALL allow users to resume a paused timer.

#### Scenario: Resume paused timer
- **WHEN** user selects a paused timer and presses `[r]` or the resume action
- **THEN** the timer state changes to "running" and continues from accumulated time

#### Scenario: Resume pauses other running timer
- **WHEN** user resumes a paused timer while another timer is running
- **THEN** the other running timer is paused and the selected timer resumes

### Requirement: Stop timer and create entry
The system SHALL stop a timer and create a time entry via the API when the user stops it.

#### Scenario: Stop timer creates entry
- **WHEN** user selects a timer and presses `[s]` or the stop action
- **THEN** the timer state changes to "stopped", the stop time is recorded, and an entry is created via `POST /entries` with format `"{duration} {description}"`

#### Scenario: Stopped timer appears in completed list
- **WHEN** a timer is stopped
- **THEN** it moves from the active timers list to the "Completed Today" section

### Requirement: Delete timer
The system SHALL allow users to delete a timer without creating an entry.

#### Scenario: Delete active timer
- **WHEN** user selects an active timer and presses `[d]` or the delete action
- **THEN** the timer is removed from the list without creating an entry

#### Scenario: Delete completed timer
- **WHEN** user selects a completed timer and presses `[d]` or the delete action
- **THEN** the timer is removed from the completed list

### Requirement: Display elapsed time
The system SHALL display the current elapsed time for running timers, updating in real-time.

#### Scenario: Running timer shows live elapsed time
- **WHEN** a timer is in "running" state
- **THEN** the elapsed time display updates every second showing format "HH:MM:SS" or "MM:SS"

#### Scenario: Paused timer shows frozen elapsed time
- **WHEN** a timer is in "paused" state
- **THEN** the elapsed time display shows the accumulated time without updating

### Requirement: Timer panel layout
The system SHALL display timers in two sections: active timers (running/paused) and completed today.

#### Scenario: Panel shows active timers section
- **WHEN** the timer panel is open and there are running or paused timers
- **THEN** they appear in an "Active Timers" section with state indicator (Running/Paused)

#### Scenario: Panel shows completed today section
- **WHEN** the timer panel is open and there are stopped timers from today
- **THEN** they appear in a "Completed Today" section with final duration

#### Scenario: Empty state display
- **WHEN** the timer panel is open and there are no timers
- **THEN** a helpful message is shown indicating how to start a timer

### Requirement: Home screen status indicator
The system SHALL display an active timer indicator on the home screen when a timer is running.

#### Scenario: Running timer shows indicator
- **WHEN** a timer is in "running" state and the user is on the home screen (panel closed)
- **THEN** a status indicator appears showing the timer description and elapsed time

#### Scenario: No indicator when no running timer
- **WHEN** no timer is running (all paused or stopped)
- **THEN** no timer indicator is displayed on the home screen

### Requirement: Local file persistence
The system SHALL persist timer state to a local JSON file at `~/.tact/timers.json`.

#### Scenario: Timers persist across sessions
- **WHEN** the TUI is closed and reopened
- **THEN** all active timers (running/paused) are restored with their accumulated time

#### Scenario: Auto-cleanup of old completed timers
- **WHEN** the TUI starts
- **THEN** completed timers from previous days are removed from the persisted file

#### Scenario: Timer state saved on change
- **WHEN** a timer is started, paused, resumed, stopped, or deleted
- **THEN** the updated timer state is immediately saved to the persistence file

### Requirement: Timer duration formatting
The system SHALL format timer durations consistently for display and API submission.

#### Scenario: Display format for elapsed time
- **WHEN** showing elapsed time in the panel
- **THEN** format is "H:MM:SS" for times over an hour, "MM:SS" for shorter times

#### Scenario: API format for entry creation
- **WHEN** creating an entry from a stopped timer
- **THEN** duration is formatted as "{N}h{M}m" (e.g., "1h30m", "45m") rounded to nearest minute

