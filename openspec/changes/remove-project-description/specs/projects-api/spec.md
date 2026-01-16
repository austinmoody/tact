## MODIFIED Requirements

### Requirement: Create Project

The API SHALL allow creating a new project via POST request. A project consists of an `id` and `name` only.

#### Scenario: Successful creation

- Given: A valid project payload with id and name
- When: POST `/projects` is called
- Then: The project is created with `active=true`
- And: HTTP 201 is returned with the created project

#### Scenario: Duplicate ID

- Given: A project with the same ID already exists
- When: POST `/projects` is called
- Then: HTTP 409 Conflict is returned

### Requirement: Update Project

The API SHALL allow updating an existing project. Only `name` and `active` fields can be updated.

#### Scenario: Successful update

- Given: A project with ID "IZG" exists
- When: PUT `/projects/IZG` is called with updated fields
- Then: The project is updated
- And: HTTP 200 is returned with the updated project

#### Scenario: Update non-existent project

- Given: No project with ID "UNKNOWN" exists
- When: PUT `/projects/UNKNOWN` is called
- Then: HTTP 404 is returned
