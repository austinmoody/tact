# time-codes-api Specification

## Purpose
TBD - created by archiving change add-timecodes-worktypes-api. Update Purpose after archive.
## Requirements
### Requirement: Create Time Code

The API SHALL allow creating a new time code via POST request.

#### Scenario: Successful creation with project

- Given: A valid time code payload with id, name, description, and project_id
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

### Requirement: List Time Codes

The API SHALL return all time codes with optional filtering.

#### Scenario: Filter by project

- Given: Time codes exist for multiple projects
- When: GET `/time-codes?project_id=IZG` is called
- Then: Only time codes for project "IZG" are returned

### Requirement: Get Single Time Code

The API SHALL return a single time code by ID.

#### Scenario: Time code exists

- Given: A time code with ID "PROJ-001" exists
- When: GET `/time-codes/PROJ-001` is called
- Then: The time code is returned

#### Scenario: Time code not found

- Given: No time code with ID "UNKNOWN" exists
- When: GET `/time-codes/UNKNOWN` is called
- Then: HTTP 404 is returned

### Requirement: Update Time Code

The API SHALL allow updating an existing time code.

#### Scenario: Change project association

- Given: A time code with ID "FEDS-163" exists in project "IZG"
- And: Project "OTHER" exists
- When: PUT `/time-codes/FEDS-163` is called with `project_id: "OTHER"`
- Then: The time code is moved to project "OTHER"
- And: HTTP 200 is returned

### Requirement: Delete Time Code

The API SHALL soft-delete a time code by setting active to false.

#### Scenario: Successful soft-delete

- Given: An active time code with ID "PROJ-001" exists
- When: DELETE `/time-codes/PROJ-001` is called
- Then: The time code's `active` field is set to false
- And: HTTP 200 is returned

#### Scenario: Delete non-existent time code

- Given: No time code with ID "UNKNOWN" exists
- When: DELETE `/time-codes/UNKNOWN` is called
- Then: HTTP 404 is returned

