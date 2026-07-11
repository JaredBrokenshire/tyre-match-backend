from database.session import Base
from database.models.file import File
from datetime import datetime, timezone
from typing import ClassVar, Optional, Dict
from sqlalchemy import Column, Integer, String, DateTime, Enum, Float
from database.models.data_types.tyre_impression_status import TyreImpressionStatus


class TyreImpression(Base):
    __tablename__ = 'tyre_impressions'

    id = Column(Integer, primary_key=True)
    uuid = Column(String(64), unique=True, nullable=False)

    status = Column(Enum(TyreImpressionStatus), default=TyreImpressionStatus.uploaded, nullable=False)

    # Metadata
    edge_density = Column(Float, nullable=True)
    void_ratio = Column(Float, nullable=True)
    groove_count = Column(Integer, nullable=True)

    created_at = Column(DateTime, default=datetime.now(timezone.utc))
    updated_at = Column(DateTime, default=datetime.now(timezone.utc), onupdate=datetime.now(timezone.utc))

    # Not table column, typing only
    files: ClassVar[Optional[Dict[str, File]]] = None


    def __repr__(self):
        return f"<TyreImpression {self.id}>"
