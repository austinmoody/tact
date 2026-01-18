# Change: Add macOS Timer App (Tact Timer)

## Why

The existing Tact system has a backend API and a terminal UI, but no native macOS experience for quick time tracking. Users want a lightweight, always-available timer that lives in the Dock and allows starting/stopping timers with minimal friction - without opening a terminal or browser.

## What Changes

- **NEW**: Native macOS app "Tact Timer" built with Swift + AppKit
- Dock-based interaction via right-click context menu
- Small popup window for entering timer descriptions
- Timer list window for viewing/managing multiple timers
- Automatic timer persistence across app restarts
- Integration with existing Tact backend API for saving time entries
- Preferences window for API URL configuration

## Impact

- Affected specs: None (new capability)
- Affected code: New `/macos` directory for the Xcode project
- No changes to existing backend API, TUI, or MCP server
