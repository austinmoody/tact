## Context

The Ollama provider (`backend/src/tact/llm/ollama.py`) currently uses `"format": "json"` when calling the Ollama API. This tells Ollama to return JSON, but doesn't constrain the structure. The model can return:
- Unexpected field names (e.g., `time_code` instead of `time_code_id`)
- Missing required fields
- Wrong types (e.g., string instead of number for confidence)

Ollama supports structured output by passing a JSON schema to the `format` parameter. When a schema is provided, the model is constrained to output only valid JSON matching that schema.

Current code:
```python
response = self.client.post(
    f"{self.base_url}/api/generate",
    json={
        "model": self.model,
        "prompt": f"{system_prompt}\n\n{user_prompt}",
        "stream": False,
        "format": "json",  # Just requests JSON, no structure enforcement
    },
)
```

## Goals / Non-Goals

**Goals:**
- Enforce exact JSON structure from Ollama responses using JSON schema
- Reduce parsing failures from malformed or unexpected responses
- Maintain backward compatibility with existing ParseResult structure

**Non-Goals:**
- Changing the ParseResult fields or types
- Adding new LLM providers
- Modifying prompt engineering

## Decisions

### 1. Schema definition approach

**Decision**: Define the schema as a Python dict constant in `ollama.py`

**Rationale**:
- Single file change, no new dependencies
- Schema stays close to the code that uses it
- Easy to update if ParseResult changes

**Alternative considered**: External JSON schema file - rejected as overkill for a single use case

### 2. Schema structure

**Decision**: Use JSON Schema draft-07 compatible format with all ParseResult fields

```python
PARSE_RESULT_SCHEMA = {
    "type": "object",
    "properties": {
        "duration_minutes": {"type": "integer"},
        "work_type_id": {"type": ["string", "null"]},
        "time_code_id": {"type": ["string", "null"]},
        "parsed_description": {"type": ["string", "null"]},
        "confidence_duration": {"type": "number"},
        "confidence_work_type": {"type": "number"},
        "confidence_time_code": {"type": "number"},
        "confidence_overall": {"type": "number"},
        "notes": {"type": ["string", "null"]}
    },
    "required": ["confidence_overall"]
}
```

**Rationale**:
- Only `confidence_overall` is truly required - other fields can be null/missing
- Using `["string", "null"]` allows nullable fields
- Confidence scores as `number` allows both int and float

### 3. Error handling

**Decision**: Keep existing error handling, schema validation happens at Ollama level

**Rationale**:
- Ollama enforces the schema during generation
- If schema validation fails in Ollama, it returns an error we already handle
- No additional client-side validation needed

## Risks / Trade-offs

**Risk**: Ollama version compatibility - older versions may not support schema format
→ **Mitigation**: Document minimum Ollama version requirement; fallback to `"format": "json"` if schema fails

**Risk**: Schema constraints could cause model to fail on edge cases
→ **Mitigation**: Use permissive schema (nullable fields, only one required field); monitor parsing success rate

**Trade-off**: Slightly more complex code vs more reliable parsing
→ **Accepted**: Reliability gain outweighs minor complexity increase
