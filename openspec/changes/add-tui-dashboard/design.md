# Design: TUI Dashboard

## Overview

A Go-based terminal user interface using Bubbletea that connects to the Tact FastAPI backend. Provides a read-only dashboard for viewing time codes and work types.

## Directory Structure

```
tui/
├── go.mod
├── go.sum
├── main.go              # Entry point, flag parsing
├── api/
│   └── client.go        # HTTP client for backend API
├── model/
│   ├── timecode.go      # TimeCode struct
│   └── worktype.go      # WorkType struct
└── ui/
    ├── dashboard.go     # Main dashboard view
    ├── styles.go        # Lip Gloss styles
    └── keys.go          # Key bindings
```

## Architecture

### Bubbletea Model-View-Update

```
┌─────────────────────────────────────────┐
│                  Model                   │
│  - timeCodes []TimeCode                 │
│  - workTypes []WorkType                 │
│  - cursor int                           │
│  - activePane (codes|types)             │
│  - loading bool                         │
│  - err error                            │
└─────────────────────────────────────────┘
           │                    ▲
           │ View()             │ Update()
           ▼                    │
┌─────────────┐          ┌─────────────┐
│   Screen    │          │   Message   │
│  (string)   │          │  (KeyMsg,   │
│             │          │   DataMsg)  │
└─────────────┘          └─────────────┘
```

### API Client

Simple HTTP client with methods:
- `FetchTimeCodes() ([]TimeCode, error)`
- `FetchWorkTypes() ([]WorkType, error)`

Configurable base URL via:
1. `--api` flag: `./tact-tui --api http://localhost:2100`
2. `TACT_API_URL` env var
3. Default: `http://localhost:2100`

## Dashboard Layout

```
┌─────────────────────────────────────────────────────────┐
│  Tact Dashboard                              [q] quit   │
├───────────────────────────┬─────────────────────────────┤
│  Time Codes               │  Work Types                 │
│  ─────────────            │  ──────────                 │
│  > PROJ-001 Project Alpha │  > development              │
│    PROJ-002 Project Beta  │    meeting                  │
│    ADMIN-01 Admin Tasks   │    code-review              │
│                           │    devops                   │
│                           │                             │
│                           │                             │
├───────────────────────────┴─────────────────────────────┤
│  PROJ-001: Project Alpha                                │
│  Keywords: alpha, main                                  │
│  Examples: "working on alpha", "alpha dev work"         │
└─────────────────────────────────────────────────────────┘
  [j/k] navigate  [h/l] switch pane  [enter] details  [q] quit
```

## Key Bindings

| Key | Action |
|-----|--------|
| `j` / `↓` | Move cursor down |
| `k` / `↑` | Move cursor up |
| `h` / `l` | Switch between panes |
| `Enter` | Toggle detail view for selected item |
| `r` | Refresh data from API |
| `q` / `Ctrl+C` | Quit |

## Styling

Using Lip Gloss for:
- Border styles (rounded corners)
- Color scheme (respects terminal theme)
- Focused vs unfocused pane highlighting
- Selected item highlighting

## Error Handling

- Connection errors shown in status bar
- Retry with `r` key
- Graceful degradation (show error, don't crash)

## Build & Run

```bash
# Build
cd tui && go build -o tact-tui .

# Run
./tact-tui                                    # Uses default localhost:2100
./tact-tui --api http://localhost:2100        # Explicit flag
TACT_API_URL=http://server:2100 ./tact-tui    # Env var
```

## Makefile Integration

```makefile
tui-build:
    cd tui && go build -o tact-tui .

tui-run:
    cd tui && go run .

tui-dev:
    cd tui && go run . --api http://localhost:2100
```
