import json
import logging
import os
from datetime import UTC, datetime

from sqlalchemy.orm import Session

from tact.db.models import Config, TimeCode, TimeEntry, WorkType
from tact.llm.provider import (
    LLMProvider,
    ParseContext,
    RAGContext,
    TimeCodeInfo,
    WorkTypeInfo,
)
from tact.rag.retrieval import retrieve_similar_contexts

logger = logging.getLogger(__name__)

DEFAULT_CONFIDENCE_THRESHOLD = 0.7


def get_confidence_threshold(session: Session) -> float:
    """Get confidence threshold from config or use default."""
    config = session.query(Config).filter(Config.key == "confidence_threshold").first()
    if config:
        try:
            return float(config.value)
        except ValueError:
            logger.warning("Invalid confidence_threshold config: %s", config.value)
    return DEFAULT_CONFIDENCE_THRESHOLD


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
        # Retrieve RAG context before building parse context
        rag_contexts = self._retrieve_rag_context(entry.raw_text, session)
        context = self._build_context(session, rag_contexts)

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
        entry.parsed_at = datetime.now(UTC)
        entry.parse_error = None

        # Set status based on required fields and confidence threshold
        # Entry is only "parsed" if both time_code and duration are set with
        # confidence above threshold. Work type is optional.
        threshold = get_confidence_threshold(session)
        has_time_code = (
            result.time_code_id is not None
            and (result.confidence_time_code or 0.0) >= threshold
        )
        has_duration = (
            result.duration_minutes is not None
            and (result.confidence_duration or 0.0) >= threshold
        )

        if has_time_code and has_duration:
            entry.status = "parsed"
        else:
            entry.status = "needs_review"

        logger.info(
            f"Entry {entry.id} parsed: duration={result.duration_minutes}, "
            f"time_code={result.time_code_id}, work_type={result.work_type_id}, "
            f"conf_duration={result.confidence_duration}, "
            f"conf_time_code={result.confidence_time_code}, status={entry.status}"
        )

        return True

    def _retrieve_rag_context(
        self, raw_text: str, session: Session
    ) -> list[RAGContext] | None:
        """Retrieve relevant context documents for the entry text."""
        try:
            similar = retrieve_similar_contexts(raw_text, session, top_k=5)
            if not similar:
                return None

            rag_contexts = [
                RAGContext(
                    content=ctx.content,
                    project_id=ctx.project_id,
                    time_code_id=ctx.time_code_id,
                    similarity=ctx.similarity,
                )
                for ctx in similar
            ]
            logger.info(
                "Retrieved %d RAG context documents for entry", len(rag_contexts)
            )
            return rag_contexts
        except Exception as e:
            logger.warning("Failed to retrieve RAG context: %s", e)
            return None

    def _build_context(
        self, session: Session, rag_contexts: list[RAGContext] | None = None
    ) -> ParseContext:
        """Build parsing context from active time codes, work types, and RAG context."""
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
            rag_contexts=rag_contexts,
        )
