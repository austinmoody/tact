# Proposal: Add LLM Entry Parsing

## Why

Entries are currently created with `status="pending"` and only `raw_text` populated. We need an LLM to parse the raw text and extract structured fields (duration, time code, work type, description) so entries can be reviewed and submitted.

## What Changes

Add a provider-agnostic LLM integration that:
1. Parses pending entries in the background
2. Supports both Ollama (local) and Anthropic (cloud) via configuration
3. Produces confidence scores for each parsed field
4. Updates entry status to `parsed` or `failed`

## Scope

### In Scope

- LLM provider abstraction with common interface
- Ollama provider implementation (local models)
- Anthropic provider implementation (Claude Haiku)
- Background parsing task that processes pending entries
- Configuration via environment variables (`TACT_LLM_PROVIDER`, `TACT_OLLAMA_URL`, etc.)
- Prompt engineering for structured extraction
- Confidence score calculation
- Docker Compose addition for optional Ollama container
- `POST /entries/{id}/reparse` endpoint to manually trigger re-parsing

### Out of Scope

- Real-time parsing (entries are parsed in background, not on create)
- Batch reparse endpoint (`POST /entries/reparse`) - future enhancement
- Fine-tuning or custom models
- OpenRouter provider (can add later if needed)
- Other providers beyond Ollama and Anthropic

## Key Decisions

1. **Start with Ollama** - Experiment with local models first (Llama 3.2 3B or Qwen 2.5 3B)
2. **Provider abstraction** - Clean interface allows swapping providers without code changes
3. **Background processing** - Parsing happens async, not blocking entry creation
4. **Confidence scores** - LLM self-reports confidence; low confidence entries flagged for review
5. **Graceful degradation** - If LLM unavailable, entries stay pending (no crashes)

## Affected Specs

- `llm-parsing` (NEW) - LLM provider abstraction and parsing logic
