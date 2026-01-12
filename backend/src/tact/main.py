import asyncio
import logging
import os
import time
from collections.abc import AsyncGenerator
from contextlib import asynccontextmanager

from fastapi import FastAPI, Request

from tact.db.migrations import run_migrations
from tact.routes.entries import router as entries_router
from tact.routes.health import router as health_router
from tact.routes.time_codes import router as time_codes_router
from tact.routes.work_types import router as work_types_router
from tact.worker import start_parser_worker

logger = logging.getLogger(__name__)


def configure_logging() -> None:
    """Configure structured logging for the application."""
    import sys

    log_format = "%(asctime)s %(levelname)s [%(name)s] %(message)s"
    date_format = "%Y-%m-%dT%H:%M:%S%z"

    # Configure root logger
    root_logger = logging.getLogger()
    root_logger.setLevel(logging.INFO)

    # Remove existing handlers and add our own
    for handler in root_logger.handlers[:]:
        root_logger.removeHandler(handler)

    handler = logging.StreamHandler(sys.stdout)
    handler.setLevel(logging.INFO)
    handler.setFormatter(logging.Formatter(log_format, date_format))
    root_logger.addHandler(handler)

    # Ensure all loggers propagate to root and are enabled
    # Include all tact.* loggers that may have been disabled
    logger_names = ["uvicorn", "uvicorn.access", "uvicorn.error"]
    for name in list(logging.Logger.manager.loggerDict.keys()):
        if name.startswith("tact"):
            logger_names.append(name)

    for name in logger_names:
        named_logger = logging.getLogger(name)
        named_logger.handlers = []
        named_logger.propagate = True
        named_logger.setLevel(logging.INFO)
        named_logger.disabled = False


@asynccontextmanager
async def lifespan(app: FastAPI) -> AsyncGenerator[None, None]:
    configure_logging()
    run_migrations()

    # Start parser worker unless disabled
    worker_task = None
    if os.getenv("TACT_DISABLE_WORKER", "").lower() != "true":
        logger.info("Starting parser worker...")
        worker_task = asyncio.create_task(start_parser_worker())

    yield

    # Cancel worker on shutdown
    if worker_task:
        worker_task.cancel()
        try:
            await worker_task
        except asyncio.CancelledError:
            pass


app = FastAPI(title="Tact", version="0.1.0", lifespan=lifespan)


@app.middleware("http")
async def log_requests(request: Request, call_next):
    """Log all HTTP requests with method, path, status, and timing."""
    # Ensure logging is configured (uvicorn may reset it after lifespan)
    root = logging.getLogger()
    if not root.handlers or root.level > logging.INFO:
        configure_logging()

    start_time = time.perf_counter()
    response = await call_next(request)
    duration_ms = (time.perf_counter() - start_time) * 1000

    logging.getLogger("tact.requests").info(
        "%s %s %d %.2fms",
        request.method,
        request.url.path,
        response.status_code,
        duration_ms,
    )
    return response


app.include_router(health_router)
app.include_router(entries_router)
app.include_router(time_codes_router)
app.include_router(work_types_router)
