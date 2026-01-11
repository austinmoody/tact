import logging
from pathlib import Path

from alembic.config import Config

from alembic import command
from tact.db.base import ensure_db_directory, get_db_path

logger = logging.getLogger(__name__)


def get_alembic_config() -> Config:
    """Get Alembic config pointing to the correct locations."""
    # Find the backend directory (where alembic.ini lives)
    backend_dir = Path(__file__).parent.parent.parent.parent
    alembic_ini = backend_dir / "alembic.ini"

    config = Config(str(alembic_ini))
    config.set_main_option("script_location", str(backend_dir / "alembic"))
    return config


def run_migrations() -> None:
    """Run all pending database migrations."""
    ensure_db_directory()
    db_path = get_db_path()

    logger.info(f"Running database migrations for {db_path}")

    config = get_alembic_config()
    command.upgrade(config, "head")

    logger.info("Database migrations complete")
