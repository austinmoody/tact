## 1. Project Structure

- [x] 1.1 Create `backend/src/tact/` directory structure
- [x] 1.2 Create `backend/tests/` directory with `__init__.py`
- [x] 1.3 Create `backend/src/tact/__init__.py` with version info

## 2. Dependency Configuration

- [x] 2.1 Create `backend/pyproject.toml` with uv configuration
- [x] 2.2 Define dependencies: fastapi, uvicorn[standard]
- [x] 2.3 Define dev dependencies: pytest, ruff
- [x] 2.4 Run `uv sync` to generate lock file

## 3. FastAPI Application

- [x] 3.1 Create `backend/src/tact/main.py` with FastAPI app instance
- [x] 3.2 Create `backend/src/tact/routes/__init__.py`
- [x] 3.3 Create `backend/src/tact/routes/health.py` with GET /health endpoint
- [x] 3.4 Wire up health router in main.py

## 4. Docker Configuration

- [x] 4.1 Create `backend/Dockerfile` with python:3.12-slim base
- [x] 4.2 Create `docker-compose.yml` in project root
- [x] 4.3 Configure volume mount for data directory
- [x] 4.4 Configure port mapping (2100:2100)

## 5. Development Tooling

- [x] 5.1 Create `Makefile` in project root with commands
- [x] 5.2 Create `backend/README.md` with setup instructions
- [x] 5.3 Add `.gitignore` entries for Python/uv artifacts

## 6. Validation

- [x] 6.1 Verify `make install` works
- [x] 6.2 Verify `make run` starts API on port 2100
- [x] 6.3 Verify `curl localhost:2100/health` returns success
- [x] 6.4 Verify `make docker-up` runs containerized API
- [x] 6.5 Verify health endpoint works in Docker
