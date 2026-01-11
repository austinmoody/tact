# Capability: Work Types API

REST API for managing work types used to categorize the type of work performed.

## ADDED Requirements

### Requirement: Create Work Type

The API SHALL allow creating a new work type via POST request. The ID is auto-generated as a URL-friendly slug from the name (e.g., "Code Review" → "code-review").

#### Scenario: Successful creation

- Given: A valid work type payload with name
- When: POST `/work-types` is called
- Then: The work type is created with `active=true` and auto-generated slug ID
- And: HTTP 201 is returned with the created work type

#### Scenario: Duplicate slug

- Given: A work type that generates the same slug already exists
- When: POST `/work-types` is called
- Then: HTTP 409 Conflict is returned

### Requirement: List Work Types

The API SHALL return all work types with optional filtering.

#### Scenario: List all work types

- Given: Multiple work types exist
- When: GET `/work-types` is called
- Then: All work types are returned

#### Scenario: Filter by active status

- Given: Both active and inactive work types exist
- When: GET `/work-types?active=true` is called
- Then: Only active work types are returned

### Requirement: Get Single Work Type

The API SHALL return a single work type by ID.

#### Scenario: Work type exists

- Given: A work type with ID "development" exists
- When: GET `/work-types/development` is called
- Then: The work type is returned

#### Scenario: Work type not found

- Given: No work type with ID "unknown" exists
- When: GET `/work-types/unknown` is called
- Then: HTTP 404 is returned

### Requirement: Update Work Type

The API SHALL allow updating an existing work type.

#### Scenario: Successful update

- Given: A work type with ID "development" exists
- When: PUT `/work-types/development` is called with updated fields
- Then: The work type is updated
- And: HTTP 200 is returned with the updated work type

#### Scenario: Update non-existent work type

- Given: No work type with ID "unknown" exists
- When: PUT `/work-types/unknown` is called
- Then: HTTP 404 is returned

### Requirement: Delete Work Type

The API SHALL soft-delete a work type by setting active to false.

#### Scenario: Successful soft-delete

- Given: An active work type with ID "development" exists
- When: DELETE `/work-types/development` is called
- Then: The work type's `active` field is set to false
- And: HTTP 200 is returned

#### Scenario: Delete non-existent work type

- Given: No work type with ID "unknown" exists
- When: DELETE `/work-types/unknown` is called
- Then: HTTP 404 is returned
