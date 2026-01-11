from collections.abc import Callable, Generator

from sqlalchemy.orm import Session, sessionmaker

from tact.db.base import engine

SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)

# Alias for external use
SessionFactory = SessionLocal

# Allow tests to override this
_session_factory: Callable[[], Session] = SessionLocal


def get_session() -> Generator[Session, None, None]:
    session = _session_factory()
    try:
        yield session
    finally:
        session.close()


def set_session_factory(factory: Callable[[], Session]) -> None:
    """Set the session factory. Used for testing."""
    global _session_factory
    _session_factory = factory


def reset_session_factory() -> None:
    """Reset to the default session factory."""
    global _session_factory
    _session_factory = SessionLocal
