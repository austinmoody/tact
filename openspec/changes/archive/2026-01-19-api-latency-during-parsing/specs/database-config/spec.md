## MODIFIED Requirements

### Requirement: SQLite Connection Settings

The SQLite connection MUST be configured for FastAPI compatibility, data integrity, and concurrent access.

#### Scenario: Connection configuration

- Given: The database connection is established
- Then: `check_same_thread=False` is set for async compatibility
- And: Foreign key constraints are enabled
- And: WAL journal mode is enabled via `PRAGMA journal_mode=WAL`
- And: Busy timeout is set via `PRAGMA busy_timeout=5000`
