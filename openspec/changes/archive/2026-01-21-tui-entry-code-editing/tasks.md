## 1. API Client Updates

- [x] 1.1 Extend `EntryUpdate` struct with `TimeCodeID *string` and `WorkTypeID *string` fields

## 2. Modal State Management

- [x] 2.1 Add fields for storing fetched time codes and work types lists
- [x] 2.2 Add `loadingCodes` bool for tracking fetch state
- [x] 2.3 Add `selectedTimeCode` and `selectedWorkType` index fields
- [x] 2.4 Add `originalTimeCodeID` and `originalWorkTypeID` for learn flag comparison
- [x] 2.5 Update `focusedField` constants to include time_code and work_type (0-3)

## 3. Edit Mode Entry

- [x] 3.1 Create async command to fetch time codes and work types when entering edit mode
- [x] 3.2 Add message types for successful fetch and fetch error
- [x] 3.3 Handle fetch response messages in Update()
- [x] 3.4 Store original time_code_id and work_type_id on edit mode entry
- [x] 3.5 Initialize selected indices based on current entry values

## 4. Dropdown Navigation

- [x] 4.1 Handle `j`/`↓` keys to select next item when time_code/work_type focused
- [x] 4.2 Handle `k`/`↑` keys to select previous item
- [x] 4.3 Implement bounds checking for list navigation
- [x] 4.4 Update Tab cycling to include all 4 fields

## 5. Rendering

- [x] 5.1 Render time_code field with current selection and focus indicator
- [x] 5.2 Render work_type field with current selection and focus indicator
- [x] 5.3 Render dropdown list when time_code focused (max 5 visible, scroll indicators)
- [x] 5.4 Render dropdown list when work_type focused
- [x] 5.5 Show loading indicator while codes are being fetched
- [x] 5.6 Update help text to include j/k navigation hint

## 6. Save Logic

- [x] 6.1 Add `shouldLearn()` method comparing current vs original code values
- [x] 6.2 Include time_code_id in update payload if selection changed
- [x] 6.3 Include work_type_id in update payload if selection changed
- [x] 6.4 Pass correct learn flag to UpdateEntry based on shouldLearn()
- [x] 6.5 Handle nil/empty selections (user can deselect)
