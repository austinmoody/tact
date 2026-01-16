# Tasks: Add Parse Notes Field

## Implementation Tasks

### Backend Database & Models
- [ ] Create Alembic migration to add `parse_notes` column to `time_entries` table
- [ ] Update `TimeEntry` model in `backend/src/tact/db/models.py`
- [ ] Update `EntryResponse` schema in `backend/src/tact/schemas/entry.py`

### Backend LLM Integration
- [ ] Update `ParseResult` in `backend/src/tact/llm/provider.py` to include `notes` field
- [ ] Update LLM prompt in `backend/src/tact/llm/prompts.py` to request reasoning explanation
- [ ] Update `OllamaProvider` in `backend/src/tact/llm/ollama.py` to parse notes from response
- [ ] Update `AnthropicProvider` in `backend/src/tact/llm/anthropic.py` to parse notes from response
- [ ] Update `EntryParser` in `backend/src/tact/llm/parser.py` to:
  - Store LLM notes in entry
  - Append closest RAG context info to notes

### Backend Tests
- [ ] Update `backend/tests/test_llm_parser.py` with parse_notes tests
- [ ] Update any other affected tests

### TUI
- [ ] Update `tui/model/entry.go`: add `ParseNotes *string` field with json tag
- [ ] Update `tui/ui/entry_detail.go`: display `parse_notes` after overall confidence section
- [ ] Consider using different styling for `needs_review` entries to highlight the notes

### Specs
- [ ] Update `openspec/specs/llm-parsing/spec.md`
- [ ] Update `openspec/specs/entries-api/spec.md`

## Verification

1. Run backend tests: `cd backend && pytest`
2. Run TUI build: `cd tui && go build ./...`
3. Manual testing:
   - Create entry that parses successfully → verify parse_notes explains the match
   - Create ambiguous entry → verify parse_notes explains why it's `needs_review`
   - Verify TUI displays parse_notes in entry detail
