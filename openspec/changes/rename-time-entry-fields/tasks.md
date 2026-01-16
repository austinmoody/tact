# Tasks: Rename Time Entry Text Fields

## Implementation Tasks

### Backend Database & Models
- [ ] Create Alembic migration to rename `raw_text` → `user_input` and `description` → `parsed_description` in `time_entries` table
- [ ] Update `TimeEntry` model in `backend/src/tact/db/models.py`
- [ ] Update schemas in `backend/src/tact/schemas/entry.py`

### Backend Routes & LLM
- [ ] Update `backend/src/tact/routes/entries.py`
- [ ] Update `backend/src/tact/llm/provider.py` (ParseResult class)
- [ ] Update `backend/src/tact/llm/prompts.py`
- [ ] Update `backend/src/tact/llm/parser.py`
- [ ] Update `backend/src/tact/llm/ollama.py`
- [ ] Update `backend/src/tact/llm/anthropic.py`

### Backend Tests
- [ ] Update `backend/tests/test_entries.py`
- [ ] Update `backend/tests/test_llm_parser.py`
- [ ] Update `backend/tests/test_db.py`
- [ ] Update `backend/tests/test_reparse.py`

### TUI
- [ ] Update `tui/model/entry.go`: rename `RawText` → `UserInput` (json tag: `user_input`) and `Description` → `ParsedDescription` (json tag: `parsed_description`)
- [ ] Update any TUI code that references these fields (search for `RawText` and `Description` usage)

### MCP
- [ ] Update `mcp/src/tact_mcp/server.py`: rename `raw_text` → `user_input` in create_entry tool schema and arguments
- [ ] Update `mcp/src/tact_mcp/server.py`: rename `description` → `parsed_description` in update_entry tool schema
- [ ] Update `mcp/src/tact_mcp/client.py`: rename `raw_text` → `user_input` parameter and JSON field in `create_entry` method

### Specs
- [ ] Update `openspec/specs/entries-api/spec.md`

## Verification

1. Run backend tests: `cd backend && pytest`
2. Run TUI build: `cd tui && go build ./...`
3. Run MCP tests: `cd mcp && pytest` (if tests exist)
4. Start all services and verify:
   - Create a new time entry via API
   - Verify parsing populates `parsed_description`
   - Verify `user_input` is preserved
   - Verify TUI displays entries correctly
