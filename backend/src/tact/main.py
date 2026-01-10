from fastapi import FastAPI

from tact.routes.health import router as health_router

app = FastAPI(title="Tact", version="0.1.0")

app.include_router(health_router)
