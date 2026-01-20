## Why

The project currently has no automated CI/CD pipeline. Without automation, code quality checks (tests, linting) must be run manually, Docker images must be built and pushed by hand, and releases require manual version management. This creates friction and risk of shipping broken code.

## What Changes

- Add GitHub Actions workflow for CI pipeline (run tests, linting, type checking on PRs and pushes)
- Add GitHub Actions workflow for Docker image builds (push to GitHub Container Registry on main branch)
- Add GitHub Actions workflow for release automation (semantic versioning, changelog generation, GitHub releases)
- Add workflow for the Go TUI component

## Capabilities

### New Capabilities

- `github-actions`: GitHub Actions workflows for CI/CD automation including:
  - CI pipeline for Python backend (pytest, ruff)
  - CI pipeline for Go TUI (go test, go vet, staticcheck)
  - Docker image builds for backend and MCP server
  - Release automation with semantic versioning

### Modified Capabilities

(none)

## Impact

- **New files**: `.github/workflows/*.yml` workflow files
- **Dependencies**: None on application code
- **Services**: GitHub Container Registry for Docker images
- **Components affected**:
  - `backend/` - Python 3.12, pytest, ruff
  - `mcp/` - Python 3.12, ruff (no tests currently)
  - `tui/` - Go 1.24, Charm libraries
