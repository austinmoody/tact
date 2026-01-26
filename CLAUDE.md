
## Critical Data Safety Rules

- **NEVER delete database files, data directories, or any files containing user data without EXPLICIT user approval.** This includes:
  - SQLite database files (*.db)
  - Data directories (data/, .data/, etc.)
  - Backup files
  - Any file that might contain user-entered information
- When a database migration fails or a database is in a bad state, **ASK the user** how they want to proceed. Suggest options like:
  - Creating a backup first
  - Attempting to fix the migration
  - Stamping the database to a known revision
- **Assume all data has value** - never assume a database is "just development data"

## Git Workflow

- **Do not commit or push without explicit user approval.** After completing code changes, wait for the user to request a commit before staging, committing, or pushing.
