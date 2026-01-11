# Proposal: Enhance TUI with Entry Management

## Why

The TUI currently displays time codes and work types in a split-pane view, but doesn't support the core workflow: creating and managing time entries. Now that the entries API and LLM parsing are complete, the TUI should become the primary interface for quick time entry.

## What Changes

Redesign the TUI with an entry-focused main screen:
1. Replace split-pane view with entries-focused home screen
2. Add quick entry input via keybinding
3. Show last 5 entries on home screen
4. Entry detail modal with reparse option
5. Menu system for accessing time codes and work types management
6. Simplified quick-add for time codes and work types

## Scope

### In Scope

- New home screen showing last 5 entries
- Quick entry input modal (keybinding: `n`)
- Entry detail modal (Enter on selected entry)
- Reparse action from detail modal
- Menu system (keybinding: `m` or `?`)
- Time codes management screen (list, edit, quick-add)
- Work types management screen (list, edit, quick-add)
- API client extensions for entries, time code mutations, work type mutations

### Out of Scope

- Full-form entry editing (just view + reparse for now)
- Entry deletion from TUI
- Filtering/searching entries
- Pagination beyond last 5 entries
- Reports or analytics views

## Key Decisions

1. **Entry-focused home** - Main screen prioritizes the primary workflow: quick entry
2. **Modal pattern** - Use modals for input, details, and quick-add rather than screen transitions
3. **Menu for management** - Time codes and work types accessible via menu, not cluttering main screen
4. **Simplified quick-add** - New time codes/work types use minimal required fields only

## Affected Specs

- `tui-dashboard` (MODIFIED) - Update requirements for new home screen and navigation
