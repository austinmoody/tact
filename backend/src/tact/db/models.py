import uuid
from datetime import UTC, date, datetime

from sqlalchemy import ForeignKey, Text
from sqlalchemy.orm import Mapped, mapped_column

from tact.db.base import Base


def utc_now() -> datetime:
    return datetime.now(UTC)


class TimeEntry(Base):
    __tablename__ = "time_entries"

    id: Mapped[str] = mapped_column(primary_key=True, default=lambda: str(uuid.uuid4()))
    raw_text: Mapped[str] = mapped_column(Text, nullable=False)

    # Parsed fields
    duration_minutes: Mapped[int | None] = mapped_column(default=None)
    work_type_id: Mapped[str | None] = mapped_column(
        ForeignKey("work_types.id"), default=None
    )
    time_code_id: Mapped[str | None] = mapped_column(
        ForeignKey("time_codes.id"), default=None
    )
    description: Mapped[str | None] = mapped_column(Text, default=None)
    entry_date: Mapped[date | None] = mapped_column(default=None)

    # Confidence scores
    confidence_duration: Mapped[float | None] = mapped_column(default=None)
    confidence_work_type: Mapped[float | None] = mapped_column(default=None)
    confidence_time_code: Mapped[float | None] = mapped_column(default=None)
    confidence_overall: Mapped[float | None] = mapped_column(default=None)

    # Status
    status: Mapped[str] = mapped_column(default="pending")
    parse_error: Mapped[str | None] = mapped_column(Text, default=None)

    # Correction tracking
    manually_corrected: Mapped[bool] = mapped_column(default=False)
    locked: Mapped[bool] = mapped_column(default=False)
    corrected_at: Mapped[datetime | None] = mapped_column(default=None)

    # Timestamps
    created_at: Mapped[datetime] = mapped_column(default=utc_now)
    parsed_at: Mapped[datetime | None] = mapped_column(default=None)
    updated_at: Mapped[datetime] = mapped_column(
        default=utc_now, onupdate=utc_now
    )


class TimeCode(Base):
    __tablename__ = "time_codes"

    id: Mapped[str] = mapped_column(primary_key=True)
    name: Mapped[str] = mapped_column(nullable=False)
    description: Mapped[str] = mapped_column(Text, nullable=False)
    keywords: Mapped[str] = mapped_column(Text, default="[]")
    examples: Mapped[str] = mapped_column(Text, default="[]")
    active: Mapped[bool] = mapped_column(default=True)
    created_at: Mapped[datetime] = mapped_column(default=utc_now)
    updated_at: Mapped[datetime] = mapped_column(
        default=utc_now, onupdate=utc_now
    )


class WorkType(Base):
    __tablename__ = "work_types"

    id: Mapped[str] = mapped_column(primary_key=True)
    name: Mapped[str] = mapped_column(nullable=False)
    description: Mapped[str | None] = mapped_column(Text, default=None)
    active: Mapped[bool] = mapped_column(default=True)
    created_at: Mapped[datetime] = mapped_column(default=utc_now)
    updated_at: Mapped[datetime] = mapped_column(
        default=utc_now, onupdate=utc_now
    )


class Config(Base):
    __tablename__ = "config"

    key: Mapped[str] = mapped_column(primary_key=True)
    value: Mapped[str] = mapped_column(Text, nullable=False)
    updated_at: Mapped[datetime] = mapped_column(
        default=utc_now, onupdate=utc_now
    )
