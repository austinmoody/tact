# database-concurrency Specification

## Purpose

Configuration and patterns for SQLite concurrent access including WAL mode, busy timeouts, and transaction scoping best practices.

## Requirements

### Requirement: WAL Mode

The system SHALL use SQLite Write-Ahead Logging (WAL) mode for improved concurrent access.

#### Scenario: WAL mode enabled on connection

- **WHEN** a database connection is established
- **THEN** the connection executes `PRAGMA journal_mode=WAL`
- **AND** concurrent readers are not blocked by writers

#### Scenario: Reads during active write transaction

- **WHEN** the parser worker is writing entry results
- **AND** an API request attempts to read entries
- **THEN** the read succeeds without blocking
- **AND** the read sees a consistent snapshot of data

### Requirement: Busy Timeout

The system SHALL configure SQLite busy timeout to handle lock contention gracefully.

#### Scenario: Busy timeout configured

- **WHEN** a database connection is established
- **THEN** the connection executes `PRAGMA busy_timeout=5000`

#### Scenario: Brief lock contention

- **WHEN** a write operation encounters a locked database
- **AND** the lock is released within 5 seconds
- **THEN** the write operation succeeds after waiting
- **AND** no error is raised

#### Scenario: Extended lock contention

- **WHEN** a write operation encounters a locked database
- **AND** the lock is not released within 5 seconds
- **THEN** SQLite raises a "database is locked" error
- **AND** the application handles the error appropriately

### Requirement: Minimal Lock Duration

Database transactions SHALL be kept as short as possible, especially during long-running operations.

#### Scenario: LLM calls outside transaction

- **WHEN** the parser worker processes an entry
- **THEN** the database session is closed before calling the LLM
- **AND** a new session is opened only to write results
- **AND** the total lock duration is under 1 second per entry

#### Scenario: Entry state validation before write

- **WHEN** the parser worker attempts to write parse results
- **AND** the entry was modified or deleted during LLM processing
- **THEN** the write is skipped
- **AND** no error is raised
- **AND** a warning is logged
