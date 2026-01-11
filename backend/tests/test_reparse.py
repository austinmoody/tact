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


def test_reparse_entry_success(client):
    # Create an entry
    create_response = client.post(
        "/entries",
        json={"raw_text": "2h coding on Project Alpha"},
    )
    assert create_response.status_code == 201
    entry_id = create_response.json()["id"]

    # Manually update it as if it was parsed
    update_response = client.patch(
        f"/entries/{entry_id}",
        json={"duration_minutes": 120, "status": "parsed"},
    )
    assert update_response.status_code == 200
    assert update_response.json()["status"] == "parsed"
    assert update_response.json()["manually_corrected"] is True

    # Reparse the entry
    reparse_response = client.post(f"/entries/{entry_id}/reparse")
    assert reparse_response.status_code == 200

    data = reparse_response.json()
    assert data["status"] == "pending"
    assert data["manually_corrected"] is False
    assert data["duration_minutes"] is None


def test_reparse_entry_not_found(client):
    response = client.post("/entries/unknown-uuid/reparse")
    assert response.status_code == 404


def test_reparse_clears_all_parsed_fields(client):
    # Create an entry
    create_response = client.post(
        "/entries",
        json={"raw_text": "test entry"},
    )
    entry_id = create_response.json()["id"]

    # Update with parsed data
    client.patch(
        f"/entries/{entry_id}",
        json={
            "duration_minutes": 60,
            "description": "Test description",
            "status": "parsed",
        },
    )

    # Reparse
    reparse_response = client.post(f"/entries/{entry_id}/reparse")
    data = reparse_response.json()

    # All parsed fields should be cleared
    assert data["duration_minutes"] is None
    assert data["work_type_id"] is None
    assert data["time_code_id"] is None
    assert data["description"] is None
    assert data["confidence_duration"] is None
    assert data["confidence_work_type"] is None
    assert data["confidence_time_code"] is None
    assert data["confidence_overall"] is None
    assert data["parsed_at"] is None
    assert data["parse_error"] is None

    # Status should be reset
    assert data["status"] == "pending"
    assert data["manually_corrected"] is False
