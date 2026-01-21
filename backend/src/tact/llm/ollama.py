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
DEFAULT_OLLAMA_PULL_TIMEOUT = 600.0  # seconds (10 minutes for model downloads)

# JSON schema for structured output - constrains LLM response to exact field structure
PARSE_RESULT_SCHEMA = {
    "type": "object",
    "properties": {
        "duration_minutes": {"type": "integer"},
        "work_type_id": {"type": ["string", "null"]},
        "time_code_id": {"type": ["string", "null"]},
        "parsed_description": {"type": ["string", "null"]},
        "confidence_duration": {"type": "number"},
        "confidence_work_type": {"type": "number"},
        "confidence_time_code": {"type": "number"},
        "confidence_overall": {"type": "number"},
        "notes": {"type": ["string", "null"]},
    },
    "required": ["confidence_overall"],
}


class OllamaProvider(LLMProvider):
    """LLM provider using Ollama for local model inference."""

    def __init__(
        self,
        base_url: str | None = None,
        model: str | None = None,
        timeout: float | None = None,
        pull_timeout: float | None = None,
    ):
        self.base_url = base_url or os.getenv("TACT_OLLAMA_URL", DEFAULT_OLLAMA_URL)
        self.model = model or os.getenv("TACT_OLLAMA_MODEL", DEFAULT_OLLAMA_MODEL)
        self.timeout = timeout or float(
            os.getenv("TACT_OLLAMA_TIMEOUT", DEFAULT_OLLAMA_TIMEOUT)
        )
        self.pull_timeout = pull_timeout or float(
            os.getenv("TACT_OLLAMA_PULL_TIMEOUT", DEFAULT_OLLAMA_PULL_TIMEOUT)
        )
        self.client = httpx.Client(timeout=self.timeout)
        self._model_verified = False

    def _ensure_model_available(self) -> str | None:
        """Check if model exists and pull if needed. Returns error message or None."""
        if self._model_verified:
            return None

        try:
            # Check available models
            response = self.client.get(f"{self.base_url}/api/tags")
            response.raise_for_status()
            tags = response.json()
            models = tags.get("models", [])

            # Check if our model is available (match by name prefix)
            model_available = any(
                m.get("name", "").startswith(self.model.split(":")[0])
                and (
                    ":" not in self.model
                    or m.get("name", "") == self.model
                    or m.get("name", "").startswith(self.model)
                )
                for m in models
            )

            if model_available:
                self._model_verified = True
                return None

            # Model not found, pull it
            logger.info(f"Model '{self.model}' not found, pulling...")
            pull_client = httpx.Client(timeout=self.pull_timeout)
            try:
                # Use streaming to handle the chunked response from Ollama pull
                with pull_client.stream(
                    "POST",
                    f"{self.base_url}/api/pull",
                    json={"name": self.model},
                ) as pull_response:
                    pull_response.raise_for_status()
                    for line in pull_response.iter_lines():
                        if line:
                            try:
                                status = json.loads(line)
                                if "status" in status:
                                    logger.debug(f"Pull status: {status['status']}")
                                if status.get("error"):
                                    return f"Model pull failed: {status['error']}"
                            except json.JSONDecodeError:
                                pass
            finally:
                pull_client.close()

            logger.info(f"Model '{self.model}' pulled successfully")
            self._model_verified = True
            return None

        except httpx.HTTPError as e:
            logger.error(f"Failed to check/pull model: {e}")
            return f"Failed to ensure model availability: {e}"
        except Exception as e:
            logger.error(f"Unexpected error checking/pulling model: {e}")
            return f"Failed to ensure model availability: {e}"

    def parse(self, user_input: str, context: ParseContext) -> ParseResult:
        """Parse user input using Ollama."""
        # Ensure model is available (auto-pull if needed)
        if error := self._ensure_model_available():
            return ParseResult(error=error)

        system_prompt = build_system_prompt(context)
        user_prompt = build_user_prompt(user_input)

        try:
            response = self.client.post(
                f"{self.base_url}/api/generate",
                json={
                    "model": self.model,
                    "prompt": f"{system_prompt}\n\n{user_prompt}",
                    "stream": False,
                    "format": PARSE_RESULT_SCHEMA,
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
                parsed_description=data.get("parsed_description"),
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
