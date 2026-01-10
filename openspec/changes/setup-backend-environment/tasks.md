## 1. Project Structure

- [ ] 1.1 Create `backend/src/tact/` directory structure
- [ ] 1.2 Create `backend/tests/` directory with `__init__.py`
- [ ] 1.3 Create `backend/src/tact/__init__.py` with version info

## 2. Dependency Configuration

- [ ] 2.1 Create `backend/pyproject.toml` with uv configuration
- [ ] 2.2 Define dependencies: fastapi, uvicorn[standard]
- [ ] 2.3 Define dev dependencies: pytest, ruff
- [ ] 2.4 Run `uv sync` to generate lock file

## 3. FastAPI Application

- [ ] 3.1 Create `backend/src/tact/main.py` with FastAPI app instance
- [ ] 3.2 Create `backend/src/tact/routes/__init__.py`
- [ ] 3.3 Create `backend/src/tact/routes/health.py` with GET /health endpoint
- [ ] 3.4 Wire up health router in main.py

## 4. Docker Configuration

- [ ] 4.1 Create `backend/Dockerfile` with python:3.12-slim base
- [ ] 4.2 Create `docker-compose.yml` in project root
- [ ] 4.3 Configure volume mount for data directory
- [ ] 4.4 Configure port mapping (2100:2100)

## 5. Development Tooling

- [ ] 5.1 Create `Makefile` in project root with commands
- [ ] 5.2 Create `backend/README.md` with setup instructions
- [ ] 5.3 Add `.gitignore` entries for Python/uv artifacts

## 6. Validation

- [ ] 6.1 Verify `make install` works
- [ ] 6.2 Verify `make run` starts API on port 2100
- [ ] 6.3 Verify `curl localhost:2100/health` returns success
- [ ] 6.4 Verify `make docker-up` runs containerized API
- [ ] 6.5 Verify health endpoint works in Docker
