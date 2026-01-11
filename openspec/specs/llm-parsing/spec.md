# llm-parsing Specification

## Purpose

Provides LLM-based parsing of time entry raw text into structured fields (duration, time code, work type, description) with confidence scores.

## Requirements

### Requirement: Provider Abstraction

The system SHALL support multiple LLM providers through a common interface.

#### Scenario: Ollama provider

- Given: `TACT_LLM_PROVIDER=ollama` is configured
- When: The parser processes an entry
- Then: The Ollama API is used for parsing

#### Scenario: Anthropic provider

- Given: `TACT_LLM_PROVIDER=anthropic` and `TACT_ANTHROPIC_API_KEY` are configured
- When: The parser processes an entry
- Then: The Anthropic API is used for parsing

### Requirement: Parse Entry Fields

The LLM SHALL extract structured fields from raw entry text.

#### Scenario: Full extraction

- Given: Raw text "2h coding on Project Alpha"
- And: Time code "PROJ-001" exists with keyword "alpha"
- And: Work type "development" exists
- When: The entry is parsed
- Then: duration_minutes is set to 120
- And: time_code_id is set to "PROJ-001"
- And: work_type_id is set to "development"
- And: description is generated

#### Scenario: Partial extraction

- Given: Raw text "meeting with team"
- And: No duration is specified
- When: The entry is parsed
- Then: duration_minutes is null
- And: confidence_duration is low (< 0.5)

#### Scenario: Vague input

- Given: Raw text with minimal information (e.g., "stuff")
- When: The entry is parsed
- Then: Fields that cannot be determined are null
- And: confidence_overall is very low (< 0.3)

### Requirement: Confidence Scores

The LLM SHALL provide confidence scores (0.0 to 1.0) for each extracted field.

#### Scenario: High confidence

- Given: Unambiguous input with clear duration and matching codes
- When: The entry is parsed
- Then: confidence_overall is high (> 0.8)

#### Scenario: Low confidence

- Given: Ambiguous input without clear matches
- When: The entry is parsed
- Then: confidence_overall is low
- And: The entry is flagged for review

### Requirement: Background Worker

The system SHALL process pending entries in the background.

#### Scenario: Automatic processing

- Given: An entry is created with status "pending"
- When: The background worker runs (every TACT_PARSER_INTERVAL seconds)
- Then: The entry is parsed
- And: status is updated to "parsed" or "failed"
- And: parsed_at is set to the current timestamp

#### Scenario: Worker disabled

- Given: `TACT_DISABLE_WORKER=true` is configured
- When: The application starts
- Then: The background worker does not start

### Requirement: Reparse Endpoint

The API SHALL allow manually triggering a reparse of an entry.

#### Scenario: Successful reparse

- Given: An entry exists with status "parsed"
- When: POST `/entries/{id}/reparse` is called
- Then: All parsed fields are cleared
- And: status is set to "pending"
- And: manually_corrected is set to false
- And: The entry will be reparsed on the next worker cycle

#### Scenario: Entry not found

- Given: No entry with the specified ID exists
- When: POST `/entries/{id}/reparse` is called
- Then: HTTP 404 is returned

### Requirement: Error Handling

The system SHALL handle parsing errors gracefully.

#### Scenario: LLM unavailable

- Given: The configured LLM provider is unavailable
- When: The worker attempts to parse an entry
- Then: The entry remains in "pending" status
- And: An error is logged
- And: Parsing is retried on the next worker cycle

#### Scenario: Invalid LLM response

- Given: The LLM returns invalid JSON
- When: The worker attempts to parse an entry
- Then: status is set to "failed"
- And: parse_error contains the error message

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `TACT_LLM_PROVIDER` | `ollama` | Provider: `ollama` or `anthropic` |
| `TACT_OLLAMA_URL` | `http://localhost:11434` | Ollama API URL |
| `TACT_OLLAMA_MODEL` | `llama3.2:3b` | Ollama model name |
| `TACT_ANTHROPIC_API_KEY` | (none) | Anthropic API key (required for anthropic provider) |
| `TACT_ANTHROPIC_MODEL` | `claude-3-haiku-20240307` | Anthropic model |
| `TACT_PARSER_INTERVAL` | `10` | Seconds between parse cycles |
| `TACT_DISABLE_WORKER` | `false` | Set to `true` to disable background parsing |
