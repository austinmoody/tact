from datetime import date

import pytest
from fastapi import FastAPI
from fastapi.testclient import TestClient
from sqlalchemy import create_engine, event
from sqlalchemy.orm import sessionmaker
from sqlalchemy.pool import StaticPool

from tact.db.base import Base
from tact.db.models import (  # noqa: F401
    Config,
    ContextDocument,
    Project,
    TimeCode,
    TimeEntry,
    WorkType,
)
from tact.db.session import get_session
from tact.routes.entries import router as entries_router


@pytest.fixture
def db_session():
    """Create an in-memory database and return a session factory."""
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
    yield TestSession
    engine.dispose()


@pytest.fixture
def client(db_session):
    """Create a test client with the shared database session."""

    def override_get_session():
        session = db_session()
        try:
            yield session
        finally:
            session.close()

    test_app = FastAPI()
    test_app.include_router(entries_router)
    test_app.dependency_overrides[get_session] = override_get_session

    with TestClient(test_app) as test_client:
        yield test_client


# Create tests

def test_create_entry_user_input_only(client):
    response = client.post(
        "/entries",
        json={"user_input": "2h coding on Project Alpha"},
    )
    assert response.status_code == 201
    data = response.json()
    assert data["user_input"] == "2h coding on Project Alpha"
    assert data["status"] == "pending"
    assert data["entry_date"] == str(date.today())
    assert data["id"] is not None
    assert data["manually_corrected"] is False


def test_create_entry_with_entry_date(client):
    response = client.post(
        "/entries",
        json={"user_input": "1h meeting", "entry_date": "2026-01-09"},
    )
    assert response.status_code == 201
    data = response.json()
    assert data["entry_date"] == "2026-01-09"


def test_create_entry_missing_user_input(client):
    response = client.post("/entries", json={})
    assert response.status_code == 422


# List tests

def test_list_entries_all(client):
    client.post("/entries", json={"user_input": "Entry 1"})
    client.post("/entries", json={"user_input": "Entry 2"})
    response = client.get("/entries")
    assert response.status_code == 200
    data = response.json()
    assert len(data) == 2


def test_list_entries_filter_status(client):
    client.post("/entries", json={"user_input": "Entry 1"})
    response = client.get("/entries?status=pending")
    assert response.status_code == 200
    data = response.json()
    assert len(data) == 1
    assert data[0]["status"] == "pending"

    response = client.get("/entries?status=parsed")
    assert response.status_code == 200
    assert len(response.json()) == 0


def test_list_entries_filter_date_range(client):
    client.post("/entries", json={"user_input": "Entry 1", "entry_date": "2026-01-10"})
    client.post("/entries", json={"user_input": "Entry 2", "entry_date": "2026-01-15"})
    client.post("/entries", json={"user_input": "Entry 3", "entry_date": "2026-01-20"})

    response = client.get("/entries?from_date=2026-01-12&to_date=2026-01-18")
    assert response.status_code == 200
    data = response.json()
    assert len(data) == 1
    assert data[0]["entry_date"] == "2026-01-15"


def test_list_entries_pagination(client):
    for i in range(5):
        client.post("/entries", json={"user_input": f"Entry {i}"})

    response = client.get("/entries?limit=2&offset=0")
    assert response.status_code == 200
    assert len(response.json()) == 2

    response = client.get("/entries?limit=2&offset=3")
    assert response.status_code == 200
    assert len(response.json()) == 2


# Get single entry tests

def test_get_entry_exists(client):
    create_response = client.post("/entries", json={"user_input": "Test entry"})
    entry_id = create_response.json()["id"]

    response = client.get(f"/entries/{entry_id}")
    assert response.status_code == 200
    assert response.json()["id"] == entry_id


def test_get_entry_not_found(client):
    response = client.get("/entries/unknown-uuid")
    assert response.status_code == 404


# Update tests

def test_update_entry_success(client):
    create_response = client.post("/entries", json={"user_input": "Original"})
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
    create_response = client.post("/entries", json={"user_input": "Original"})
    entry_id = create_response.json()["id"]
    assert create_response.json()["manually_corrected"] is False

    response = client.patch(
        f"/entries/{entry_id}",
        json={"parsed_description": "Updated description"},
    )
    assert response.status_code == 200
    assert response.json()["manually_corrected"] is True


def test_update_entry_not_found(client):
    response = client.patch("/entries/unknown-uuid", json={"duration_minutes": 60})
    assert response.status_code == 404


# Delete tests

def test_delete_entry_success(client):
    create_response = client.post("/entries", json={"user_input": "To be deleted"})
    entry_id = create_response.json()["id"]

    response = client.delete(f"/entries/{entry_id}")
    assert response.status_code == 204

    # Verify it's gone
    get_response = client.get(f"/entries/{entry_id}")
    assert get_response.status_code == 404


def test_delete_entry_not_found(client):
    response = client.delete("/entries/unknown-uuid")
    assert response.status_code == 404


# Learning tests


def test_update_entry_with_learn_creates_context(client, db_session):
    """Test that updating entry with learn=true creates context document."""
    # Create a project, time code, and work type first
    session = db_session()
    project = Project(id="test-project", name="Test Project")
    time_code = TimeCode(
        id="TC-001",
        project_id="test-project",
        name="Test Time Code",
            )
    work_type = WorkType(id="meetings", name="Meetings")
    session.add(project)
    session.add(time_code)
    session.add(work_type)
    session.commit()
    session.close()

    # Create an entry
    create_response = client.post("/entries", json={"user_input": "2h standup meeting"})
    entry_id = create_response.json()["id"]

    # Update with time_code_id (learn=true by default)
    response = client.patch(
        f"/entries/{entry_id}",
        json={
            "time_code_id": "TC-001",
            "duration_minutes": 120,
            "work_type_id": "meetings",
        },
    )
    assert response.status_code == 200

    # Verify context document was created
    session = db_session()
    contexts = session.query(ContextDocument).filter(
        ContextDocument.time_code_id == "TC-001"
    ).all()
    assert len(contexts) == 1
    assert 'Example: "2h standup meeting"' in contexts[0].content
    assert "120 minutes" in contexts[0].content
    assert "work_type: meetings" in contexts[0].content
    session.close()


def test_update_entry_with_learn_false_no_context(client, db_session):
    """Test that updating entry with learn=false does not create context."""
    # Create a project and time code first
    session = db_session()
    project = Project(id="test-project", name="Test Project")
    time_code = TimeCode(
        id="TC-001",
        project_id="test-project",
        name="Test Time Code",
            )
    session.add(project)
    session.add(time_code)
    session.commit()
    session.close()

    # Create an entry
    create_response = client.post("/entries", json={"user_input": "2h coding"})
    entry_id = create_response.json()["id"]

    # Update with learn=false
    response = client.patch(
        f"/entries/{entry_id}?learn=false",
        json={
            "time_code_id": "TC-001",
            "duration_minutes": 120,
        },
    )
    assert response.status_code == 200

    # Verify no context document was created
    session = db_session()
    contexts = session.query(ContextDocument).filter(
        ContextDocument.time_code_id == "TC-001"
    ).all()
    assert len(contexts) == 0
    session.close()


def test_update_entry_without_time_code_no_context(client, db_session):
    """Test that updating entry without time_code_id does not create context."""
    # Create an entry
    create_response = client.post("/entries", json={"user_input": "2h work"})
    entry_id = create_response.json()["id"]

    # Update without time_code_id
    response = client.patch(
        f"/entries/{entry_id}",
        json={"duration_minutes": 120},
    )
    assert response.status_code == 200

    # Verify no context document was created
    session = db_session()
    contexts = session.query(ContextDocument).all()
    assert len(contexts) == 0
    session.close()


def test_update_entry_context_format_duration_only(client, db_session):
    """Test context format when only duration is set."""
    # Create a project and time code first
    session = db_session()
    project = Project(id="test-project", name="Test Project")
    time_code = TimeCode(
        id="TC-001",
        project_id="test-project",
        name="Test Time Code",
            )
    session.add(project)
    session.add(time_code)
    session.commit()
    session.close()

    # Create an entry
    create_response = client.post("/entries", json={"user_input": "30m quick fix"})
    entry_id = create_response.json()["id"]

    # Update with only duration (no work_type_id)
    response = client.patch(
        f"/entries/{entry_id}",
        json={
            "time_code_id": "TC-001",
            "duration_minutes": 30,
        },
    )
    assert response.status_code == 200

    # Verify context format
    session = db_session()
    contexts = session.query(ContextDocument).filter(
        ContextDocument.time_code_id == "TC-001"
    ).all()
    assert len(contexts) == 1
    assert 'Example: "30m quick fix"' in contexts[0].content
    assert "30 minutes" in contexts[0].content
    assert "work_type" not in contexts[0].content
    session.close()
