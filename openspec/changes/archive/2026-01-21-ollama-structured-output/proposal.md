## Why

The current Ollama provider uses `"format": "json"` which requests JSON output but doesn't constrain the structure. This can lead to parsing failures when the model returns unexpected field names, missing fields, or incorrect types. Using Ollama's structured output feature with a JSON schema will enforce the exact response structure, improving parsing reliability.

## What Changes

- Update OllamaProvider to pass a JSON schema to the `format` parameter instead of just `"json"`
- Define a schema that matches the expected ParseResult fields with correct types
- Add schema validation to ensure responses conform to the expected structure

## Capabilities

### New Capabilities
- `ollama-structured-output`: JSON schema-constrained output from Ollama for reliable parsing

### Modified Capabilities
None - this is an implementation enhancement that doesn't change spec-level behavior

## Impact

- **Backend**: `backend/src/tact/llm/ollama.py` - Update the `parse()` method to use JSON schema format
- **API**: No changes - same ParseResult structure returned
- **Dependencies**: None - Ollama already supports this feature
