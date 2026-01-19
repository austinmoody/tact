## 1. Database Migration

- [ ] 1.1 Create Alembic migration to drop `description`, `keywords`, `examples` columns from `time_codes` table

## 2. Backend Model & Schema Updates

- [ ] 2.1 Remove fields from `TimeCode` model in `db/models.py`
- [ ] 2.2 Update `TimeCodeCreate` schema (remove description, keywords, examples)
- [ ] 2.3 Update `TimeCodeUpdate` schema (remove description, keywords, examples)
- [ ] 2.4 Update `TimeCodeResponse` schema (remove description, keywords, examples)

## 3. Backend Route Updates

- [ ] 3.1 Update `routes/time_codes.py` - remove field handling in create/update/response

## 4. Backend LLM Prompt Updates

- [ ] 4.1 Update `llm/prompts.py` - simplify time code display (remove description and keywords)
- [ ] 4.2 Update `llm/provider.py` - remove keywords/description from `TimeCodeInfo` dataclass

## 5. Backend Test Updates

- [ ] 5.1 Update `tests/test_time_codes.py` - remove fields from test payloads
- [ ] 5.2 Update `tests/test_llm_parser.py` - remove keywords/description from fixtures
- [ ] 5.3 Update any other tests that reference these fields

## 6. TUI Updates (Go)

- [ ] 6.1 Update `tui/model/timecode.go` - remove Description, Keywords, Examples from struct
- [ ] 6.2 Update `tui/api/client.go` - remove fields from CreateTimeCodeRequest/UpdateTimeCodeRequest
- [ ] 6.3 Update `tui/api/client.go` - simplify CreateTimeCode function signature
- [ ] 6.4 Update `tui/ui/timecode_edit.go` - remove description, keywords, examples input fields and related logic

## 7. MCP Updates (Python)

- [ ] 7.1 Update `mcp/src/tact_mcp/client.py` - remove description, keywords, examples from create_time_code
- [ ] 7.2 Update `mcp/src/tact_mcp/server.py` - remove fields from create_time_code tool definition
- [ ] 7.3 Update `mcp/src/tact_mcp/server.py` - remove fields from update_time_code tool definition

## 8. Documentation

- [ ] 8.1 Update backend README if it mentions time code fields

## Verification

1. Run migrations successfully
2. Create a time code via API with just id, name, project_id
3. Create a time code via TUI with simplified form
4. Create a time code via MCP tool
5. Verify LLM parsing still works (uses RAG context for matching)
6. All backend tests pass
