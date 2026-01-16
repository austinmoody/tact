# Tasks

## 1. Backend

- [ ] 1.1 Remove `description` field from `Project` model in `models.py`
- [ ] 1.2 Remove `description` from `ProjectCreate`, `ProjectUpdate`, `ProjectResponse` schemas
- [ ] 1.3 Remove `description` handling from `routes/projects.py` (create and update)
- [ ] 1.4 Create Alembic migration to drop `description` column from `projects` table
- [ ] 1.5 Update tests in `test_projects.py` to remove description references
- [ ] 1.6 Update tests in `test_db.py`, `test_time_codes.py`, `test_llm_parser.py` to remove description from Project fixtures

## 2. TUI

- [ ] 2.1 Remove `Description` field from `Project` struct in `model/project.go`
- [ ] 2.2 Remove `Description` from project request/update structs in `api/client.go`
- [ ] 2.3 Remove `description` parameter from `CreateProject` function in `api/client.go`
- [ ] 2.4 Remove description input field and related logic from `ui/project_edit.go`

## 3. MCP

- [ ] 3.1 Remove `description` parameter from `create_project` method in `client.py`
- [ ] 3.2 Remove `description` from `create_project` tool schema in `server.py`
- [ ] 3.3 Remove `description` from `update_project` tool schema in `server.py`

## Dependencies

- Task 1.4 (migration) must run after 1.1 (model change)
- TUI and MCP changes (2.x, 3.x) can be done in parallel with backend changes
- All components should be updated together before deployment

## Verification

1. Run backend tests: `cd backend && make test`
2. Run TUI build: `cd tui && go build ./...`
3. Run MCP tests: `cd mcp && uv run pytest`
4. Start the stack and verify:
   - `curl -X POST http://localhost:2100/projects -H "Content-Type: application/json" -d '{"id": "test", "name": "Test Project"}'` succeeds
   - Response does not contain `description` field
   - TUI can create/edit projects without description field
