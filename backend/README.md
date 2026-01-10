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
