import logging
import os
from pathlib import Path

from sqlalchemy import create_engine, event
from sqlalchemy.orm import DeclarativeBase

logger = logging.getLogger(__name__)

DEFAULT_DB_PATH = "./data/tact.db"


def get_db_path() -> str:
    return os.environ.get("TACT_DB_PATH", DEFAULT_DB_PATH)


def get_database_url() -> str:
    db_path = get_db_path()
    return f"sqlite:///{db_path}"


def ensure_db_directory() -> None:
    db_path = Path(get_db_path())
    db_path.parent.mkdir(parents=True, exist_ok=True)


class Base(DeclarativeBase):
    pass


def _set_sqlite_pragma(dbapi_conn, connection_record):
    cursor = dbapi_conn.cursor()
    cursor.execute("PRAGMA foreign_keys=ON")
    cursor.execute("PRAGMA journal_mode=WAL")
    cursor.execute("PRAGMA busy_timeout=5000")
    logger.debug(
        "SQLite pragmas set: foreign_keys=ON, journal_mode=WAL, busy_timeout=5000"
    )
    cursor.close()


def create_db_engine():
    ensure_db_directory()
    db_engine = create_engine(
        get_database_url(),
        connect_args={"check_same_thread": False},
    )
    event.listen(db_engine, "connect", _set_sqlite_pragma)
    return db_engine


engine = create_db_engine()
