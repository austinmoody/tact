import logging
import re

from fastapi import APIRouter, Depends, HTTPException, Query
from sqlalchemy.orm import Session

from tact.db.models import WorkType
from tact.db.session import get_session
from tact.schemas.work_type import WorkTypeCreate, WorkTypeResponse, WorkTypeUpdate

logger = logging.getLogger(__name__)

router = APIRouter(prefix="/work-types", tags=["work-types"])


def slugify(name: str) -> str:
    """Convert a name to a URL-friendly slug."""
    slug = name.lower()
    slug = re.sub(r"[^\w\s-]", "", slug)  # Remove special characters
    slug = re.sub(r"[\s_]+", "-", slug)  # Replace spaces/underscores with hyphens
    slug = re.sub(r"-+", "-", slug)  # Collapse multiple hyphens
    return slug.strip("-")


@router.post("", response_model=WorkTypeResponse, status_code=201)
def create_work_type(
    data: WorkTypeCreate,
    session: Session = Depends(get_session),
) -> WorkTypeResponse:
    work_type_id = slugify(data.name)
    existing = session.query(WorkType).filter(WorkType.id == work_type_id).first()
    if existing:
        raise HTTPException(status_code=409, detail="Work type already exists")

    work_type = WorkType(
        id=work_type_id,
        name=data.name,
    )
    session.add(work_type)
    session.commit()
    session.refresh(work_type)
    logger.info("Created work type %s", work_type_id)
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
    logger.info("Listed work types: active=%s count=%d", active, len(work_types))
    return [WorkTypeResponse.model_validate(wt) for wt in work_types]


@router.get("/{work_type_id}", response_model=WorkTypeResponse)
def get_work_type(
    work_type_id: str,
    session: Session = Depends(get_session),
) -> WorkTypeResponse:
    work_type = session.query(WorkType).filter(WorkType.id == work_type_id).first()
    if not work_type:
        raise HTTPException(status_code=404, detail="Work type not found")
    logger.info("Retrieved work type %s", work_type_id)
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
    if data.active is not None:
        work_type.active = data.active

    session.commit()
    session.refresh(work_type)
    logger.info("Updated work type %s", work_type_id)
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
    logger.info("Deleted work type %s", work_type_id)
    return WorkTypeResponse.model_validate(work_type)
