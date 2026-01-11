import pytest
from fastapi import FastAPI
from fastapi.testclient import TestClient
from sqlalchemy import create_engine, event
from sqlalchemy.orm import sessionmaker
from sqlalchemy.pool import StaticPool

from tact.db.base import Base
from tact.db.models import Config, TimeCode, TimeEntry, WorkType  # noqa: F401
from tact.db.session import get_session
from tact.routes.work_types import router as work_types_router


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
    test_app.include_router(work_types_router)
    test_app.dependency_overrides[get_session] = override_get_session

    with TestClient(test_app) as test_client:
        yield test_client

    engine.dispose()


def test_create_work_type(client):
    response = client.post(
        "/work-types",
        json={"name": "Development"},
    )
    assert response.status_code == 201
    data = response.json()
    assert data["id"] == "development"
    assert data["name"] == "Development"
    assert data["active"] is True


def test_create_work_type_slug_generation(client):
    response = client.post(
        "/work-types",
        json={"name": "Code Review"},
    )
    assert response.status_code == 201
    data = response.json()
    assert data["id"] == "code-review"
    assert data["name"] == "Code Review"


def test_create_duplicate_work_type(client):
    client.post("/work-types", json={"name": "Development"})
    response = client.post("/work-types", json={"name": "Development"})
    assert response.status_code == 409


def test_create_duplicate_work_type_same_slug(client):
    # "DevOps" and "DEVOPS" both generate "devops" slug
    client.post("/work-types", json={"name": "DevOps"})
    response = client.post("/work-types", json={"name": "DEVOPS"})
    assert response.status_code == 409


def test_list_work_types(client):
    client.post("/work-types", json={"name": "Development"})
    client.post("/work-types", json={"name": "Meeting"})
    response = client.get("/work-types")
    assert response.status_code == 200
    data = response.json()
    assert len(data) == 2


def test_list_work_types_filter_active(client):
    client.post("/work-types", json={"name": "Development"})
    client.post("/work-types", json={"name": "Meeting"})
    client.delete("/work-types/meeting")

    response = client.get("/work-types?active=true")
    assert response.status_code == 200
    data = response.json()
    assert len(data) == 1
    assert data[0]["id"] == "development"


def test_get_work_type(client):
    client.post("/work-types", json={"name": "Development"})
    response = client.get("/work-types/development")
    assert response.status_code == 200
    assert response.json()["id"] == "development"


def test_get_work_type_not_found(client):
    response = client.get("/work-types/unknown")
    assert response.status_code == 404


def test_update_work_type(client):
    client.post("/work-types", json={"name": "Original"})
    response = client.put("/work-types/original", json={"name": "Updated"})
    assert response.status_code == 200
    assert response.json()["name"] == "Updated"


def test_update_work_type_not_found(client):
    response = client.put("/work-types/unknown", json={"name": "Test"})
    assert response.status_code == 404


def test_delete_work_type(client):
    client.post("/work-types", json={"name": "Development"})
    response = client.delete("/work-types/development")
    assert response.status_code == 200
    assert response.json()["active"] is False

    # Verify it's still there but inactive
    get_response = client.get("/work-types/development")
    assert get_response.json()["active"] is False


def test_delete_work_type_not_found(client):
    response = client.delete("/work-types/unknown")
    assert response.status_code == 404
