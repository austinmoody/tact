# TACT

**T**racking **A**nd **C**alculating **T**ime

*Set out runnin', but I track my time.*

## Components

- **backend/** - FastAPI REST API (Python)
- **tui/** - Terminal UI dashboard (Go/Bubbletea)

## Quick Start

```bash
# Start the backend
make run

# In another terminal, run the TUI
make tui-run
```

## TUI

The terminal user interface provides a read-only dashboard for viewing time codes and work types.

### Build & Run

```bash
# Build binary
make tui-build

# Run directly
make tui-run

# Run with explicit API URL
make tui-dev
```

### Configuration

The TUI connects to the backend API. Configure the URL via:

1. `--api` flag: `./tui/tact-tui --api http://localhost:2100`
2. `TACT_API_URL` env var: `TACT_API_URL=http://server:2100 ./tui/tact-tui`
3. Default: `http://localhost:2100`

### Key Bindings

| Key | Action |
|-----|--------|
| `j` / `↓` | Move cursor down |
| `k` / `↑` | Move cursor up |
| `h` / `l` | Switch between panes |
| `Enter` | Toggle detail view |
| `r` | Refresh data |
| `q` / `Ctrl+C` | Quit |
