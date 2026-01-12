FROM python:3.12-slim

WORKDIR /app

COPY --from=ghcr.io/astral-sh/uv:latest /uv /uvx /bin/

COPY pyproject.toml ./
RUN uv sync --frozen --no-dev --no-install-project 2>/dev/null || uv sync --no-dev --no-install-project

COPY src/ ./src/
RUN uv sync --frozen --no-dev 2>/dev/null || uv sync --no-dev

# MCP servers use stdio transport - no ports needed
CMD ["uv", "run", "tact-mcp"]
