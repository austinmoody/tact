## MODIFIED Requirements

### Requirement: Create Time Code

The API SHALL allow creating a new time code via POST request.

#### Scenario: Successful creation with project

- Given: A valid time code payload with id, name, and project_id
- And: The referenced project exists
- When: POST `/time-codes` is called
- Then: The time code is created with `active=true`
- And: The time code is associated with the specified project
- And: HTTP 201 is returned with the created time code

#### Scenario: Project not found

- Given: A time code payload references a non-existent project_id
- When: POST `/time-codes` is called
- Then: HTTP 400 Bad Request is returned
- And: Error message indicates project not found

## NEW Requirements

### Requirement: Time Code Fields

The TimeCode model SHALL have the following core fields: id, project_id, name, active, created_at, updated_at.

#### Scenario: Minimal field set

- Given: A time code is created or retrieved
- Then: The time code has fields: id, project_id, name, active, created_at, updated_at
- And: No description, keywords, or examples fields exist
