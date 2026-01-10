# Change: Setup Backend Environment

## Why

The project needs a foundational backend structure before any features can be implemented. This establishes the development environment, Docker configuration, and basic project scaffolding that all future work will build upon.

## What Changes

- Create `backend/` directory with standard Python package structure (`backend/src/tact/`)
- Configure `uv` for dependency management with `pyproject.toml`
- Set up FastAPI application with a single `/health` endpoint
- Create Docker and Docker Compose configuration for containerized execution
- Add Makefile with common development commands
- Document setup and usage in README

## Impact

- Affected specs: New `health-check` capability
- Affected code: Creates new `backend/` directory (no existing code affected)
- This is the foundation for all future backend development
