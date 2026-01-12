## 1. Project Setup
- [x] 1.1 Create `mcp/` directory structure
- [x] 1.2 Create `pyproject.toml` with dependencies (`mcp`, `httpx`, `pydantic`)
- [x] 1.3 Create basic `src/tact_mcp/__init__.py`

## 2. API Client
- [x] 2.1 Create `client.py` with Tact API client class
- [x] 2.2 Implement entry endpoints (create, list, get, update, delete, reparse)
- [x] 2.3 Implement time code endpoints (list, get, create, update, delete)
- [x] 2.4 Implement work type endpoints (list, get, create, update, delete)
- [x] 2.5 Implement reports endpoint (get_summary)

## 3. MCP Server
- [x] 3.1 Create `server.py` with MCP server skeleton
- [x] 3.2 Tools implemented inline in server.py (no separate tools.py needed)
- [x] 3.3 Register entry tools (create_entry, list_entries, get_entry, update_entry, delete_entry, reparse_entry)
- [x] 3.4 Register time code tools (list_time_codes, get_time_code, create_time_code, update_time_code, delete_time_code)
- [x] 3.5 Register work type tools (list_work_types, get_work_type, create_work_type, update_work_type, delete_work_type)
- [x] 3.6 Register report tool (get_summary)

## 4. Docker & Deployment
- [x] 4.1 Create `Dockerfile` for MCP server
- [x] 4.2 Add MCP service to `docker-compose.yml` with `profiles: [tools]` (builds with `docker compose build`, but doesn't start with `docker compose up`)
- [x] 4.3 Test Docker build and container startup

## 5. Documentation
- [x] 5.1 Create `mcp/README.md` with overview
- [x] 5.2 Document Claude Desktop configuration
- [x] 5.3 Document Goose configuration
- [x] 5.4 Document development/local run instructions

## Verification

Steps to verify the implementation works:

1. Build the Docker image:
   ```bash
   docker compose build mcp
   ```

2. Start the Tact API (if not running):
   ```bash
   docker compose up -d api
   ```

3. Test MCP server starts without errors:
   ```bash
   echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}' | docker run -i --rm --network tact_default -e TACT_API_URL=http://api:2100 tact-mcp
   ```
   Expected: JSON response with server capabilities

4. Configure Claude Desktop and verify tools appear:
   - Add configuration to `~/Library/Application Support/Claude/claude_desktop_config.json`
   - Restart Claude Desktop
   - Check that Tact tools appear in the tools list

5. Test creating an entry via Claude Desktop:
   - Ask Claude: "Log 2 hours of development work on Project Alpha"
   - Verify entry appears in TUI or API
