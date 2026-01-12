## Context

Tact needs an MCP server to expose its API to AI assistants. The server must:
- Work with multiple MCP clients (Claude Desktop, Goose, GitHub Copilot, Cursor)
- Run in Docker alongside the existing API container
- Expose full CRUD operations for entries, time codes, and work types

### Constraints
- No JavaScript/TypeScript/Node.js (project constraint)
- Single user focus (no auth needed)
- Local-first deployment via Docker

## Goals / Non-Goals

**Goals:**
- Expose full Tact API as MCP tools
- Docker deployment with clear client configuration docs
- stdio transport for universal client compatibility

**Non-Goals:**
- HTTP/SSE transport (overkill for personal use)
- Multi-user authentication
- Rate limiting or quota management

## Decisions

### Language: Python

**Decision:** Use Python with the official Anthropic MCP SDK (`mcp` package).

**Rationale:**
- Backend is already Python (FastAPI) - consistent tech stack
- Official Anthropic SDK is well-maintained and up-to-date
- Can share Pydantic models with backend if needed
- Fits project constraint: no JavaScript/TypeScript/Node.js

**Alternatives considered:**
- Go: Community SDK only (not official), less mature
- Java: Heavy runtime, no clear advantage

### Transport: stdio

**Decision:** Use stdio (stdin/stdout) transport exclusively.

**Rationale:**
- Supported by all major MCP clients
- Simpler than HTTP - no networking, ports, or auth needed
- Docker-friendly with `-i` flag
- Sufficient for single-user personal use

**Alternatives considered:**
- HTTP/SSE: More complex, requires auth, overkill for personal use

### Project Structure

```
mcp/
├── pyproject.toml          # Dependencies and build config
├── Dockerfile              # Container build
├── README.md               # Client configuration docs
└── src/
    └── tact_mcp/
        ├── __init__.py
        ├── server.py       # MCP server setup and tool registration
        ├── tools.py        # Tool implementations
        └── client.py       # HTTP client for Tact API
```

### Tool Design

Each API endpoint maps to an MCP tool:

| Category | Tools |
|----------|-------|
| Entries | `create_entry`, `list_entries`, `get_entry`, `update_entry`, `delete_entry`, `reparse_entry` |
| Time Codes | `list_time_codes`, `get_time_code`, `create_time_code`, `update_time_code`, `delete_time_code` |
| Work Types | `list_work_types`, `get_work_type`, `create_work_type`, `update_work_type`, `delete_work_type` |
| Reports | `get_summary` |

### Docker Integration

The MCP server runs as a separate container that connects to the API container:

```yaml
mcp:
  build:
    context: ./mcp
    dockerfile: Dockerfile
  environment:
    - TACT_API_URL=http://api:2100
  stdin_open: true
  depends_on:
    - api
```

Clients spawn the container directly:
```bash
docker run -i --rm --network tact_default -e TACT_API_URL=http://api:2100 tact-mcp
```

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| API changes break MCP server | Version API, update MCP client accordingly |
| Docker networking complexity | Document network setup clearly, test on fresh install |
| MCP SDK breaking changes | Pin SDK version, test before upgrading |
| Container startup latency | Accept for now; could add warm-start optimization later |

## Migration Plan

Not applicable - new capability with no existing state.

## Open Questions

None - requirements are clear.
