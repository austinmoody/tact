## Context

The tact backend uses SQLite for persistence and a background worker for LLM-based parsing of time entries. Currently, the parser worker holds a database session open throughout the entire parsing operation, including during LLM API calls that can take 2-180+ seconds.

SQLite uses database-level write locking - when any connection holds a write lock, all other write operations (and reads in default mode) block until the lock is released. This causes the API to become unresponsive during parsing.

**Current flow in `parser_worker.py`:**
```
session = SessionFactory()           # Opens connection
pending = session.query(...).all()  # DB read
for entry in pending:
    parser.parse_entry(entry, session)  # Includes LLM call (2-180s) - HOLDS LOCK
    session.commit()                     # Releases lock
session.close()
```

The LLM call inside `parse_entry` includes:
- RAG context retrieval (DB queries)
- Building parse context (DB queries)
- Calling the LLM provider (network I/O - the slow part)
- Updating entry fields (in-memory)

## Goals / Non-Goals

**Goals:**
- API remains responsive during parse operations (< 500ms response time for reads)
- Parser worker continues to function correctly with same external behavior
- SQLite remains the database (no migration to PostgreSQL)
- Minimize code changes while solving the blocking issue

**Non-Goals:**
- Full async rewrite of the application
- Moving to a different database
- Adding a job queue system (Celery, etc.)
- Real-time parsing feedback to users

## Decisions

### Decision 1: Enable WAL Mode for SQLite

**Choice:** Enable Write-Ahead Logging (WAL) mode via PRAGMA.

**Rationale:** WAL mode allows concurrent readers during writes. Readers see a consistent snapshot while a write is in progress. This is a simple configuration change that immediately improves concurrency.

**Alternatives considered:**
- *Do nothing* - Rejected because it doesn't solve the core blocking issue
- *Use PostgreSQL* - Rejected as overkill; SQLite is sufficient for this use case and WAL addresses the concurrency issue

**Implementation:**
```python
# In _set_sqlite_pragma:
cursor.execute("PRAGMA journal_mode=WAL")
```

### Decision 2: Configure Busy Timeout

**Choice:** Set SQLite busy timeout to 5000ms (5 seconds).

**Rationale:** Instead of immediately failing on lock contention, SQLite will retry for up to 5 seconds. This handles brief contentions gracefully while not hanging indefinitely.

**Alternatives considered:**
- *No timeout (default ~0ms)* - Rejected; causes immediate failures on any contention
- *Longer timeout (30s+)* - Rejected; would mask the underlying issue and still cause poor UX

**Implementation:**
```python
cursor.execute("PRAGMA busy_timeout=5000")
```

### Decision 3: Restructure Parser with Thread Pool Execution

**Choice:** Split parsing into three phases AND run all blocking operations in a thread pool via `asyncio.to_thread()`:
1. **Fetch phase**: Query pending entries and build contexts (in thread pool)
2. **Parse phase**: Call LLM with no DB connection held (in thread pool)
3. **Write phase**: Write results with optimistic concurrency check (in thread pool)

**Rationale:** The parser worker is an `async` function, but calls synchronous blocking code (embedding model, LLM HTTP calls). Calling synchronous code directly from an async function blocks the entire event loop, preventing other requests from being handled. Using `asyncio.to_thread()` offloads blocking work to a thread pool.

**Key insight discovered during implementation:** Phase separation alone wasn't enough. Even with phases, synchronous calls like `embed_text()` (11s to load model) and `provider.parse()` (44s for LLM) blocked the event loop. The fix required wrapping all blocking operations in `asyncio.to_thread()`.

**Alternatives considered:**
- *Phase separation without thread pool* - Insufficient; still blocks event loop
- *Full async rewrite with async HTTP client* - Too invasive; thread pool achieves the goal
- *Move to background job queue* - Overkill for this use case

**Implementation approach:**
```python
def _fetch_entries_and_build_contexts(parser):
    """Runs in thread pool."""
    with get_session_context() as session:
        # ... fetch and build contexts
    return entries_to_process

def _parse_single_entry(parser, entry_data):
    """Runs in thread pool."""
    return parser.parse_text(entry_data.user_input, entry_data.context)

async def process_pending_entries(parser: EntryParser) -> int:
    # Phase 1: In thread pool
    entries = await asyncio.to_thread(_fetch_entries_and_build_contexts, parser)

    # Phase 2: In thread pool (sequential to avoid overwhelming LLM)
    for entry_data in entries:
        result = await asyncio.to_thread(_parse_single_entry, parser, entry_data)

    # Phase 3: In thread pool
    for parsed in parsed_entries:
        await asyncio.to_thread(_write_parse_result, parser, parsed)
```

### Decision 4: Handle Race Conditions with Optimistic Approach

**Choice:** Use optimistic concurrency - re-check entry status before writing.

**Rationale:** Between fetching and writing, an entry could be:
- Deleted by the user
- Manually edited (no longer "pending")
- Already parsed by another process

We handle this by re-querying and checking status before writing. If the entry changed, we skip it silently.

**Alternatives considered:**
- *Pessimistic locking with SELECT FOR UPDATE* - SQLite doesn't support row-level locking
- *Version column* - More complex, optimistic check is sufficient

### Decision 5: Separate Context Building from LLM Calls

**Choice:** Refactor `EntryParser` to have two methods:
- `build_parse_context(session)` - Queries DB for time codes, work types, RAG context
- `parse_text(user_input, context)` - Pure function, calls LLM with pre-built context

**Rationale:** This allows fetching all DB data upfront in the fetch phase, then running LLM calls with no DB dependency.

**Implementation:**
```python
class EntryParser:
    def build_parse_context(self, user_input: str, session: Session) -> ParseContext:
        """Fetch all context needed for parsing. Requires active session."""
        rag_contexts = self._retrieve_rag_context(user_input, session)
        return self._build_context(session, rag_contexts)

    def parse_text(self, user_input: str, context: ParseContext) -> ParseResult:
        """Parse text using pre-built context. No DB access."""
        return self.provider.parse(user_input, context)
```

## Risks / Trade-offs

**[Risk] Entry modified between fetch and write**
→ Mitigation: Re-check entry exists and is still "pending" before applying results. Skip silently if changed.

**[Risk] Context becomes stale during long parse**
→ Mitigation: Acceptable trade-off. Time codes/work types rarely change. If they do, entry can be reparsed.

**[Risk] WAL mode uses more disk space**
→ Mitigation: WAL files auto-checkpoint. For this app's scale, disk usage is negligible.

**[Risk] Multiple workers could pick same entry**
→ Mitigation: Currently single worker, but if scaling: use `status = "parsing"` intermediate state set in fetch phase.

## Open Questions

1. **Should we add a "parsing" intermediate status?** - Would prevent duplicate work if we ever scale to multiple workers. Low priority for now.

2. **Should context routes also defer embedding?** - Embedding generation in `/projects/{id}/context` routes is synchronous but typically fast (~1-2s). Could defer to background, but may be over-engineering for now.
