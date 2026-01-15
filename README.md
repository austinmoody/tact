# TACT

**T**racking **A**nd **C**alculating **T**ime

*Set out runnin', but I track my time.*

## Components

- **backend/** - FastAPI REST API (Python)
- **tui/** - Terminal UI dashboard (Go/Bubbletea)
- **mcp/** - MCP server for AI assistants (Python)

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
