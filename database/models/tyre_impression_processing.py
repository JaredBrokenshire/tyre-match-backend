from datetime import datetime
from database.session import Base
from sqlalchemy.orm import relationship
from sqlalchemy.dialects.mysql import LONGTEXT
from sqlalchemy import Column, Integer, ForeignKey, Float, DateTime


class TyreImpressionProcessing(Base):
    __tablename__ = 'tyre_impression_processing'

    id = Column(Integer, primary_key=True)
    tyre_impression_id = Column(
        Integer,
        ForeignKey('tyre_impressions.id'),
        nullable=False,
        unique=True
    )

    original_file_id = Column(Integer, ForeignKey('files.id'), unique=True)
    normalised_file_id = Column(Integer, ForeignKey('files.id'), unique=True)
    enhanced_file_id = Column(Integer, ForeignKey('files.id'), unique=True)
    binary_file_id = Column(Integer, ForeignKey('files.id'), unique=True)
    clean_file_id = Column(Integer, ForeignKey('files.id'), unique=True)
    skeleton_file_id = Column(Integer, ForeignKey('files.id'), unique=True)

    original_file = relationship("File", foreign_keys="[TyreImpressionProcessing.original_file_id]")
    normalised_file = relationship("File", foreign_keys="[TyreImpressionProcessing.normalised_file_id]")
    enhanced_file = relationship("File", foreign_keys="[TyreImpressionProcessing.enhanced_file_id]")
    binary_file = relationship("File", foreign_keys="[TyreImpressionProcessing.binary_file_id]")
    clean_file = relationship("File", foreign_keys="[TyreImpressionProcessing.clean_file_id]")
    skeleton_file = relationship("File", foreign_keys="[TyreImpressionProcessing.skeleton_file_id]")

    edge_density = Column(Float, nullable=True)
    void_ratio = Column(Float, nullable=True)
    groove_count = Column(Integer, nullable=True)

    feature_vector_json = Column(LONGTEXT, nullable=True)
    match_results_json = Column(LONGTEXT, nullable=True)

    pipeline_version = Column(Integer, default=1, nullable=False)

    created_at = Column(DateTime, nullable=True, default=datetime.now)

    # Define relationships
    tyre_impression = relationship(
        "TyreImpression",
        uselist=False,
        single_parent=True,
    )

    def __repr__(self):
        return f"<TyreImpressionProcessing {self.id} v{self.pipeline_version}>"