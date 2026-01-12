# api-logging Specification

## Purpose

Provides structured logging for all API operations, enabling visibility into system behavior for troubleshooting and monitoring. Logs are output to stdout in a consistent format with timestamps, log levels, and module names for container compatibility.
## Requirements
### Requirement: Structured Log Format

The API SHALL output logs in a consistent structured format including timestamp, log level, module name, and message.

#### Scenario: Log entry format

- **WHEN** any log entry is written
- **THEN** the log includes ISO-8601 timestamp, log level (INFO/WARNING/ERROR), module name, and message
- **AND** logs are written to stdout for container compatibility

### Requirement: Request Logging

The API SHALL log all incoming HTTP requests with method, path, response status, and timing.

#### Scenario: Successful request logging

- **WHEN** an HTTP request is processed successfully
- **THEN** a log entry is written with method, path, status code, and duration in milliseconds

#### Scenario: Failed request logging

- **WHEN** an HTTP request fails with an error
- **THEN** a log entry is written with method, path, error status code, and duration

### Requirement: Entry Operation Logging

The API SHALL log all entry CRUD operations with relevant context.

#### Scenario: Entry creation logged

- **WHEN** POST `/entries` creates an entry
- **THEN** a log entry is written with the new entry ID and raw_text preview

#### Scenario: Entry list logged

- **WHEN** GET `/entries` is called
- **THEN** a log entry is written with filter parameters and result count

#### Scenario: Entry update logged

- **WHEN** PATCH `/entries/{id}` updates an entry
- **THEN** a log entry is written with entry ID and updated field names

#### Scenario: Entry deletion logged

- **WHEN** DELETE `/entries/{id}` removes an entry
- **THEN** a log entry is written with the deleted entry ID

#### Scenario: Entry reparse logged

- **WHEN** POST `/entries/{id}/reparse` is called
- **THEN** a log entry is written with the entry ID being reparsed

### Requirement: Work Type Operation Logging

The API SHALL log all work type CRUD operations.

#### Scenario: Work type create/update/delete logged

- **WHEN** a work type is created, updated, or deleted
- **THEN** a log entry is written with the operation type and work type ID

### Requirement: Time Code Operation Logging

The API SHALL log all time code CRUD operations.

#### Scenario: Time code create/update/delete logged

- **WHEN** a time code is created, updated, or deleted
- **THEN** a log entry is written with the operation type and time code ID

