from datetime import date, datetime

from pydantic import BaseModel


class EntryCreate(BaseModel):
    raw_text: str
    entry_date: date | None = None


class EntryUpdate(BaseModel):
    raw_text: str | None = None
    duration_minutes: int | None = None
    work_type_id: str | None = None
    time_code_id: str | None = None
    description: str | None = None
    entry_date: date | None = None
    status: str | None = None


class EntryResponse(BaseModel):
    id: str
    raw_text: str
    duration_minutes: int | None
    work_type_id: str | None
    time_code_id: str | None
    description: str | None
    entry_date: date | None
    confidence_duration: float | None
    confidence_work_type: float | None
    confidence_time_code: float | None
    confidence_overall: float | None
    status: str
    parse_error: str | None
    parse_notes: str | None
    manually_corrected: bool
    locked: bool
    corrected_at: datetime | None
    created_at: datetime
    parsed_at: datetime | None
    updated_at: datetime

    model_config = {"from_attributes": True}
