## MODIFIED Requirements

### Requirement: Time Code Tools

The `create_time_code` tool SHALL accept ID, name, and project_id (defaults to "default") and create a new time code.

The `update_time_code` tool SHALL accept a time code ID and update fields including name, active, and project_id.

#### Scenario: Create time code

- WHEN: user invokes `create_time_code` with id="PROJ-BETA", name="Project Beta"
- THEN: a new time code is created and returned

#### Scenario: Create time code in project

- WHEN: user invokes `create_time_code` with id="IZG-001", name="IZG Task", project_id="izg"
- THEN: a new time code is created under the specified project

#### Scenario: Update time code name

- WHEN: user invokes `update_time_code` with time_code_id="IZG-001" and name="Updated Name"
- THEN: the time code name is updated

## REMOVED Requirements

### Requirement: Time Code Examples Update

~~The `update_time_code` tool accepted examples parameter.~~

Removed: Examples field removed from time code model. Use context documents instead.
