from datetime import datetime

from pydantic import BaseModel


class TimeCodeCreate(BaseModel):
    id: str
    project_id: str
    name: str


class TimeCodeUpdate(BaseModel):
    project_id: str | None = None
    name: str | None = None
    active: bool | None = None


class TimeCodeResponse(BaseModel):
    id: str
    project_id: str
    name: str
    active: bool
    created_at: datetime
    updated_at: datetime

    model_config = {"from_attributes": True}
