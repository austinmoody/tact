# Tact Backend

Personal time-tracking tool backend built with FastAPI.

## Prerequisites

- Python 3.12+
- [uv](https://docs.astral.sh/uv/) for dependency management
- Docker (optional, for containerized execution)

## Local Development

> **Note:** See the [Makefile](../Makefile) for underlying commands if you prefer running them directly.

### Install dependencies

```bash
make install
```

### Run the API

```bash
make run
```

The API will be available at `http://localhost:2100`.

### Run tests

```bash
make test
```

### Linting and formatting

```bash
make lint    # Check for issues
make format  # Auto-format code
```

## Database

The backend uses SQLite with Alembic for migrations. The database is automatically created and migrated on application startup.

### Configuration

Set the database path via environment variable:

```bash
export TACT_DB_PATH=/path/to/tact.db
```

Default: `./data/tact.db`

### Manual migration commands

```bash
make migrate                    # Run pending migrations
make db-revision msg="add foo"  # Create new migration
```

## Docker

### Build and run

```bash
make docker-build
make docker-up
```

### Stop

```bash
make docker-down
```

## API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/health` | GET | Health check - returns `{"status": "healthy"}` |
| `/time-codes` | GET | List all time codes (optional `?active=true/false` filter) |
| `/time-codes` | POST | Create a time code (user-specified ID) |
| `/time-codes/{id}` | GET | Get a single time code |
| `/time-codes/{id}` | PUT | Update a time code |
| `/time-codes/{id}` | DELETE | Soft-delete a time code (sets `active=false`) |
| `/work-types` | GET | List all work types (optional `?active=true/false` filter) |
| `/work-types` | POST | Create a work type (ID auto-generated from name as slug) |
| `/work-types/{id}` | GET | Get a single work type |
| `/work-types/{id}` | PUT | Update a work type |
| `/work-types/{id}` | DELETE | Soft-delete a work type (sets `active=false`) |
| `/entries` | GET | List entries (filters: `status`, `time_code_id`, `work_type_id`, `from_date`, `to_date`, `limit`, `offset`) |
| `/entries` | POST | Create an entry (`raw_text` required, `entry_date` optional - defaults to today) |
| `/entries/{id}` | GET | Get a single entry |
| `/entries/{id}` | PATCH | Update an entry (sets `manually_corrected=true`) |
| `/entries/{id}` | DELETE | Hard-delete an entry (returns 204) |
| `/entries/{id}/reparse` | POST | Reset entry to pending for re-parsing |

## LLM Integration

The backend uses an LLM to parse time entries. It supports two providers:

### Ollama (Local)

1. **Install Ollama:**
   ```bash
   # macOS
   brew install ollama

   # Or download from https://ollama.ai
   ```

2. **Pull a model:**
   ```bash
   ollama pull llama3.2:3b
   ```

3. **Start Ollama:**
   ```bash
   ollama serve
   ```

4. **Configure (optional):**
   ```bash
   export TACT_LLM_PROVIDER=ollama           # Default
   export TACT_OLLAMA_URL=http://localhost:11434  # Default
   export TACT_OLLAMA_MODEL=llama3.2:3b      # Default
   ```

### Anthropic (Cloud)

1. **Get an API key** from https://console.anthropic.com

2. **Configure:**
   ```bash
   export TACT_LLM_PROVIDER=anthropic
   export TACT_ANTHROPIC_API_KEY=sk-ant-...
   export TACT_ANTHROPIC_MODEL=claude-3-haiku-20240307  # Default
   ```

### Docker with Ollama

Both the API and Ollama start together:

```bash
make docker-up
docker compose exec ollama ollama pull llama3.2:3b
```

### Configuration Reference

| Variable | Default | Description |
|----------|---------|-------------|
| `TACT_LLM_PROVIDER` | `ollama` | Provider: `ollama` or `anthropic` |
| `TACT_OLLAMA_URL` | `http://localhost:11434` | Ollama API URL |
| `TACT_OLLAMA_MODEL` | `llama3.2:3b` | Ollama model name |
| `TACT_ANTHROPIC_API_KEY` | (none) | Anthropic API key |
| `TACT_ANTHROPIC_MODEL` | `claude-3-haiku-20240307` | Anthropic model |
| `TACT_PARSER_INTERVAL` | `10` | Seconds between parse cycles |
| `TACT_DISABLE_WORKER` | `false` | Set to `true` to disable background parsing |
