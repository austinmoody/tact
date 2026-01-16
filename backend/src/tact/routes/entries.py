import logging
from datetime import date

from fastapi import APIRouter, Depends, HTTPException, Query, Response
from sqlalchemy.orm import Session

from tact.db.models import ContextDocument, TimeEntry
from tact.db.session import get_session
from tact.rag.embeddings import embed_text
from tact.schemas.entry import EntryCreate, EntryResponse, EntryUpdate

logger = logging.getLogger(__name__)

router = APIRouter(prefix="/entries", tags=["entries"])


@router.post("", response_model=EntryResponse, status_code=201)
def create_entry(
    data: EntryCreate,
    session: Session = Depends(get_session),
) -> EntryResponse:
    entry = TimeEntry(
        user_input=data.user_input,
        entry_date=data.entry_date if data.entry_date else date.today(),
        status="pending",
    )
    session.add(entry)
    session.commit()
    session.refresh(entry)
    logger.info("Created entry %s: %s", entry.id, data.user_input[:50])
    return EntryResponse.model_validate(entry)


@router.get("", response_model=list[EntryResponse])
def list_entries(
    status: str | None = Query(None),
    time_code_id: str | None = Query(None),
    work_type_id: str | None = Query(None),
    from_date: date | None = Query(None),
    to_date: date | None = Query(None),
    limit: int = Query(100, ge=1, le=1000),
    offset: int = Query(0, ge=0),
    session: Session = Depends(get_session),
) -> list[EntryResponse]:
    query = session.query(TimeEntry)

    if status is not None:
        query = query.filter(TimeEntry.status == status)
    if time_code_id is not None:
        query = query.filter(TimeEntry.time_code_id == time_code_id)
    if work_type_id is not None:
        query = query.filter(TimeEntry.work_type_id == work_type_id)
    if from_date is not None:
        query = query.filter(TimeEntry.entry_date >= from_date)
    if to_date is not None:
        query = query.filter(TimeEntry.entry_date <= to_date)

    query = query.offset(offset).limit(limit)
    entries = query.all()

    filters = {
        k: v
        for k, v in [
            ("status", status),
            ("time_code_id", time_code_id),
            ("work_type_id", work_type_id),
            ("from_date", from_date),
            ("to_date", to_date),
        ]
        if v is not None
    }
    logger.info("Listed entries: filters=%s count=%d", filters, len(entries))
    return [EntryResponse.model_validate(e) for e in entries]


@router.get("/{entry_id}", response_model=EntryResponse)
def get_entry(
    entry_id: str,
    session: Session = Depends(get_session),
) -> EntryResponse:
    entry = session.query(TimeEntry).filter(TimeEntry.id == entry_id).first()
    if not entry:
        raise HTTPException(status_code=404, detail="Entry not found")
    logger.info("Retrieved entry %s", entry_id)
    return EntryResponse.model_validate(entry)


@router.patch("/{entry_id}", response_model=EntryResponse)
def update_entry(
    entry_id: str,
    data: EntryUpdate,
    learn: bool = Query(True, description="Create context document from correction"),
    session: Session = Depends(get_session),
) -> EntryResponse:
    entry = session.query(TimeEntry).filter(TimeEntry.id == entry_id).first()
    if not entry:
        raise HTTPException(status_code=404, detail="Entry not found")

    update_data = data.model_dump(exclude_unset=True)
    if update_data:
        for field, value in update_data.items():
            setattr(entry, field, value)
        entry.manually_corrected = True

    session.commit()
    session.refresh(entry)

    # Create context document for learning if enabled and entry has time_code_id
    if learn and entry.time_code_id:
        _create_learned_context(entry, session)

    logger.info("Updated entry %s: fields=%s", entry_id, list(update_data.keys()))
    return EntryResponse.model_validate(entry)


def _create_learned_context(entry: TimeEntry, session: Session) -> None:
    """Create a context document from a manually corrected entry."""
    # Build the content from the entry
    parts = [f'Example: "{entry.user_input}"']

    parsed_parts = []
    if entry.duration_minutes is not None:
        parsed_parts.append(f"{entry.duration_minutes} minutes")
    if entry.work_type_id is not None:
        parsed_parts.append(f"work_type: {entry.work_type_id}")

    if parsed_parts:
        parts.append(f"Parsed as: {', '.join(parsed_parts)}")

    content = "\n".join(parts)

    # Create the context document
    context = ContextDocument(
        time_code_id=entry.time_code_id,
        content=content,
        embedding=embed_text(content),
    )
    session.add(context)
    session.commit()

    logger.info(
        "Created learned context for time code %s from entry %s",
        entry.time_code_id,
        entry.id,
    )


@router.delete("/{entry_id}", status_code=204)
def delete_entry(
    entry_id: str,
    session: Session = Depends(get_session),
) -> Response:
    entry = session.query(TimeEntry).filter(TimeEntry.id == entry_id).first()
    if not entry:
        raise HTTPException(status_code=404, detail="Entry not found")

    session.delete(entry)
    session.commit()
    logger.info("Deleted entry %s", entry_id)
    return Response(status_code=204)


@router.post("/{entry_id}/reparse", response_model=EntryResponse)
def reparse_entry(
    entry_id: str,
    session: Session = Depends(get_session),
) -> EntryResponse:
    """Reset an entry to pending status for re-parsing."""
    entry = session.query(TimeEntry).filter(TimeEntry.id == entry_id).first()
    if not entry:
        raise HTTPException(status_code=404, detail="Entry not found")

    # Clear parsed fields
    entry.duration_minutes = None
    entry.work_type_id = None
    entry.time_code_id = None
    entry.parsed_description = None
    entry.confidence_duration = None
    entry.confidence_work_type = None
    entry.confidence_time_code = None
    entry.confidence_overall = None
    entry.parsed_at = None
    entry.parse_error = None

    # Reset status
    entry.status = "pending"
    entry.manually_corrected = False

    session.commit()
    session.refresh(entry)
    logger.info("Reparse requested for entry %s", entry_id)
    return EntryResponse.model_validate(entry)
