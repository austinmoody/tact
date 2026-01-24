## Context

TACT has a Python/FastAPI backend serving REST APIs, a Go-based TUI, and a Go-based MCP server. The web UI will be a new Go service that consumes the existing backend APIs, following the same pattern as the TUI.

**Current Architecture:**
```
backend/ (Python/FastAPI) - REST API on port 2100
tui/ (Go) - Terminal UI, calls backend API
mcp/ (Go) - MCP server for AI assistants, calls backend API
```

**Proposed Addition:**
```
webui/ (Go) - Web UI on port 2200, calls backend API
```

## Goals / Non-Goals

**Goals:**
- Full feature parity with TUI (entries, timers, projects, time codes, work types, context)
- Real-time timer updates without page refresh
- Responsive design that works on desktop and mobile browsers
- Simple development workflow (single `go run` command)
- Minimal JavaScript - rely on HTMX for interactivity

**Non-Goals:**
- User authentication (single-user app, same as TUI)
- Offline support / PWA capabilities
- Mobile native app
- Replacing the TUI (both will coexist)

## Decisions

### Decision 1: Go + Templ + HTMX Stack

**Choice**: Use Go with Templ templates and HTMX for interactivity.

**Rationale**:
- Same language as TUI allows code/pattern sharing
- Templ provides type-safe templates with Go syntax
- HTMX enables SPA-like UX with minimal JS (~14kb)
- No Node.js/npm build pipeline required

**Alternatives Considered**:
- Python + Jinja2 + HTMX: Would work but adds another language to maintain
- Go + html/template: Less ergonomic than Templ, no type safety
- Full SPA (React/Vue): Explicitly ruled out per requirements

### Decision 2: Server-Sent Events for Timer Updates

**Choice**: Use SSE for real-time timer display updates.

**Rationale**:
- Simpler than WebSockets for one-way server→client updates
- Native browser support, no JS library needed
- HTMX has built-in SSE support (`hx-ext="sse"`)
- Perfect for timer tick updates (1/second)

**Implementation**:
```go
// SSE endpoint streams timer state
GET /timer/stream
Content-Type: text/event-stream

event: tick
data: {"elapsed": "01:23:45", "running": true}
```

### Decision 3: Separate Service from Backend

**Choice**: Run web UI as a separate Go service on port 2200.

**Rationale**:
- Clean separation of concerns
- Can be deployed independently
- Follows existing pattern (TUI and MCP are also separate)
- Allows different scaling if needed

**Alternative Considered**:
- Embed in backend: Would mix Python and Go, complicate deployment

### Decision 4: CSS Framework - PicoCSS

**Choice**: Use PicoCSS for styling.

**Rationale**:
- Classless CSS - semantic HTML gets styled automatically
- Minimal size (~10kb)
- No build step required (just a CSS file)
- Provides dark mode support
- Clean, modern appearance

**Alternatives Considered**:
- TailwindCSS: Powerful but requires build step, utility classes everywhere
- No framework: Too much custom CSS to write
- Bootstrap: Heavy, dated appearance

### Decision 5: Project Structure

```
webui/
├── main.go              # Entry point, server setup
├── go.mod
├── handlers/            # HTTP handlers
│   ├── entries.go
│   ├── timer.go
│   ├── projects.go
│   ├── timecodes.go
│   └── worktypes.go
├── templates/           # Templ files
│   ├── layouts/
│   │   └── base.templ   # Main layout with nav
│   ├── pages/
│   │   ├── home.templ
│   │   ├── entries.templ
│   │   └── ...
│   └── components/
│       ├── entry_row.templ
│       ├── timer.templ
│       └── ...
├── api/                 # Backend API client (shared with TUI?)
│   └── client.go
├── static/
│   ├── css/
│   │   └── pico.min.css
│   └── htmx.min.js
└── Makefile
```

### Decision 6: API Client Sharing

**Choice**: Create new API client in webui/, consider extracting shared client later.

**Rationale**:
- TUI already has its own API client
- Web UI may need different error handling (HTTP responses vs TUI messages)
- Can refactor to shared module after initial implementation if warranted

## Risks / Trade-offs

**[Risk]** Templ is relatively new, smaller community than established frameworks
→ **Mitigation**: Templ is stable, well-documented, and backed by active development. Fallback to html/template if needed.

**[Risk]** SSE connection management for timers (disconnects, multiple tabs)
→ **Mitigation**: Implement reconnection logic in HTMX, timer state lives in backend so reconnects restore state.

**[Risk]** Two separate services to run during development
→ **Mitigation**: Provide `make dev` command that starts both, or document workflow clearly.

**[Trade-off]** PicoCSS limits customization compared to utility-first frameworks
→ **Accepted**: Clean defaults are sufficient for this app. Can add custom CSS for specific needs.

## Open Questions

1. **Shared API client**: Should we extract `tui/api/` to a shared Go module, or keep separate clients?
2. **Timer state**: Should web UI maintain its own timer state, or sync with TUI timer state via backend?
3. **Configuration**: Use same env vars as TUI (`TACT_API_URL`) or separate config?
