# Tact MCP Server

MCP (Model Context Protocol) server that exposes the Tact time tracking API to AI clients like Claude Desktop, Goose, and GitHub Copilot.

## Features

- Full CRUD operations for time entries
- Time code management
- Work type management
- Report generation
- Natural language time entry creation

## Available Tools

### Entries
- `create_entry` - Create a new time entry using natural language
- `list_entries` - List entries with optional filters
- `get_entry` - Get a single entry by ID
- `update_entry` - Update an entry
- `delete_entry` - Delete an entry
- `reparse_entry` - Trigger re-parsing

### Time Codes
- `list_time_codes` - List all time codes
- `get_time_code` - Get a single time code
- `create_time_code` - Create a new time code
- `update_time_code` - Update a time code
- `delete_time_code` - Deactivate a time code

### Work Types
- `list_work_types` - List all work types
- `get_work_type` - Get a single work type
- `create_work_type` - Create a new work type
- `update_work_type` - Update a work type
- `delete_work_type` - Deactivate a work type

### Reports
- `get_summary` - Get time summary report

## Installation

### Prerequisites

1. The Tact API must be running:
   ```bash
   docker compose up -d api
   ```

2. Build the MCP server image:
   ```bash
   docker compose build mcp
   ```

## Client Configuration

### Claude Desktop

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

After saving, restart Claude Desktop. The Tact tools will appear in the tools menu.

### Goose

Add to `.goose/config.yaml`:

```yaml
mcp:
  tact:
    command: docker
    args:
      - run
      - -i
      - --rm
      - --network=tact_default
      - -e
      - TACT_API_URL=http://api:2100
      - tact-mcp
```

### VS Code / Cursor

Add to your MCP configuration:

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

## Development

### Running Locally (without Docker)

1. Install dependencies:
   ```bash
   cd mcp
   uv sync
   ```

2. Set the API URL:
   ```bash
   export TACT_API_URL=http://localhost:2100
   ```

3. Run the server:
   ```bash
   uv run tact-mcp
   ```

### Local Development with Claude Desktop

For local development, update the Claude Desktop config:

```json
{
  "mcpServers": {
    "tact": {
      "command": "uv",
      "args": ["run", "--directory", "/path/to/tact/mcp", "tact-mcp"],
      "env": {
        "TACT_API_URL": "http://localhost:2100"
      }
    }
  }
}
```

## Usage Examples

Once connected, you can ask your AI assistant:

- "Log 2 hours of development work on Project Alpha"
- "Show me my time entries for this week"
- "List all active time codes"
- "Create a new work type called 'Code Review'"
- "What's my time summary for this month?"

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `TACT_API_URL` | `http://localhost:2100` | URL of the Tact API |

## Troubleshooting

### "Connection refused" errors

Make sure the Tact API is running and accessible from the MCP container:

```bash
# Check API is running
docker compose ps

# Test connectivity
docker run --rm --network tact_default curlimages/curl http://api:2100/health
```

### Tools not appearing in Claude Desktop

1. Check Claude Desktop logs for errors
2. Verify the config file syntax is valid JSON
3. Restart Claude Desktop after config changes
4. Ensure Docker is running

### Network issues

The MCP container needs to be on the same Docker network as the API:

```bash
# Check network exists
docker network ls | grep tact

# If not, start the API first
docker compose up -d api
```
