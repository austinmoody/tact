import asyncio
import logging
import os
from dataclasses import dataclass

from tact.db.models import TimeEntry
from tact.db.session import get_session_context
from tact.llm.parser import EntryParser, ParseContextWithRAG
from tact.llm.provider import ParseResult

logger = logging.getLogger(__name__)

DEFAULT_PARSER_INTERVAL = 10


@dataclass
class EntryToProcess:
    """Entry data extracted for processing without holding DB connection."""

    id: str
    user_input: str
    context: ParseContextWithRAG


@dataclass
class ParsedEntry:
    """Parse result ready to be written back to the database."""

    entry_id: str
    result: ParseResult
    rag_contexts: list | None


async def start_parser_worker() -> None:
    """Start the background parser worker."""
    interval = int(os.getenv("TACT_PARSER_INTERVAL", DEFAULT_PARSER_INTERVAL))
    logger.info(f"Starting parser worker with {interval}s interval")

    parser = EntryParser()

    while True:
        try:
            await process_pending_entries(parser)
        except Exception as e:
            logger.error(f"Parser worker error: {e}")

        await asyncio.sleep(interval)


def _fetch_entries_and_build_contexts(parser: EntryParser) -> list[EntryToProcess]:
    """Synchronous function to fetch entries and build contexts.

    Runs in thread pool to avoid blocking the event loop.
    """
    entries_to_process: list[EntryToProcess] = []

    with get_session_context() as session:
        pending = (
            session.query(TimeEntry)
            .filter(TimeEntry.status == "pending")
            .limit(10)  # Process in batches
            .all()
        )

        if not pending:
            return []

        logger.info(f"Processing {len(pending)} pending entries")

        for entry in pending:
            try:
                context = parser.build_parse_context(entry.user_input, session)
                entries_to_process.append(
                    EntryToProcess(
                        id=entry.id,
                        user_input=entry.user_input,
                        context=context,
                    )
                )
            except Exception as e:
                logger.error(f"Failed to build context for entry {entry.id}: {e}")

    return entries_to_process


def _parse_single_entry(
    parser: EntryParser, entry_data: EntryToProcess
) -> ParsedEntry | None:
    """Synchronous function to parse a single entry.

    Runs in thread pool to avoid blocking the event loop.
    """
    try:
        logger.info(f"Parsing entry {entry_data.id}: {entry_data.user_input[:50]}...")
        result = parser.parse_text(
            entry_data.user_input, entry_data.context.context
        )
        return ParsedEntry(
            entry_id=entry_data.id,
            result=result,
            rag_contexts=entry_data.context.rag_contexts,
        )
    except Exception as e:
        logger.error(f"Failed to parse entry {entry_data.id}: {e}")
        return None


def _write_parse_result(parser: EntryParser, parsed: ParsedEntry) -> bool:
    """Synchronous function to write parse result to database.

    Runs in thread pool to avoid blocking the event loop.
    """
    with get_session_context() as session:
        try:
            # Re-fetch entry and verify it still exists and is pending
            entry = session.get(TimeEntry, parsed.entry_id)

            if entry is None:
                logger.warning(
                    f"Entry {parsed.entry_id} was deleted during parsing, "
                    "discarding results"
                )
                return False

            if entry.status != "pending":
                logger.warning(
                    f"Entry {parsed.entry_id} status changed to '{entry.status}' "
                    "during parsing, discarding results"
                )
                return False

            # Apply parse results
            parser.apply_parse_result(
                entry, parsed.result, parsed.rag_contexts, session
            )
            session.commit()
            return True

        except Exception as e:
            logger.error(f"Failed to save parse result for entry {parsed.entry_id}: {e}")
            session.rollback()
            return False


async def process_pending_entries(parser: EntryParser) -> int:
    """Process all pending entries using three-phase approach.

    All blocking operations run in thread pool to avoid blocking the event loop.

    Phase 1: Fetch pending entries and build contexts (in thread pool)
    Phase 2: Call LLM for each entry (in thread pool, no DB connection held)
    Phase 3: Write results back (in thread pool, with optimistic concurrency check)

    Returns:
        Number of entries successfully processed
    """
    # Phase 1: Fetch pending entries and build contexts (in thread pool)
    entries_to_process = await asyncio.to_thread(
        _fetch_entries_and_build_contexts, parser
    )

    if not entries_to_process:
        return 0

    # Phase 2: Call LLM for each entry (in thread pool)
    # Note: We process sequentially to avoid overwhelming the LLM
    parsed_entries: list[ParsedEntry] = []

    for entry_data in entries_to_process:
        result = await asyncio.to_thread(_parse_single_entry, parser, entry_data)
        if result:
            parsed_entries.append(result)

    # Phase 3: Write results back (in thread pool)
    processed = 0

    for parsed in parsed_entries:
        success = await asyncio.to_thread(_write_parse_result, parser, parsed)
        if success:
            processed += 1

    return processed
