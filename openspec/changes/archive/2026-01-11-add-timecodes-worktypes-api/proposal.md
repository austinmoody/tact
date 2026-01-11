# Change: Add Time Codes and Work Types API

## Why

The database models exist but there's no way to manage time codes or work types. These are foundational resources needed before implementing time entries. Adding these CRUD endpoints establishes the API patterns for the rest of the application.

## What Changes

- Add Time Codes API: POST, GET (list), GET (single), PUT, DELETE
- Add Work Types API: POST, GET (list), GET (single), PUT, DELETE
- Add Pydantic schemas for request/response validation
- Direct SQLAlchemy queries in route handlers (simple, no repository layer)
- Soft-delete behavior (set `active=false` instead of hard delete)

## Impact

- Affected specs: New `time-codes-api` and `work-types-api` capabilities
- Affected code: Adds routes, schemas, and tests
- Establishes CRUD patterns for future endpoints (entries, config)
