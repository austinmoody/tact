# Change: Setup Data Model

## Why

The backend needs a database layer before any features can be implemented. This establishes the SQLite schema, migration infrastructure, and automatic initialization so the app is ready to persist time entries, time codes, work types, and configuration.

## What Changes

- Add Alembic for database migrations with auto-run on startup
- Create initial migration with all four tables: `time_entries`, `time_codes`, `work_types`, `config`
- Configure database path via `TACT_DB_PATH` environment variable with fallback default
- Add SQLAlchemy models matching the schema defined in project.md
- Integrate database initialization into FastAPI startup lifecycle

## Impact

- Affected specs: New `database-migrations` and `database-config` capabilities
- Affected code: Adds `backend/src/tact/db/` module, Alembic configuration, SQLAlchemy models
- Foundation for all data persistence in future features
