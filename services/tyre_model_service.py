from flask import current_app
from domain import DatabaseError
from sqlalchemy.exc import IntegrityError
from database.repositories import TyreModelRepository


class TyreModelService:
    def __init__(self):
        self.repo = TyreModelRepository()

    def create(self, model):
        try:
            tyre_model = self.repo.create(
                manufacturer=model["manufacturer"],
                model_name=model["model_name"],
                category=model["category"],
                vehicle_type=model["vehicle_type"],
                width_mm=model["width_mm"],
                aspect_ratio=model["aspect_ratio"],
                rim_diameter_inches=model["rim_diameter_inches"],
                groove_count=model["groove_count"],
                pattern_type=model["pattern_type"],
                tread_pitch_length_mm=model["tread_pitch_length_mm"],
                dataset_source=model["dataset_source"],
                notes=model["notes"],
            )
        except IntegrityError as e:
            current_app.logger.error(f"Error creating tyre_model record: {e}")
            raise DatabaseError(f"Error creating tyre_model record: {e}")

        return tyre_model