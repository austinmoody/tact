## Why

TACT currently offers a TUI for desktop terminal users and an MCP server for AI assistant integration, but lacks a web-based interface for browser access. A web UI would enable access from any device with a browser, provide a more visual experience for managing time entries, and support real-time timer functionality without requiring terminal access.

## What Changes

- Add new `webui/` directory alongside existing `backend/`, `tui/`, and `mcp/` directories
- Implement Go-based web server using Templ for type-safe HTML templating
- Use HTMX for dynamic interactions without heavy JavaScript frameworks
- Implement Server-Sent Events (SSE) for real-time timer updates
- Full feature parity with TUI: entries, timers, projects, time codes, work types, context management

## Capabilities

### New Capabilities

- `web-ui-core`: Core web application structure, routing, authentication (if needed), and shared layouts/components using Go + Templ + HTMX
- `web-ui-entries`: Web interface for listing, viewing, creating, and editing time entries with status display and filtering
- `web-ui-timer`: Real-time timer functionality with start/stop/pause controls, elapsed time display via SSE, and timer-to-entry conversion
- `web-ui-management`: Web interfaces for managing projects, time codes, work types, and their associated context documents

### Modified Capabilities

None - this is a new frontend that consumes existing backend APIs.

## Impact

- **New Code**: `webui/` directory with Go web server, Templ templates, static assets
- **Dependencies**: Go modules for Templ, PicoCSS for styling, HTMX for interactivity
- **Infrastructure**: New service to run on port 2200 (separate from backend API on 2100)
- **Build**: Templ compilation step via `templ generate`
- **Documentation**: Setup instructions, development workflow

## DevOps

### Makefile Targets

Added to root Makefile:
- `make webui-generate` - Generate templ files
- `make webui-build` - Build webui binary (runs templ generate first)
- `make webui-run` - Run webui server (port 2200)
- `make webui-dev` - Run webui in dev mode (localhost:2100 API)

### Docker

- `webui/Dockerfile` - Multi-stage build using Go 1.24-alpine
- Added `webui` service to `docker-compose.yml` with:
  - Port mapping: 2200:2200
  - Environment: `TACT_API_URL=http://api:2100`
  - Depends on: api service

### GitHub Actions

Updated `.github/workflows/ci-go.yml`:
- Added `build-and-test-webui` job (parallel to TUI job)
- Installs templ, generates files, builds, tests, vets, runs staticcheck

Updated `.github/workflows/docker.yml`:
- Added `build-webui` job for Docker image publishing
- Publishes to `ghcr.io/<repo>/webui` on releases
- Multi-platform: linux/amd64, linux/arm64
