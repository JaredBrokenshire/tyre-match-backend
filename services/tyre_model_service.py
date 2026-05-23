from flask import current_app
from domain import DatabaseError
from database.extensions import db
from database.unit_of_work import UnitOfWork
from database.repositories import TyreModelRepository


class TyreModelService:
    def __init__(self):
        self.repo = TyreModelRepository()

    def create(self, dto):
        with UnitOfWork(db.session):
            try:
                tyre_model = self.repo.create(
                    manufacturer=dto.get("manufacturer", "Temp Manufacturer"),
                    model_name=dto.get("model_name", "Temp Model None"),
                    category=dto.get("category", None),
                    vehicle_type=dto.get("vehicle_type", None),
                    width_mm=dto.get("width_mm", None),
                    aspect_ratio=dto.get("aspect_ratio", None),
                    rim_diameter_inches=dto.get("rim_diameter_inches", None),
                    groove_count=dto.get("groove_count", None),
                    pattern_type=dto.get("pattern_type", None),
                    tread_pitch_length_mm=dto.get("tread_pitch_length_mm", None),
                    dataset_source=dto.get("dataset_source", None),
                    notes=dto.get("notes", None),
                )
            except DatabaseError as e:
                current_app.logger.error(f"Error creating tyre_model record: {e}")
                raise DatabaseError(f"Error creating tyre_model record: {e}")

        return tyre_model
