# Change: Add Structured API Logging

## Why

The API currently provides minimal visibility into operations. When running in Docker or elsewhere, there are no log entries for CRUD operations on entries, work types, or time codes. This makes troubleshooting issues difficult since there's no indication of what's happening in the system.

**Current state:**
- API routes for entries, work types, time codes, and health check have zero logging
- Only the parser worker and LLM providers log their operations
- No request/response logging middleware
- No structured log format (timestamps, request IDs, etc.)

## What Changes

- Add structured logging to all API route operations (entries, work types, time codes)
- Add request logging middleware to capture request method, path, status code, and timing
- Configure consistent log format with timestamps, log levels, and module names
- Use correlation IDs to trace requests across the system

## Impact

- Affected specs: New `api-logging` capability
- Affected code:
  - `backend/src/tact/main.py` - logging configuration
  - `backend/src/tact/routes/entries.py` - add logging to all operations
  - `backend/src/tact/routes/work_types.py` - add logging to all operations
  - `backend/src/tact/routes/time_codes.py` - add logging to all operations
