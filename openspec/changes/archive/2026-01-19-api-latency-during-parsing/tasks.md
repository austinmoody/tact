## 1. Database Configuration

- [x] 1.1 Add WAL mode pragma to `_set_sqlite_pragma` in `backend/src/tact/db/base.py`
- [x] 1.2 Add busy timeout pragma (5000ms) to `_set_sqlite_pragma`
- [x] 1.3 Verify pragmas are applied by adding debug logging on connection

## 2. Refactor EntryParser

- [x] 2.1 Add `build_parse_context(user_input, session)` method that fetches RAG context, time codes, and work types
- [x] 2.2 Add `parse_text(user_input, context)` method that calls LLM provider without DB access
- [x] 2.3 Keep existing `parse_entry()` method working for backwards compatibility (can be deprecated later)

## 3. Restructure Parser Worker

- [x] 3.1 Refactor `process_pending_entries` to use three-phase approach:
  - Phase 1: Fetch pending entries and build contexts, then close session
  - Phase 2: Call LLM for each entry (no DB connection held)
  - Phase 3: Write results with new sessions, checking entry still pending
- [x] 3.2 Add optimistic concurrency check before writing (verify entry exists and status == "pending")
- [x] 3.3 Add warning logging when entry was modified/deleted during parsing
- [x] 3.4 Update `SessionFactory` usage to use context manager pattern (`with SessionFactory() as session:`)

## 4. Testing

- [x] 4.1 Test API responsiveness during active parsing (manual or integration test)
- [x] 4.2 Test race condition: modify entry while parsing is in progress
- [x] 4.3 Test race condition: delete entry while parsing is in progress
- [x] 4.4 Verify WAL mode is active (`PRAGMA journal_mode` returns "wal")
- [x] 4.5 Verify existing parse behavior is unchanged (entries still parse correctly)
