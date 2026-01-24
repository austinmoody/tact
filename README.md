# TACT

**T**racking **A**nd **C**alculating **T**ime

*Set out runnin', but I track my time.*

## Components

- **backend/** - FastAPI REST API (Python)
- **tui/** - Terminal UI dashboard (Go/Bubbletea)
- **webui/** - Web UI dashboard (Go/Templ/HTMX)
- **mcp/** - MCP server for AI assistants (Python)
- **macos/** - Native macOS timer app (Swift/AppKit)

## Quick Start

```bash
# Start the backend
make run

# In another terminal, run the TUI
make tui-run
```

## MCP Server

The MCP (Model Context Protocol) server exposes the Tact API to AI assistants like Claude Desktop, Goose, and GitHub Copilot for natural language time tracking.

### Setup

```bash
# Start the API
docker compose up -d api

# Build the MCP server
make docker-build
```

### Claude Desktop Configuration

Add to `~/Library/Application Support/Claude/claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "tact": {
      "command": "docker",
      "args": [
        "run", "-i", "--rm",
        "--network", "tact_default",
        "-e", "TACT_API_URL=http://api:2100",
        "tact-mcp"
      ]
    }
  }
}
```

Once configured, you can ask Claude things like:
- "Log 2 hours of development work on Project Alpha"
- "Show me my time entries for this week"
- "List all active time codes"

See [mcp/README.md](mcp/README.md) for full documentation including Goose/VS Code configuration and troubleshooting.

## TUI

The terminal user interface provides a dashboard for managing time entries, projects, time codes, work types, and RAG context documents.

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

### Features

- **Home Screen**: View and manage time entries, create new entries with natural language
- **Timer**: Track time with a floating timer panel, supports multiple timers (one active at a time), auto-creates entries on stop
- **Projects**: Create and manage projects for grouping time codes
- **Time Codes**: Manage billing codes with project associations
- **Work Types**: Manage categories of work (development, meetings, etc.)
- **Context Management**: Add RAG context documents to projects and time codes to improve LLM parsing

### Key Bindings

#### Global

| Key | Action |
|-----|--------|
| `j` / `↓` | Move cursor down |
| `k` / `↑` | Move cursor up |
| `m` | Open menu |
| `t` | Open timer panel |
| `r` | Refresh data |
| `q` / `Ctrl+C` | Quit |
| `Esc` | Close modal / Go back |

#### List Screens (Projects, Time Codes, Work Types)

| Key | Action |
|-----|--------|
| `a` | Add new item |
| `e` | Edit selected item |
| `d` | Delete selected item |
| `c` | Manage context (Projects, Time Codes) |

#### Modals

| Key | Action |
|-----|--------|
| `Tab` | Next field |
| `Shift+Tab` | Previous field |
| `Enter` | Save / Submit |
| `Ctrl+S` | Save (in context editor) |
| `Esc` | Cancel / Close |

#### Timer Panel

| Key | Action |
|-----|--------|
| `n` | New timer |
| `p` | Pause selected timer |
| `r` | Resume selected timer |
| `s` | Stop timer (creates entry) |
| `d` | Delete timer |
| `t` / `Esc` | Close panel |

## Web UI

A browser-based dashboard built with Go, Templ templates, and HTMX for a responsive, server-rendered experience.

### Build & Run

```bash
# Build binary
make webui-build

# Run directly (connects to API at http://localhost:2100)
make webui-run

# Run with explicit API URL
make webui-dev
```

### Docker

```bash
# Build the Docker image
docker build -t tact-webui ./webui

# Run with docker-compose
docker compose up webui
```

### Configuration

The Web UI connects to the backend API. Configure the URL via:

1. `--api` flag: `./webui/tact-webui --api http://localhost:2100`
2. `TACT_API_URL` env var: `TACT_API_URL=http://server:2100 ./webui/tact-webui`
3. Default: `http://localhost:2100`

The Web UI runs on port **2200** by default.

### Features

- **Entries**: View and manage time entries with filtering by status and date range
- **Timer**: Track time with start/pause/resume/stop controls, live elapsed time display via SSE
- **Projects**: Create and manage projects with search/filter functionality
- **Time Codes**: Manage billing codes with project associations
- **Work Types**: Manage categories of work
- **Context Management**: Add RAG context documents to projects and time codes

### Technology Stack

- **Go** - Backend server
- **Templ** - Type-safe HTML templating
- **HTMX** - Dynamic interactions without JavaScript frameworks
- **Pico CSS** - Minimal, classless CSS framework
- **SSE** - Server-Sent Events for real-time timer updates

## macOS App (Tact Timer)

A native macOS app for quick timer-based time tracking. Lives in the Dock for easy access.

### Requirements

- macOS 13.0+
- Xcode 15+ (for building)

### Build & Run

```bash
# Open in Xcode
open macos/TactTimer/TactTimer.xcodeproj

# Build and run with Cmd+R
```

Or build from command line:

```bash
cd macos/TactTimer
xcodebuild -scheme TactTimer -configuration Debug build
```

### Features

- **Timer List Window**: Shows on launch, displays all active timers with live elapsed time
- **Multiple Timers**: Run multiple concurrent timers (only one active at a time, others paused)
- **Quick Start**: Click "+ New Timer" or use Cmd+N to start tracking
- **Dock Menu**: Right-click the dock icon for quick access to timers and preferences
- **Auto-Pause**: Starting a new timer automatically pauses the current one
- **API Integration**: Stopping a timer saves the entry to the Tact backend
- **Completed Today**: Stopped timers remain visible for quick restart with "Start New"
- **Persistence**: Timers survive app restarts; old completed timers auto-cleanup daily

### Configuration

Open Preferences (Cmd+,) to set the API URL. Default: `http://localhost:2100`

### Usage

1. Start the backend API (`make run`)
2. Launch Tact Timer
3. Click "+ New Timer" and enter a description
4. Work on your task
5. Click "Stop" when done - the entry is saved to the backend with the elapsed duration

## Potential Future Enhancements

### Ollama Structured Output Schema Support

The current Ollama integration uses basic `"format": "json"` mode which asks the model to return JSON but doesn't enforce a specific structure. Ollama now supports **schema-constrained structured output** where you pass a full JSON schema to the `format` parameter, guaranteeing the exact output format.

**Current implementation** (`backend/src/tact/llm/ollama.py`):
```python
"format": "json",
```

**Enhanced implementation** would pass the expected schema:
```python
"format": {
    "type": "object",
    "properties": {
        "duration_minutes": {"type": ["integer", "null"]},
        "time_code_id": {"type": ["string", "null"]},
        "work_type_id": {"type": ["string", "null"]},
        "parsed_description": {"type": ["string", "null"]},
        "confidence_duration": {"type": "number"},
        "confidence_time_code": {"type": "number"},
        "confidence_work_type": {"type": "number"},
        "confidence_overall": {"type": "number"},
        "notes": {"type": "string"}
    },
    "required": ["duration_minutes", "time_code_id", "confidence_duration", "confidence_time_code"]
}
```

This would improve parsing reliability across all Ollama models by ensuring the model always returns the expected fields, rather than relying on prompt instructions alone.

**References:**
- [Ollama Structured Outputs Blog](https://ollama.com/blog/structured-outputs)
- [Ollama Structured Outputs Docs](https://docs.ollama.com/capabilities/structured-outputs)

