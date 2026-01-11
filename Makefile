.PHONY: help install run test lint format migrate db-revision docker-build docker-up docker-down tui-build tui-run tui-dev

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

install: ## Install backend dependencies
	cd backend && uv sync --extra dev

run: ## Run backend API (port 2100)
	cd backend && uv run uvicorn tact.main:app --reload --port 2100

test: ## Run backend tests
	cd backend && uv run pytest

lint: ## Check code with ruff
	cd backend && uv run ruff check src tests

format: ## Format code with ruff
	cd backend && uv run ruff format src tests

migrate: ## Run database migrations
	cd backend && uv run alembic upgrade head

db-revision: ## Create migration (usage: make db-revision msg="description")
	cd backend && uv run alembic revision --autogenerate -m "$(msg)"

docker-build: ## Build Docker images
	docker compose build

docker-up: ## Start Docker containers
	docker compose up -d

docker-down: ## Stop Docker containers
	docker compose down

tui-build: ## Build TUI binary
	cd tui && go build -o tact-tui .

tui-run: ## Run TUI
	cd tui && go run .

tui-dev: ## Run TUI (dev mode, localhost:2100)
	cd tui && go run . --api http://localhost:2100
