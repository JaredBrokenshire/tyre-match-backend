from database.session import Base
from sqlalchemy import Column, Integer, String, Float


class TyreModel(Base):
    __tablename__ = 'tyre_models'

    # Basic identifiers
    id = Column(Integer, primary_key=True)
    manufacturer = Column(String(100), nullable=False)
    model_name = Column(String(100), nullable=False)

    # Tyre classifications
    category = Column(String(50), nullable=True)  # e.g. Summer, Winter, Performance, etc.
    vehicle_type = Column(String(50), nullable=True)  # e.g Passenger Car, SUV, Light Truck (Note, may remove depending on sample diversity)

    # Size specifications
    width_mm = Column(Integer, nullable=True)
    aspect_ratio = Column(Integer, nullable=True)
    rim_diameter_inches = Column(Integer, nullable=True)

    # Tread Pattern Metadata
    groove_count = Column(Integer, nullable=True)  # Number of circumferential grooves
    pattern_type = Column(String(50), nullable=True)  # e.g. symmetric, asymmetric, directional
    tread_pitch_length_mm = Column(Float, nullable=True)  # Estimated repeat unit size

    # Research Metadata
    dataset_source = Column(String(100), nullable=True)
    notes = Column(String(255), nullable=True)

    def __repr__(self):
        return f"<TyreModel {self.manufacturer} {self.model_name}>"