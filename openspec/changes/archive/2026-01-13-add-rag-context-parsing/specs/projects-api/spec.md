# projects-api Specification Delta

## ADDED Requirements

### Requirement: Create Project

The API SHALL allow creating a new project via POST request.

#### Scenario: Successful creation

- Given: A valid project payload with id and name
- When: POST `/projects` is called
- Then: The project is created with `active=true`
- And: HTTP 201 is returned with the created project

#### Scenario: Duplicate ID

- Given: A project with the same ID already exists
- When: POST `/projects` is called
- Then: HTTP 409 Conflict is returned

### Requirement: List Projects

The API SHALL return all projects with optional filtering.

#### Scenario: List all projects

- Given: Multiple projects exist
- When: GET `/projects` is called
- Then: All projects are returned

#### Scenario: Filter by active status

- Given: Both active and inactive projects exist
- When: GET `/projects?active=true` is called
- Then: Only active projects are returned

### Requirement: Get Single Project

The API SHALL return a single project by ID.

#### Scenario: Project exists

- Given: A project with ID "IZG" exists
- When: GET `/projects/IZG` is called
- Then: The project is returned

#### Scenario: Project not found

- Given: No project with ID "UNKNOWN" exists
- When: GET `/projects/UNKNOWN` is called
- Then: HTTP 404 is returned

### Requirement: Update Project

The API SHALL allow updating an existing project.

#### Scenario: Successful update

- Given: A project with ID "IZG" exists
- When: PUT `/projects/IZG` is called with updated fields
- Then: The project is updated
- And: HTTP 200 is returned with the updated project

#### Scenario: Update non-existent project

- Given: No project with ID "UNKNOWN" exists
- When: PUT `/projects/UNKNOWN` is called
- Then: HTTP 404 is returned

### Requirement: Delete Project

The API SHALL soft-delete a project by setting active to false.

#### Scenario: Successful soft-delete

- Given: An active project with ID "IZG" exists
- When: DELETE `/projects/IZG` is called
- Then: The project's `active` field is set to false
- And: HTTP 200 is returned

#### Scenario: Delete non-existent project

- Given: No project with ID "UNKNOWN" exists
- When: DELETE `/projects/UNKNOWN` is called
- Then: HTTP 404 is returned

#### Scenario: Delete project with time codes

- Given: A project with associated time codes exists
- When: DELETE `/projects/{id}` is called
- Then: The project is soft-deleted
- And: Associated time codes remain but reference an inactive project
