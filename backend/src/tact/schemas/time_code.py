from datetime import datetime

from pydantic import BaseModel


class TimeCodeCreate(BaseModel):
    id: str
    name: str
    description: str
    keywords: list[str] = []
    examples: list[str] = []


class TimeCodeUpdate(BaseModel):
    name: str | None = None
    description: str | None = None
    keywords: list[str] | None = None
    examples: list[str] | None = None
    active: bool | None = None


class TimeCodeResponse(BaseModel):
    id: str
    name: str
    description: str
    keywords: list[str]
    examples: list[str]
    active: bool
    created_at: datetime
    updated_at: datetime

    model_config = {"from_attributes": True}
