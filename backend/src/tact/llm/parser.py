import logging
import os
from datetime import UTC, datetime

from sqlalchemy.orm import Session

from tact.db.models import Config, TimeCode, TimeEntry, WorkType
from tact.llm.provider import (
    LLMProvider,
    ParseContext,
    ParseResult,
    RAGContext,
    TimeCodeInfo,
    WorkTypeInfo,
)
from tact.rag.retrieval import retrieve_similar_contexts
from tact.utils.duration import round_duration

logger = logging.getLogger(__name__)

from dataclasses import dataclass

DEFAULT_CONFIDENCE_THRESHOLD = 0.7


@dataclass
class ParseContextWithRAG:
    """Parse context bundled with RAG contexts for later use in apply_parse_result."""

    context: ParseContext
    rag_contexts: list[RAGContext] | None


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

    def build_parse_context(
        self, user_input: str, session: Session
    ) -> ParseContextWithRAG:
        """Build parsing context from database. Requires active session.

        This method fetches all data needed for parsing (RAG context, time codes,
        work types) so the actual LLM call can be made without a database connection.

        Args:
            user_input: The raw text to parse (used for RAG retrieval)
            session: Active database session

        Returns:
            ParseContextWithRAG containing ParseContext and RAG contexts
        """
        rag_contexts = self._retrieve_rag_context(user_input, session)
        context = self._build_context(session, rag_contexts)
        return ParseContextWithRAG(context=context, rag_contexts=rag_contexts)

    def parse_text(self, user_input: str, context: ParseContext) -> ParseResult:
        """Parse text using pre-built context. No database access required.

        Args:
            user_input: The raw text to parse
            context: Pre-built ParseContext from build_parse_context()

        Returns:
            ParseResult with extracted fields and confidence scores
        """
        return self.provider.parse(user_input, context)

    def apply_parse_result(
        self,
        entry: TimeEntry,
        result: ParseResult,
        rag_contexts: list[RAGContext] | None,
        session: Session,
    ) -> bool:
        """Apply a ParseResult to an entry. Requires active session for threshold lookup.

        Args:
            entry: The entry to update
            result: ParseResult from parse_text()
            rag_contexts: RAG contexts used for building parse notes
            session: Active database session (for threshold config lookup)

        Returns:
            True if parsing succeeded, False if there was an error
        """
        if result.error:
            entry.status = "failed"
            entry.parse_error = result.error
            logger.warning(f"Entry {entry.id} parse failed: {result.error}")
            return False

        # Update entry with parsed results (apply duration rounding if configured)
        entry.duration_minutes = round_duration(result.duration_minutes)
        entry.work_type_id = result.work_type_id
        entry.time_code_id = result.time_code_id
        entry.parsed_description = result.parsed_description
        entry.confidence_duration = result.confidence_duration
        entry.confidence_work_type = result.confidence_work_type
        entry.confidence_time_code = result.confidence_time_code
        entry.confidence_overall = result.confidence_overall
        entry.parsed_at = datetime.now(UTC)
        entry.parse_error = None

        # Build parse_notes from LLM reasoning and RAG context info
        entry.parse_notes = self._build_parse_notes(result.notes, rag_contexts)

        # Set status based on required fields and confidence threshold
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

    def parse_entry(self, entry: TimeEntry, session: Session) -> bool:
        """Parse a single entry and update it with results.

        Args:
            entry: The entry to parse
            session: Database session

        Returns:
            True if parsing succeeded, False otherwise
        """
        # Retrieve RAG context before building parse context
        rag_contexts = self._retrieve_rag_context(entry.user_input, session)
        context = self._build_context(session, rag_contexts)

        logger.info(f"Parsing entry {entry.id}: {entry.user_input[:50]}...")

        result = self.provider.parse(entry.user_input, context)

        if result.error:
            entry.status = "failed"
            entry.parse_error = result.error
            logger.warning(f"Entry {entry.id} parse failed: {result.error}")
            return False

        # Update entry with parsed results (apply duration rounding if configured)
        entry.duration_minutes = round_duration(result.duration_minutes)
        entry.work_type_id = result.work_type_id
        entry.time_code_id = result.time_code_id
        entry.parsed_description = result.parsed_description
        entry.confidence_duration = result.confidence_duration
        entry.confidence_work_type = result.confidence_work_type
        entry.confidence_time_code = result.confidence_time_code
        entry.confidence_overall = result.confidence_overall
        entry.parsed_at = datetime.now(UTC)
        entry.parse_error = None

        # Build parse_notes from LLM reasoning and RAG context info
        entry.parse_notes = self._build_parse_notes(result.notes, rag_contexts)

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
                TimeCodeInfo(id=tc.id, name=tc.name) for tc in time_codes
            ],
            work_types=[
                WorkTypeInfo(id=wt.id, name=wt.name) for wt in work_types
            ],
            rag_contexts=rag_contexts,
        )

    def _build_parse_notes(
        self, llm_notes: str | None, rag_contexts: list[RAGContext] | None
    ) -> str | None:
        """Build parse notes from LLM reasoning and RAG context info.

        Args:
            llm_notes: Notes/reasoning from the LLM
            rag_contexts: RAG contexts that were used for parsing

        Returns:
            Combined notes string or None if no notes available
        """
        parts = []

        # Add LLM reasoning
        if llm_notes:
            parts.append(llm_notes)

        # Add info about closest RAG context if available
        if rag_contexts:
            # Get the highest similarity context
            best_context = max(rag_contexts, key=lambda c: c.similarity)
            source = (
                f"time_code:{best_context.time_code_id}"
                if best_context.time_code_id
                else f"project:{best_context.project_id}"
            )
            context_info = (
                f"[Context used: '{best_context.content[:100]}' "
                f"({source}, similarity: {best_context.similarity:.2f})]"
            )
            parts.append(context_info)

        if not parts:
            return None

        return " ".join(parts)
