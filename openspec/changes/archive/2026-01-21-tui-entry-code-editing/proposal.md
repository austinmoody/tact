## Why

When the AI parses a time entry incorrectly (wrong time_code or work_type), users need a way to correct it. These corrections should trigger the "learn" feature so the AI improves over time. The current entry editing only supports user_input and entry_date fields, which don't constitute AI parsing corrections.

## What Changes

- Add time_code selection to entry edit mode with dropdown showing available time codes
- Add work_type selection to entry edit mode with dropdown showing available work types
- Use `?learn=true` when saving if time_code or work_type was changed (these are AI parsing corrections)
- Keep `?learn=false` for user_input/entry_date-only changes (not AI corrections)
- Fetch available time codes and work types when entering edit mode

## Capabilities

### New Capabilities
- `entry-code-selection`: Dropdown selection UI for time_code and work_type fields in entry edit mode, with learn flag logic

### Modified Capabilities
- `entry-editing`: Add time_code and work_type fields to edit mode, conditional learn flag based on which fields changed

## Impact

- **TUI**: `tui/ui/entry_detail.go` - Add time_code/work_type selection fields
- **TUI**: `tui/api/client.go` - Update `EntryUpdate` struct to include time_code_id and work_type_id
- **API**: No backend changes needed - PATCH /entries already supports these fields and learn flag
