## 1. Python CI Workflow

- [x] 1.1 Create `.github/workflows/ci-python.yml` workflow file
- [x] 1.2 Configure trigger on push to main/develop and pull requests
- [x] 1.3 Set up Python 3.12 with `actions/setup-python@v5`
- [x] 1.4 Add pip caching with `actions/cache@v4`
- [x] 1.5 Install backend dependencies and run pytest
- [x] 1.6 Install backend dependencies and run ruff (lint + format check)
- [x] 1.7 Install mcp dependencies and run ruff (lint + format check)

## 2. Go CI Workflow

- [x] 2.1 Create `.github/workflows/ci-go.yml` workflow file
- [x] 2.2 Configure trigger on push to main/develop and pull requests
- [x] 2.3 Set up Go 1.24 with `actions/setup-go@v5`
- [x] 2.4 Enable Go module caching (built into setup-go)
- [x] 2.5 Run `go build` to verify compilation
- [x] 2.6 Run `go test ./...` for tests
- [x] 2.7 Run `go vet ./...` for static analysis
- [x] 2.8 Install and run `staticcheck`

## 3. Docker Build Workflow

- [x] 3.1 Create `.github/workflows/docker.yml` workflow file
- [x] 3.2 Configure trigger on GitHub release published
- [x] 3.3 Set up Docker Buildx with `docker/setup-buildx-action@v3`
- [x] 3.4 Configure GHCR login with `docker/login-action@v3`
- [x] 3.5 Build and push backend image with version and latest tags
- [x] 3.6 Build and push mcp image with version and latest tags
- [x] 3.7 Configure multi-platform builds (linux/amd64, linux/arm64)
- [x] 3.8 Enable Docker layer caching

## 4. Release Automation

- [x] 4.1 Create `.github/workflows/release.yml` workflow file
- [x] 4.2 Configure `google-github-actions/release-please-action@v4`
- [x] 4.3 Set release type to `simple` (or `node` if using package.json)
- [x] 4.4 Configure to create releases on main branch

## 5. Verification

- [x] 5.1 Push a test commit to verify Python CI runs
- [x] 5.2 Push a test commit to verify Go CI runs
- [x] 5.3 Verify release-please creates a release PR (requires merge to main)
- [x] 5.4 Test Docker build workflow (requires a release)
