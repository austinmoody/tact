## 1. Schema Definition

- [x] 1.1 Define `PARSE_RESULT_SCHEMA` constant with JSON schema for ParseResult fields
- [x] 1.2 Include all fields: duration_minutes, work_type_id, time_code_id, parsed_description, confidence scores, notes
- [x] 1.3 Set appropriate types (integer for duration, number for confidences, string/null for IDs)
- [x] 1.4 Mark only `confidence_overall` as required

## 2. API Integration

- [x] 2.1 Update `parse()` method to pass schema to `format` parameter instead of `"json"`
- [x] 2.2 Verify request payload structure matches Ollama API expectations

## 3. Testing

- [x] 3.1 Test parsing with schema enforcement enabled
- [x] 3.2 Verify response fields match expected structure
- [x] 3.3 Test edge cases (null values, missing optional fields)
