"""drop time code description keywords examples

Revision ID: 1c02d24f00be
Revises: 342356b8e3e0
Create Date: 2026-01-18 21:33:53.294106

"""

from typing import Sequence, Union

from alembic import op
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision: str = "1c02d24f00be"
down_revision: Union[str, Sequence[str], None] = "342356b8e3e0"
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None


def upgrade() -> None:
    """Upgrade schema."""
    op.drop_column("time_codes", "description")
    op.drop_column("time_codes", "keywords")
    op.drop_column("time_codes", "examples")


def downgrade() -> None:
    """Downgrade schema."""
    op.add_column(
        "time_codes",
        sa.Column("description", sa.Text(), nullable=False, server_default=""),
    )
    op.add_column(
        "time_codes",
        sa.Column("keywords", sa.Text(), nullable=True, server_default="[]"),
    )
    op.add_column(
        "time_codes",
        sa.Column("examples", sa.Text(), nullable=True, server_default="[]"),
    )
