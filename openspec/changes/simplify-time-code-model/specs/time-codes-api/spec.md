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

### Requirement: Time Code Fields

The TimeCode model SHALL have the following fields only:

#### Scenario: Minimal field set

- Given: A time code is created or retrieved
- Then: The time code has fields: id, project_id, name, active, created_at, updated_at
- And: No description, keywords, or examples fields exist

## REMOVED Requirements

### Requirement: Time Code Description

~~The TimeCode model included a description field for longer explanations.~~

Removed: Description was redundant with name field. RAG context handles matching.

### Requirement: Time Code Keywords

~~The TimeCode model included keywords for LLM matching hints.~~

Removed: RAG context documents provide better matching through learned examples.

### Requirement: Time Code Examples

~~The TimeCode model included examples for LLM context.~~

Removed: Field was never used in LLM prompt. RAG context serves this purpose.
