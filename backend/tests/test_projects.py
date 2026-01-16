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
    test_app.dependency_overrides[get_session] = override_get_session

    with TestClient(test_app) as test_client:
        yield test_client

    engine.dispose()


def test_create_project(client):
    response = client.post(
        "/projects",
        json={
            "id": "IZG",
            "name": "IZ Gateway",
        },
    )
    assert response.status_code == 201
    data = response.json()
    assert data["id"] == "IZG"
    assert data["name"] == "IZ Gateway"
    assert data["active"] is True


def test_create_duplicate_project(client):
    client.post(
        "/projects",
        json={"id": "IZG", "name": "IZ Gateway"},
    )
    response = client.post(
        "/projects",
        json={"id": "IZG", "name": "Another"},
    )
    assert response.status_code == 409


def test_list_projects(client):
    client.post("/projects", json={"id": "IZG", "name": "IZ Gateway"})
    client.post("/projects", json={"id": "TESTME", "name": "TestMe"})
    response = client.get("/projects")
    assert response.status_code == 200
    data = response.json()
    assert len(data) == 2


def test_list_projects_filter_active(client):
    client.post("/projects", json={"id": "IZG", "name": "Active"})
    client.post("/projects", json={"id": "OLD", "name": "Will be inactive"})
    client.delete("/projects/OLD")

    response = client.get("/projects?active=true")
    assert response.status_code == 200
    data = response.json()
    assert len(data) == 1
    assert data[0]["id"] == "IZG"


def test_get_project(client):
    client.post("/projects", json={"id": "IZG", "name": "IZ Gateway"})
    response = client.get("/projects/IZG")
    assert response.status_code == 200
    assert response.json()["id"] == "IZG"


def test_get_project_not_found(client):
    response = client.get("/projects/UNKNOWN")
    assert response.status_code == 404


def test_update_project(client):
    client.post("/projects", json={"id": "IZG", "name": "Original"})
    response = client.put(
        "/projects/IZG",
        json={"name": "Updated"},
    )
    assert response.status_code == 200
    data = response.json()
    assert data["name"] == "Updated"


def test_update_project_not_found(client):
    response = client.put("/projects/UNKNOWN", json={"name": "Test"})
    assert response.status_code == 404


def test_delete_project(client):
    client.post("/projects", json={"id": "IZG", "name": "IZ Gateway"})
    response = client.delete("/projects/IZG")
    assert response.status_code == 200
    assert response.json()["active"] is False

    # Verify it's still there but inactive
    get_response = client.get("/projects/IZG")
    assert get_response.json()["active"] is False


def test_delete_project_not_found(client):
    response = client.delete("/projects/UNKNOWN")
    assert response.status_code == 404
