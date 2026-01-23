## Context

The TUI home screen currently displays up to 15 entries (`entriesLimit = 15`) regardless of terminal height. The View() function renders:
1. Title bar (1 line) + blank line
2. Entries section with date headers (variable lines)
3. Timer status (optional, 2 lines if running)
4. Help bar (2 lines including blank line before)

When the terminal is short, the help bar gets cut off because entries render without considering available vertical space.

The `Home` struct already tracks `height` via `SetSize(width, height)`, but this value is not used to limit entry display.

## Goals / Non-Goals

**Goals:**
- Calculate available lines for entries based on terminal height
- Dynamically limit displayed entries to fit within available space
- Always keep the help bar visible
- Show scroll indicator when entries are truncated

**Non-Goals:**
- Scrollable entry list with viewport (too complex for this change)
- Pagination with explicit page controls
- Persisting scroll position across refreshes

## Decisions

### Decision 1: Calculate available entry lines dynamically

**Approach**: Calculate `maxVisibleEntries` at render time in View() based on `h.height` minus fixed UI elements.

Fixed UI overhead:
- Title bar + newline: 2 lines
- Help bar + newline before: 2 lines
- Timer status (if running): 2 lines
- Minimum padding: 1 line

**Formula**: `availableLines = height - 4 (or -6 if timer running)`

Each date group uses:
- Date header: 1 line
- Separator: 1 line
- Entries: 1 line each
- Space between groups: 1 line (except first)

**Rationale**: Calculating at render time ensures the display adapts immediately to window resizes without needing to re-fetch data.

**Alternative considered**: Pre-calculate in SetSize() - rejected because timer visibility changes dynamically and affects available space.

### Decision 2: Limit entries displayed, not entries fetched

**Approach**: Continue fetching `entriesLimit` (15) entries from API, but only display entries that fit in available space.

**Rationale**:
- Keeps cursor navigation working within fetched data
- Allows showing more entries when terminal grows without re-fetch
- Simpler than dynamic API limit calculation

### Decision 3: Show scroll indicator when entries are hidden

**Approach**: When there are more entries than can be displayed, show a subtle indicator like `↓ 5 more entries` after the last visible entry.

**Rationale**: Users need to know there's more content they can't see. The indicator uses minimal vertical space (1 line) and communicates the count of hidden entries.

### Decision 4: Cursor bounds during resize

**Approach**: If cursor is beyond visible range after resize, clamp cursor to last visible entry.

**Rationale**: Prevents cursor pointing at invisible entry, which would confuse users.

## Risks / Trade-offs

**[Risk]** Date headers consume variable space, making exact calculation complex
→ **Mitigation**: Track lines used while rendering; stop when limit reached. May cut mid-group but ensures help bar visibility.

**[Risk]** Very short terminals (< 10 lines) may show no entries
→ **Mitigation**: Set minimum of 1 entry regardless of calculated space. Help bar less critical than showing at least one entry.

**[Trade-off]** Users in short terminals see fewer entries
→ **Accepted**: This is the intended behavior. Users can resize terminal for more entries.

## Open Questions

None - implementation is straightforward given the existing height tracking infrastructure.
