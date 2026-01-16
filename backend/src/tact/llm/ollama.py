import json
import logging
import os

import httpx

from tact.llm.prompts import build_system_prompt, build_user_prompt
from tact.llm.provider import LLMProvider, ParseContext, ParseResult

logger = logging.getLogger(__name__)

DEFAULT_OLLAMA_URL = "http://localhost:11434"
DEFAULT_OLLAMA_MODEL = "llama3.2:3b"
DEFAULT_OLLAMA_TIMEOUT = 180.0  # seconds


class OllamaProvider(LLMProvider):
    """LLM provider using Ollama for local model inference."""

    def __init__(
        self,
        base_url: str | None = None,
        model: str | None = None,
        timeout: float | None = None,
    ):
        self.base_url = base_url or os.getenv("TACT_OLLAMA_URL", DEFAULT_OLLAMA_URL)
        self.model = model or os.getenv("TACT_OLLAMA_MODEL", DEFAULT_OLLAMA_MODEL)
        self.timeout = timeout or float(
            os.getenv("TACT_OLLAMA_TIMEOUT", DEFAULT_OLLAMA_TIMEOUT)
        )
        self.client = httpx.Client(timeout=self.timeout)

    def parse(self, raw_text: str, context: ParseContext) -> ParseResult:
        """Parse raw text using Ollama."""
        system_prompt = build_system_prompt(context)
        user_prompt = build_user_prompt(raw_text)

        try:
            response = self.client.post(
                f"{self.base_url}/api/generate",
                json={
                    "model": self.model,
                    "prompt": f"{system_prompt}\n\n{user_prompt}",
                    "stream": False,
                    "format": "json",
                },
            )
            response.raise_for_status()
            result = response.json()
            response_text = result.get("response", "")

            return self._parse_response(response_text)

        except httpx.HTTPError as e:
            logger.error(f"Ollama HTTP error: {e}")
            return ParseResult(error=f"HTTP error: {e}")
        except Exception as e:
            logger.error(f"Ollama error: {e}")
            return ParseResult(error=str(e))

    def _parse_response(self, response_text: str) -> ParseResult:
        """Parse the JSON response from Ollama."""
        try:
            data = json.loads(response_text)
            # Convert duration to int if present (LLM may return float like 127.5)
            duration = data.get("duration_minutes")
            if duration is not None:
                duration = round(duration)
            return ParseResult(
                duration_minutes=duration,
                work_type_id=data.get("work_type_id"),
                time_code_id=data.get("time_code_id"),
                description=data.get("description"),
                confidence_duration=float(data.get("confidence_duration", 0)),
                confidence_work_type=float(data.get("confidence_work_type", 0)),
                confidence_time_code=float(data.get("confidence_time_code", 0)),
                confidence_overall=float(data.get("confidence_overall", 0)),
                notes=data.get("notes"),
            )
        except json.JSONDecodeError as e:
            logger.error(f"Failed to parse Ollama response: {e}")
            logger.debug(f"Response text: {response_text}")
            return ParseResult(error=f"Invalid JSON response: {e}")
        except (KeyError, TypeError, ValueError) as e:
            logger.error(f"Failed to extract fields from response: {e}")
            return ParseResult(error=f"Failed to extract fields: {e}")
