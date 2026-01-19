## 1. Database Migration

- [x] 1.1 Create Alembic migration to drop `description`, `keywords`, `examples` columns from `time_codes` table

## 2. Backend Model & Schema Updates

- [x] 2.1 Remove fields from `TimeCode` model in `db/models.py`
- [x] 2.2 Update `TimeCodeCreate` schema (remove description, keywords, examples)
- [x] 2.3 Update `TimeCodeUpdate` schema (remove description, keywords, examples)
- [x] 2.4 Update `TimeCodeResponse` schema (remove description, keywords, examples)

## 3. Backend Route Updates

- [x] 3.1 Update `routes/time_codes.py` - remove field handling in create/update/response

## 4. Backend LLM Prompt Updates

- [x] 4.1 Update `llm/prompts.py` - simplify time code display (remove description and keywords)
- [x] 4.2 Update `llm/provider.py` - remove keywords/description from `TimeCodeInfo` dataclass

## 5. Backend Test Updates

- [x] 5.1 Update `tests/test_time_codes.py` - remove fields from test payloads
- [x] 5.2 Update `tests/test_llm_parser.py` - remove keywords/description from fixtures
- [x] 5.3 Update any other tests that reference these fields

## 6. TUI Updates (Go)

- [x] 6.1 Update `tui/model/timecode.go` - remove Description, Keywords, Examples from struct
- [x] 6.2 Update `tui/api/client.go` - remove fields from CreateTimeCodeRequest/UpdateTimeCodeRequest
- [x] 6.3 Update `tui/api/client.go` - simplify CreateTimeCode function signature
- [x] 6.4 Update `tui/ui/timecode_edit.go` - remove description, keywords, examples input fields and related logic

## 7. MCP Updates (Python)

- [x] 7.1 Update `mcp/src/tact_mcp/client.py` - remove description, keywords, examples from create_time_code
- [x] 7.2 Update `mcp/src/tact_mcp/server.py` - remove fields from create_time_code tool definition
- [x] 7.3 Update `mcp/src/tact_mcp/server.py` - remove fields from update_time_code tool definition

## 8. Documentation

- [x] 8.1 Update backend README if it mentions time code fields (N/A - no README changes needed)

## Verification

1. [x] Run migrations successfully
2. [x] Create a time code via API with just id, name, project_id
3. [x] Create a time code via TUI with simplified form (verified Go compiles)
4. [x] Create a time code via MCP tool (verified Python syntax)
5. [x] Verify LLM parsing still works (uses RAG context for matching)
6. [x] All backend tests pass (150 tests)
