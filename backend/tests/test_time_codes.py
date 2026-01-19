import pytest
from fastapi import FastAPI
from fastapi.testclient import TestClient
from sqlalchemy import create_engine, event
from sqlalchemy.orm import sessionmaker
from sqlalchemy.pool import StaticPool

from tact.db.base import Base
from tact.db.models import Config, Project, TimeCode, TimeEntry, WorkType  # noqa: F401
from tact.db.session import get_session
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


def test_create_time_code(client, project):
    response = client.post(
        "/time-codes",
        json={
            "id": "PROJ-001",
            "project_id": "IZG",
            "name": "Project Alpha",
        },
    )
    assert response.status_code == 201
    data = response.json()
    assert data["id"] == "PROJ-001"
    assert data["project_id"] == "IZG"
    assert data["name"] == "Project Alpha"
    assert data["active"] is True


def test_create_time_code_project_not_found(client):
    response = client.post(
        "/time-codes",
        json={
            "id": "PROJ-001",
            "project_id": "UNKNOWN",
            "name": "Project",
        },
    )
    assert response.status_code == 400
    assert response.json()["detail"] == "Project not found"


def test_create_duplicate_time_code(client, project):
    client.post(
        "/time-codes",
        json={
            "id": "PROJ-001",
            "project_id": "IZG",
            "name": "Project",
        },
    )
    response = client.post(
        "/time-codes",
        json={
            "id": "PROJ-001",
            "project_id": "IZG",
            "name": "Another",
        },
    )
    assert response.status_code == 409


def test_list_time_codes(client, project):
    client.post(
        "/time-codes",
        json={
            "id": "PROJ-001",
            "project_id": "IZG",
            "name": "Project 1",
        },
    )
    client.post(
        "/time-codes",
        json={
            "id": "PROJ-002",
            "project_id": "IZG",
            "name": "Project 2",
        },
    )
    response = client.get("/time-codes")
    assert response.status_code == 200
    data = response.json()
    assert len(data) == 2


def test_list_time_codes_filter_active(client, project):
    client.post(
        "/time-codes",
        json={
            "id": "PROJ-001",
            "project_id": "IZG",
            "name": "Active",
        },
    )
    client.post(
        "/time-codes",
        json={
            "id": "PROJ-002",
            "project_id": "IZG",
            "name": "Will be inactive",
        },
    )
    client.delete("/time-codes/PROJ-002")

    response = client.get("/time-codes?active=true")
    assert response.status_code == 200
    data = response.json()
    assert len(data) == 1
    assert data[0]["id"] == "PROJ-001"


def test_list_time_codes_filter_by_project(client):
    # Create two projects
    client.post("/projects", json={"id": "IZG", "name": "IZ Gateway"})
    client.post("/projects", json={"id": "TESTME", "name": "TestMe"})

    # Create time codes in each project
    client.post(
        "/time-codes",
        json={
            "id": "IZG-001",
            "project_id": "IZG",
            "name": "IZG Code",
        },
    )
    client.post(
        "/time-codes",
        json={
            "id": "TESTME-001",
            "project_id": "TESTME",
            "name": "TestMe Code",
        },
    )

    # Filter by project
    response = client.get("/time-codes?project_id=IZG")
    assert response.status_code == 200
    data = response.json()
    assert len(data) == 1
    assert data[0]["id"] == "IZG-001"


def test_get_time_code(client, project):
    client.post(
        "/time-codes",
        json={
            "id": "PROJ-001",
            "project_id": "IZG",
            "name": "Project",
        },
    )
    response = client.get("/time-codes/PROJ-001")
    assert response.status_code == 200
    data = response.json()
    assert data["id"] == "PROJ-001"
    assert data["project_id"] == "IZG"


def test_get_time_code_not_found(client):
    response = client.get("/time-codes/UNKNOWN")
    assert response.status_code == 404


def test_update_time_code(client, project):
    client.post(
        "/time-codes",
        json={
            "id": "PROJ-001",
            "project_id": "IZG",
            "name": "Original",
        },
    )
    response = client.put(
        "/time-codes/PROJ-001",
        json={"name": "Updated"},
    )
    assert response.status_code == 200
    assert response.json()["name"] == "Updated"


def test_update_time_code_change_project(client):
    # Create two projects
    client.post("/projects", json={"id": "IZG", "name": "IZ Gateway"})
    client.post("/projects", json={"id": "OTHER", "name": "Other Project"})

    # Create time code in first project
    client.post(
        "/time-codes",
        json={
            "id": "PROJ-001",
            "project_id": "IZG",
            "name": "Original",
        },
    )

    # Move to second project
    response = client.put(
        "/time-codes/PROJ-001",
        json={"project_id": "OTHER"},
    )
    assert response.status_code == 200
    assert response.json()["project_id"] == "OTHER"


def test_update_time_code_invalid_project(client, project):
    client.post(
        "/time-codes",
        json={
            "id": "PROJ-001",
            "project_id": "IZG",
            "name": "Original",
        },
    )
    response = client.put(
        "/time-codes/PROJ-001",
        json={"project_id": "UNKNOWN"},
    )
    assert response.status_code == 400
    assert response.json()["detail"] == "Project not found"


def test_update_time_code_not_found(client):
    response = client.put("/time-codes/UNKNOWN", json={"name": "Test"})
    assert response.status_code == 404


def test_delete_time_code(client, project):
    client.post(
        "/time-codes",
        json={
            "id": "PROJ-001",
            "project_id": "IZG",
            "name": "Project",
        },
    )
    response = client.delete("/time-codes/PROJ-001")
    assert response.status_code == 200
    assert response.json()["active"] is False

    # Verify it's still there but inactive
    get_response = client.get("/time-codes/PROJ-001")
    assert get_response.json()["active"] is False


def test_delete_time_code_not_found(client):
    response = client.delete("/time-codes/UNKNOWN")
    assert response.status_code == 404
