from datetime import date, datetime

from pydantic import BaseModel


class EntryCreate(BaseModel):
    user_input: str
    entry_date: date | None = None


class EntryUpdate(BaseModel):
    user_input: str | None = None
    duration_minutes: int | None = None
    work_type_id: str | None = None
    time_code_id: str | None = None
    parsed_description: str | None = None
    entry_date: date | None = None
    status: str | None = None


class EntryResponse(BaseModel):
    id: str
    user_input: str
    duration_minutes: int | None
    work_type_id: str | None
    time_code_id: str | None
    parsed_description: str | None
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
