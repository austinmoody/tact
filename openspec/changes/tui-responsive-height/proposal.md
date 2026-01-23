## Why

When the terminal window is short, the TUI's bottom help bar gets cut off because the home screen displays a fixed number of entries (15) regardless of available vertical space. Users cannot see essential navigation hints like `[n] New  [t] Timer  [Enter] Details  [m] Menu  [r] Refresh  [q] Quit` when the terminal height is limited.

## What Changes

- Calculate maximum displayable entries based on terminal height minus fixed UI elements (title bar, help bar, spacing)
- Dynamically limit the number of entries shown to fit within available space
- Ensure the help bar is always visible at the bottom of the screen
- Add scroll indicators when there are more entries than can be displayed

## Capabilities

### New Capabilities

None - this is a modification to existing TUI dashboard behavior.

### Modified Capabilities

- `tui-dashboard`: Add responsive height handling to dynamically limit entries displayed based on terminal height, ensuring the help bar remains visible.

## Impact

- **Code**: `tui/ui/home.go` - Replace hardcoded `entriesLimit` with dynamic calculation based on `height`
- **User Experience**: Help bar will always be visible; users in short terminals will see fewer entries but can still navigate
- **Behavior Change**: Entry count adapts to terminal size instead of always showing up to 15 entries
