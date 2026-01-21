## 1. Fix Backend API Sorting

- [x] 1.1 Add `order_by(entry_date DESC, created_at DESC)` to entries endpoint in `backend/src/tact/routes/entries.py`

## 2. Fix TUI Home Screen Entry Sorting

- [x] 2.1 Update sort logic in `tui/ui/home.go` to sort by `entry_date` (descending) then `created_at` (descending)
- [x] 2.2 Verify entries appear under correct date headers with mixed entry_date/created_at values

## 3. Fix Timer Timezone Handling

- [x] 3.1 Add `localMidnight()` helper function to `tui/timer/manager.go`
- [x] 3.2 Update `CompletedToday()` to use `localMidnight()` instead of UTC truncation
- [x] 3.3 Update `cleanupOldCompleted()` to use `localMidnight()` for consistency

## 4. Verification

- [x] 4.1 Restart backend to pick up API changes
- [x] 4.2 Test home screen shows today's entries when entries exist for today
- [x] 4.3 Test timer panel "Completed Today" only shows timers stopped after local midnight
- [x] 4.4 Test old completed timers are cleaned up on app restart
