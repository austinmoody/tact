## 1. Add Time Code Caching to Home

- [x] 1.1 Add `timeCodes []model.TimeCode` field to Home struct
- [x] 1.2 Add `timeCodeNames map[string]string` lookup map to Home struct
- [x] 1.3 Create `fetchTimeCodes()` command that calls `client.FetchTimeCodes()`
- [x] 1.4 Update `Init()` to fetch time codes alongside entries (use `tea.Batch`)
- [x] 1.5 Handle `timeCodesMsg` to populate timeCodes and build lookup map
- [x] 1.6 Update `Refresh()` to also re-fetch time codes

## 2. Update Entry List Display

- [x] 2.1 Create helper function `getTimeCodeDisplay(id string) string` that returns "ID Name" or "ID" if name not found
- [x] 2.2 Update `renderEntryLine()` to include time code column between user input and status
- [x] 2.3 Adjust user input max length calculation to account for time code column width
- [x] 2.4 Truncate time code name if combined "ID Name" exceeds ~25 chars
- [x] 2.5 Handle entries without time code (show empty space to maintain alignment)

## 3. Update Entry Detail Display

- [x] 3.1 Add `timeCodeNames map[string]string` parameter or field to EntryDetailModal
- [x] 3.2 Pass time code names map from App to EntryDetailModal when opening
- [x] 3.3 Update parsed fields section to show "Time Code: ID - Name (confidence%)"
- [x] 3.4 Truncate name if too long, keeping ID and confidence visible

## 4. Testing

- [x] 4.1 Verify entry list shows time code ID and name for parsed entries
- [x] 4.2 Verify entry list handles entries without time code gracefully
- [x] 4.3 Verify entry detail shows time code name in parsed fields
- [x] 4.4 Verify refresh updates time code names
- [x] 4.5 Test with long time code names to verify truncation works
