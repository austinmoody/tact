import json
from unittest.mock import MagicMock, patch

import anthropic
import pytest

from tact.llm.parser import EntryParser
from tact.llm.prompts import build_system_prompt, build_user_prompt
from tact.llm.provider import ParseContext, ParseResult, TimeCodeInfo, WorkTypeInfo


@pytest.fixture
def sample_context():
    return ParseContext(
        time_codes=[
            TimeCodeInfo(id="PROJ-001", name="Project Alpha"),
            TimeCodeInfo(id="ADMIN-01", name="Admin Tasks"),
        ],
        work_types=[
            WorkTypeInfo(id="development", name="Development"),
            WorkTypeInfo(id="meeting", name="Meeting"),
        ],
    )


@pytest.fixture
def empty_context():
    return ParseContext(time_codes=[], work_types=[])


class TestPrompts:
    def test_build_system_prompt_includes_time_codes(self, sample_context):
        prompt = build_system_prompt(sample_context)
        assert "PROJ-001" in prompt
        assert "Project Alpha" in prompt
        assert "ADMIN-01" in prompt

    def test_build_system_prompt_includes_work_types(self, sample_context):
        prompt = build_system_prompt(sample_context)
        assert "development" in prompt
        assert "Development" in prompt
        assert "meeting" in prompt

    def test_build_system_prompt_empty_context(self, empty_context):
        prompt = build_system_prompt(empty_context)
        assert "(none defined)" in prompt

    def test_build_user_prompt(self):
        user_input = "2h coding on alpha"
        prompt = build_user_prompt(user_input)
        assert user_input in prompt
        assert "Parse this time entry" in prompt


class TestParseResult:
    def test_parse_result_defaults(self):
        result = ParseResult()
        assert result.duration_minutes is None
        assert result.work_type_id is None
        assert result.time_code_id is None
        assert result.parsed_description is None
        assert result.confidence_duration == 0.0
        assert result.confidence_overall == 0.0
        assert result.notes is None
        assert result.error is None

    def test_parse_result_with_values(self):
        result = ParseResult(
            duration_minutes=120,
            time_code_id="PROJ-001",
            work_type_id="development",
            parsed_description="Coding",
            confidence_duration=0.95,
            confidence_overall=0.90,
            notes="Matched based on 'alpha' keyword",
        )
        assert result.duration_minutes == 120
        assert result.time_code_id == "PROJ-001"
        assert result.confidence_overall == 0.90
        assert result.notes == "Matched based on 'alpha' keyword"


class TestOllamaProvider:
    def test_parse_success(self, sample_context):
        from tact.llm.ollama import OllamaProvider

        provider = OllamaProvider()
        provider._model_verified = True  # Skip model check for this test

        mock_response = MagicMock()
        mock_response.json.return_value = {
            "response": json.dumps(
                {
                    "duration_minutes": 120,
                    "time_code_id": "PROJ-001",
                    "work_type_id": "development",
                    "parsed_description": "Working on alpha",
                    "confidence_duration": 0.95,
                    "confidence_time_code": 0.85,
                    "confidence_work_type": 0.90,
                    "confidence_overall": 0.85,
                    "notes": "Matched to PROJ-001 based on 'alpha' keyword",
                }
            )
        }

        with patch.object(provider.client, "post", return_value=mock_response):
            result = provider.parse("2h coding on alpha", sample_context)

        assert result.duration_minutes == 120
        assert result.time_code_id == "PROJ-001"
        assert result.work_type_id == "development"
        assert result.notes == "Matched to PROJ-001 based on 'alpha' keyword"
        assert result.error is None

    def test_parse_float_duration_rounds_to_int(self, sample_context):
        """Test that float duration like 127.5 gets rounded to int."""
        from tact.llm.ollama import OllamaProvider

        provider = OllamaProvider()
        provider._model_verified = True  # Skip model check for this test

        mock_response = MagicMock()
        mock_response.json.return_value = {
            "response": json.dumps(
                {
                    "duration_minutes": 127.5,
                    "time_code_id": "PROJ-001",
                    "work_type_id": "development",
                    "parsed_description": "Working on alpha",
                    "confidence_duration": 0.95,
                    "confidence_time_code": 0.85,
                    "confidence_work_type": 0.90,
                    "confidence_overall": 0.85,
                }
            )
        }

        with patch.object(provider.client, "post", return_value=mock_response):
            result = provider.parse("2h coding on alpha", sample_context)

        # 127.5 rounds to 128
        assert result.duration_minutes == 128
        assert isinstance(result.duration_minutes, int)
        assert result.error is None

    def test_parse_invalid_json(self, sample_context):
        from tact.llm.ollama import OllamaProvider

        provider = OllamaProvider()
        provider._model_verified = True  # Skip model check for this test

        mock_response = MagicMock()
        mock_response.json.return_value = {"response": "not valid json"}

        with patch.object(provider.client, "post", return_value=mock_response):
            result = provider.parse("test", sample_context)

        assert result.error is not None
        assert "Invalid JSON" in result.error

    def test_parse_http_error(self, sample_context):
        from httpx import HTTPError

        from tact.llm.ollama import OllamaProvider

        provider = OllamaProvider()
        provider._model_verified = True  # Skip model check for this test

        with patch.object(
            provider.client, "post", side_effect=HTTPError("Connection refused")
        ):
            result = provider.parse("test", sample_context)

        assert result.error is not None
        assert "HTTP error" in result.error

    def test_ensure_model_available_when_model_exists(self):
        """Test that model check passes when model is already available."""
        from tact.llm.ollama import OllamaProvider

        provider = OllamaProvider(model="llama3.2:3b")

        mock_tags_response = MagicMock()
        mock_tags_response.json.return_value = {
            "models": [{"name": "llama3.2:3b"}, {"name": "mistral:latest"}]
        }

        with patch.object(provider.client, "get", return_value=mock_tags_response):
            error = provider._ensure_model_available()

        assert error is None
        assert provider._model_verified is True

    def test_ensure_model_available_pulls_missing_model(self):
        """Test that missing model triggers a pull."""
        from tact.llm.ollama import OllamaProvider

        provider = OllamaProvider(model="llama3.2:3b")

        mock_tags_response = MagicMock()
        mock_tags_response.json.return_value = {"models": []}  # No models

        # Mock the streaming pull response
        mock_pull_response = MagicMock()
        mock_pull_response.__enter__ = MagicMock(return_value=mock_pull_response)
        mock_pull_response.__exit__ = MagicMock(return_value=False)
        mock_pull_response.iter_lines.return_value = [
            '{"status": "pulling manifest"}',
            '{"status": "downloading"}',
            '{"status": "success"}',
        ]

        with (
            patch.object(provider.client, "get", return_value=mock_tags_response),
            patch("tact.llm.ollama.httpx.Client") as mock_client_class,
        ):
            mock_pull_client = MagicMock()
            mock_pull_client.stream.return_value = mock_pull_response
            mock_client_class.return_value = mock_pull_client

            error = provider._ensure_model_available()

        assert error is None
        assert provider._model_verified is True
        mock_pull_client.stream.assert_called_once()

    def test_ensure_model_available_skips_when_verified(self):
        """Test that model check is skipped when already verified."""
        from tact.llm.ollama import OllamaProvider

        provider = OllamaProvider()
        provider._model_verified = True

        with patch.object(provider.client, "get") as mock_get:
            error = provider._ensure_model_available()

        assert error is None
        mock_get.assert_not_called()

    def test_ensure_model_available_returns_error_on_pull_failure(self):
        """Test that pull error is returned properly."""
        from tact.llm.ollama import OllamaProvider

        provider = OllamaProvider(model="nonexistent:model")

        mock_tags_response = MagicMock()
        mock_tags_response.json.return_value = {"models": []}

        mock_pull_response = MagicMock()
        mock_pull_response.__enter__ = MagicMock(return_value=mock_pull_response)
        mock_pull_response.__exit__ = MagicMock(return_value=False)
        mock_pull_response.iter_lines.return_value = [
            '{"error": "model not found"}',
        ]

        with (
            patch.object(provider.client, "get", return_value=mock_tags_response),
            patch("tact.llm.ollama.httpx.Client") as mock_client_class,
        ):
            mock_pull_client = MagicMock()
            mock_pull_client.stream.return_value = mock_pull_response
            mock_client_class.return_value = mock_pull_client

            error = provider._ensure_model_available()

        assert error is not None
        assert "model not found" in error
        assert provider._model_verified is False


class TestAnthropicProvider:
    def test_parse_success(self, sample_context):
        from tact.llm.anthropic import AnthropicProvider

        mock_message = MagicMock()
        mock_message.content = [
            MagicMock(
                text=json.dumps(
                    {
                        "duration_minutes": 120,
                        "time_code_id": "PROJ-001",
                        "work_type_id": "development",
                        "parsed_description": "Working on alpha",
                        "confidence_duration": 0.95,
                        "confidence_time_code": 0.85,
                        "confidence_work_type": 0.90,
                        "confidence_overall": 0.85,
                        "notes": "Matched to PROJ-001 based on 'alpha' keyword",
                    }
                )
            )
        ]

        with (
            patch.object(anthropic.Anthropic, "__init__", return_value=None),
            patch.object(anthropic.Anthropic, "messages", create=True) as mock_messages,
        ):
            mock_messages.create.return_value = mock_message
            provider = AnthropicProvider(api_key="test-key")
            provider.client = MagicMock()
            provider.client.messages.create.return_value = mock_message

            result = provider.parse("2h coding on alpha", sample_context)

        assert result.duration_minutes == 120
        assert result.time_code_id == "PROJ-001"
        assert result.work_type_id == "development"
        assert result.notes == "Matched to PROJ-001 based on 'alpha' keyword"
        assert result.error is None

    def test_parse_float_duration_rounds_to_int(self, sample_context):
        """Test that float duration like 127.5 gets rounded to int."""
        from tact.llm.anthropic import AnthropicProvider

        mock_message = MagicMock()
        mock_message.content = [
            MagicMock(
                text=json.dumps(
                    {
                        "duration_minutes": 127.5,
                        "time_code_id": "PROJ-001",
                        "work_type_id": "development",
                        "parsed_description": "Working on alpha",
                        "confidence_duration": 0.95,
                        "confidence_time_code": 0.85,
                        "confidence_work_type": 0.90,
                        "confidence_overall": 0.85,
                    }
                )
            )
        ]

        with patch.object(anthropic.Anthropic, "__init__", return_value=None):
            provider = AnthropicProvider(api_key="test-key")
            provider.client = MagicMock()
            provider.client.messages.create.return_value = mock_message

            result = provider.parse("2h coding on alpha", sample_context)

        # 127.5 rounds to 128
        assert result.duration_minutes == 128
        assert isinstance(result.duration_minutes, int)
        assert result.error is None

    def test_parse_invalid_json(self, sample_context):
        from tact.llm.anthropic import AnthropicProvider

        mock_message = MagicMock()
        mock_message.content = [MagicMock(text="not valid json")]

        with patch.object(anthropic.Anthropic, "__init__", return_value=None):
            provider = AnthropicProvider(api_key="test-key")
            provider.client = MagicMock()
            provider.client.messages.create.return_value = mock_message

            result = provider.parse("test", sample_context)

        assert result.error is not None
        assert "Invalid JSON" in result.error

    def test_parse_api_error(self, sample_context):
        from tact.llm.anthropic import AnthropicProvider

        with patch.object(anthropic.Anthropic, "__init__", return_value=None):
            provider = AnthropicProvider(api_key="test-key")
            provider.client = MagicMock()
            provider.client.messages.create.side_effect = anthropic.APIConnectionError(
                request=MagicMock()
            )

            result = provider.parse("test", sample_context)

        assert result.error is not None
        assert "Connection error" in result.error

    def test_missing_api_key_raises(self):
        from tact.llm.anthropic import AnthropicProvider

        with (
            patch.dict("os.environ", {}, clear=True),
            patch.object(anthropic.Anthropic, "__init__", return_value=None),
        ):
            with pytest.raises(ValueError, match="API key required"):
                AnthropicProvider()


class TestEntryParser:
    def test_parse_entry_success(self):
        mock_provider = MagicMock()
        mock_provider.parse.return_value = ParseResult(
            duration_minutes=120,
            time_code_id="PROJ-001",
            work_type_id="development",
            parsed_description="Coding",
            confidence_duration=0.95,
            confidence_time_code=0.85,
            confidence_work_type=0.90,
            confidence_overall=0.85,
        )

        mock_entry = MagicMock()
        mock_entry.id = "entry-123"
        mock_entry.user_input = "2h coding on alpha"

        mock_session = MagicMock()
        mock_query_result = MagicMock()
        mock_query_result.filter.return_value.all.return_value = []
        mock_query_result.filter.return_value.first.return_value = None
        mock_session.query.return_value = mock_query_result

        parser = EntryParser(provider=mock_provider)
        result = parser.parse_entry(mock_entry, mock_session)

        assert result is True
        assert mock_entry.duration_minutes == 120
        assert mock_entry.time_code_id == "PROJ-001"
        assert mock_entry.status == "parsed"

    def test_parse_entry_failure(self):
        mock_provider = MagicMock()
        mock_provider.parse.return_value = ParseResult(error="Parse failed")

        mock_entry = MagicMock()
        mock_entry.id = "entry-123"
        mock_entry.user_input = "invalid"

        mock_session = MagicMock()
        mock_query_result = MagicMock()
        mock_query_result.filter.return_value.all.return_value = []
        mock_query_result.filter.return_value.first.return_value = None
        mock_session.query.return_value = mock_query_result

        parser = EntryParser(provider=mock_provider)
        result = parser.parse_entry(mock_entry, mock_session)

        assert result is False
        assert mock_entry.status == "failed"
        assert mock_entry.parse_error == "Parse failed"

    def test_parse_entry_missing_time_code_needs_review(self):
        """Test that entries without time_code get needs_review status."""
        mock_provider = MagicMock()
        mock_provider.parse.return_value = ParseResult(
            duration_minutes=60,
            time_code_id=None,  # Missing time code
            work_type_id="development",
            parsed_description="Some work",
            confidence_duration=0.9,
            confidence_time_code=0.3,
            confidence_work_type=0.8,
            confidence_overall=0.5,
        )

        mock_entry = MagicMock()
        mock_entry.id = "entry-456"
        mock_entry.user_input = "did some work"

        mock_session = MagicMock()
        mock_query_result = MagicMock()
        mock_query_result.filter.return_value.all.return_value = []
        mock_query_result.filter.return_value.first.return_value = None
        mock_session.query.return_value = mock_query_result

        parser = EntryParser(provider=mock_provider)
        result = parser.parse_entry(mock_entry, mock_session)

        assert result is True
        assert mock_entry.status == "needs_review"

    def test_parse_entry_missing_duration_needs_review(self):
        """Test that entries without duration get needs_review status."""
        mock_provider = MagicMock()
        mock_provider.parse.return_value = ParseResult(
            duration_minutes=None,  # Missing duration
            time_code_id="PROJ-001",
            work_type_id="development",
            parsed_description="Some work",
            confidence_duration=0.3,
            confidence_time_code=0.9,
            confidence_work_type=0.8,
            confidence_overall=0.5,
        )

        mock_entry = MagicMock()
        mock_entry.id = "entry-789"
        mock_entry.user_input = "worked on project"

        mock_session = MagicMock()
        mock_query_result = MagicMock()
        mock_query_result.filter.return_value.all.return_value = []
        mock_query_result.filter.return_value.first.return_value = None
        mock_session.query.return_value = mock_query_result

        parser = EntryParser(provider=mock_provider)
        result = parser.parse_entry(mock_entry, mock_session)

        assert result is True
        assert mock_entry.status == "needs_review"

    def test_parse_entry_low_time_code_confidence_needs_review(self):
        """Test that low time_code confidence gets needs_review even with value."""
        mock_provider = MagicMock()
        mock_provider.parse.return_value = ParseResult(
            duration_minutes=60,
            time_code_id="PROJ-001",  # Has value but low confidence
            work_type_id="development",
            parsed_description="Some work",
            confidence_duration=0.9,
            confidence_time_code=0.5,  # Below 0.7 threshold
            confidence_work_type=0.8,
            confidence_overall=0.6,
        )

        mock_entry = MagicMock()
        mock_entry.id = "entry-low-conf"
        mock_entry.user_input = "maybe worked on project"

        mock_session = MagicMock()
        mock_query_result = MagicMock()
        mock_query_result.filter.return_value.all.return_value = []
        mock_query_result.filter.return_value.first.return_value = None
        mock_session.query.return_value = mock_query_result

        parser = EntryParser(provider=mock_provider)
        result = parser.parse_entry(mock_entry, mock_session)

        assert result is True
        assert mock_entry.status == "needs_review"

    def test_parse_entry_low_duration_confidence_needs_review(self):
        """Test that low duration confidence gets needs_review even with value."""
        mock_provider = MagicMock()
        mock_provider.parse.return_value = ParseResult(
            duration_minutes=60,  # Has value but low confidence
            time_code_id="PROJ-001",
            work_type_id="development",
            parsed_description="Some work",
            confidence_duration=0.5,  # Below 0.7 threshold
            confidence_time_code=0.9,
            confidence_work_type=0.8,
            confidence_overall=0.6,
        )

        mock_entry = MagicMock()
        mock_entry.id = "entry-low-dur-conf"
        mock_entry.user_input = "worked some time on project"

        mock_session = MagicMock()
        mock_query_result = MagicMock()
        mock_query_result.filter.return_value.all.return_value = []
        mock_query_result.filter.return_value.first.return_value = None
        mock_session.query.return_value = mock_query_result

        parser = EntryParser(provider=mock_provider)
        result = parser.parse_entry(mock_entry, mock_session)

        assert result is True
        assert mock_entry.status == "needs_review"

    def test_parse_entry_includes_notes(self):
        """Test that parse_notes is set from LLM notes."""
        mock_provider = MagicMock()
        mock_provider.parse.return_value = ParseResult(
            duration_minutes=120,
            time_code_id="PROJ-001",
            work_type_id="development",
            parsed_description="Coding",
            confidence_duration=0.95,
            confidence_time_code=0.85,
            confidence_work_type=0.90,
            confidence_overall=0.85,
            notes="Matched PROJ-001 based on 'alpha' keyword",
        )

        mock_entry = MagicMock()
        mock_entry.id = "entry-with-notes"
        mock_entry.user_input = "2h coding on alpha"

        mock_session = MagicMock()
        mock_query_result = MagicMock()
        mock_query_result.filter.return_value.all.return_value = []
        mock_query_result.filter.return_value.first.return_value = None
        mock_session.query.return_value = mock_query_result

        parser = EntryParser(provider=mock_provider)
        result = parser.parse_entry(mock_entry, mock_session)

        assert result is True
        assert mock_entry.parse_notes == "Matched PROJ-001 based on 'alpha' keyword"


class TestBuildParseNotes:
    """Tests for the _build_parse_notes helper method."""

    def test_build_parse_notes_with_llm_notes_only(self):
        """Test parse notes with only LLM reasoning."""
        parser = EntryParser(provider=MagicMock())
        notes = parser._build_parse_notes("Matched based on keyword", None)
        assert notes == "Matched based on keyword"

    def test_build_parse_notes_with_rag_context_only(self):
        """Test parse notes with only RAG context."""
        from tact.llm.provider import RAGContext

        parser = EntryParser(provider=MagicMock())
        rag_contexts = [
            RAGContext(
                content="ALL APHL meetings go to FEDS-163",
                project_id=None,
                time_code_id="FEDS-163",
                similarity=0.85,
            )
        ]
        notes = parser._build_parse_notes(None, rag_contexts)
        assert "[Context used:" in notes
        assert "FEDS-163" in notes
        assert "0.85" in notes

    def test_build_parse_notes_with_both(self):
        """Test parse notes with both LLM notes and RAG context."""
        from tact.llm.provider import RAGContext

        parser = EntryParser(provider=MagicMock())
        rag_contexts = [
            RAGContext(
                content="ALL APHL meetings go to FEDS-163",
                project_id=None,
                time_code_id="FEDS-163",
                similarity=0.85,
            )
        ]
        notes = parser._build_parse_notes("Used APHL rule", rag_contexts)
        assert "Used APHL rule" in notes
        assert "[Context used:" in notes

    def test_build_parse_notes_empty(self):
        """Test parse notes returns None when no info available."""
        parser = EntryParser(provider=MagicMock())
        notes = parser._build_parse_notes(None, None)
        assert notes is None

    def test_build_parse_notes_picks_highest_similarity(self):
        """Test that highest similarity context is used."""
        from tact.llm.provider import RAGContext

        parser = EntryParser(provider=MagicMock())
        rag_contexts = [
            RAGContext(
                content="Low similarity context",
                project_id="proj-1",
                time_code_id=None,
                similarity=0.5,
            ),
            RAGContext(
                content="High similarity context",
                project_id=None,
                time_code_id="TC-001",
                similarity=0.9,
            ),
        ]
        notes = parser._build_parse_notes(None, rag_contexts)
        assert "High similarity context" in notes
        assert "0.90" in notes
