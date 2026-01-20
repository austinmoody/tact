"""rename_time_entry_fields

Revision ID: 8809038f19d7
Revises: f79f252be3d6
Create Date: 2026-01-16 06:43:18.501677

"""

from collections.abc import Sequence

from alembic import op

# revision identifiers, used by Alembic.
revision: str = "8809038f19d7"
down_revision: str | Sequence[str] | None = "f79f252be3d6"
branch_labels: str | Sequence[str] | None = None
depends_on: str | Sequence[str] | None = None


def upgrade() -> None:
    """Rename raw_text -> user_input and description -> parsed_description."""
    # Disable foreign keys during batch operation (required for SQLite)
    op.execute("PRAGMA foreign_keys=OFF")
    with op.batch_alter_table("time_entries") as batch_op:
        batch_op.alter_column("raw_text", new_column_name="user_input")
        batch_op.alter_column("description", new_column_name="parsed_description")
    op.execute("PRAGMA foreign_keys=ON")


def downgrade() -> None:
    """Revert user_input -> raw_text and parsed_description -> description."""
    op.execute("PRAGMA foreign_keys=OFF")
    with op.batch_alter_table("time_entries") as batch_op:
        batch_op.alter_column("user_input", new_column_name="raw_text")
        batch_op.alter_column("parsed_description", new_column_name="description")
    op.execute("PRAGMA foreign_keys=ON")
