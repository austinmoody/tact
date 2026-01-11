from datetime import date

from fastapi import APIRouter, Depends, HTTPException, Query, Response
from sqlalchemy.orm import Session

from tact.db.models import TimeEntry
from tact.db.session import get_session
from tact.schemas.entry import EntryCreate, EntryResponse, EntryUpdate

router = APIRouter(prefix="/entries", tags=["entries"])


@router.post("", response_model=EntryResponse, status_code=201)
def create_entry(
    data: EntryCreate,
    session: Session = Depends(get_session),
) -> EntryResponse:
    entry = TimeEntry(
        raw_text=data.raw_text,
        entry_date=data.entry_date if data.entry_date else date.today(),
        status="pending",
    )
    session.add(entry)
    session.commit()
    session.refresh(entry)
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
    return [EntryResponse.model_validate(e) for e in entries]


@router.get("/{entry_id}", response_model=EntryResponse)
def get_entry(
    entry_id: str,
    session: Session = Depends(get_session),
) -> EntryResponse:
    entry = session.query(TimeEntry).filter(TimeEntry.id == entry_id).first()
    if not entry:
        raise HTTPException(status_code=404, detail="Entry not found")
    return EntryResponse.model_validate(entry)


@router.patch("/{entry_id}", response_model=EntryResponse)
def update_entry(
    entry_id: str,
    data: EntryUpdate,
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
    return EntryResponse.model_validate(entry)


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
    return Response(status_code=204)
