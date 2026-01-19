# Tasks

## Implementation
- [x] Add `TACT_OLLAMA_PULL_TIMEOUT` configuration constant
- [x] Add `pull_timeout` parameter to `OllamaProvider.__init__`
- [x] Add `_model_verified` flag to track model availability
- [x] Implement `_ensure_model_available()` method
- [x] Call `_ensure_model_available()` at start of `parse()` method

## Testing
- [x] Add test for model check when model exists
- [x] Add test for auto-pull when model is missing
- [x] Add test for skipping check when already verified
- [x] Add test for error handling on pull failure
- [x] Update existing tests to set `_model_verified = True`

## Documentation
- [x] Update backend README with new environment variable
- [x] Update Docker with Ollama section to mention auto-pull
- [x] Add note about auto-pull in Ollama setup instructions
