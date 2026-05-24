from sqlalchemy import or_
from flask import current_app
from database.extensions import db
from database.models import TyreModel
from database.unit_of_work import UnitOfWork
from domain import DatabaseError, ModelNotFoundError
from database.repositories import TyreModelRepository


class TyreModelService:
    def __init__(self):
        self.repo = TyreModelRepository()

    def get_all(self, page=1, page_size=20, search="") -> (list[TyreModel], int):
        filters = None

        if search != "":
            search_term = f"%{search}%"
            filters = or_(
                TyreModel.manufacturer.ilike(search_term),
                TyreModel.model_name.ilike(search_term),
            )

        return self.repo.get_all(page=page, page_size=page_size, filters=filters)

    def get_by_id(self, id_: int) -> TyreModel:
        tyre_model =  self.repo.get_by_id(id_)

        if not tyre_model:
            current_app.logger.error(f"Error getting tyre model by id: {id_}")
            raise ModelNotFoundError(f"Error getting tyre model by id: {id_}")

        return tyre_model


    def create(self, dto) -> TyreModel:
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
                current_app.logger.error(f"Error creating tyre model record: {e}")
                raise DatabaseError(f"Error creating tyre model record: {e}")

        return tyre_model

    def update(self, id_, dto) -> TyreModel:
        with UnitOfWork(db.session):
            tyre_model = self.repo.get_by_id(id_)

            if not tyre_model:
                current_app.logger.error(f"Error getting tyre model by id: {id_}")
                raise ModelNotFoundError(f"Error getting tyre model by id: {id_}")

            try:
                allowed_fields = [
                    "manufacturer",
                    "model_name",
                    "category",
                    "vehicle_type",
                    "width_mm",
                    "aspect_ratio",
                    "rim_diameter_inches",
                    "groove_count",
                    "pattern_type",
                    "tread_pitch_length_mm",
                    "dataset_source",
                    "notes",
                ]

                update_data = {
                    k: v for k, v in dto.items()
                    if k in allowed_fields
                }

                updated_tyre_model = self.repo.update(
                    entity=tyre_model,
                    **update_data
                )
            except DatabaseError as e:
                current_app.logger.error(f"Error updating tyre model record: {e}")
                raise DatabaseError(f"Error updating tyre model record: {e}")

        return updated_tyre_model

    def delete(self, id_: int) -> bool:
        try:
            res = self.repo.delete(id_)
        except ModelNotFoundError as e:
            current_app.logger.error(f"Tyre model with id {id_} not found: {e}")
            raise ModelNotFoundError(f"Tyre model with id {id_} not found")
        except DatabaseError as e:
            current_app.logger.error(f"Error deleting tyre model record: {e}")
            raise DatabaseError(f"Error deleting tyre model record: {e}")

        return res