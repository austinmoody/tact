"""merge_multiple_heads

Revision ID: 342356b8e3e0
Revises: 8809038f19d7, 930ed8fe7fa7, a3c3608e750f
Create Date: 2026-01-18 11:32:02.376272

"""

from collections.abc import Sequence

# revision identifiers, used by Alembic.
revision: str = "342356b8e3e0"
down_revision: str | Sequence[str] | None = (
    "8809038f19d7",
    "930ed8fe7fa7",
    "a3c3608e750f",
)
branch_labels: str | Sequence[str] | None = None
depends_on: str | Sequence[str] | None = None


def upgrade() -> None:
    """Upgrade schema."""
    pass


def downgrade() -> None:
    """Downgrade schema."""
    pass
