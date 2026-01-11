import json
import logging
import os
from datetime import UTC, datetime

from sqlalchemy.orm import Session

from tact.db.models import TimeCode, TimeEntry, WorkType
from tact.llm.provider import LLMProvider, ParseContext, TimeCodeInfo, WorkTypeInfo

logger = logging.getLogger(__name__)


def get_provider() -> LLMProvider:
    """Get the configured LLM provider."""
    provider_name = os.getenv("TACT_LLM_PROVIDER", "ollama").lower()

    if provider_name == "ollama":
        from tact.llm.ollama import OllamaProvider

        return OllamaProvider()
    elif provider_name == "anthropic":
        from tact.llm.anthropic import AnthropicProvider

        return AnthropicProvider()
    else:
        raise ValueError(f"Unknown LLM provider: {provider_name}")


class EntryParser:
    """Parses time entries using an LLM provider."""

    def __init__(self, provider: LLMProvider | None = None):
        self.provider = provider or get_provider()

    def parse_entry(self, entry: TimeEntry, session: Session) -> bool:
        """Parse a single entry and update it with results.

        Args:
            entry: The entry to parse
            session: Database session

        Returns:
            True if parsing succeeded, False otherwise
        """
        context = self._build_context(session)

        logger.info(f"Parsing entry {entry.id}: {entry.raw_text[:50]}...")

        result = self.provider.parse(entry.raw_text, context)

        if result.error:
            entry.status = "failed"
            entry.parse_error = result.error
            logger.warning(f"Entry {entry.id} parse failed: {result.error}")
            return False

        # Update entry with parsed results
        entry.duration_minutes = result.duration_minutes
        entry.work_type_id = result.work_type_id
        entry.time_code_id = result.time_code_id
        entry.description = result.description
        entry.confidence_duration = result.confidence_duration
        entry.confidence_work_type = result.confidence_work_type
        entry.confidence_time_code = result.confidence_time_code
        entry.confidence_overall = result.confidence_overall
        entry.status = "parsed"
        entry.parsed_at = datetime.now(UTC)
        entry.parse_error = None

        logger.info(
            f"Entry {entry.id} parsed: duration={result.duration_minutes}, "
            f"time_code={result.time_code_id}, work_type={result.work_type_id}, "
            f"confidence={result.confidence_overall}"
        )

        return True

    def _build_context(self, session: Session) -> ParseContext:
        """Build parsing context from active time codes and work types."""
        time_codes = session.query(TimeCode).filter(TimeCode.active.is_(True)).all()
        work_types = session.query(WorkType).filter(WorkType.active.is_(True)).all()

        return ParseContext(
            time_codes=[
                TimeCodeInfo(
                    id=tc.id,
                    name=tc.name,
                    description=tc.description,
                    keywords=json.loads(tc.keywords) if tc.keywords else [],
                )
                for tc in time_codes
            ],
            work_types=[
                WorkTypeInfo(id=wt.id, name=wt.name) for wt in work_types
            ],
        )
