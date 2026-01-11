# Capability: LLM Entry Parsing

LLM-powered parsing of raw time entry text into structured fields.

## ADDED Requirements

### Requirement: Provider Abstraction

The system SHALL support multiple LLM providers via a common interface.

#### Scenario: Ollama provider configured

- Given: `TACT_LLM_PROVIDER=ollama` is set
- When: The parser is initialized
- Then: The Ollama provider is used for parsing

#### Scenario: Anthropic provider configured

- Given: `TACT_LLM_PROVIDER=anthropic` is set
- When: The parser is initialized
- Then: The Anthropic provider is used for parsing

#### Scenario: Default provider

- Given: No `TACT_LLM_PROVIDER` is set
- When: The parser is initialized
- Then: The Ollama provider is used by default

### Requirement: Entry Parsing

The LLM SHALL parse raw text into structured entry fields.

#### Scenario: Successful parse

- Given: A pending entry with raw_text "2h coding on Project Alpha"
- And: Time code "PROJ-001" exists with name "Project Alpha"
- And: Work type "development" exists
- When: The entry is parsed
- Then: `duration_minutes` is set to 120
- And: `time_code_id` is set to "PROJ-001"
- And: `work_type_id` is set to "development"
- And: `status` is set to "parsed"
- And: `parsed_at` is set to current timestamp
- And: Confidence scores are populated

#### Scenario: Parse with low confidence

- Given: A pending entry with ambiguous raw_text
- When: The entry is parsed
- And: `confidence_overall` is below 0.7
- Then: The entry is marked as parsed
- And: The low confidence is preserved for review

#### Scenario: Parse failure

- Given: A pending entry
- When: The LLM returns invalid output
- Then: `status` is set to "failed"
- And: `parse_error` contains the error message

### Requirement: Background Worker

The system SHALL process pending entries automatically.

#### Scenario: Worker processes pending entries

- Given: Entries exist with `status="pending"`
- When: The background worker runs
- Then: Each pending entry is parsed via the LLM
- And: Entry status is updated accordingly

#### Scenario: Worker handles LLM unavailability

- Given: The LLM provider is unreachable
- When: The background worker attempts to parse
- Then: Entries remain with `status="pending"`
- And: No errors are raised to crash the application

### Requirement: Reparse Entry

The API SHALL allow manually triggering a reparse of an entry.

#### Scenario: Successful reparse trigger

- Given: An entry exists with `status="parsed"`
- When: POST `/entries/{id}/reparse` is called
- Then: The entry's parsed fields are cleared
- And: `status` is set to "pending"
- And: `manually_corrected` is set to false
- And: HTTP 200 is returned

#### Scenario: Reparse non-existent entry

- Given: No entry with ID "unknown-uuid" exists
- When: POST `/entries/unknown-uuid/reparse` is called
- Then: HTTP 404 is returned

### Requirement: Configuration

The system SHALL be configurable via environment variables.

#### Scenario: Custom Ollama URL

- Given: `TACT_OLLAMA_URL=http://ollama:11434` is set
- When: The Ollama provider makes requests
- Then: Requests are sent to the specified URL

#### Scenario: Custom Ollama model

- Given: `TACT_OLLAMA_MODEL=qwen2.5:3b` is set
- When: The Ollama provider parses entries
- Then: The specified model is used

#### Scenario: Custom parse interval

- Given: `TACT_PARSER_INTERVAL=30` is set
- When: The background worker runs
- Then: It waits 30 seconds between parse cycles
