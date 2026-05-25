from datetime import datetime

import sqlalchemy

from database.session import Base
from sqlalchemy.orm import Relationship, relationship

from sqlalchemy import Column, Integer, ForeignKey, String, Float, DateTime


class TyreImpressionProcessing(Base):
    __tablename__ = 'tyre_impression_processing'

    id = Column(Integer, primary_key=True)
    tyre_impression_id = Column(
        Integer,
        ForeignKey('tyre_impressions.id'),
        nullable=False,
        unique=True
    )

    grayscale_path = Column(String(255), nullable=True)
    binary_path = Column(String(255), nullable=True)
    skeleton_path = Column(String(255), nullable=True)

    edge_density = Column(Float, nullable=True)
    void_ratio = Column(Float, nullable=True)
    groove_count = Column(Integer, nullable=True)

    preprocessing_version = Column(Integer, default=1, nullable=False)

    created_at = Column(DateTime, nullable=True, default=datetime.now)

    # Define relationships
    tyre_impression = relationship(
        "TyreImpression",
        uselist=False,
        single_parent=True,
    )

    def __repr__(self):
        return f"<TyreImpressionProcessing {self.id} v{self.preprocessing_version}>"