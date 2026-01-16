# Tasks: Rename Time Entry Text Fields

## Implementation Tasks

### Backend Database & Models
- [x] Create Alembic migration to rename `raw_text` → `user_input` and `description` → `parsed_description` in `time_entries` table
- [x] Update `TimeEntry` model in `backend/src/tact/db/models.py`
- [x] Update schemas in `backend/src/tact/schemas/entry.py`

### Backend Routes & LLM
- [x] Update `backend/src/tact/routes/entries.py`
- [x] Update `backend/src/tact/llm/provider.py` (ParseResult class)
- [x] Update `backend/src/tact/llm/prompts.py`
- [x] Update `backend/src/tact/llm/parser.py`
- [x] Update `backend/src/tact/llm/ollama.py`
- [x] Update `backend/src/tact/llm/anthropic.py`

### Backend Tests
- [x] Update `backend/tests/test_entries.py`
- [x] Update `backend/tests/test_llm_parser.py`
- [x] Update `backend/tests/test_db.py`
- [x] Update `backend/tests/test_reparse.py`

### TUI
- [x] Update `tui/model/entry.go`: rename `RawText` → `UserInput` (json tag: `user_input`) and `Description` → `ParsedDescription` (json tag: `parsed_description`)
- [x] Update any TUI code that references these fields (search for `RawText` and `Description` usage)

### MCP
- [x] Update `mcp/src/tact_mcp/server.py`: rename `raw_text` → `user_input` in create_entry tool schema and arguments
- [x] Update `mcp/src/tact_mcp/server.py`: rename `description` → `parsed_description` in update_entry tool schema
- [x] Update `mcp/src/tact_mcp/client.py`: rename `raw_text` → `user_input` parameter and JSON field in `create_entry` method

### Specs
- [x] Update `openspec/specs/entries-api/spec.md`

## Verification

1. Run backend tests: `cd backend && pytest`
2. Run TUI build: `cd tui && go build ./...`
3. Run MCP tests: `cd mcp && pytest` (if tests exist)
4. Start all services and verify:
   - Create a new time entry via API
   - Verify parsing populates `parsed_description`
   - Verify `user_input` is preserved
   - Verify TUI displays entries correctly
