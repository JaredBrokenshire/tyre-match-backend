"""Updated tyre impression status enum

Revision ID: f45d3b90fa8e
Revises: e8059bdc7f45
Create Date: 2026-05-25 08:44:27.528809

"""
from alembic import op
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision = 'f45d3b90fa8e'
down_revision = 'e8059bdc7f45'
branch_labels = None
depends_on = None


def upgrade():
    op.execute("""
        ALTER TABLE tyre_impressions
        MODIFY status ENUM(
            'uploaded',
            'queued',
            'preprocessing',
            'preprocessed',
            'extracting_features',
            'matching',
            'matched',
            'failed',
            'archived'
        ) NOT NULL;  
    """)


def downgrade():
    op.execute("""
        ALTER TABLE tyre_impressions
        MODIFY status ENUM(
            'uploaded',
            'processing',
            'processed',
            'matched',
            'failed'
        ) NOT NULL;  
    """)
