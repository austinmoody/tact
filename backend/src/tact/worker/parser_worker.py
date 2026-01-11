import asyncio
import logging
import os

from tact.db.models import TimeEntry
from tact.db.session import SessionFactory
from tact.llm.parser import EntryParser

logger = logging.getLogger(__name__)

DEFAULT_PARSER_INTERVAL = 10


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


async def process_pending_entries(parser: EntryParser) -> int:
    """Process all pending entries.

    Returns:
        Number of entries processed
    """
    session = SessionFactory()
    try:
        pending = (
            session.query(TimeEntry)
            .filter(TimeEntry.status == "pending")
            .limit(10)  # Process in batches
            .all()
        )

        if not pending:
            return 0

        logger.info(f"Processing {len(pending)} pending entries")

        processed = 0
        for entry in pending:
            try:
                parser.parse_entry(entry, session)
                session.commit()
                processed += 1
            except Exception as e:
                logger.error(f"Failed to parse entry {entry.id}: {e}")
                session.rollback()

        return processed

    finally:
        session.close()
