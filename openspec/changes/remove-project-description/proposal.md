# Change: Remove description field from Projects

## Why

The `description` field on Projects is redundant. Projects only need an `id` and `name` - there's no need for a more verbose description. Removing it simplifies the data model and reduces unnecessary fields in the API, TUI, and MCP.

## What Changes

- **BREAKING**: Remove `description` column from `projects` table
- Remove `description` from Project API schemas (create, update, response)
- Remove description input from TUI project edit screen
- Remove description parameter from MCP project tools

## Impact

- Affected specs: `projects-api`
- Affected code:
  - Backend: `models.py`, `schemas/project.py`, `routes/projects.py`, tests
  - TUI: `model/project.go`, `api/client.go`, `ui/project_edit.go`
  - MCP: `client.py`, `server.py`
- Database migration required to drop column
- Any existing description data will be lost (acceptable per user)
