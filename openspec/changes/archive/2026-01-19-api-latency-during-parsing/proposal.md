## Why

The API becomes unresponsive during parse/reparse operations because SQLite's database-level locking blocks all other requests while the parser holds a transaction open during long-running LLM calls (up to 180-600 seconds). This creates a poor user experience where listing entries, fetching projects, or any other API call hangs until parsing completes.

## What Changes

- **Restructure parser worker to minimize DB lock duration**: Move LLM calls outside database transactions so the lock is only held during the brief read/write operations, not during the entire LLM processing time.
- **Enable SQLite WAL mode**: Switch from rollback journal to Write-Ahead Logging to allow concurrent reads during writes.
- **Add busy timeout configuration**: Configure SQLite's busy timeout to handle lock contention gracefully instead of failing immediately.
- **Make embedding generation non-blocking for API routes**: Defer embedding generation so context creation doesn't block API responses.

## Capabilities

### New Capabilities

- `database-concurrency`: Configuration and patterns for SQLite concurrent access including WAL mode, busy timeouts, and transaction scoping best practices.

### Modified Capabilities

- `llm-parsing`: Restructure the background worker to release database locks before calling the LLM. The worker will: (1) fetch pending entries and release lock, (2) call LLM without holding lock, (3) acquire lock only to write results. This changes internal implementation while preserving all external behavior.
- `database-config`: Add WAL mode and busy timeout configuration to improve concurrent access.

## Impact

**Code affected:**
- `backend/src/tact/worker/parser_worker.py` - Restructure to minimize lock duration
- `backend/src/tact/db/base.py` - Add WAL mode and busy timeout configuration
- `backend/src/tact/llm/parser.py` - Separate data fetching from LLM calls
- `backend/src/tact/routes/context.py` - Consider async embedding or background processing

**APIs:**
- No external API changes - all changes are internal to improve concurrency

**Dependencies:**
- No new dependencies required

**Risk areas:**
- Must ensure data consistency when entry state can change between read and write
- Need to handle race conditions if an entry is deleted or modified during parsing
