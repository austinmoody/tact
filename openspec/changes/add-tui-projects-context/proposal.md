# Proposal: Add TUI Support for Projects and Context

## Summary

Add TUI screens and modals to manage Projects and Context Documents, enabling users to configure RAG-based parsing context through the terminal interface.

## Motivation

The backend now supports Projects (for grouping time codes) and Context Documents (for RAG-enhanced parsing). Users need a way to manage these entities through the TUI without relying on curl commands.

## Scope

### In Scope

- Projects management screen (list, add, edit, delete)
- Context management for Projects (via `c` key)
- Context management for Time Codes (via `c` key)
- API client methods for Projects and Context endpoints
- New models for Project and ContextDocument

### Out of Scope

- Bulk import of context documents
- Context document search/filtering
- Showing which context was used during parsing

## Design Decisions

1. **Projects as top-level menu item** - Projects appear alongside Time Codes and Work Types in the main menu
2. **Context managed inline** - Press `c` on a selected Project or Time Code to view/manage its context documents
3. **Truncated display** - Context shown truncated in lists, full content in detail/edit modals

## Affected Specs

- `tui-dashboard` - Add Projects screen, Context management, menu updates
