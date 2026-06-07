import logging
from sqlalchemy import or_
from database.extensions import db
from database.unit_of_work import UnitOfWork
from policies.file_naming import uuid_filename
from werkzeug.datastructures import FileStorage
from database.models.tyre_model import TyreModel
from services.file_service import FileService, FileSaveRequest
from database.models.data_types.files import FileType, FileModel
from database.repositories.tyre_model_repository import TyreModelRepository
from domain.exceptions import DatabaseError, ModelNotFoundError, InvalidFileTypeError, InvalidFileError, FileSaveError

logger = logging.getLogger(__name__)


class TyreModelService:
    def __init__(self):
        self.repo = TyreModelRepository()
        self.file_service = FileService()

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
        try:
            tyre_model = self.repo.get_by_id(id_)
        except ModelNotFoundError as e:
            logger.error(f"TyreModel with id {id_} not found: {e}")
            raise ModelNotFoundError(f"TyreModel with id {id_} not found: {e}") from e

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
                logger.error(f"Error creating tyre model record: {e}")
                raise DatabaseError(f"Error creating tyre model record: {e}")

        return tyre_model

    def upload_image(self, tyre_model: TyreModel, file: FileStorage) -> TyreModel:
        if not file:
            logger.error("No file provided in tyre model service")
            raise InvalidFileTypeError("No file provided")

        uuid, filename = uuid_filename(file)
        file.filename = filename

        with UnitOfWork(db.session):
            try:
                file_record = self.file_service.handle_file(
                    FileSaveRequest(
                        file=file,
                        upload_directory=f"/tyre_match/files/tyre_models/{tyre_model.id}",
                        valid_extensions=["png", "jpg", "jpeg"],
                        model=FileModel.tyre_model,
                        model_id=tyre_model.id,
                        file_type=FileType.original
                    )
                )
            except InvalidFileError as e:
                logger.error(f"Invalid file error from file service in tyre model service: {e}")
                raise FileSaveError(f"Invalid file error from file service in tyre model service: {e}") from e
            except InvalidFileTypeError as e:
                logger.error(f"Invalid file type error from file service in tyre model service: {e}")
                raise FileSaveError(f"Invalid file type error from file service in tyre model service: {e}") from e
            except (PermissionError, OSError) as e:
                logger.error(f"Permission or OS error from file service in tyre model service: {e}")
                raise FileSaveError(f"Permission or OS error from file service in tyre model service: {e}") from e
            except DatabaseError as e:
                logger.error(f"Database error from file service in tyre model service: {e}")
                raise FileSaveError(f"Database error from file service in tyre model service: {e}") from e

            tyre_model.files.append(file_record)

        return tyre_model

    def update(self, id_, dto) -> TyreModel:
        with UnitOfWork(db.session):
            try:
                tyre_model = self.repo.get_by_id(id_)
            except ModelNotFoundError as e:
                logger.error(f"TyreModel with id {id_} not found: {e}")
                raise ModelNotFoundError(f"TyreModel with id {id_} not found: {e}") from e

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
                logger.error(f"Error updating tyre model record: {e}")
                raise DatabaseError(f"Error updating tyre model record: {e}")

        return updated_tyre_model

    def delete(self, id_: int) -> bool:
        with UnitOfWork(db.session):
            try:
                res = self.repo.delete(id_)
            except ModelNotFoundError as e:
                logger.error(f"Tyre model with id {id_} not found: {e}")
                raise ModelNotFoundError(f"Tyre model with id {id_} not found")
            except DatabaseError as e:
                logger.error(f"Error deleting tyre model record: {e}")
                raise DatabaseError(f"Error deleting tyre model record: {e}")

            return res
