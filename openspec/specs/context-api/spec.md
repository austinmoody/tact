# context-api Specification

## Purpose
TBD - created by archiving change add-rag-context-parsing. Update Purpose after archive.
## Requirements
### Requirement: Add Project Context

The API SHALL allow adding context documents to a project.

#### Scenario: Successful creation

- Given: A project with ID "IZG" exists
- And: A valid context payload with content
- When: POST `/projects/IZG/context` is called
- Then: The context document is created
- And: The content is embedded for vector search
- And: HTTP 201 is returned with the created context document

#### Scenario: Project not found

- Given: No project with ID "UNKNOWN" exists
- When: POST `/projects/UNKNOWN/context` is called
- Then: HTTP 404 is returned

### Requirement: Add Time Code Context

The API SHALL allow adding context documents to a time code.

#### Scenario: Successful creation

- Given: A time code with ID "FEDS-163" exists
- And: A valid context payload with content
- When: POST `/time-codes/FEDS-163/context` is called
- Then: The context document is created
- And: The content is embedded for vector search
- And: HTTP 201 is returned with the created context document

#### Scenario: Time code not found

- Given: No time code with ID "UNKNOWN" exists
- When: POST `/time-codes/UNKNOWN/context` is called
- Then: HTTP 404 is returned

### Requirement: List Project Context

The API SHALL return all context documents for a project.

#### Scenario: List context documents

- Given: A project with ID "IZG" exists
- And: Multiple context documents are associated with it
- When: GET `/projects/IZG/context` is called
- Then: All context documents for the project are returned

#### Scenario: Empty context

- Given: A project with ID "IZG" exists
- And: No context documents are associated with it
- When: GET `/projects/IZG/context` is called
- Then: An empty list is returned

### Requirement: List Time Code Context

The API SHALL return all context documents for a time code.

#### Scenario: List context documents

- Given: A time code with ID "FEDS-163" exists
- And: Multiple context documents are associated with it
- When: GET `/time-codes/FEDS-163/context` is called
- Then: All context documents for the time code are returned

### Requirement: Delete Context Document

The API SHALL allow deleting a context document.

#### Scenario: Successful deletion

- Given: A context document with a known ID exists
- When: DELETE `/context/{id}` is called
- Then: The context document is deleted
- And: The embedding is removed from the vector index
- And: HTTP 200 is returned

#### Scenario: Context not found

- Given: No context document with the specified ID exists
- When: DELETE `/context/{id}` is called
- Then: HTTP 404 is returned

### Requirement: Update Context Document

The API SHALL allow updating a context document's content.

#### Scenario: Successful update

- Given: A context document with a known ID exists
- When: PUT `/context/{id}` is called with new content
- Then: The content is updated
- And: The embedding is regenerated
- And: HTTP 200 is returned with the updated document

#### Scenario: Context not found

- Given: No context document with the specified ID exists
- When: PUT `/context/{id}` is called
- Then: HTTP 404 is returned

