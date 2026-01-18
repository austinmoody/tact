# Change: Add Completed Timers Section to macOS App

## Why

Users often need to resume work on tasks they completed earlier in the day. Currently, when a timer is stopped, it's removed from the list entirely after being saved to the API. This forces users to retype the description if they want to start tracking the same task again later.

## What Changes

- **MODIFIED**: Timer List window now shows a "Completed Today" section below active timers
- **MODIFIED**: Stopped timers are retained locally (with `stopped` state) instead of being removed
- **MODIFIED**: TactTimer model gains `stoppedAt` timestamp and `stopped` state
- **NEW**: "Start New" button on completed timers to create a fresh timer with the same description
- **NEW**: Auto-cleanup removes completed timers from previous days on app launch

## Impact

- Affected specs: `macos-timer-app`
- Affected code: `TimerModel.swift`, `TimerManager.swift`, `TimerListWindowController.swift`
- No backend changes required
