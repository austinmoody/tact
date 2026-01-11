## 1. Dependencies

- [x] 1.1 Add `sqlalchemy>=2.0` to pyproject.toml dependencies
- [x] 1.2 Add `alembic>=1.13` to pyproject.toml dependencies
- [x] 1.3 Run `uv sync --extra dev` to update lock file

## 2. Database Module

- [x] 2.1 Create `backend/src/tact/db/__init__.py`
- [x] 2.2 Create `backend/src/tact/db/base.py` with engine and Base setup
- [x] 2.3 Create `backend/src/tact/db/models.py` with all four models
- [x] 2.4 Create `backend/src/tact/db/session.py` with session factory

## 3. Alembic Setup

- [x] 3.1 Initialize Alembic in backend directory
- [x] 3.2 Configure `alembic.ini` for SQLite and project structure
- [x] 3.3 Update `alembic/env.py` to use app models and config
- [x] 3.4 Create initial migration `001_initial_schema.py`

## 4. Auto-Migration

- [x] 4.1 Create `backend/src/tact/db/migrations.py` with run_migrations function
- [x] 4.2 Add lifespan handler to FastAPI app for startup migration
- [x] 4.3 Add `TACT_DB_PATH` environment variable handling with default

## 5. Docker Configuration

- [x] 5.1 Update `docker-compose.yml` to set `TACT_DB_PATH=/data/tact.db`
- [x] 5.2 Ensure `/data` volume mount exists (already configured)

## 6. Developer Tooling

- [x] 6.1 Add `make migrate` command for manual migration runs
- [x] 6.2 Add `make db-revision` command for creating new migrations
- [x] 6.3 Update `backend/README.md` with database documentation

## 7. Testing

- [x] 7.1 Create test fixture for in-memory database
- [x] 7.2 Add test for database initialization
- [x] 7.3 Add test for each model can be created and queried
- [x] 7.4 Verify `make test` passes

## Verification

Steps to independently verify the implementation:

1. **Start the app and confirm migrations run:**
   ```bash
   make run
   ```
   Look for migration log output before the server starts.

2. **Verify database file was created:**
   ```bash
   ls -la data/tact.db
   ```

3. **Check all tables exist:**
   ```bash
   sqlite3 data/tact.db ".tables"
   ```
   Expected: `alembic_version  config  time_codes  time_entries  work_types`

4. **Verify the API still works:**
   ```bash
   curl http://localhost:2100/health
   ```
   Expected: `{"status":"healthy"}`

5. **Test manual migration command:**
   ```bash
   make migrate
   ```
   Should complete without errors.
