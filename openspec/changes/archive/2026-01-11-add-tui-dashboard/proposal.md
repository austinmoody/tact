# Proposal: Add TUI Dashboard

## Summary

Add a terminal user interface (TUI) built with Go and Bubbletea that provides a read-only dashboard view of time codes and work types. This introduces Go as a second language for client applications while keeping the Python backend.

## Motivation

- Provide a quick, keyboard-driven way to view time tracking data
- Learn Go and the Bubbletea TUI framework
- Establish patterns for future TUI features (entry submission, etc.)
- Complement the future web UI and MCP server as another client

## Scope

### In Scope

- New `tui/` directory with Go module
- Bubbletea-based TUI application
- Combined dashboard view showing time codes and work types
- Vim-like keyboard navigation (j/k, enter, q to quit)
- Configurable backend URL via `--api` flag or `TACT_API_URL` env var
- Makefile targets for building and running the TUI

### Out of Scope

- Entry management (future enhancement after entries API)
- Write operations (creating/updating/deleting)
- Authentication (single-user tool)

## Key Decisions

1. **Go + Bubbletea** - Modern TUI framework with Elm Architecture
2. **Lip Gloss for styling** - Charm's styling library for colors/borders
3. **Separate Go module** - `tui/` directory with its own `go.mod`
4. **API client** - Simple HTTP client calling the FastAPI backend
5. **Vim-style keys** - j/k navigation, enter to select/expand, q to quit

## Affected Specs

- `tui-dashboard` (NEW) - Dashboard view capabilities
