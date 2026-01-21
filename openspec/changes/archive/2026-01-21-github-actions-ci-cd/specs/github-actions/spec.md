## ADDED Requirements

### Requirement: Python CI Workflow

The system SHALL run CI checks on Python components (backend, mcp) for every pull request and push to main/develop branches.

#### Scenario: PR triggers CI

- **WHEN** a pull request is opened or updated
- **THEN** the Python CI workflow runs
- **AND** pytest runs on the backend component
- **AND** ruff lint and format checks run on backend and mcp

#### Scenario: Push to main triggers CI

- **WHEN** code is pushed to the main or develop branch
- **THEN** the Python CI workflow runs

#### Scenario: CI uses Python 3.12

- **WHEN** the Python CI workflow runs
- **THEN** Python 3.12 is used for all checks

#### Scenario: Dependencies are cached

- **WHEN** the Python CI workflow runs
- **THEN** pip dependencies are cached between runs
- **AND** subsequent runs are faster

### Requirement: Go CI Workflow

The system SHALL run CI checks on the Go TUI component for every pull request and push to main/develop branches.

#### Scenario: PR triggers Go CI

- **WHEN** a pull request is opened or updated
- **THEN** the Go CI workflow runs
- **AND** go build verifies compilation
- **AND** go test runs all tests
- **AND** go vet checks for common issues
- **AND** staticcheck runs static analysis

#### Scenario: Go modules are cached

- **WHEN** the Go CI workflow runs
- **THEN** Go module dependencies are cached between runs

#### Scenario: CI uses Go 1.24

- **WHEN** the Go CI workflow runs
- **THEN** Go 1.24 is used

### Requirement: Docker Build Workflow

The system SHALL build and push Docker images to GitHub Container Registry on releases.

#### Scenario: Release triggers Docker build

- **WHEN** a GitHub release is published
- **THEN** Docker images are built for backend and mcp components
- **AND** images are pushed to ghcr.io

#### Scenario: Images are tagged with version

- **WHEN** a release v1.2.3 is published
- **THEN** Docker images are tagged with `1.2.3`
- **AND** Docker images are tagged with `latest`

#### Scenario: Docker layer caching

- **WHEN** Docker images are built
- **THEN** Docker layer caching is used to speed up builds

#### Scenario: Multi-platform builds

- **WHEN** Docker images are built
- **THEN** images are built for linux/amd64 and linux/arm64

### Requirement: Release Automation Workflow

The system SHALL automate releases using release-please with conventional commits.

#### Scenario: Conventional commit triggers release PR

- **WHEN** commits with conventional prefixes (feat:, fix:) are pushed to main
- **THEN** release-please creates or updates a release PR
- **AND** the PR contains a generated changelog

#### Scenario: Merging release PR creates release

- **WHEN** a release-please PR is merged
- **THEN** a GitHub release is created
- **AND** the release is tagged with the new version
- **AND** the Docker build workflow is triggered

#### Scenario: Version follows semver

- **WHEN** a release is created
- **THEN** the version follows semantic versioning
- **AND** feat: commits bump minor version
- **AND** fix: commits bump patch version
- **AND** BREAKING CHANGE commits bump major version

### Requirement: CI Status Checks

The system SHALL enforce CI checks before merging pull requests.

#### Scenario: Required status checks

- **WHEN** a pull request is opened
- **THEN** the PR cannot be merged until Python CI passes
- **AND** the PR cannot be merged until Go CI passes

#### Scenario: CI failure blocks merge

- **WHEN** any CI check fails
- **THEN** the pull request shows a failed status
- **AND** merge is blocked (if branch protection is enabled)
