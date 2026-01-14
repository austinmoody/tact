import pytest
from fastapi import FastAPI
from fastapi.testclient import TestClient
from sqlalchemy import create_engine, event
from sqlalchemy.orm import sessionmaker
from sqlalchemy.pool import StaticPool

from tact.db.base import Base
from tact.db.session import get_session
from tact.routes.context import router as context_router
from tact.routes.projects import router as projects_router
from tact.routes.time_codes import router as time_codes_router


@pytest.fixture
def client():
    """Create a test client with an in-memory database."""
    engine = create_engine(
        "sqlite:///:memory:",
        connect_args={"check_same_thread": False},
        poolclass=StaticPool,
    )

    @event.listens_for(engine, "connect")
    def set_sqlite_pragma(dbapi_conn, connection_record):
        cursor = dbapi_conn.cursor()
        cursor.execute("PRAGMA foreign_keys=ON")
        cursor.close()

    Base.metadata.create_all(bind=engine)

    TestSession = sessionmaker(autocommit=False, autoflush=False, bind=engine)

    def override_get_session():
        session = TestSession()
        try:
            yield session
        finally:
            session.close()

    test_app = FastAPI()
    test_app.include_router(context_router)
    test_app.include_router(projects_router)
    test_app.include_router(time_codes_router)
    test_app.dependency_overrides[get_session] = override_get_session

    with TestClient(test_app) as test_client:
        yield test_client

    engine.dispose()


@pytest.fixture
def project(client):
    """Create a default project for tests."""
    response = client.post("/projects", json={"id": "IZG", "name": "IZ Gateway"})
    assert response.status_code == 201
    return response.json()


@pytest.fixture
def time_code(client, project):
    """Create a default time code for tests."""
    response = client.post(
        "/time-codes",
        json={
            "id": "FEDS-163",
            "project_id": "IZG",
            "name": "Development",
            "description": "Development work",
        },
    )
    assert response.status_code == 201
    return response.json()


# Project context tests
def test_add_project_context(client, project):
    response = client.post(
        "/projects/IZG/context",
        json={"content": "ALL meetings with APHL go to FEDS-163"},
    )
    assert response.status_code == 201
    data = response.json()
    assert data["project_id"] == "IZG"
    assert data["time_code_id"] is None
    assert data["content"] == "ALL meetings with APHL go to FEDS-163"
    assert "id" in data


def test_add_project_context_not_found(client):
    response = client.post(
        "/projects/UNKNOWN/context",
        json={"content": "Some context"},
    )
    assert response.status_code == 404


def test_list_project_context(client, project):
    # Add multiple context documents
    client.post("/projects/IZG/context", json={"content": "Context 1"})
    client.post("/projects/IZG/context", json={"content": "Context 2"})

    response = client.get("/projects/IZG/context")
    assert response.status_code == 200
    data = response.json()
    assert len(data) == 2


def test_list_project_context_empty(client, project):
    response = client.get("/projects/IZG/context")
    assert response.status_code == 200
    data = response.json()
    assert len(data) == 0


def test_list_project_context_not_found(client):
    response = client.get("/projects/UNKNOWN/context")
    assert response.status_code == 404


# Time code context tests
def test_add_time_code_context(client, time_code):
    response = client.post(
        "/time-codes/FEDS-163/context",
        json={"content": "ALL deployments go to this code"},
    )
    assert response.status_code == 201
    data = response.json()
    assert data["project_id"] is None
    assert data["time_code_id"] == "FEDS-163"
    assert data["content"] == "ALL deployments go to this code"


def test_add_time_code_context_not_found(client):
    response = client.post(
        "/time-codes/UNKNOWN/context",
        json={"content": "Some context"},
    )
    assert response.status_code == 404


def test_list_time_code_context(client, time_code):
    client.post("/time-codes/FEDS-163/context", json={"content": "Context 1"})
    client.post("/time-codes/FEDS-163/context", json={"content": "Context 2"})

    response = client.get("/time-codes/FEDS-163/context")
    assert response.status_code == 200
    data = response.json()
    assert len(data) == 2


def test_list_time_code_context_not_found(client):
    response = client.get("/time-codes/UNKNOWN/context")
    assert response.status_code == 404


# Context document CRUD tests
def test_get_context(client, project):
    create_response = client.post(
        "/projects/IZG/context",
        json={"content": "Test content"},
    )
    context_id = create_response.json()["id"]

    response = client.get(f"/context/{context_id}")
    assert response.status_code == 200
    assert response.json()["content"] == "Test content"


def test_get_context_not_found(client):
    response = client.get("/context/nonexistent-id")
    assert response.status_code == 404


def test_update_context(client, project):
    create_response = client.post(
        "/projects/IZG/context",
        json={"content": "Original content"},
    )
    context_id = create_response.json()["id"]

    response = client.put(
        f"/context/{context_id}",
        json={"content": "Updated content"},
    )
    assert response.status_code == 200
    assert response.json()["content"] == "Updated content"


def test_update_context_not_found(client):
    response = client.put(
        "/context/nonexistent-id",
        json={"content": "Updated content"},
    )
    assert response.status_code == 404


def test_delete_context(client, project):
    create_response = client.post(
        "/projects/IZG/context",
        json={"content": "To be deleted"},
    )
    context_id = create_response.json()["id"]

    response = client.delete(f"/context/{context_id}")
    assert response.status_code == 200

    # Verify it's gone
    get_response = client.get(f"/context/{context_id}")
    assert get_response.status_code == 404


def test_delete_context_not_found(client):
    response = client.delete("/context/nonexistent-id")
    assert response.status_code == 404


# Test isolation between project and time code contexts
def test_contexts_are_separate(client, project, time_code):
    # Add context to project
    client.post("/projects/IZG/context", json={"content": "Project context"})

    # Add context to time code
    client.post("/time-codes/FEDS-163/context", json={"content": "Time code context"})

    # Verify they're separate
    project_contexts = client.get("/projects/IZG/context").json()
    time_code_contexts = client.get("/time-codes/FEDS-163/context").json()

    assert len(project_contexts) == 1
    assert len(time_code_contexts) == 1
    assert project_contexts[0]["content"] == "Project context"
    assert time_code_contexts[0]["content"] == "Time code context"
