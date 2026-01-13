from datetime import datetime

from pydantic import BaseModel


class ContextCreate(BaseModel):
    content: str


class ContextUpdate(BaseModel):
    content: str


class ContextResponse(BaseModel):
    id: str
    project_id: str | None
    time_code_id: str | None
    content: str
    created_at: datetime
    updated_at: datetime

    model_config = {"from_attributes": True}
