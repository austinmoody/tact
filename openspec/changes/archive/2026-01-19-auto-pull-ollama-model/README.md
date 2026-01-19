# Auto-pull Ollama Model on First Use

## Summary

When the Ollama LLM provider is first used and the configured model is not available locally, automatically pull the model from the Ollama registry before proceeding. This enables a "just works" experience for fresh deployments without requiring manual model setup.

## Motivation

Currently, if `TACT_OLLAMA_MODEL` specifies a model that hasn't been pulled to the Ollama container, requests fail. This requires manual intervention to run `ollama pull <model>` before the system works. Auto-pulling removes this friction for new deployments.

## Approach

**Lazy loading on first AI use:**
- When `OllamaProvider.parse()` is first called, check if the model exists
- If not, trigger `ollama pull` and wait for completion
- Cache the "model verified" state to avoid repeated checks
- Subsequent calls skip the check entirely

This approach was chosen because:
- Parsing is a background task - users aren't waiting for immediate results
- API startup remains fast
- Only pulls if AI features are actually used

## Implementation

### Changes to `backend/src/tact/llm/ollama.py`

1. Add `_model_verified: bool` instance variable (default `False`)

2. Add `_ensure_model_available()` method:
   - Call `GET /api/tags` to list available models
   - If configured model not found, call `POST /api/pull` with `{"name": model}`
   - The pull endpoint streams progress; wait for completion
   - Set `_model_verified = True` on success
   - Log progress/status for visibility

3. Modify `parse()` to call `_ensure_model_available()` at the start if `_model_verified` is `False`

### Error Handling

- If model pull fails (network error, invalid model name), return a `ParseResult` with an appropriate error message
- Log the error for debugging

### Timeout Considerations

- Model pulls can take several minutes for large models
- Use a separate longer timeout for the pull operation (configurable via `TACT_OLLAMA_PULL_TIMEOUT`, default 600 seconds / 10 minutes)
- The regular `TACT_OLLAMA_TIMEOUT` continues to apply to inference requests

## Affected Specs

- `ollama-provider` - Adding model auto-pull capability

## Testing

- Unit test: Mock Ollama API to verify pull is triggered when model missing
- Unit test: Verify pull is skipped when model already exists
- Unit test: Verify `_model_verified` flag prevents repeated checks
- Integration test: Fresh Ollama container successfully pulls model on first parse request
