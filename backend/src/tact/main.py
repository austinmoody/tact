from collections.abc import AsyncGenerator
from contextlib import asynccontextmanager

from fastapi import FastAPI

from tact.db.migrations import run_migrations
from tact.routes.health import router as health_router
from tact.routes.time_codes import router as time_codes_router
from tact.routes.work_types import router as work_types_router


@asynccontextmanager
async def lifespan(app: FastAPI) -> AsyncGenerator[None, None]:
    run_migrations()
    yield


app = FastAPI(title="Tact", version="0.1.0", lifespan=lifespan)

app.include_router(health_router)
app.include_router(time_codes_router)
app.include_router(work_types_router)
