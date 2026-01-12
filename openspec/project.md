# Project Context

## Purpose

Tact is a personal time tracking tool designed for logging work hours using natural language input. The primary goal is to make time entry effortless by allowing entries like:

- "Development on Project Alpha PROJ-2345 for 45m"
- "Daily standup meeting 30 minutes"
- "Corporate training 1 hour"
- "Worked on deployment scripts yesterday for 1 hour"

The system parses these entries to extract:
- **Duration** - normalized to minutes for calculations
- **Work type** - coding, research, troubleshooting, meeting, training, etc.
- **Time code** - maps to organizational billing categories
- **Entry date** - supports relative dates ("yesterday", "last Friday")

This is API-first to support multiple future clients (web UI, desktop apps, MCP server for chat).

## Tech Stack

- **Language**: Python 3.12+
- **Framework**: FastAPI
- **Database**: SQLite (file mounted via Docker volume)
- **Deployment**: Docker Compose
- **LLM Integration**: External API for natural language parsing (provider TBD)
- **Future UI**: HTMX + Jinja2 templates (no JavaScript frameworks)

## Project Conventions

### Code Style

- Follow PEP 8 and use `ruff` for linting and formatting
- Type hints required on all function signatures
- Use `pydantic` for request/response models and validation
- Prefer explicit over implicit

### Architecture Patterns

- **API-first design** - all functionality exposed via REST endpoints
- **Service layer** - business logic separated from API routes
- **Repository pattern** - database access abstracted behind repositories
- Keep it simple - avoid over-engineering for hypothetical future needs

### Testing Strategy

- `pytest` for all tests
- Unit tests for parsing logic and business rules
- Integration tests for API endpoints
- Test database uses in-memory SQLite

### Git Workflow

- `main` branch is stable/release-ready
- `develop` branch for active development
- Feature branches off `develop`
- Conventional commits preferred (feat:, fix:, docs:, etc.)

## Domain Context

### Time Codes

Time codes represent organizational billing categories. They are:
- **Dynamic** - new codes are added when contracts are won, removed when contracts end
- **Remappable** - the same type of work may need to map to different codes based on current business direction
- Managed via API endpoints (CRUD operations)

### Work Types

Categories of work being performed:
- Development / Coding
- Research
- Troubleshooting
- Meetings
- Training
- Documentation
- Code Review
- (extensible list)

### Entry Parsing

Natural language entries need to extract:
1. Duration (various formats: "45m", "1.5 hours", "1h30m", "90 minutes")
2. Relative or absolute dates ("yesterday", "last Friday", "2024-01-15")
3. Work type classification
4. Time code assignment based on current mappings

Original entry text is always preserved for potential re-parsing if mappings change.

### Parsing Flow

Parsing happens **asynchronously** via background processing:

1. User submits natural language entry
2. Entry is saved immediately with status `pending`
3. API returns entry ID right away (fast response for quick entry workflow)
4. Background worker picks up pending entries
5. LLM parses entry using time code context (descriptions + examples)
6. Entry updated with parsed fields, confidence scores, and new status

### Entry Statuses

| Status | Meaning |
|--------|---------|
| `pending` | Submitted, awaiting parsing |
| `parsed` | Successfully parsed, confidence met threshold |
| `needs_review` | Parsed but confidence below threshold |
| `failed` | LLM call failed, needs retry or manual entry |

### Confidence Scoring

Each parsed field receives a confidence score (0.0 - 1.0):
- Duration confidence
- Work type confidence
- Time code confidence

Overall entry confidence = minimum of all field confidences.

A configurable **threshold** (default: 0.7) determines whether an entry is marked `parsed` or `needs_review`.

### Time Code Categorization

Time codes use a hybrid approach for LLM classification:

```
TimeCode:
  id: str              # Unique identifier (e.g., "PROJ-ALPHA")
  name: str            # Human-readable name
  description: str     # Rich description for LLM context
  keywords: list[str]  # Pattern matches (project names, prefixes)
  examples: list[str]  # Example entries that map to this code
  active: bool         # Can be deactivated without deletion
```

The LLM prompt includes all active time codes with their descriptions and examples to make classification decisions.

### Manual Corrections

When a user manually corrects a parsed entry:
- Entry is marked `manually_corrected = true`
- User chooses whether to **lock** the entry (skip on future re-parses) or **use as example** (correction improves future classification)
- If not locked, the corrected entry's raw_text + corrected values can be fed back as training examples

## Data Model

### TimeEntry

```
TimeEntry:
  id: uuid
  raw_text: str                     # Original input (immutable)

  # Parsed fields
  duration_minutes: int | None
  work_type_id: str | None          # FK to WorkType
  time_code_id: str | None          # FK to TimeCode (soft-delete only)
  description: str | None
  entry_date: date | None

  # Confidence scores (0.0 - 1.0)
  confidence_duration: float | None
  confidence_work_type: float | None
  confidence_time_code: float | None
  confidence_overall: float | None

  # Status
  status: str                       # pending, parsed, needs_review, failed
  parse_error: str | None

  # Correction tracking
  manually_corrected: bool          # Has user edited this entry?
  locked: bool                      # Skip on re-parse?
  corrected_at: datetime | None

  # Timestamps
  created_at: datetime
  parsed_at: datetime | None
  updated_at: datetime
```

### TimeCode

```
TimeCode:
  id: str                           # PK (e.g., "PROJ-ALPHA")
  name: str
  description: str
  keywords: list[str]               # JSON array in SQLite
  examples: list[str]               # JSON array in SQLite
  active: bool                      # Soft-delete (never hard delete)
  created_at: datetime
  updated_at: datetime
```

### WorkType

```
WorkType:
  id: str                           # PK (e.g., "development")
  name: str                         # "Development"
  description: str | None           # Optional context for LLM
  active: bool
  created_at: datetime
  updated_at: datetime
```

### Config

```
Config:
  key: str                          # PK (e.g., "confidence_threshold")
  value: str
  updated_at: datetime
```

## API Design

### Entries

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/entries` | Create entry (returns immediately, parsing in background) |
| `GET` | `/entries` | List entries with filters |
| `GET` | `/entries/{id}` | Get single entry |
| `PATCH` | `/entries/{id}` | Update/correct an entry |
| `DELETE` | `/entries/{id}` | Delete entry |
| `POST` | `/entries/{id}/reparse` | Re-trigger parsing for this entry |
| `POST` | `/entries/reparse` | Bulk re-parse (all unlocked entries) |

**Query params for `GET /entries`:**
- `status` - filter by pending, parsed, needs_review, failed
- `time_code_id` - filter by time code
- `work_type_id` - filter by work type
- `from_date` / `to_date` - date range
- `limit` / `offset` - pagination

### Time Codes

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/time-codes` | Create time code |
| `GET` | `/time-codes` | List all (optionally filter active only) |
| `GET` | `/time-codes/{id}` | Get single time code |
| `PUT` | `/time-codes/{id}` | Update time code |
| `DELETE` | `/time-codes/{id}` | Soft-delete (deactivate) |

### Work Types

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/work-types` | Create work type |
| `GET` | `/work-types` | List all |
| `GET` | `/work-types/{id}` | Get single |
| `PUT` | `/work-types/{id}` | Update |
| `DELETE` | `/work-types/{id}` | Soft-delete |

### Reports

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/reports/summary` | Aggregated time by time code, work type, date range |

### Config

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/config` | Get all config values |
| `GET` | `/config/{key}` | Get specific config |
| `PUT` | `/config/{key}` | Update config value |

### Health

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/health` | Basic health check |

## Important Constraints

- **No JavaScript/TypeScript/Node.js** - This project explicitly avoids the JavaScript ecosystem. Backend is Python, TUI is Go, and any other components must use languages other than JavaScript/TypeScript.
- **Single user focus** - no auth/multi-tenancy needed initially
- **Local-first** - designed to run via `docker compose up` on a personal machine
- **Portable data** - SQLite file can be easily backed up or moved
- **External LLM dependency** - parsing requires network access to LLM API

## External Dependencies

- **LLM API** (TBD) - for natural language parsing of time entries
- No other external services required for core functionality
