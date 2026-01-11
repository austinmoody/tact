import json

from fastapi import APIRouter, Depends, HTTPException, Query
from sqlalchemy.orm import Session

from tact.db.models import TimeCode
from tact.db.session import get_session
from tact.schemas.time_code import TimeCodeCreate, TimeCodeResponse, TimeCodeUpdate

router = APIRouter(prefix="/time-codes", tags=["time-codes"])


def _model_to_response(time_code: TimeCode) -> TimeCodeResponse:
    """Convert SQLAlchemy model to response, parsing JSON fields."""
    return TimeCodeResponse(
        id=time_code.id,
        name=time_code.name,
        description=time_code.description,
        keywords=json.loads(time_code.keywords),
        examples=json.loads(time_code.examples),
        active=time_code.active,
        created_at=time_code.created_at,
        updated_at=time_code.updated_at,
    )


@router.post("", response_model=TimeCodeResponse, status_code=201)
def create_time_code(
    data: TimeCodeCreate,
    session: Session = Depends(get_session),
) -> TimeCodeResponse:
    existing = session.query(TimeCode).filter(TimeCode.id == data.id).first()
    if existing:
        raise HTTPException(status_code=409, detail="Time code already exists")

    time_code = TimeCode(
        id=data.id,
        name=data.name,
        description=data.description,
        keywords=json.dumps(data.keywords),
        examples=json.dumps(data.examples),
    )
    session.add(time_code)
    session.commit()
    session.refresh(time_code)
    return _model_to_response(time_code)


@router.get("", response_model=list[TimeCodeResponse])
def list_time_codes(
    active: bool | None = Query(None),
    session: Session = Depends(get_session),
) -> list[TimeCodeResponse]:
    query = session.query(TimeCode)
    if active is not None:
        query = query.filter(TimeCode.active == active)
    time_codes = query.all()
    return [_model_to_response(tc) for tc in time_codes]


@router.get("/{time_code_id}", response_model=TimeCodeResponse)
def get_time_code(
    time_code_id: str,
    session: Session = Depends(get_session),
) -> TimeCodeResponse:
    time_code = session.query(TimeCode).filter(TimeCode.id == time_code_id).first()
    if not time_code:
        raise HTTPException(status_code=404, detail="Time code not found")
    return _model_to_response(time_code)


@router.put("/{time_code_id}", response_model=TimeCodeResponse)
def update_time_code(
    time_code_id: str,
    data: TimeCodeUpdate,
    session: Session = Depends(get_session),
) -> TimeCodeResponse:
    time_code = session.query(TimeCode).filter(TimeCode.id == time_code_id).first()
    if not time_code:
        raise HTTPException(status_code=404, detail="Time code not found")

    if data.name is not None:
        time_code.name = data.name
    if data.description is not None:
        time_code.description = data.description
    if data.keywords is not None:
        time_code.keywords = json.dumps(data.keywords)
    if data.examples is not None:
        time_code.examples = json.dumps(data.examples)
    if data.active is not None:
        time_code.active = data.active

    session.commit()
    session.refresh(time_code)
    return _model_to_response(time_code)


@router.delete("/{time_code_id}", response_model=TimeCodeResponse)
def delete_time_code(
    time_code_id: str,
    session: Session = Depends(get_session),
) -> TimeCodeResponse:
    time_code = session.query(TimeCode).filter(TimeCode.id == time_code_id).first()
    if not time_code:
        raise HTTPException(status_code=404, detail="Time code not found")

    time_code.active = False
    session.commit()
    session.refresh(time_code)
    return _model_to_response(time_code)
