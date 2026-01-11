from collections.abc import Generator

from sqlalchemy.orm import Session, sessionmaker

from tact.db.base import engine

SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)


def get_session() -> Generator[Session, None, None]:
    session = SessionLocal()
    try:
        yield session
    finally:
        session.close()
