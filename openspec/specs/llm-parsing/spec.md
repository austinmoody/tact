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

The LLM SHALL extract structured fields from raw entry text using RAG-retrieved context.

#### Scenario: Context-aware matching

- Given: Raw text "2h APHL meeting about UI"
- And: Context document exists: "ALL meetings with APHL go to FEDS-163 regardless of topic"
- And: Context document exists: "ALL UI work goes to FEDS-167"
- When: The entry is parsed
- Then: Relevant context is retrieved via vector similarity
- And: time_code_id is set to "FEDS-163" (APHL rule overrides UI)
- And: The LLM uses context to make the disambiguation

#### Scenario: Acronym expansion

- Given: Raw text "1h IZG deployment"
- And: Project context exists: "IZG = IZ Gateway"
- And: Time code context exists for FEDS-163: "ALL deployments go to this code"
- When: The entry is parsed
- Then: The LLM understands IZG refers to IZ Gateway
- And: time_code_id is set to "FEDS-163"

#### Scenario: No relevant context

- Given: Raw text "30m general admin work"
- And: No specific context matches this entry
- When: The entry is parsed
- Then: The LLM falls back to time code descriptions and keywords
- And: A reasonable match is made based on available information

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

### Requirement: RAG Context Retrieval

The parser SHALL retrieve relevant context documents before calling the LLM.

#### Scenario: Retrieve similar context

- Given: An entry "2h security scan review"
- And: Context documents exist with various content
- When: The parser prepares the LLM prompt
- Then: The entry text is embedded using the same model as context docs
- And: Top-k most similar context chunks are retrieved
- And: Retrieved chunks are included in the LLM prompt

#### Scenario: Context includes source

- Given: Context is retrieved for an entry
- When: The LLM prompt is built
- Then: Each context chunk is labeled with its source (project or time code)
- And: The LLM can see which time code each rule applies to

#### Scenario: Empty context store

- Given: No context documents exist in the system
- When: An entry is parsed
- Then: Parsing proceeds without RAG context
- And: The LLM uses only time code descriptions and keywords

### Requirement: Local Embeddings

The system SHALL generate embeddings locally without external API calls.

#### Scenario: Embed entry text

- Given: An entry is submitted for parsing
- When: The parser retrieves context
- Then: The entry text is embedded using a local model
- And: No external embedding API is called

#### Scenario: Embed context document

- Given: A context document is created
- When: The document is saved
- Then: The content is embedded using a local model
- And: The embedding is stored with the document

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
