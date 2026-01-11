import pytest
from fastapi import FastAPI
from fastapi.testclient import TestClient
from sqlalchemy import create_engine, event
from sqlalchemy.orm import sessionmaker
from sqlalchemy.pool import StaticPool

from tact.db.base import Base
from tact.db.models import Config, TimeCode, TimeEntry, WorkType  # noqa: F401
from tact.db.session import get_session
from tact.routes.time_codes import router as time_codes_router


@pytest.fixture
def client():
    """Create a test client with an in-memory database."""
    # Use StaticPool to ensure all connections use the same in-memory db
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

    # Create tables
    Base.metadata.create_all(bind=engine)

    TestSession = sessionmaker(autocommit=False, autoflush=False, bind=engine)

    def override_get_session():
        session = TestSession()
        try:
            yield session
        finally:
            session.close()

    # Create a test app without lifespan
    test_app = FastAPI()
    test_app.include_router(time_codes_router)
    test_app.dependency_overrides[get_session] = override_get_session

    with TestClient(test_app) as test_client:
        yield test_client

    engine.dispose()


def test_create_time_code(client):
    response = client.post(
        "/time-codes",
        json={
            "id": "PROJ-001",
            "name": "Project Alpha",
            "description": "Main project",
            "keywords": ["alpha", "main"],
            "examples": ["working on alpha"],
        },
    )
    assert response.status_code == 201
    data = response.json()
    assert data["id"] == "PROJ-001"
    assert data["name"] == "Project Alpha"
    assert data["active"] is True
    assert data["keywords"] == ["alpha", "main"]


def test_create_duplicate_time_code(client):
    client.post(
        "/time-codes",
        json={"id": "PROJ-001", "name": "Project", "description": "Desc"},
    )
    response = client.post(
        "/time-codes",
        json={"id": "PROJ-001", "name": "Another", "description": "Desc"},
    )
    assert response.status_code == 409


def test_list_time_codes(client):
    client.post(
        "/time-codes",
        json={"id": "PROJ-001", "name": "Project 1", "description": "Desc 1"},
    )
    client.post(
        "/time-codes",
        json={"id": "PROJ-002", "name": "Project 2", "description": "Desc 2"},
    )
    response = client.get("/time-codes")
    assert response.status_code == 200
    data = response.json()
    assert len(data) == 2


def test_list_time_codes_filter_active(client):
    client.post(
        "/time-codes",
        json={"id": "PROJ-001", "name": "Active", "description": "Desc"},
    )
    client.post(
        "/time-codes",
        json={"id": "PROJ-002", "name": "Will be inactive", "description": "Desc"},
    )
    client.delete("/time-codes/PROJ-002")

    response = client.get("/time-codes?active=true")
    assert response.status_code == 200
    data = response.json()
    assert len(data) == 1
    assert data[0]["id"] == "PROJ-001"


def test_get_time_code(client):
    client.post(
        "/time-codes",
        json={"id": "PROJ-001", "name": "Project", "description": "Desc"},
    )
    response = client.get("/time-codes/PROJ-001")
    assert response.status_code == 200
    assert response.json()["id"] == "PROJ-001"


def test_get_time_code_not_found(client):
    response = client.get("/time-codes/UNKNOWN")
    assert response.status_code == 404


def test_update_time_code(client):
    client.post(
        "/time-codes",
        json={"id": "PROJ-001", "name": "Original", "description": "Desc"},
    )
    response = client.put(
        "/time-codes/PROJ-001",
        json={"name": "Updated"},
    )
    assert response.status_code == 200
    assert response.json()["name"] == "Updated"


def test_update_time_code_not_found(client):
    response = client.put("/time-codes/UNKNOWN", json={"name": "Test"})
    assert response.status_code == 404


def test_delete_time_code(client):
    client.post(
        "/time-codes",
        json={"id": "PROJ-001", "name": "Project", "description": "Desc"},
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
