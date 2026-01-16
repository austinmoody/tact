import logging

from fastapi import APIRouter, Depends, HTTPException, Query
from sqlalchemy.orm import Session

from tact.db.models import Project
from tact.db.session import get_session
from tact.schemas.project import ProjectCreate, ProjectResponse, ProjectUpdate

logger = logging.getLogger(__name__)

router = APIRouter(prefix="/projects", tags=["projects"])


@router.post("", response_model=ProjectResponse, status_code=201)
def create_project(
    data: ProjectCreate,
    session: Session = Depends(get_session),
) -> ProjectResponse:
    existing = session.query(Project).filter(Project.id == data.id).first()
    if existing:
        raise HTTPException(status_code=409, detail="Project already exists")

    project = Project(
        id=data.id,
        name=data.name,
    )
    session.add(project)
    session.commit()
    session.refresh(project)
    logger.info("Created project %s", data.id)
    return ProjectResponse.model_validate(project)


@router.get("", response_model=list[ProjectResponse])
def list_projects(
    active: bool | None = Query(None),
    session: Session = Depends(get_session),
) -> list[ProjectResponse]:
    query = session.query(Project)
    if active is not None:
        query = query.filter(Project.active == active)
    projects = query.all()
    logger.info("Listed projects: active=%s count=%d", active, len(projects))
    return [ProjectResponse.model_validate(p) for p in projects]


@router.get("/{project_id}", response_model=ProjectResponse)
def get_project(
    project_id: str,
    session: Session = Depends(get_session),
) -> ProjectResponse:
    project = session.query(Project).filter(Project.id == project_id).first()
    if not project:
        raise HTTPException(status_code=404, detail="Project not found")
    logger.info("Retrieved project %s", project_id)
    return ProjectResponse.model_validate(project)


@router.put("/{project_id}", response_model=ProjectResponse)
def update_project(
    project_id: str,
    data: ProjectUpdate,
    session: Session = Depends(get_session),
) -> ProjectResponse:
    project = session.query(Project).filter(Project.id == project_id).first()
    if not project:
        raise HTTPException(status_code=404, detail="Project not found")

    if data.name is not None:
        project.name = data.name
    if data.active is not None:
        project.active = data.active

    session.commit()
    session.refresh(project)
    logger.info("Updated project %s", project_id)
    return ProjectResponse.model_validate(project)


@router.delete("/{project_id}", response_model=ProjectResponse)
def delete_project(
    project_id: str,
    session: Session = Depends(get_session),
) -> ProjectResponse:
    project = session.query(Project).filter(Project.id == project_id).first()
    if not project:
        raise HTTPException(status_code=404, detail="Project not found")

    project.active = False
    session.commit()
    session.refresh(project)
    logger.info("Deleted project %s", project_id)
    return ProjectResponse.model_validate(project)
