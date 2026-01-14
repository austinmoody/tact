# Tasks

## 1. Add Client Methods

- [x] 1.1 Add Project methods to `client.py` (list, get, create, update, delete)
- [x] 1.2 Add Project Context methods (`list_project_context`, `create_project_context`)
- [x] 1.3 Add Time Code Context methods (`list_time_code_context`, `create_time_code_context`)
- [x] 1.4 Add generic Context methods (`get_context`, `update_context`, `delete_context`)
- [x] 1.5 Update `list_time_codes` client method to support `project_id` filter
- [x] 1.6 Update `create_time_code` client method to support `project_id` parameter
- [x] 1.7 Update `update_time_code` client method to support `project_id` parameter

## 2. Add Project Tools

- [x] 2.1 Add `list_projects` tool definition and handler
- [x] 2.2 Add `get_project` tool definition and handler
- [x] 2.3 Add `create_project` tool definition and handler
- [x] 2.4 Add `update_project` tool definition and handler
- [x] 2.5 Add `delete_project` tool definition and handler

## 3. Add Context Tools

- [x] 3.1 Add `list_project_context` tool definition and handler
- [x] 3.2 Add `list_time_code_context` tool definition and handler
- [x] 3.3 Add `create_project_context` tool definition and handler
- [x] 3.4 Add `create_time_code_context` tool definition and handler
- [x] 3.5 Add `get_context` tool definition and handler
- [x] 3.6 Add `update_context` tool definition and handler
- [x] 3.7 Add `delete_context` tool definition and handler

## 4. Update Existing Time Code Tools

- [x] 4.1 Add `project_id` filter to `list_time_codes` tool
- [x] 4.2 Add `project_id` parameter to `create_time_code` tool (default: "default")
- [x] 4.3 Add `project_id` parameter to `update_time_code` tool

## 5. Testing

- [x] 5.1 Manual test with Claude Desktop or MCP inspector
- [x] 5.2 Verify all tools appear in tool list
- [x] 5.3 Test CRUD operations for projects and context
- [x] 5.4 Test time code project_id filtering and assignment

## Dependencies

- Task 2 depends on Task 1.1 (need client methods for project tools)
- Task 3 depends on Tasks 1.2-1.4 (need client methods for context tools)
- Task 4 depends on Tasks 1.5-1.7 (need updated client methods)

## Notes

- Follow existing patterns in server.py for tool definitions
- All tools return JSON responses via `json_response()` helper
- Error handling via `error_response()` helper
