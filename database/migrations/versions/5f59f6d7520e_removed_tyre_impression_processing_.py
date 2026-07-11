"""Removed tyre impression processing object and associated artifacts

Revision ID: 5f59f6d7520e
Revises: 6044c53dc9d8
Create Date: 2026-07-11 11:02:04.126473

"""
from alembic import op
import sqlalchemy as sa
from sqlalchemy.dialects import mysql

# revision identifiers, used by Alembic.
revision = '5f59f6d7520e'
down_revision = '6044c53dc9d8'
branch_labels = None
depends_on = None


def upgrade():
    op.drop_table('tyre_impression_processing')
    with op.batch_alter_table('tyre_impressions', schema=None) as batch_op:
        batch_op.add_column(sa.Column('edge_density', sa.Float(), nullable=True))
        batch_op.add_column(sa.Column('void_ratio', sa.Float(), nullable=True))
        batch_op.add_column(sa.Column('groove_count', sa.Integer(), nullable=True))
        batch_op.add_column(sa.Column('updated_at', sa.DateTime(), nullable=True))


def downgrade():
    with op.batch_alter_table('tyre_impressions', schema=None) as batch_op:
        batch_op.drop_column('updated_at')
        batch_op.drop_column('groove_count')
        batch_op.drop_column('void_ratio')
        batch_op.drop_column('edge_density')

    op.create_table('tyre_impression_processing',
    sa.Column('id', mysql.INTEGER(), autoincrement=True, nullable=False),
    sa.Column('tyre_impression_id', mysql.INTEGER(), autoincrement=False, nullable=False),
    sa.Column('edge_density', mysql.FLOAT(), nullable=True),
    sa.Column('void_ratio', mysql.FLOAT(), nullable=True),
    sa.Column('groove_count', mysql.INTEGER(), autoincrement=False, nullable=True),
    sa.Column('created_at', mysql.DATETIME(), nullable=True),
    sa.Column('feature_vector_json', mysql.LONGTEXT(), nullable=True),
    sa.Column('match_results_json', mysql.LONGTEXT(), nullable=True),
    sa.Column('pipeline_version', mysql.INTEGER(), autoincrement=False, nullable=False),
    sa.ForeignKeyConstraint(['tyre_impression_id'], ['tyre_impressions.id'], name=op.f('tyre_impression_processing_ibfk_1')),
    sa.PrimaryKeyConstraint('id'),
    mysql_collate='utf8mb4_0900_ai_ci',
    mysql_default_charset='utf8mb4',
    mysql_engine='InnoDB'
    )
    with op.batch_alter_table('tyre_impression_processing', schema=None) as batch_op:
        batch_op.create_index(batch_op.f('tyre_impression_id'), ['tyre_impression_id'], unique=True)
