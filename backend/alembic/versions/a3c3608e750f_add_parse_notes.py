"""add_parse_notes

Revision ID: a3c3608e750f
Revises: f79f252be3d6
Create Date: 2026-01-15

"""

from collections.abc import Sequence

import sqlalchemy as sa

from alembic import op

# revision identifiers, used by Alembic.
revision: str = "a3c3608e750f"
down_revision: str | Sequence[str] | None = "f79f252be3d6"
branch_labels: str | Sequence[str] | None = None
depends_on: str | Sequence[str] | None = None


def upgrade() -> None:
    """Add parse_notes column to time_entries table."""
    op.add_column(
        "time_entries",
        sa.Column("parse_notes", sa.Text(), nullable=True),
    )


def downgrade() -> None:
    """Remove parse_notes column from time_entries table."""
    op.drop_column("time_entries", "parse_notes")
