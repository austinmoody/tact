# database-migrations Specification

## Purpose
TBD - created by archiving change setup-data-model. Update Purpose after archive.
## Requirements
### Requirement: Automatic Migration Execution

The application SHALL automatically run pending database migrations when starting.

#### Scenario: Fresh database initialization

- Given: No database file exists at the configured path
- When: The application starts
- Then: Parent directories are created if needed
- And: All migrations are applied
- And: The database is ready for use

#### Scenario: Existing database with pending migrations

- Given: A database file exists with some migrations applied
- When: The application starts
- Then: Only pending migrations are applied
- And: Previously applied migrations are skipped

#### Scenario: Database is up to date

- Given: A database file exists with all migrations applied
- When: The application starts
- Then: No migrations are executed
- And: The application starts normally

### Requirement: Initial Schema

The initial migration MUST create all core tables required by the application.

#### Scenario: Tables created by initial migration

- Given: The initial migration runs
- Then: The `time_entries` table exists with all columns from TimeEntry model
- And: The `time_codes` table exists with all columns from TimeCode model
- And: The `work_types` table exists with all columns from WorkType model
- And: The `config` table exists with all columns from Config model
- And: Foreign key constraints are enforced

### Requirement: Manual Migration Commands

The project SHALL provide Makefile commands for manual migration management.

#### Scenario: Run pending migrations

- Given: The developer runs `make migrate`
- Then: All pending migrations are applied to the configured database

#### Scenario: Create new migration

- Given: The developer runs `make db-revision` with a message
- Then: A new migration file is created in `alembic/versions/`

