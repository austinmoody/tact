"""remove_project_description

Revision ID: 930ed8fe7fa7
Revises: f79f252be3d6
Create Date: 2026-01-15 20:56:48.796257

"""

from typing import Sequence, Union

from alembic import op
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision: str = "930ed8fe7fa7"
down_revision: Union[str, Sequence[str], None] = "f79f252be3d6"
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None


def upgrade() -> None:
    """Drop description column from projects table."""
    # Disable foreign keys during batch operation (required for SQLite)
    op.execute("PRAGMA foreign_keys=OFF")
    with op.batch_alter_table("projects") as batch_op:
        batch_op.drop_column("description")
    op.execute("PRAGMA foreign_keys=ON")


def downgrade() -> None:
    """Re-add description column to projects table."""
    op.execute("PRAGMA foreign_keys=OFF")
    with op.batch_alter_table("projects") as batch_op:
        batch_op.add_column(sa.Column("description", sa.TEXT(), nullable=True))
    op.execute("PRAGMA foreign_keys=ON")
