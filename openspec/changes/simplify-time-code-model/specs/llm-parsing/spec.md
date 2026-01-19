## MODIFIED Requirements

### Requirement: Parse Entry Fields

The LLM SHALL extract structured fields from raw entry text using RAG-retrieved context.

#### Scenario: No relevant context

- Given: Raw text "30m general admin work"
- And: No specific context matches this entry
- When: The entry is parsed
- Then: The LLM falls back to time code names
- And: A reasonable match is made based on available information

### Requirement: Parse Notes

The LLM SHALL provide reasoning notes explaining how parsing decisions were made.

#### Scenario: No matching context

- Given: Raw text "30m general admin"
- And: No relevant context documents match
- When: The entry is parsed
- Then: `parse_notes` explains matching was based on time code names
- And: `parse_notes` indicates no specific context rules applied

### Requirement: RAG Context Retrieval

The parser SHALL retrieve relevant context documents before calling the LLM.

#### Scenario: Empty context store

- Given: No context documents exist in the system
- When: An entry is parsed
- Then: Parsing proceeds without RAG context
- And: The LLM uses only time code names for matching
