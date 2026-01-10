## 1. Dependencies

- [ ] 1.1 Add `sqlalchemy>=2.0` to pyproject.toml dependencies
- [ ] 1.2 Add `alembic>=1.13` to pyproject.toml dependencies
- [ ] 1.3 Run `uv sync --extra dev` to update lock file

## 2. Database Module

- [ ] 2.1 Create `backend/src/tact/db/__init__.py`
- [ ] 2.2 Create `backend/src/tact/db/base.py` with engine and Base setup
- [ ] 2.3 Create `backend/src/tact/db/models.py` with all four models
- [ ] 2.4 Create `backend/src/tact/db/session.py` with session factory

## 3. Alembic Setup

- [ ] 3.1 Initialize Alembic in backend directory
- [ ] 3.2 Configure `alembic.ini` for SQLite and project structure
- [ ] 3.3 Update `alembic/env.py` to use app models and config
- [ ] 3.4 Create initial migration `001_initial_schema.py`

## 4. Auto-Migration

- [ ] 4.1 Create `backend/src/tact/db/migrations.py` with run_migrations function
- [ ] 4.2 Add lifespan handler to FastAPI app for startup migration
- [ ] 4.3 Add `TACT_DB_PATH` environment variable handling with default

## 5. Docker Configuration

- [ ] 5.1 Update `docker-compose.yml` to set `TACT_DB_PATH=/data/tact.db`
- [ ] 5.2 Ensure `/data` volume mount exists (already configured)

## 6. Developer Tooling

- [ ] 6.1 Add `make migrate` command for manual migration runs
- [ ] 6.2 Add `make db-revision` command for creating new migrations
- [ ] 6.3 Update `backend/README.md` with database documentation

## 7. Testing

- [ ] 7.1 Create test fixture for in-memory database
- [ ] 7.2 Add test for database initialization
- [ ] 7.3 Add test for each model can be created and queried
- [ ] 7.4 Verify `make test` passes
