from datetime import date

import pytest
from fastapi import FastAPI
from fastapi.testclient import TestClient
from sqlalchemy import create_engine, event
from sqlalchemy.orm import sessionmaker
from sqlalchemy.pool import StaticPool

from tact.db.base import Base
from tact.db.models import Config, TimeCode, TimeEntry, WorkType  # noqa: F401
from tact.db.session import get_session
from tact.routes.entries import router as entries_router


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
    test_app.include_router(entries_router)
    test_app.dependency_overrides[get_session] = override_get_session

    with TestClient(test_app) as test_client:
        yield test_client

    engine.dispose()


# Create tests

def test_create_entry_raw_text_only(client):
    response = client.post(
        "/entries",
        json={"raw_text": "2h coding on Project Alpha"},
    )
    assert response.status_code == 201
    data = response.json()
    assert data["raw_text"] == "2h coding on Project Alpha"
    assert data["status"] == "pending"
    assert data["entry_date"] == str(date.today())
    assert data["id"] is not None
    assert data["manually_corrected"] is False


def test_create_entry_with_entry_date(client):
    response = client.post(
        "/entries",
        json={"raw_text": "1h meeting", "entry_date": "2026-01-09"},
    )
    assert response.status_code == 201
    data = response.json()
    assert data["entry_date"] == "2026-01-09"


def test_create_entry_missing_raw_text(client):
    response = client.post("/entries", json={})
    assert response.status_code == 422


# List tests

def test_list_entries_all(client):
    client.post("/entries", json={"raw_text": "Entry 1"})
    client.post("/entries", json={"raw_text": "Entry 2"})
    response = client.get("/entries")
    assert response.status_code == 200
    data = response.json()
    assert len(data) == 2


def test_list_entries_filter_status(client):
    client.post("/entries", json={"raw_text": "Entry 1"})
    response = client.get("/entries?status=pending")
    assert response.status_code == 200
    data = response.json()
    assert len(data) == 1
    assert data[0]["status"] == "pending"

    response = client.get("/entries?status=parsed")
    assert response.status_code == 200
    assert len(response.json()) == 0


def test_list_entries_filter_date_range(client):
    client.post("/entries", json={"raw_text": "Entry 1", "entry_date": "2026-01-10"})
    client.post("/entries", json={"raw_text": "Entry 2", "entry_date": "2026-01-15"})
    client.post("/entries", json={"raw_text": "Entry 3", "entry_date": "2026-01-20"})

    response = client.get("/entries?from_date=2026-01-12&to_date=2026-01-18")
    assert response.status_code == 200
    data = response.json()
    assert len(data) == 1
    assert data[0]["entry_date"] == "2026-01-15"


def test_list_entries_pagination(client):
    for i in range(5):
        client.post("/entries", json={"raw_text": f"Entry {i}"})

    response = client.get("/entries?limit=2&offset=0")
    assert response.status_code == 200
    assert len(response.json()) == 2

    response = client.get("/entries?limit=2&offset=3")
    assert response.status_code == 200
    assert len(response.json()) == 2


# Get single entry tests

def test_get_entry_exists(client):
    create_response = client.post("/entries", json={"raw_text": "Test entry"})
    entry_id = create_response.json()["id"]

    response = client.get(f"/entries/{entry_id}")
    assert response.status_code == 200
    assert response.json()["id"] == entry_id


def test_get_entry_not_found(client):
    response = client.get("/entries/unknown-uuid")
    assert response.status_code == 404


# Update tests

def test_update_entry_success(client):
    create_response = client.post("/entries", json={"raw_text": "Original"})
    entry_id = create_response.json()["id"]

    response = client.patch(
        f"/entries/{entry_id}",
        json={"duration_minutes": 120},
    )
    assert response.status_code == 200
    data = response.json()
    assert data["duration_minutes"] == 120
    assert data["manually_corrected"] is True


def test_update_entry_sets_manually_corrected(client):
    create_response = client.post("/entries", json={"raw_text": "Original"})
    entry_id = create_response.json()["id"]
    assert create_response.json()["manually_corrected"] is False

    response = client.patch(
        f"/entries/{entry_id}",
        json={"description": "Updated description"},
    )
    assert response.status_code == 200
    assert response.json()["manually_corrected"] is True


def test_update_entry_not_found(client):
    response = client.patch("/entries/unknown-uuid", json={"duration_minutes": 60})
    assert response.status_code == 404


# Delete tests

def test_delete_entry_success(client):
    create_response = client.post("/entries", json={"raw_text": "To be deleted"})
    entry_id = create_response.json()["id"]

    response = client.delete(f"/entries/{entry_id}")
    assert response.status_code == 204

    # Verify it's gone
    get_response = client.get(f"/entries/{entry_id}")
    assert get_response.status_code == 404


def test_delete_entry_not_found(client):
    response = client.delete("/entries/unknown-uuid")
    assert response.status_code == 404
