## Context

This is the initial project setup establishing patterns that will be used throughout development. The backend is a FastAPI application for a personal time-tracking tool. It needs to support both local development (via `uv`) and containerized execution (via Docker).

**Constraints from project.md:**
- Python 3.12+
- FastAPI framework
- SQLite database (file mounted via Docker volume)
- Docker Compose deployment
- No JavaScript
- Single user focus

## Goals / Non-Goals

**Goals:**
- Establish a clean, standard Python package structure
- Enable running the API locally with `uv run` or via Docker
- Provide a Makefile for common operations
- Set up minimal dependencies (FastAPI + uvicorn only for now)

**Non-Goals:**
- Database setup (future change)
- Authentication (not needed per project.md)
- CI/CD configuration (future change)
- Testing infrastructure beyond basic pytest setup

## Decisions

### Package Structure

```
backend/
├── src/
│   └── tact/
│       ├── __init__.py
│       ├── main.py          # FastAPI app entry point
│       └── routes/
│           ├── __init__.py
│           └── health.py    # Health check endpoint
├── tests/
│   └── __init__.py
├── pyproject.toml           # uv/pip configuration
├── Dockerfile
└── README.md
```

**Rationale:** Using `src/` layout is Python packaging best practice. It prevents accidental imports of the local package during development and makes the package structure explicit.

### Dependency Management: uv

Using `uv` for:
- Fast dependency resolution and installation
- `pyproject.toml` based configuration (PEP 621 compliant)
- Virtual environment management via `uv sync`
- Running commands via `uv run`

### Docker Configuration

**Base image:** `python:3.12-slim`
- Official image, well-maintained
- Slim variant reduces image size
- Python 3.12 matches project requirement

**Strategy:** Multi-stage build not needed initially (simple app), but Dockerfile structured to allow easy addition later.

**Port:** 2100 (as specified)

### Docker Compose

Single service initially:
- `api` service running the FastAPI application
- Volume mount for future SQLite database (`./data:/app/data`)
- Port mapping: `2100:2100`

### Makefile Commands

| Command | Description |
|---------|-------------|
| `make install` | Install dependencies with uv |
| `make run` | Run API locally |
| `make docker-build` | Build Docker image |
| `make docker-up` | Start with Docker Compose |
| `make docker-down` | Stop Docker Compose |
| `make test` | Run tests |
| `make lint` | Run ruff linter |
| `make format` | Format code with ruff |

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| uv is relatively new | uv is stable, actively maintained, and generates standard `pyproject.toml` |
| No database in initial setup | Explicit non-goal; separate change will add SQLite |

## Open Questions

None - all decisions confirmed with user.
