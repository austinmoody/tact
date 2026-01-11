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
