## MODIFIED Requirements

### Requirement: Time Code Tools

The system SHALL provide MCP tools for managing time codes.

The `list_time_codes` tool SHALL return all time codes, optionally filtered by active status or project_id.

The `get_time_code` tool SHALL accept a time code ID and return full details.

The `create_time_code` tool SHALL accept ID, name, and project_id (defaults to "default") and create a new time code.

The `update_time_code` tool SHALL accept a time code ID and update fields including name, active, and project_id.

The `delete_time_code` tool SHALL accept a time code ID and deactivate it (soft delete).

#### Scenario: List active time codes
- **WHEN** user invokes `list_time_codes` with active_only=true
- **THEN** only active time codes are returned

#### Scenario: List time codes by project
- **WHEN** user invokes `list_time_codes` with project_id="izg"
- **THEN** only time codes belonging to that project are returned

#### Scenario: Create time code
- **WHEN** user invokes `create_time_code` with id="PROJ-BETA", name="Project Beta"
- **THEN** a new time code is created and returned

#### Scenario: Create time code in project
- **WHEN** user invokes `create_time_code` with id="IZG-001", name="IZG Task", project_id="izg"
- **THEN** a new time code is created under the specified project

#### Scenario: Update time code name
- **WHEN** user invokes `update_time_code` with time_code_id="IZG-001" and name="Updated Name"
- **THEN** the time code name is updated

#### Scenario: Move time code to different project
- **WHEN** user invokes `update_time_code` with time_code_id="IZG-001" and project_id="other-project"
- **THEN** the time code is reassigned to the new project

#### Scenario: Deactivate time code
- **WHEN** user invokes `delete_time_code` with a valid time code ID
- **THEN** the time code is marked inactive
