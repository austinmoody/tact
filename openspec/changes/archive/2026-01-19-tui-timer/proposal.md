## Why

The TUI currently lacks timer functionality, requiring users to switch to the Mac app to track time while working. Adding a timer to the TUI enables a complete terminal-based workflow for time tracking without context switching.

## What Changes

- Add a floating timer panel overlay accessible via `[t]` key from the home screen
- Display active timer status indicator on the home screen (visible when panel is closed)
- Support multiple timers with only one running at a time (pause existing when starting new)
- Timer states: running, paused, stopped
- Local file-based persistence (`~/.tact/timers.json`) - independent from Mac app
- Timer creation with simple description input (like Mac app)
- Completed timers from today shown in the panel
- On stop: create entry via existing `/entries` API with formatted duration + description

## Capabilities

### New Capabilities
- `tui-timer`: Timer management within the TUI including panel UI, timer lifecycle (start/pause/resume/stop), local persistence, and home screen status indicator

### Modified Capabilities
<!-- None - this is additive functionality that uses existing entry creation API -->

## Impact

- **TUI code**: New timer panel component, timer manager, persistence layer, home screen modifications for status indicator
- **Key bindings**: Add `[t]` for timer panel toggle, panel-specific bindings for timer actions
- **File system**: Creates `~/.tact/timers.json` for local timer state
- **API**: Uses existing `POST /entries` endpoint (no backend changes needed)
