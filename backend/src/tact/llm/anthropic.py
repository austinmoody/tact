import json
import logging
import os

import anthropic

from tact.llm.prompts import build_system_prompt, build_user_prompt
from tact.llm.provider import LLMProvider, ParseContext, ParseResult

logger = logging.getLogger(__name__)

DEFAULT_ANTHROPIC_MODEL = "claude-3-haiku-20240307"


class AnthropicProvider(LLMProvider):
    """LLM provider using Anthropic's Claude API."""

    def __init__(
        self,
        api_key: str | None = None,
        model: str | None = None,
    ):
        self.api_key = api_key or os.getenv("TACT_ANTHROPIC_API_KEY")
        if not self.api_key:
            raise ValueError(
                "Anthropic API key required. Set TACT_ANTHROPIC_API_KEY environment "
                "variable or pass api_key parameter."
            )
        self.model = model or os.getenv(
            "TACT_ANTHROPIC_MODEL", DEFAULT_ANTHROPIC_MODEL
        )
        self.client = anthropic.Anthropic(api_key=self.api_key)

    def parse(self, user_input: str, context: ParseContext) -> ParseResult:
        """Parse user input using Anthropic's Claude API."""
        system_prompt = build_system_prompt(context)
        user_prompt = build_user_prompt(user_input)

        try:
            message = self.client.messages.create(
                model=self.model,
                max_tokens=1024,
                system=system_prompt,
                messages=[
                    {"role": "user", "content": user_prompt}
                ],
            )

            response_text = message.content[0].text
            return self._parse_response(response_text)

        except anthropic.APIConnectionError as e:
            logger.error(f"Anthropic connection error: {e}")
            return ParseResult(error=f"Connection error: {e}")
        except anthropic.RateLimitError as e:
            logger.error(f"Anthropic rate limit error: {e}")
            return ParseResult(error=f"Rate limit exceeded: {e}")
        except anthropic.APIStatusError as e:
            logger.error(f"Anthropic API error: {e}")
            return ParseResult(error=f"API error: {e}")
        except Exception as e:
            logger.error(f"Anthropic error: {e}")
            return ParseResult(error=str(e))

    def _parse_response(self, response_text: str) -> ParseResult:
        """Parse the JSON response from Anthropic."""
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
            )
        except json.JSONDecodeError as e:
            logger.error(f"Failed to parse Anthropic response: {e}")
            logger.debug(f"Response text: {response_text}")
            return ParseResult(error=f"Invalid JSON response: {e}")
        except (KeyError, TypeError, ValueError) as e:
            logger.error(f"Failed to extract fields from response: {e}")
            return ParseResult(error=f"Failed to extract fields: {e}")
