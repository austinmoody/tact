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
            TimeCodeInfo(
                id="PROJ-001",
                name="Project Alpha",
                description="Main project",
                keywords=["alpha", "main"],
            ),
            TimeCodeInfo(
                id="ADMIN-01",
                name="Admin Tasks",
                description="Administrative work",
                keywords=["admin", "paperwork"],
            ),
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
        assert "alpha, main" in prompt
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
        raw_text = "2h coding on alpha"
        prompt = build_user_prompt(raw_text)
        assert raw_text in prompt
        assert "Parse this time entry" in prompt


class TestParseResult:
    def test_parse_result_defaults(self):
        result = ParseResult()
        assert result.duration_minutes is None
        assert result.work_type_id is None
        assert result.time_code_id is None
        assert result.description is None
        assert result.confidence_duration == 0.0
        assert result.confidence_overall == 0.0
        assert result.error is None

    def test_parse_result_with_values(self):
        result = ParseResult(
            duration_minutes=120,
            time_code_id="PROJ-001",
            work_type_id="development",
            description="Coding",
            confidence_duration=0.95,
            confidence_overall=0.90,
        )
        assert result.duration_minutes == 120
        assert result.time_code_id == "PROJ-001"
        assert result.confidence_overall == 0.90


class TestOllamaProvider:
    def test_parse_success(self, sample_context):
        from tact.llm.ollama import OllamaProvider

        provider = OllamaProvider()

        mock_response = MagicMock()
        mock_response.json.return_value = {
            "response": json.dumps({
                "duration_minutes": 120,
                "time_code_id": "PROJ-001",
                "work_type_id": "development",
                "description": "Working on alpha",
                "confidence_duration": 0.95,
                "confidence_time_code": 0.85,
                "confidence_work_type": 0.90,
                "confidence_overall": 0.85,
            })
        }

        with patch.object(provider.client, "post", return_value=mock_response):
            result = provider.parse("2h coding on alpha", sample_context)

        assert result.duration_minutes == 120
        assert result.time_code_id == "PROJ-001"
        assert result.work_type_id == "development"
        assert result.error is None

    def test_parse_invalid_json(self, sample_context):
        from tact.llm.ollama import OllamaProvider

        provider = OllamaProvider()

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

        with patch.object(
            provider.client, "post", side_effect=HTTPError("Connection refused")
        ):
            result = provider.parse("test", sample_context)

        assert result.error is not None
        assert "HTTP error" in result.error


class TestAnthropicProvider:
    def test_parse_success(self, sample_context):
        from tact.llm.anthropic import AnthropicProvider

        mock_message = MagicMock()
        mock_message.content = [
            MagicMock(
                text=json.dumps({
                    "duration_minutes": 120,
                    "time_code_id": "PROJ-001",
                    "work_type_id": "development",
                    "description": "Working on alpha",
                    "confidence_duration": 0.95,
                    "confidence_time_code": 0.85,
                    "confidence_work_type": 0.90,
                    "confidence_overall": 0.85,
                })
            )
        ]

        with patch.object(
            anthropic.Anthropic, "__init__", return_value=None
        ), patch.object(
            anthropic.Anthropic, "messages", create=True
        ) as mock_messages:
            mock_messages.create.return_value = mock_message
            provider = AnthropicProvider(api_key="test-key")
            provider.client = MagicMock()
            provider.client.messages.create.return_value = mock_message

            result = provider.parse("2h coding on alpha", sample_context)

        assert result.duration_minutes == 120
        assert result.time_code_id == "PROJ-001"
        assert result.work_type_id == "development"
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

        with patch.dict("os.environ", {}, clear=True), patch.object(
            anthropic.Anthropic, "__init__", return_value=None
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
            description="Coding",
            confidence_duration=0.95,
            confidence_time_code=0.85,
            confidence_work_type=0.90,
            confidence_overall=0.85,
        )

        mock_entry = MagicMock()
        mock_entry.id = "entry-123"
        mock_entry.raw_text = "2h coding on alpha"

        mock_session = MagicMock()
        mock_session.query.return_value.filter.return_value.all.return_value = []

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
        mock_entry.raw_text = "invalid"

        mock_session = MagicMock()
        mock_session.query.return_value.filter.return_value.all.return_value = []

        parser = EntryParser(provider=mock_provider)
        result = parser.parse_entry(mock_entry, mock_session)

        assert result is False
        assert mock_entry.status == "failed"
        assert mock_entry.parse_error == "Parse failed"
