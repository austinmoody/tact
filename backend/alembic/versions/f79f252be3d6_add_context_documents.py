"""add_context_documents

Revision ID: f79f252be3d6
Revises: 2e1d8609b084
Create Date: 2026-01-13 15:06:06.542745

"""
from typing import Sequence, Union

from alembic import op
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision: str = 'f79f252be3d6'
down_revision: Union[str, Sequence[str], None] = '2e1d8609b084'
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None


def upgrade() -> None:
    """Upgrade schema."""
    op.create_table(
        "context_documents",
        sa.Column("id", sa.String(), nullable=False),
        sa.Column("project_id", sa.String(), nullable=True),
        sa.Column("time_code_id", sa.String(), nullable=True),
        sa.Column("content", sa.Text(), nullable=False),
        sa.Column("embedding", sa.LargeBinary(), nullable=True),
        sa.Column("created_at", sa.DateTime(), nullable=False),
        sa.Column("updated_at", sa.DateTime(), nullable=False),
        sa.PrimaryKeyConstraint("id"),
        sa.ForeignKeyConstraint(["project_id"], ["projects.id"]),
        sa.ForeignKeyConstraint(["time_code_id"], ["time_codes.id"]),
        sa.CheckConstraint(
            "(project_id IS NOT NULL AND time_code_id IS NULL) OR "
            "(project_id IS NULL AND time_code_id IS NOT NULL)",
            name="exactly_one_parent",
        ),
    )


def downgrade() -> None:
    """Downgrade schema."""
    op.drop_table("context_documents")
