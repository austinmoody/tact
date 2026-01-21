## Why

Time code IDs like "FEDS-165" are difficult to remember without seeing the associated name. Users currently have to memorize what each code means or navigate to the time codes list to look them up. Showing the time code name alongside the ID improves usability and reduces cognitive load.

## What Changes

- Display time code ID and name in the main entry list (between user input and status columns)
- Show time code name in the entry detail "Parsed Fields" section alongside the ID and confidence percentage
- Truncate long time code names as needed to fit available space

## Capabilities

### New Capabilities
None - this enhances existing display functionality

### Modified Capabilities
- `tui-dashboard`: Update entry list to show time code ID + name; update entry detail parsed fields to include time code name

## Impact

- **TUI**: `tui/ui/home.go` - Update entry row rendering to include time code column
- **TUI**: `tui/ui/entry_detail.go` - Update parsed fields section to show time code name
- **API**: May need to ensure time code name is included in entry list response (check if already present)
