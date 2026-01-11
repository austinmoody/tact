import asyncio
import logging
import os
from collections.abc import AsyncGenerator
from contextlib import asynccontextmanager

from fastapi import FastAPI

from tact.db.migrations import run_migrations
from tact.routes.entries import router as entries_router
from tact.routes.health import router as health_router
from tact.routes.time_codes import router as time_codes_router
from tact.routes.work_types import router as work_types_router
from tact.worker import start_parser_worker

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)


@asynccontextmanager
async def lifespan(app: FastAPI) -> AsyncGenerator[None, None]:
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

app.include_router(health_router)
app.include_router(entries_router)
app.include_router(time_codes_router)
app.include_router(work_types_router)
