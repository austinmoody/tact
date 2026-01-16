from abc import ABC, abstractmethod
from dataclasses import dataclass


@dataclass
class ParseResult:
    """Result of parsing a time entry."""

    duration_minutes: int | None = None
    work_type_id: str | None = None
    time_code_id: str | None = None
    parsed_description: str | None = None
    confidence_duration: float = 0.0
    confidence_work_type: float = 0.0
    confidence_time_code: float = 0.0
    confidence_overall: float = 0.0
    notes: str | None = None
    error: str | None = None


@dataclass
class TimeCodeInfo:
    """Time code information for context."""

    id: str
    name: str
    description: str
    keywords: list[str]


@dataclass
class WorkTypeInfo:
    """Work type information for context."""

    id: str
    name: str


@dataclass
class RAGContext:
    """RAG-retrieved context for parsing."""

    content: str
    project_id: str | None
    time_code_id: str | None
    similarity: float


@dataclass
class ParseContext:
    """Context provided to the LLM for parsing."""

    time_codes: list[TimeCodeInfo]
    work_types: list[WorkTypeInfo]
    rag_contexts: list[RAGContext] | None = None


class LLMProvider(ABC):
    """Abstract base class for LLM providers."""

    @abstractmethod
    def parse(self, user_input: str, context: ParseContext) -> ParseResult:
        """Parse user input into structured entry fields.

        Args:
            user_input: The user's natural language time entry text to parse
            context: Available time codes and work types

        Returns:
            ParseResult with extracted fields and confidence scores
        """
        pass
