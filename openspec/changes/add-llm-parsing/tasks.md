## 1. Provider Abstraction

- [ ] 1.1 Create `backend/src/tact/llm/__init__.py`
- [ ] 1.2 Create `backend/src/tact/llm/provider.py` with:
  - `ParseResult` dataclass (fields + confidence scores)
  - `ParseContext` dataclass (time_codes, work_types lists)
  - `LLMProvider` abstract base class
- [ ] 1.3 Create `backend/src/tact/llm/prompts.py` with system prompt template

## 2. Ollama Provider

- [ ] 2.1 Create `backend/src/tact/llm/ollama.py` with `OllamaProvider` class
- [ ] 2.2 Implement HTTP calls to Ollama API (`/api/generate`)
- [ ] 2.3 Implement JSON response parsing with error handling
- [ ] 2.4 Add configuration support (`TACT_OLLAMA_URL`, `TACT_OLLAMA_MODEL`)

## 3. Entry Parser

- [ ] 3.1 Create `backend/src/tact/llm/parser.py` with `EntryParser` class
- [ ] 3.2 Implement `parse_entry(entry, session)` method:
  - Fetch active time codes and work types
  - Build context for LLM
  - Call provider
  - Update entry with results
- [ ] 3.3 Add provider factory based on `TACT_LLM_PROVIDER` config

## 4. Background Worker

- [ ] 4.1 Create `backend/src/tact/worker/__init__.py`
- [ ] 4.2 Create `backend/src/tact/worker/parser_worker.py`:
  - Query pending entries
  - Parse each entry
  - Handle errors gracefully
  - Sleep for `TACT_PARSER_INTERVAL`
- [ ] 4.3 Integrate worker startup in `main.py` lifespan

## 5. Reparse Endpoint

- [ ] 5.1 Add `POST /entries/{id}/reparse` to entries router
- [ ] 5.2 Implement endpoint logic (clear fields, set pending)

## 6. Docker Integration

- [ ] 6.1 Add Ollama service to `docker-compose.yml` with profile
- [ ] 6.2 Update README with Ollama setup instructions

## 7. Testing

- [ ] 7.1 Create `backend/tests/test_llm_parser.py` with unit tests:
  - Test prompt generation
  - Test response parsing
  - Test error handling
- [ ] 7.2 Create `backend/tests/test_reparse.py` with endpoint tests
- [ ] 7.3 Verify `make test` passes
- [ ] 7.4 Verify `make lint` passes

## 8. Anthropic Provider

- [ ] 8.1 Add `anthropic` dependency to pyproject.toml
- [ ] 8.2 Create `backend/src/tact/llm/anthropic.py` with `AnthropicProvider` class
- [ ] 8.3 Implement Anthropic API calls (Messages API)
- [ ] 8.4 Add configuration support (`TACT_ANTHROPIC_API_KEY`, `TACT_ANTHROPIC_MODEL`)
- [ ] 8.5 Add tests for Anthropic provider

## Verification

1. **Start Ollama locally:**
   ```bash
   # If not using Docker
   ollama pull llama3.2:3b
   ollama serve

   # Or with Docker
   docker compose --profile ollama up -d
   ```

2. **Start the backend:**
   ```bash
   make run
   ```

3. **Create a test entry:**
   ```bash
   # First create a time code and work type
   curl -X POST http://localhost:2100/time-codes \
     -H "Content-Type: application/json" \
     -d '{"id": "PROJ-001", "name": "Project Alpha", "description": "Main project", "keywords": ["alpha", "main"]}'

   curl -X POST http://localhost:2100/work-types \
     -H "Content-Type: application/json" \
     -d '{"name": "Development"}'

   # Create an entry
   curl -X POST http://localhost:2100/entries \
     -H "Content-Type: application/json" \
     -d '{"raw_text": "2h coding on Project Alpha"}'
   ```

4. **Wait for parsing (or check logs):**
   ```bash
   # After ~10 seconds, check the entry
   curl http://localhost:2100/entries
   # Should show status="parsed" with duration_minutes=120
   ```

5. **Test reparse:**
   ```bash
   curl -X POST http://localhost:2100/entries/{id}/reparse
   # Entry should reset to pending
   ```
