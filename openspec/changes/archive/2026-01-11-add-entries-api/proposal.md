# Proposal: Add Entries API

## Summary

Add CRUD API endpoints for time entries. Entries are saved with `status="pending"` and remain uncategorized until LLM parsing is implemented in a future change. This establishes the core entry management functionality independent of the parsing system.

## Motivation

- Enable users to submit and manage time entries before LLM integration is ready
- Entries accumulate with `raw_text` preserved for future parsing
- Provides the foundation for the primary user workflow (quick entry submission)

## Scope

### In Scope

- `POST /entries` - Create entry (raw_text required, optional entry_date defaults to today)
- `GET /entries` - List entries with filters (status, time_code_id, work_type_id, date range, pagination)
- `GET /entries/{id}` - Get single entry
- `PATCH /entries/{id}` - Update/correct entry fields
- `DELETE /entries/{id}` - Hard delete entry

### Out of Scope

- `POST /entries/{id}/reparse` - LLM-related, future change
- `POST /entries/reparse` - LLM-related, future change
- Background parsing worker - future change
- Confidence score calculation - set by future LLM parsing

## Key Decisions

1. **Entry date defaults to today** - If not provided in request, `entry_date` is set to current date
2. **Status starts as pending** - All new entries have `status="pending"` until parsed
3. **Hard delete** - Unlike time codes/work types, entries are hard deleted (not soft-deleted)
4. **PATCH for updates** - Partial updates only; changed fields trigger `manually_corrected=true`
5. **No validation of FK references** - work_type_id and time_code_id are optional and not validated on create (will be set by parser)

## Affected Specs

- `entries-api` (NEW) - CRUD operations for time entries
