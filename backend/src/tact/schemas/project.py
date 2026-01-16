from datetime import datetime

from pydantic import BaseModel


class ProjectCreate(BaseModel):
    id: str
    name: str


class ProjectUpdate(BaseModel):
    name: str | None = None
    active: bool | None = None


class ProjectResponse(BaseModel):
    id: str
    name: str
    active: bool
    created_at: datetime
    updated_at: datetime

    model_config = {"from_attributes": True}
