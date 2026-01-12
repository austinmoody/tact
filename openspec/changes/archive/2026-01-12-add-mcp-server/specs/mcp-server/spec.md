## ADDED Requirements

### Requirement: MCP Server Core

The system SHALL provide an MCP server that exposes the Tact API to MCP-compatible clients.

The server SHALL use stdio transport (stdin/stdout) for communication.

The server SHALL connect to the Tact API using the `TACT_API_URL` environment variable.

#### Scenario: Server initialization
- **WHEN** an MCP client connects to the server
- **THEN** the server responds with its capabilities and available tools

#### Scenario: API connection configured via environment
- **WHEN** the server starts with `TACT_API_URL=http://api:2100`
- **THEN** all tool calls are routed to that API URL

### Requirement: Entry Tools

The system SHALL provide MCP tools for managing time entries.

The `create_entry` tool SHALL accept a natural language text parameter and create a new entry.

The `list_entries` tool SHALL accept optional filter parameters (status, time_code_id, work_type_id, from_date, to_date, limit, offset) and return matching entries.

The `get_entry` tool SHALL accept an entry ID and return the full entry details.

The `update_entry` tool SHALL accept an entry ID and update fields and modify the entry.

The `delete_entry` tool SHALL accept an entry ID and delete the entry.

The `reparse_entry` tool SHALL accept an entry ID and trigger re-parsing of the entry.

#### Scenario: Create entry via natural language
- **WHEN** user invokes `create_entry` with text "2 hours development on Project Alpha"
- **THEN** a new entry is created with the raw text
- **AND** the entry ID is returned

#### Scenario: List entries with filters
- **WHEN** user invokes `list_entries` with status="parsed" and limit=10
- **THEN** up to 10 parsed entries are returned

#### Scenario: Get entry details
- **WHEN** user invokes `get_entry` with a valid entry ID
- **THEN** the full entry details are returned including parsed fields and confidence scores

#### Scenario: Update entry
- **WHEN** user invokes `update_entry` with entry ID and new time_code_id
- **THEN** the entry is updated with the new time code

#### Scenario: Delete entry
- **WHEN** user invokes `delete_entry` with a valid entry ID
- **THEN** the entry is deleted

#### Scenario: Reparse entry
- **WHEN** user invokes `reparse_entry` with a valid entry ID
- **THEN** the entry is queued for re-parsing

### Requirement: Time Code Tools

The system SHALL provide MCP tools for managing time codes.

The `list_time_codes` tool SHALL return all time codes, optionally filtered by active status.

The `get_time_code` tool SHALL accept a time code ID and return full details.

The `create_time_code` tool SHALL accept ID, name, and optional description/keywords/examples and create a new time code.

The `update_time_code` tool SHALL accept a time code ID and update fields.

The `delete_time_code` tool SHALL accept a time code ID and deactivate it (soft delete).

#### Scenario: List active time codes
- **WHEN** user invokes `list_time_codes` with active_only=true
- **THEN** only active time codes are returned

#### Scenario: Create time code
- **WHEN** user invokes `create_time_code` with id="PROJ-BETA", name="Project Beta"
- **THEN** a new time code is created and returned

#### Scenario: Update time code with examples
- **WHEN** user invokes `update_time_code` with examples=["2h coding", "bug fix"]
- **THEN** the time code examples are updated

#### Scenario: Deactivate time code
- **WHEN** user invokes `delete_time_code` with a valid time code ID
- **THEN** the time code is marked inactive

### Requirement: Work Type Tools

The system SHALL provide MCP tools for managing work types.

The `list_work_types` tool SHALL return all work types.

The `get_work_type` tool SHALL accept a work type ID and return full details.

The `create_work_type` tool SHALL accept a name and create a new work type.

The `update_work_type` tool SHALL accept a work type ID and update fields.

The `delete_work_type` tool SHALL accept a work type ID and deactivate it (soft delete).

#### Scenario: List work types
- **WHEN** user invokes `list_work_types`
- **THEN** all work types are returned

#### Scenario: Create work type
- **WHEN** user invokes `create_work_type` with name="Code Review"
- **THEN** a new work type is created and returned

#### Scenario: Update work type
- **WHEN** user invokes `update_work_type` with id and name="Peer Review"
- **THEN** the work type name is updated

#### Scenario: Deactivate work type
- **WHEN** user invokes `delete_work_type` with a valid work type ID
- **THEN** the work type is marked inactive

### Requirement: Report Tools

The system SHALL provide MCP tools for generating reports.

The `get_summary` tool SHALL accept optional parameters (time_code_id, work_type_id, from_date, to_date) and return aggregated time data.

#### Scenario: Get weekly summary
- **WHEN** user invokes `get_summary` with from_date and to_date for the current week
- **THEN** aggregated time by time code and work type is returned

### Requirement: Docker Deployment

The MCP server SHALL be deployable as a Docker container.

The container SHALL support stdio transport via `docker run -i`.

The container SHALL connect to the Tact API via Docker networking.

#### Scenario: Run MCP server in Docker
- **WHEN** user runs `docker run -i --network tact_default -e TACT_API_URL=http://api:2100 tact-mcp`
- **THEN** the MCP server starts and accepts connections via stdin/stdout

#### Scenario: Claude Desktop configuration
- **WHEN** user configures Claude Desktop with the Docker command
- **THEN** Claude Desktop can spawn the MCP server and invoke tools
