## 1. Model Changes
- [ ] 1.1 Add `stopped` case to `TimerState` enum
- [ ] 1.2 Add `stoppedAt: Date?` property to `TactTimer` struct
- [ ] 1.3 Add computed property `formattedFinalDuration` for completed timer display

## 2. Timer Manager Changes
- [ ] 2.1 Modify `stopTimer()` to set state to `.stopped` and `stoppedAt` instead of removing
- [ ] 2.2 Add computed properties `activeTimers` and `completedTodayTimers` for filtering
- [ ] 2.3 Add cleanup logic in `load()` to remove completed timers from previous days
- [ ] 2.4 Add `startNewFromCompleted(id:)` method (or reuse existing `startNewTimer`)

## 3. Timer List Window Changes
- [ ] 3.1 Update table view to show two sections (active and completed)
- [ ] 3.2 Add "Completed Today" section header with separator
- [ ] 3.3 Create `CompletedTimerCellView` with description, duration, and "Start New" button
- [ ] 3.4 Wire "Start New" button to create new timer with copied description
- [ ] 3.5 Hide completed section when no completed timers exist
- [ ] 3.6 Update empty state to only show when both sections are empty

## 4. Testing & Verification
- [ ] 4.1 Verify stopping a timer moves it to completed section
- [ ] 4.2 Verify "Start New" creates fresh timer with same description
- [ ] 4.3 Verify completed timers persist across app restart (same day)
- [ ] 4.4 Verify old completed timers are cleaned up on launch (next day)
- [ ] 4.5 Verify active timer count in dock menu excludes completed timers

## Verification

1. Build and run Tact Timer from Xcode
2. Start a timer, let it run briefly, then Stop it
3. Verify timer appears in "Completed Today" section with final duration
4. Click "Start New" on the completed timer
5. Verify new timer starts with same description
6. Verify original completed timer remains in completed section
7. Quit and relaunch app → verify completed timer still shows
8. (Optional) Change system date to tomorrow, relaunch → verify cleanup
