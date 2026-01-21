## Context

The TUI displays time entries grouped by date on the home screen and shows completed timers in a "Completed Today" section in the timer panel. Both features have timezone/date handling bugs that cause incorrect displays.

**Current State:**
- Home screen sorts entries by `created_at` but groups them by `entry_date`, causing misalignment
- Timer manager uses `time.Now().Truncate(24 * time.Hour)` which truncates to UTC midnight, not local timezone midnight

## Goals / Non-Goals

**Goals:**
- Entries appear under the correct date header on the home screen
- Timer "Completed Today" accurately reflects timers stopped today in the user's local timezone
- Old completed timers are properly cleaned up based on local timezone

**Non-Goals:**
- Supporting explicit timezone configuration (use system timezone)
- Changing how `entry_date` is determined by the backend
- Modifying the API

## Decisions

### 1. Home Screen: Sort by `entry_date` then `created_at`

**Decision:** Change the sort in `home.go` to sort primarily by `entry_date` (descending), then by `created_at` (descending) as a tiebreaker.

**Rationale:** Since date headers are based on `entry_date`, the sort order must match. Using `created_at` as a secondary sort preserves chronological order within each date group.

**Alternative considered:** Sort only by `entry_date` - rejected because entries on the same date would have arbitrary order.

**Implementation:**
```go
sort.Slice(h.entries, func(i, j int) bool {
    if h.entries[i].EntryDate != h.entries[j].EntryDate {
        return h.entries[i].EntryDate > h.entries[j].EntryDate
    }
    return h.entries[i].CreatedAt.After(h.entries[j].CreatedAt.Time)
})
```

### 2. Timer: Use Local Timezone for "Today" Calculation

**Decision:** Create a helper function `localMidnight()` that returns the start of today in local timezone, replacing `time.Now().Truncate(24 * time.Hour)`.

**Rationale:** `Truncate(24 * time.Hour)` operates in UTC, not local time. Go's `time.Date()` function respects the local timezone when using `time.Local`.

**Alternative considered:** Store timezone with each timer - rejected as over-engineering; system timezone is the expected behavior.

**Implementation:**
```go
func localMidnight() time.Time {
    now := time.Now()
    return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
}
```

### 3. Apply Same Fix to `cleanupOldCompleted()`

**Decision:** Use the same `localMidnight()` helper in `cleanupOldCompleted()` for consistency.

**Rationale:** This function uses the same flawed UTC truncation logic. Both functions need to agree on what "today" means.

## Risks / Trade-offs

**[Risk] Daylight Saving Time edge cases** → The `time.Local` approach handles DST correctly because `time.Date()` with `time.Local` accounts for DST transitions. No additional mitigation needed.

**[Risk] Timers created before this fix may behave unexpectedly during the fix rollout** → Minor risk; timers are ephemeral and cleaned up daily. No migration needed.

**[Trade-off] No timezone configuration** → Users who want a different timezone must change their system timezone. This matches typical CLI tool behavior and avoids configuration complexity.
