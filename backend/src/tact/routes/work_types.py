from fastapi import APIRouter, Depends, HTTPException, Query
from sqlalchemy.orm import Session

from tact.db.models import WorkType
from tact.db.session import get_session
from tact.schemas.work_type import WorkTypeCreate, WorkTypeResponse, WorkTypeUpdate

router = APIRouter(prefix="/work-types", tags=["work-types"])


@router.post("", response_model=WorkTypeResponse, status_code=201)
def create_work_type(
    data: WorkTypeCreate,
    session: Session = Depends(get_session),
) -> WorkTypeResponse:
    existing = session.query(WorkType).filter(WorkType.id == data.id).first()
    if existing:
        raise HTTPException(status_code=409, detail="Work type already exists")

    work_type = WorkType(
        id=data.id,
        name=data.name,
        description=data.description,
    )
    session.add(work_type)
    session.commit()
    session.refresh(work_type)
    return WorkTypeResponse.model_validate(work_type)


@router.get("", response_model=list[WorkTypeResponse])
def list_work_types(
    active: bool | None = Query(None),
    session: Session = Depends(get_session),
) -> list[WorkTypeResponse]:
    query = session.query(WorkType)
    if active is not None:
        query = query.filter(WorkType.active == active)
    work_types = query.all()
    return [WorkTypeResponse.model_validate(wt) for wt in work_types]


@router.get("/{work_type_id}", response_model=WorkTypeResponse)
def get_work_type(
    work_type_id: str,
    session: Session = Depends(get_session),
) -> WorkTypeResponse:
    work_type = session.query(WorkType).filter(WorkType.id == work_type_id).first()
    if not work_type:
        raise HTTPException(status_code=404, detail="Work type not found")
    return WorkTypeResponse.model_validate(work_type)


@router.put("/{work_type_id}", response_model=WorkTypeResponse)
def update_work_type(
    work_type_id: str,
    data: WorkTypeUpdate,
    session: Session = Depends(get_session),
) -> WorkTypeResponse:
    work_type = session.query(WorkType).filter(WorkType.id == work_type_id).first()
    if not work_type:
        raise HTTPException(status_code=404, detail="Work type not found")

    if data.name is not None:
        work_type.name = data.name
    if data.description is not None:
        work_type.description = data.description
    if data.active is not None:
        work_type.active = data.active

    session.commit()
    session.refresh(work_type)
    return WorkTypeResponse.model_validate(work_type)


@router.delete("/{work_type_id}", response_model=WorkTypeResponse)
def delete_work_type(
    work_type_id: str,
    session: Session = Depends(get_session),
) -> WorkTypeResponse:
    work_type = session.query(WorkType).filter(WorkType.id == work_type_id).first()
    if not work_type:
        raise HTTPException(status_code=404, detail="Work type not found")

    work_type.active = False
    session.commit()
    session.refresh(work_type)
    return WorkTypeResponse.model_validate(work_type)
