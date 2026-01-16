# entries-api Specification Delta

## MODIFIED Requirements

### Requirement: Get Single Entry

The API SHALL return a single entry by ID, including parse notes.

#### Scenario: Entry with parse notes

- Given: An entry exists that has been parsed
- When: GET `/entries/{id}` is called
- Then: The response includes `parse_notes` field
- And: `parse_notes` contains the LLM reasoning and context info from parsing

### Requirement: List Entries

The API SHALL return entries with parse notes included.

#### Scenario: List includes parse notes

- Given: Multiple parsed entries exist
- When: GET `/entries` is called
- Then: Each entry in the response includes `parse_notes` field
