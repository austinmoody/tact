# Design: Setup Data Model

## Overview

This design covers the database layer architecture including SQLAlchemy models, Alembic migrations, and automatic database initialization.

## Directory Structure

```
backend/
├── alembic.ini
├── alembic/
│   ├── env.py
│   ├── script.py.mako
│   └── versions/
│       └── 001_initial_schema.py
└── src/tact/
    └── db/
        ├── __init__.py
        ├── base.py          # SQLAlchemy Base and engine setup
        ├── models.py        # All table models
        ├── session.py       # Session management
        └── migrations.py    # Auto-migration runner
```

## Database Configuration

### Environment Variable

- **Variable**: `TACT_DB_PATH`
- **Default**: `./data/tact.db` (relative to working directory)
- **Docker default**: `/data/tact.db` (mounted volume)

The path is resolved at startup. Parent directories are created if they don't exist.

### Connection Settings

SQLite-specific settings:
- `check_same_thread=False` for FastAPI async compatibility
- WAL mode enabled for better concurrency
- Foreign keys enforced

## SQLAlchemy Models

### Base Configuration

Using SQLAlchemy 2.0 style with `mapped_column` and type annotations.

### TimeEntry

```python
class TimeEntry(Base):
    __tablename__ = "time_entries"

    id: Mapped[uuid.UUID] = mapped_column(primary_key=True, default=uuid.uuid4)
    raw_text: Mapped[str] = mapped_column(Text, nullable=False)

    # Parsed fields
    duration_minutes: Mapped[int | None]
    work_type_id: Mapped[str | None] = mapped_column(ForeignKey("work_types.id"))
    time_code_id: Mapped[str | None] = mapped_column(ForeignKey("time_codes.id"))
    description: Mapped[str | None] = mapped_column(Text)
    entry_date: Mapped[date | None]

    # Confidence scores
    confidence_duration: Mapped[float | None]
    confidence_work_type: Mapped[float | None]
    confidence_time_code: Mapped[float | None]
    confidence_overall: Mapped[float | None]

    # Status
    status: Mapped[str] = mapped_column(default="pending")
    parse_error: Mapped[str | None] = mapped_column(Text)

    # Correction tracking
    manually_corrected: Mapped[bool] = mapped_column(default=False)
    locked: Mapped[bool] = mapped_column(default=False)
    corrected_at: Mapped[datetime | None]

    # Timestamps
    created_at: Mapped[datetime] = mapped_column(default=datetime.utcnow)
    parsed_at: Mapped[datetime | None]
    updated_at: Mapped[datetime] = mapped_column(default=datetime.utcnow, onupdate=datetime.utcnow)
```

### TimeCode

```python
class TimeCode(Base):
    __tablename__ = "time_codes"

    id: Mapped[str] = mapped_column(primary_key=True)
    name: Mapped[str] = mapped_column(nullable=False)
    description: Mapped[str] = mapped_column(Text, nullable=False)
    keywords: Mapped[str] = mapped_column(Text, default="[]")  # JSON array
    examples: Mapped[str] = mapped_column(Text, default="[]")  # JSON array
    active: Mapped[bool] = mapped_column(default=True)
    created_at: Mapped[datetime] = mapped_column(default=datetime.utcnow)
    updated_at: Mapped[datetime] = mapped_column(default=datetime.utcnow, onupdate=datetime.utcnow)
```

### WorkType

```python
class WorkType(Base):
    __tablename__ = "work_types"

    id: Mapped[str] = mapped_column(primary_key=True)
    name: Mapped[str] = mapped_column(nullable=False)
    description: Mapped[str | None] = mapped_column(Text)
    active: Mapped[bool] = mapped_column(default=True)
    created_at: Mapped[datetime] = mapped_column(default=datetime.utcnow)
    updated_at: Mapped[datetime] = mapped_column(default=datetime.utcnow, onupdate=datetime.utcnow)
```

### Config

```python
class Config(Base):
    __tablename__ = "config"

    key: Mapped[str] = mapped_column(primary_key=True)
    value: Mapped[str] = mapped_column(Text, nullable=False)
    updated_at: Mapped[datetime] = mapped_column(default=datetime.utcnow, onupdate=datetime.utcnow)
```

## Migration Strategy

### Alembic Configuration

- Migrations stored in `backend/alembic/versions/`
- Uses SQLAlchemy metadata for autogenerate support
- Configured to work with SQLite

### Auto-Migration on Startup

On FastAPI startup:
1. Check if database file exists
2. If not, create parent directories
3. Run all pending Alembic migrations
4. Log migration status

This ensures a fresh checkout with `make run` "just works" without manual migration steps.

### Migration Naming

Format: `NNN_description.py` (e.g., `001_initial_schema.py`)

## Testing Considerations

- Tests use in-memory SQLite (`:memory:`)
- Each test gets a fresh database via fixture
- Migrations are applied to test database
