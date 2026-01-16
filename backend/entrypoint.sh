#!/bin/sh
set -e

# Run database migrations
echo "Running database migrations..."
uv run alembic upgrade head

# Start the server
echo "Starting server..."
exec uv run uvicorn tact.main:app --host 0.0.0.0 --port 2100
