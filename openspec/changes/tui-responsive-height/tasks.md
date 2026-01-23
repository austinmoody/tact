## 1. Calculate Available Lines

- [x] 1.1 Add `calculateAvailableEntryLines()` method to Home struct that computes available lines based on height
- [x] 1.2 Account for fixed UI elements: title bar (2 lines), help bar (2 lines), optional timer status (2 lines)
- [x] 1.3 Set minimum of 1 entry line regardless of calculated space

## 2. Dynamic Entry Rendering

- [x] 2.1 Modify View() to track lines used while rendering entries
- [x] 2.2 Stop rendering entries when available lines are exhausted
- [x] 2.3 Account for date headers (2 lines per group) and spacing between groups (1 line)

## 3. Scroll Indicator

- [x] 3.1 Track how many entries are hidden when display limit is reached
- [x] 3.2 Render scroll indicator (e.g., `↓ 5 more entries`) when entries are truncated
- [x] 3.3 Style the indicator using existing helpStyle

## 4. Cursor Management

- [x] 4.1 Add method to clamp cursor to visible range
- [x] 4.2 Call cursor clamping after height changes or entry list updates
- [x] 4.3 Ensure cursor stays on a visible entry during window resize

## 5. Testing

- [ ] 5.1 Manually verify help bar stays visible in short terminal
- [ ] 5.2 Verify scroll indicator appears when entries are hidden
- [ ] 5.3 Verify cursor clamping works on resize
- [ ] 5.4 Verify behavior with timer running (reduces available space)
