# Design: LLM Entry Parsing

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                     Entry Created                        │
│                   (status=pending)                       │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│                  Background Worker                       │
│              (polls for pending entries)                 │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│                   LLM Parser                             │
│  ┌─────────────────────────────────────────────────┐    │
│  │              Provider Interface                  │    │
│  │  - parse(raw_text, context) -> ParseResult      │    │
│  └─────────────────────────────────────────────────┘    │
│           │                         │                    │
│           ▼                         ▼                    │
│  ┌─────────────────┐     ┌─────────────────────┐        │
│  │ OllamaProvider  │     │  AnthropicProvider  │        │
│  │ (local, free)   │     │  (cloud, paid)      │        │
│  └─────────────────┘     └─────────────────────────┘    │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│                   Entry Updated                          │
│  - duration_minutes, work_type_id, time_code_id         │
│  - description, confidence scores                        │
│  - status = "parsed" | "failed"                         │
│  - parsed_at = now()                                    │
└─────────────────────────────────────────────────────────┘
```

## Directory Structure

```
backend/src/tact/
├── llm/
│   ├── __init__.py
│   ├── provider.py      # Abstract base class
│   ├── ollama.py        # Ollama implementation
│   ├── anthropic.py     # Anthropic implementation
│   ├── parser.py        # Entry parsing logic
│   └── prompts.py       # System prompts
└── worker/
    └── parser_worker.py # Background task
```

## Provider Interface

```python
class ParseResult:
    duration_minutes: int | None
    work_type_id: str | None
    time_code_id: str | None
    description: str | None
    confidence_duration: float
    confidence_work_type: float
    confidence_time_code: float
    confidence_overall: float

class LLMProvider(ABC):
    @abstractmethod
    def parse(self, raw_text: str, context: ParseContext) -> ParseResult:
        """Parse raw text into structured entry fields."""
        pass
```

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `TACT_LLM_PROVIDER` | `ollama` | Provider: `ollama` or `anthropic` |
| `TACT_OLLAMA_URL` | `http://localhost:11434` | Ollama API URL |
| `TACT_OLLAMA_MODEL` | `llama3.2:3b` | Ollama model name |
| `TACT_ANTHROPIC_API_KEY` | (none) | Anthropic API key |
| `TACT_ANTHROPIC_MODEL` | `claude-3-haiku-20240307` | Anthropic model |
| `TACT_PARSER_INTERVAL` | `10` | Seconds between parse attempts |

## Prompt Design

The prompt includes:
1. **System context**: Available time codes and work types with descriptions/keywords
2. **Task instruction**: Parse the raw text into structured JSON
3. **Output format**: Strict JSON schema with confidence scores

Example prompt context:
```
Available Time Codes:
- PROJ-001: Project Alpha (keywords: alpha, main project)
- ADMIN-01: Administrative Tasks (keywords: admin, paperwork)

Available Work Types:
- development: Software development
- meeting: Meetings and calls
- code-review: Code review activities

Parse this time entry and return JSON:
"2h working on the main project, code review for the auth PR"
```

Expected output:
```json
{
  "duration_minutes": 120,
  "time_code_id": "PROJ-001",
  "work_type_id": "code-review",
  "description": "Code review for auth PR",
  "confidence_duration": 0.95,
  "confidence_time_code": 0.85,
  "confidence_work_type": 0.90,
  "confidence_overall": 0.85
}
```

## Background Worker

Simple polling approach:
1. Query for entries with `status="pending"`
2. For each entry, call LLM parser
3. Update entry with parsed fields and new status
4. Sleep for `TACT_PARSER_INTERVAL` seconds
5. Repeat

The worker runs as a background task started on app startup.

## Docker Integration

Add optional Ollama service to docker-compose.yml:

```yaml
services:
  ollama:
    image: ollama/ollama
    ports:
      - "11434:11434"
    volumes:
      - ollama_data:/root/.ollama
    profiles:
      - ollama  # Only starts with --profile ollama
```

## Error Handling

- **LLM unavailable**: Entry stays `pending`, retry on next cycle
- **Parse failure**: Set `status="failed"`, store error in `parse_error`
- **Low confidence**: Entry marked `parsed` but flagged for review (confidence < 0.7)
- **Invalid JSON response**: Retry once, then mark as `failed`

## Reparse Endpoint

`POST /entries/{id}/reparse` allows manually triggering a reparse:
- Clears existing parsed fields
- Sets `status="pending"`
- Clears `manually_corrected` flag
- Entry will be picked up by next worker cycle
