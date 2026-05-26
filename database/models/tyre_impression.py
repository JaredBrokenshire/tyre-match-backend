from database.session import Base
from sqlalchemy.orm import relationship
from datetime import datetime, timezone
from database.models.data_types import TyreImpressionStatus
from sqlalchemy import Column, Integer, String, DateTime, Enum


class TyreImpression(Base):
    __tablename__ = 'tyre_impressions'

    id = Column(Integer, primary_key=True)
    uuid = Column(String(64), unique=True, nullable=False)

    status = Column(Enum(TyreImpressionStatus), default=TyreImpressionStatus.uploaded, nullable=False)

    created_at = Column(DateTime, default=datetime.now(timezone.utc))

    processing = relationship(
        "TyreImpressionProcessing",
        back_populates="tyre_impression",
        uselist=False,
        cascade="all, delete-orphan"
    )

    def __repr__(self):
        return f"<TyreImpression {self.id}>"
