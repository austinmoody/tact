from datetime import datetime

from pydantic import BaseModel


class WorkTypeCreate(BaseModel):
    name: str


class WorkTypeUpdate(BaseModel):
    name: str | None = None
    active: bool | None = None


class WorkTypeResponse(BaseModel):
    id: str
    name: str
    active: bool
    created_at: datetime
    updated_at: datetime

    model_config = {"from_attributes": True}
