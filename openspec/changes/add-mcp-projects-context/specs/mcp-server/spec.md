# mcp-server Spec Delta

## MODIFIED Requirements

### Requirement: Time Code Tools - Project Support

The `list_time_codes` tool SHALL accept an optional project_id parameter to filter time codes by project.

The `create_time_code` tool SHALL accept a project_id parameter (defaults to "default").

The `update_time_code` tool SHALL accept an optional project_id parameter to reassign a time code to a different project.

#### Scenario: List time codes by project

- **WHEN** user invokes `list_time_codes` with project_id="izg"
- **THEN** only time codes belonging to that project are returned

#### Scenario: Create time code in project

- **WHEN** user invokes `create_time_code` with id="IZG-001", name="IZG Task", project_id="izg"
- **THEN** a new time code is created under the specified project

#### Scenario: Move time code to different project

- **WHEN** user invokes `update_time_code` with time_code_id="IZG-001" and project_id="other-project"
- **THEN** the time code is reassigned to the new project

## ADDED Requirements

### Requirement: Project Tools

The system SHALL provide MCP tools for managing projects.

The `list_projects` tool SHALL return all projects, optionally filtered by active status.

The `get_project` tool SHALL accept a project ID and return full details.

The `create_project` tool SHALL accept ID, name, and optional description and create a new project.

The `update_project` tool SHALL accept a project ID and update fields.

The `delete_project` tool SHALL accept a project ID and deactivate it (soft delete).

#### Scenario: List active projects

- **WHEN** user invokes `list_projects` with active_only=true
- **THEN** only active projects are returned

#### Scenario: Create project

- **WHEN** user invokes `create_project` with id="izg", name="IZG Hub"
- **THEN** a new project is created and returned

#### Scenario: Update project

- **WHEN** user invokes `update_project` with id and description="Updated description"
- **THEN** the project description is updated

#### Scenario: Deactivate project

- **WHEN** user invokes `delete_project` with a valid project ID
- **THEN** the project is marked inactive

### Requirement: Context Tools

The system SHALL provide MCP tools for managing context documents.

The `list_project_context` tool SHALL accept a project ID and return all context documents for that project.

The `list_time_code_context` tool SHALL accept a time code ID and return all context documents for that time code.

The `create_project_context` tool SHALL accept a project ID and content string and create a new context document.

The `create_time_code_context` tool SHALL accept a time code ID and content string and create a new context document.

The `get_context` tool SHALL accept a context ID and return full details.

The `update_context` tool SHALL accept a context ID and new content and update the document.

The `delete_context` tool SHALL accept a context ID and delete the document.

#### Scenario: List project context

- **WHEN** user invokes `list_project_context` with project_id="izg"
- **THEN** all context documents for that project are returned

#### Scenario: List time code context

- **WHEN** user invokes `list_time_code_context` with time_code_id="FEDS-165"
- **THEN** all context documents for that time code are returned

#### Scenario: Create project context

- **WHEN** user invokes `create_project_context` with project_id="izg" and content="IZG Hub is the integration gateway..."
- **THEN** a new context document is created with embedding
- **AND** the context document is returned

#### Scenario: Create time code context

- **WHEN** user invokes `create_time_code_context` with time_code_id="FEDS-165" and content="Use for Help Desk activities..."
- **THEN** a new context document is created with embedding
- **AND** the context document is returned

#### Scenario: Update context

- **WHEN** user invokes `update_context` with context_id and new content
- **THEN** the context document content is updated
- **AND** the embedding is regenerated

#### Scenario: Delete context

- **WHEN** user invokes `delete_context` with a valid context ID
- **THEN** the context document is deleted
