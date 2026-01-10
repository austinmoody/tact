# Capability: Database Configuration

Configures SQLite database location via environment variable with sensible defaults for local and containerized execution.

## ADDED Requirements

### Requirement: Database Path Configuration

The database file path MUST be configurable via environment variable with a fallback default.

#### Scenario: Custom path via environment variable

- Given: `TACT_DB_PATH` environment variable is set to `/custom/path/db.sqlite`
- When: The application starts
- Then: The database is created at `/custom/path/db.sqlite`

#### Scenario: Default path when not configured

- Given: `TACT_DB_PATH` environment variable is not set
- When: The application starts
- Then: The database is created at `./data/tact.db` relative to working directory

#### Scenario: Docker default path

- Given: The application runs in Docker via docker-compose
- Then: `TACT_DB_PATH` is set to `/data/tact.db`
- And: The `/data` directory is mounted as a volume for persistence

### Requirement: Directory Creation

The application SHALL create parent directories automatically if they don't exist.

#### Scenario: Parent directory does not exist

- Given: `TACT_DB_PATH` is set to `/data/subdir/tact.db`
- And: The `/data/subdir/` directory does not exist
- When: The application starts
- Then: The `/data/subdir/` directory is created
- And: The database file is created successfully

### Requirement: SQLite Connection Settings

The SQLite connection MUST be configured for FastAPI compatibility and data integrity.

#### Scenario: Connection configuration

- Given: The database connection is established
- Then: `check_same_thread=False` is set for async compatibility
- And: Foreign key constraints are enabled
