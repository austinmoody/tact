import logging

from fastapi import APIRouter, Depends, HTTPException
from sqlalchemy.orm import Session

from tact.db.models import ContextDocument, Project, TimeCode
from tact.db.session import get_session
from tact.rag.embeddings import embed_text
from tact.schemas.context import ContextCreate, ContextResponse, ContextUpdate

logger = logging.getLogger(__name__)

router = APIRouter(tags=["context"])


@router.post(
    "/projects/{project_id}/context",
    response_model=ContextResponse,
    status_code=201,
)
def add_project_context(
    project_id: str,
    data: ContextCreate,
    session: Session = Depends(get_session),
) -> ContextResponse:
    project = session.query(Project).filter(Project.id == project_id).first()
    if not project:
        raise HTTPException(status_code=404, detail="Project not found")

    context = ContextDocument(
        project_id=project_id,
        content=data.content,
        embedding=embed_text(data.content),
    )
    session.add(context)
    session.commit()
    session.refresh(context)
    logger.info("Added context to project %s (with embedding)", project_id)
    return ContextResponse.model_validate(context)


@router.get("/projects/{project_id}/context", response_model=list[ContextResponse])
def list_project_context(
    project_id: str,
    session: Session = Depends(get_session),
) -> list[ContextResponse]:
    project = session.query(Project).filter(Project.id == project_id).first()
    if not project:
        raise HTTPException(status_code=404, detail="Project not found")

    contexts = (
        session.query(ContextDocument)
        .filter(ContextDocument.project_id == project_id)
        .all()
    )
    logger.info("Listed %d context documents for project %s", len(contexts), project_id)
    return [ContextResponse.model_validate(c) for c in contexts]


@router.post(
    "/time-codes/{time_code_id}/context",
    response_model=ContextResponse,
    status_code=201,
)
def add_time_code_context(
    time_code_id: str,
    data: ContextCreate,
    session: Session = Depends(get_session),
) -> ContextResponse:
    time_code = session.query(TimeCode).filter(TimeCode.id == time_code_id).first()
    if not time_code:
        raise HTTPException(status_code=404, detail="Time code not found")

    context = ContextDocument(
        time_code_id=time_code_id,
        content=data.content,
        embedding=embed_text(data.content),
    )
    session.add(context)
    session.commit()
    session.refresh(context)
    logger.info("Added context to time code %s (with embedding)", time_code_id)
    return ContextResponse.model_validate(context)


@router.get("/time-codes/{time_code_id}/context", response_model=list[ContextResponse])
def list_time_code_context(
    time_code_id: str,
    session: Session = Depends(get_session),
) -> list[ContextResponse]:
    time_code = session.query(TimeCode).filter(TimeCode.id == time_code_id).first()
    if not time_code:
        raise HTTPException(status_code=404, detail="Time code not found")

    contexts = (
        session.query(ContextDocument)
        .filter(ContextDocument.time_code_id == time_code_id)
        .all()
    )
    logger.info(
        "Listed %d context documents for time code %s", len(contexts), time_code_id
    )
    return [ContextResponse.model_validate(c) for c in contexts]


@router.get("/context/{context_id}", response_model=ContextResponse)
def get_context(
    context_id: str,
    session: Session = Depends(get_session),
) -> ContextResponse:
    context = (
        session.query(ContextDocument).filter(ContextDocument.id == context_id).first()
    )
    if not context:
        raise HTTPException(status_code=404, detail="Context document not found")
    logger.info("Retrieved context %s", context_id)
    return ContextResponse.model_validate(context)


@router.put("/context/{context_id}", response_model=ContextResponse)
def update_context(
    context_id: str,
    data: ContextUpdate,
    session: Session = Depends(get_session),
) -> ContextResponse:
    context = (
        session.query(ContextDocument).filter(ContextDocument.id == context_id).first()
    )
    if not context:
        raise HTTPException(status_code=404, detail="Context document not found")

    context.content = data.content
    # Regenerate embedding for the new content
    context.embedding = embed_text(data.content)

    session.commit()
    session.refresh(context)
    logger.info("Updated context %s (embedding regenerated)", context_id)
    return ContextResponse.model_validate(context)


@router.delete("/context/{context_id}", response_model=ContextResponse)
def delete_context(
    context_id: str,
    session: Session = Depends(get_session),
) -> ContextResponse:
    context = (
        session.query(ContextDocument).filter(ContextDocument.id == context_id).first()
    )
    if not context:
        raise HTTPException(status_code=404, detail="Context document not found")

    session.delete(context)
    session.commit()
    logger.info("Deleted context %s", context_id)
    return ContextResponse.model_validate(context)
