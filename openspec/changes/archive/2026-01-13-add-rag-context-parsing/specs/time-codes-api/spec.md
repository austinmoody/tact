# time-codes-api Specification Delta

## MODIFIED Requirements

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

### Requirement: Update Time Code

The API SHALL allow updating an existing time code.

#### Scenario: Change project association

- Given: A time code with ID "FEDS-163" exists in project "IZG"
- And: Project "OTHER" exists
- When: PUT `/time-codes/FEDS-163` is called with `project_id: "OTHER"`
- Then: The time code is moved to project "OTHER"
- And: HTTP 200 is returned
