from database.session import Base
from datetime import datetime, timezone
from database.models.data_types import FileModel, FileType
from sqlalchemy import Column, Integer, DateTime, Enum, String, Index


class File(Base):
    __tablename__ = 'files'

    id = Column(Integer, primary_key=True)

    model = Column(Enum(FileModel), nullable=False)
    model_id = Column(Integer, nullable=False)
    file_type = Column(Enum(FileType), default=FileType.original, nullable=False)

    file_name = Column(String(255), nullable=False)
    file_location = Column(String(512), nullable=False)
    mime_type = Column(String(128), nullable=False)

    created_at = Column(DateTime, default=datetime.now(timezone.utc))
    updated_at = Column(DateTime, default=datetime.now, onupdate=datetime.now(timezone.utc))

    __table_args__ = (
        Index(
            'index_file_owner',
            'model',
            'model_id'
        ),
    )