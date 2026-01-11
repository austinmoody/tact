# Capability: Time Codes API

REST API for managing time codes used to categorize time entries.

## ADDED Requirements

### Requirement: Create Time Code

The API SHALL allow creating a new time code via POST request.

#### Scenario: Successful creation

- Given: A valid time code payload with id, name, and description
- When: POST `/time-codes` is called
- Then: The time code is created with `active=true`
- And: HTTP 201 is returned with the created time code

#### Scenario: Duplicate ID

- Given: A time code with the same ID already exists
- When: POST `/time-codes` is called
- Then: HTTP 409 Conflict is returned

### Requirement: List Time Codes

The API SHALL return all time codes with optional filtering.

#### Scenario: List all time codes

- Given: Multiple time codes exist
- When: GET `/time-codes` is called
- Then: All time codes are returned

#### Scenario: Filter by active status

- Given: Both active and inactive time codes exist
- When: GET `/time-codes?active=true` is called
- Then: Only active time codes are returned

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

#### Scenario: Successful update

- Given: A time code with ID "PROJ-001" exists
- When: PUT `/time-codes/PROJ-001` is called with updated fields
- Then: The time code is updated
- And: HTTP 200 is returned with the updated time code

#### Scenario: Update non-existent time code

- Given: No time code with ID "UNKNOWN" exists
- When: PUT `/time-codes/UNKNOWN` is called
- Then: HTTP 404 is returned

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
