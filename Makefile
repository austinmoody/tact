.PHONY: install run test lint format migrate db-revision docker-build docker-up docker-down

install:
	cd backend && uv sync --extra dev

run:
	cd backend && uv run uvicorn tact.main:app --reload --port 2100

test:
	cd backend && uv run pytest

lint:
	cd backend && uv run ruff check src tests

format:
	cd backend && uv run ruff format src tests

migrate:
	cd backend && uv run alembic upgrade head

db-revision:
	cd backend && uv run alembic revision --autogenerate -m "$(msg)"

docker-build:
	docker compose build

docker-up:
	docker compose up -d

docker-down:
	docker compose down
