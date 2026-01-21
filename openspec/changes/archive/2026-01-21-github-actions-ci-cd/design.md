## Context

The tact project consists of three components:
- **backend/** - Python 3.12 FastAPI application with pytest tests and ruff linting
- **mcp/** - Python 3.12 MCP server (no tests, uses ruff)
- **tui/** - Go 1.24 terminal UI using Charm libraries

Currently there are no GitHub Actions workflows. The `.github/agents/` directory contains OpenSpec agent definitions but no CI/CD automation.

## Goals / Non-Goals

**Goals:**
- Automated CI on all PRs and pushes to main branches
- Docker images automatically built and pushed to GHCR on releases
- Semantic versioning with automated changelog generation
- Fast feedback loops (parallel jobs, caching)

**Non-Goals:**
- Deployment automation (this is a local-first tool)
- Complex release branching strategies
- Publishing to PyPI or other package registries
- macOS app builds (handled separately)

## Decisions

### 1. Workflow Structure

**Decision:** Separate workflows for CI, Docker builds, and releases.

**Rationale:** Separation of concerns allows each workflow to have appropriate triggers and independent failure modes. CI runs on every push/PR, Docker builds on releases only.

**Alternatives considered:**
- Single monolithic workflow: Rejected because it would run Docker builds unnecessarily on PRs
- Matrix builds for all components: Rejected because Python and Go have different tooling

### 2. CI Triggers

**Decision:** Run CI on:
- Push to `main` and `develop` branches
- All pull requests

**Rationale:** Catches issues before merge while avoiding noise on feature branches.

### 3. Docker Registry

**Decision:** Use GitHub Container Registry (ghcr.io).

**Rationale:** Native GitHub integration, no additional credentials to manage, free for public repos.

**Alternatives considered:**
- Docker Hub: Requires separate account and credentials
- Self-hosted registry: Unnecessary complexity

### 4. Release Strategy

**Decision:** Use `release-please` for automated releases with conventional commits.

**Rationale:** Automatically generates changelogs from commit messages, creates release PRs, and handles version bumping. Works well with the existing conventional commit style (e.g., `feat(tui):`, `fix(parser):`).

**Alternatives considered:**
- Manual releases: Error-prone and tedious
- semantic-release: More complex setup, less transparent
- Manual version tags: Loses automatic changelog generation

### 5. Caching Strategy

**Decision:** Use GitHub Actions cache for:
- Python pip dependencies
- Go modules
- Docker layer caching

**Rationale:** Significantly reduces CI time on repeat runs.

### 6. Go CI Tools

**Decision:** Use `go vet` and `staticcheck` for Go linting.

**Rationale:** `go vet` catches common mistakes, `staticcheck` provides deeper analysis. Both are standard in the Go ecosystem.

**Alternatives considered:**
- golangci-lint: More comprehensive but heavier, overkill for a small TUI

## Risks / Trade-offs

**Risk:** Release-please may create unexpected version bumps.
→ Mitigation: Review release PRs before merging. Use conventional commit prefixes carefully.

**Risk:** Docker builds on every release may be slow.
→ Mitigation: Use Docker layer caching. Only build on actual releases, not PRs.

**Risk:** Go 1.24 is very new, may have CI action compatibility issues.
→ Mitigation: Use `setup-go@v5` which supports Go 1.24. Pin version explicitly.

**Trade-off:** Running CI on all PRs increases GitHub Actions minutes usage.
→ Acceptable for a personal project. Can add path filters later if needed.
