## Context

The TUI already supports editing `user_input` and `entry_date` fields in `EntryDetailModal` (`tui/ui/entry_detail.go`). The backend PATCH `/entries/{id}` endpoint supports `time_code_id` and `work_type_id` fields, plus a `?learn=true` flag that triggers context document creation when correcting AI parsing mistakes.

Current entry edit mode:
- Press `e` to enter edit mode
- Text inputs for `user_input` and `entry_date`
- `Tab` to switch between fields
- `Enter` to save with `learn=false` (these aren't AI corrections)
- `Esc` to cancel

The API client already has `FetchTimeCodes()` and `FetchWorkTypes()` methods that return available options.

## Goals / Non-Goals

**Goals:**
- Add time_code selection with dropdown showing available time codes
- Add work_type selection with dropdown showing available work types
- Use `?learn=true` when time_code or work_type is changed (AI parsing corrections)
- Keep `?learn=false` when only user_input/entry_date change
- Fetch available options when entering edit mode

**Non-Goals:**
- Creating new time codes or work types from the edit modal
- Inline search/filter for large lists (simple list selection for now)
- Editing duration directly (recalculated from time_code)

## Decisions

### 1. Dropdown selection UI approach

**Decision**: Use simple list selection with `j/k` navigation when a dropdown field is focused

**Rationale**:
- Consistent with TUI navigation patterns (j/k for up/down)
- No external dependencies needed
- Can show code ID and description in each row

**Alternative considered**: Separate modal for selection - rejected as adds context switching

### 2. Field layout in edit mode

**Decision**: Add time_code and work_type as additional focusable fields below user_input and date

**Rationale**:
- Natural top-to-bottom field order
- Tab cycles through all 4 fields: user_input → date → time_code → work_type → user_input
- Each code field shows current selection inline

### 3. Dropdown list behavior

**Decision**: Show dropdown list inline below the focused code field, max 5 visible items with scroll indicators

**Rationale**:
- Keeps user in context (no separate modal)
- 5 items is enough to see options while not overwhelming the view
- Scroll indicators (`↑`/`↓`) show more items available

### 4. Learn flag logic

**Decision**: Track original values on edit mode entry, compare on save:
- If `time_code_id` or `work_type_id` changed → `learn=true`
- If only `user_input` or `entry_date` changed → `learn=false`

**Rationale**:
- Clear distinction: code field changes are AI corrections, text/date changes are not
- Matches backend expectation that learn flag means "use this correction for training"

### 5. API struct extension

**Decision**: Add `TimeCodeID` and `WorkTypeID` to `EntryUpdate` struct

**Rationale**:
- Backend already accepts these in PATCH body
- Follows existing pattern for optional fields with `omitempty`

### 6. Fetching options

**Decision**: Fetch time codes and work types asynchronously when entering edit mode

**Rationale**:
- Data may have changed since TUI started
- Non-blocking fetch with loading indicator
- Cache in modal until edit mode exits

## Risks / Trade-offs

**Risk**: Many time codes/work types makes dropdown unwieldy
→ **Mitigation**: 5-item visible window with scroll; future enhancement could add search/filter

**Risk**: Slow API response when fetching options
→ **Mitigation**: Show "Loading..." in dropdown area, disable save until loaded

**Risk**: Time code changed but work type wasn't (or vice versa) - should learn still trigger?
→ **Mitigation**: Yes, any code field change triggers learn. Partial corrections still valuable.

**Trade-off**: Fetching options every time edit mode is entered vs caching globally
→ **Accepted**: Per-edit-session fetch ensures fresh data; overhead is minimal

## Implementation Components

```
tui/
├── api/
│   └── client.go          # Extend EntryUpdate struct with TimeCodeID, WorkTypeID
└── ui/
    └── entry_detail.go    # Add dropdown selection UI, learn flag logic
```

### API Client Changes

```go
type EntryUpdate struct {
    UserInput  *string `json:"user_input,omitempty"`
    EntryDate  *string `json:"entry_date,omitempty"`
    TimeCodeID *string `json:"time_code_id,omitempty"`
    WorkTypeID *string `json:"work_type_id,omitempty"`
}
```

### Modal State Changes

```go
type EntryDetailModal struct {
    // existing fields...

    // Code selection
    timeCodes        []model.TimeCode
    workTypes        []model.WorkType
    loadingCodes     bool
    selectedTimeCode int  // index in timeCodes slice, -1 for none
    selectedWorkType int  // index in workTypes slice, -1 for none

    // For learn flag logic
    originalTimeCodeID *string
    originalWorkTypeID *string
}
```

### Focus Field Values

```go
const (
    fieldUserInput = 0
    fieldDate      = 1
    fieldTimeCode  = 2
    fieldWorkType  = 3
)
```

### Key Bindings (Edit Mode)

- `Tab` / `Shift+Tab` - Cycle through all 4 fields
- `j` / `↓` - Next item in dropdown (when time_code or work_type focused)
- `k` / `↑` - Previous item in dropdown
- `Enter` - Save changes
- `Esc` - Cancel edit mode

### Dropdown Rendering

```
> Time Code:
  [ABC-001] Project Alpha - Development
  ┌─────────────────────────────────────┐
  │ ABC-001  Project Alpha - Dev    ✓   │
  │ ABC-002  Project Alpha - Test       │
  │ XYZ-100  Client Work               │
  │ ↓ 3 more                           │
  └─────────────────────────────────────┘
```

### Learn Flag Logic

```go
func (m *EntryDetailModal) shouldLearn() bool {
    currentTimeCode := m.getSelectedTimeCodeID()
    currentWorkType := m.getSelectedWorkTypeID()

    timeCodeChanged := !ptrEqual(m.originalTimeCodeID, currentTimeCode)
    workTypeChanged := !ptrEqual(m.originalWorkTypeID, currentWorkType)

    return timeCodeChanged || workTypeChanged
}
```
