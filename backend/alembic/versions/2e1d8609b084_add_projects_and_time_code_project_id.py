"""add_projects_and_time_code_project_id

Revision ID: 2e1d8609b084
Revises: 64619729f440
Create Date: 2026-01-13 15:03:28.964334

"""

from typing import Sequence, Union
import logging
import sys

from alembic import op
import sqlalchemy as sa

logger = logging.getLogger(__name__)

# revision identifiers, used by Alembic.
revision: str = "2e1d8609b084"
down_revision: Union[str, Sequence[str], None] = "64619729f440"
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None


def log_step(msg: str) -> None:
    """Log migration step progress."""
    print(f"[MIGRATION] {msg}", flush=True)
    sys.stdout.flush()


def upgrade() -> None:
    """Upgrade schema."""
    log_step("Starting migration 2e1d8609b084")

    # Create projects table
    log_step("Step 1: Creating projects table...")
    op.create_table(
        "projects",
        sa.Column("id", sa.String(), nullable=False),
        sa.Column("name", sa.String(), nullable=False),
        sa.Column("description", sa.Text(), nullable=True),
        sa.Column("active", sa.Boolean(), nullable=False, server_default="1"),
        sa.Column("created_at", sa.DateTime(), nullable=False),
        sa.Column("updated_at", sa.DateTime(), nullable=False),
        sa.PrimaryKeyConstraint("id"),
    )
    log_step("Step 1: DONE - projects table created")

    # Create default project for existing time codes
    log_step("Step 2: Inserting default project...")
    op.execute(
        """
        INSERT INTO projects (id, name, description, active, created_at, updated_at)
        VALUES ('default', 'Default Project', 'Default project for existing time codes', 1, datetime('now'), datetime('now'))
        """
    )
    log_step("Step 2: DONE - default project inserted")

    # Add project_id column to time_codes (nullable - no table recreation needed)
    log_step("Step 3: Adding project_id column to time_codes...")
    op.add_column(
        "time_codes",
        sa.Column("project_id", sa.String(), nullable=True),
    )
    log_step("Step 3: DONE - project_id column added")

    # Backfill existing time codes to default project
    log_step("Step 4: Backfilling time_codes with default project_id...")
    op.execute("UPDATE time_codes SET project_id = 'default'")
    log_step("Step 4: DONE - time_codes backfilled")

    log_step("Migration 2e1d8609b084 COMPLETE")


def downgrade() -> None:
    """Downgrade schema."""
    # Drop project_id column from time_codes
    op.drop_column("time_codes", "project_id")

    # Drop projects table
    op.drop_table("projects")
