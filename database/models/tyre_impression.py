from datetime import datetime
from database.session import Base
from database.models.data_types import TyreImpressionStatus
from sqlalchemy import Column, Integer, String, DateTime, Enum


class TyreImpression(Base):
    __tablename__ = 'tyre_impressions'

    id = Column(Integer, primary_key=True)
    uuid = Column(String(64), unique=True, nullable=False)

    file_path = Column(String(255), nullable=False)

    status = Column(Enum(TyreImpressionStatus), default=TyreImpressionStatus.uploaded, nullable=False)

    created_at = Column(DateTime, default=datetime.now)

    def __repr__(self):
        return f"<TyreImpression {self.file_path} {self.status}>"
