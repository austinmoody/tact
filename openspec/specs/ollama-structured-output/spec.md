## ADDED Requirements

### Requirement: Ollama uses JSON schema for structured output
The OllamaProvider SHALL use a JSON schema to constrain LLM output structure instead of simple JSON format mode.

#### Scenario: Schema enforces field structure
- **WHEN** OllamaProvider calls the Ollama API
- **THEN** the `format` parameter SHALL contain a JSON schema defining expected fields
- **AND** the schema SHALL include: duration_minutes, work_type_id, time_code_id, parsed_description, confidence_duration, confidence_work_type, confidence_time_code, confidence_overall, notes

#### Scenario: Schema allows nullable fields
- **WHEN** the LLM cannot determine a field value
- **THEN** the schema SHALL allow null values for optional fields (work_type_id, time_code_id, parsed_description, notes)
- **AND** confidence scores SHALL be required as numbers

#### Scenario: Response conforms to schema
- **WHEN** Ollama returns a response with schema enforcement
- **THEN** the response SHALL contain only fields defined in the schema
- **AND** field types SHALL match the schema definitions

### Requirement: Ollama structured output improves parsing reliability
The structured output feature SHALL reduce parsing failures from malformed responses.

#### Scenario: No unexpected field names
- **WHEN** the LLM generates a response with schema enforcement
- **THEN** field names SHALL exactly match the schema (e.g., `time_code_id` not `time_code`)

#### Scenario: Required fields always present
- **WHEN** the LLM generates a response with schema enforcement
- **THEN** the `confidence_overall` field SHALL always be present
- **AND** confidence fields SHALL be numeric values
