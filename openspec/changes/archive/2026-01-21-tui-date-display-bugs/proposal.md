## Why

The TUI has two date-related display bugs that cause confusing user experiences:
1. Time entries appear under incorrect date headers on the home screen
2. The timer panel's "Completed Today" section shows timers from previous days

These bugs make the TUI unreliable for tracking work, as users cannot trust the displayed dates.

## What Changes

- **Fix backend API sorting**: Add `order_by(entry_date DESC, created_at DESC)` to the entries endpoint so newest entries are returned first when using `limit`
- **Fix TUI home screen sorting**: Sort entries by `entry_date` (descending) then `created_at` (descending) instead of only by `created_at`, ensuring entries appear under the correct date header
- **Fix timer "Completed Today" timezone handling**: Use local timezone midnight truncation instead of UTC midnight when determining which timers completed today
- **Fix timer cleanup logic**: Apply the same timezone fix to `cleanupOldCompleted()` so old timers are properly removed

## Capabilities

### New Capabilities
- None

### Modified Capabilities
- None (these are bug fixes to existing behavior, not spec-level requirement changes)

## Impact

**Affected Code:**
- `backend/src/tact/routes/entries.py:58` - Missing order_by clause in list_entries endpoint
- `tui/ui/home.go:80-83` - Entry sorting logic
- `tui/timer/manager.go:126-138` - `CompletedToday()` function
- `tui/timer/manager.go:159-180` - `cleanupOldCompleted()` function

**Root Causes Identified:**
1. **API bug**: The `/entries` endpoint had no `order_by` clause, returning entries in arbitrary order. With `limit=15`, newest entries were cut off before reaching the TUI.
2. **TUI sort bug**: Entries were sorted by `created_at` but grouped by `entry_date`, causing misalignment when an entry's `entry_date` differs from its `created_at` date
3. **Timer bug**: `time.Now().Truncate(24 * time.Hour)` truncates to UTC midnight, not local midnight. A timer stopped at 9:37 PM EST on Jan 20 (which is 2:37 AM UTC on Jan 21) incorrectly appears as "completed today" on Jan 21 local time
