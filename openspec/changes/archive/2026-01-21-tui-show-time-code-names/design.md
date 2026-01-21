## Context

The TUI currently shows time code IDs (like "FEDS-165") without the human-readable name. Users must memorize what each code means.

Current entry list layout in `home.go`:
```
> 2h worked on admin tasks              parsed
  meeting with team about planning      pending
```

Current entry detail "Parsed Fields" in `entry_detail.go`:
```
Parsed Fields:
  Duration: 2h (95%)
  Time Code: FEDS-165 (85%)
  Work Type: development (90%)
```

The API response (`EntryResponse`) only includes `time_code_id`, not the name. However, the TUI already fetches time codes list when entering edit mode (`model.TimeCode` has ID and Name fields).

## Goals / Non-Goals

**Goals:**
- Show time code name alongside ID in the main entry list
- Show time code name in entry detail parsed fields section
- Keep the display compact with reasonable truncation

**Non-Goals:**
- Modifying the backend API response (avoid backend changes for UI enhancement)
- Showing work type names (IDs are already human-readable like "development", "meeting")

## Decisions

### 1. Data fetching approach

**Decision**: Fetch time codes once on startup and cache them in Home component

**Rationale**:
- Time codes rarely change during a session
- Avoids N+1 requests (one per entry)
- Home component already has `*api.Client`, can easily fetch time codes
- Entry detail already has this pattern for edit mode

**Alternative considered**: Add `time_code_name` to API response - rejected to avoid backend changes for a UI-only enhancement

### 2. Entry list layout

**Decision**: Add time code column between user input and status

Current:
```
> 2h worked on admin tasks              parsed
```

New:
```
> 2h worked on admin tasks        FEDS-165 Admin Work     parsed
```

**Layout formula** (with truncation):
- Cursor: 2 chars
- User input: flexible, max `width - 45` chars
- Gap: 2 chars
- Time code + name: ~25 chars (truncate name if needed)
- Gap: 2 chars
- Status: 7 chars

**Rationale**:
- Places time code info in natural reading flow (left to right: what → where → status)
- Truncation keeps layout stable across different terminal widths

### 3. Entry detail display

**Decision**: Show "Time Code: ID - Name (confidence%)"

Current:
```
  Time Code: FEDS-165 (85%)
```

New:
```
  Time Code: FEDS-165 - Admin Work (85%)
```

**Rationale**:
- Minimal change, keeps existing format
- Dash separator clearly distinguishes ID from name
- Name truncated to ~20 chars if needed

### 4. Time code name lookup

**Decision**: Create a lookup map `map[string]string` from time code ID to name

**Rationale**:
- O(1) lookup per entry
- Simple implementation
- Time codes list is small (typically <100), so memory is negligible

## Risks / Trade-offs

**Risk**: Time code list might be stale if codes are added/modified
→ **Mitigation**: Refresh time codes on manual refresh (r key). Also refreshed when entering edit mode.

**Trade-off**: Extra API call on startup vs backend API change
→ **Accepted**: Single extra request is negligible; keeps changes TUI-only
