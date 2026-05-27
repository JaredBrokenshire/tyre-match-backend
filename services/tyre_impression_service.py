import policies
from flask import current_app
from utils import validate_file
from database.extensions import db
from database.models import TyreImpression
from database.unit_of_work import UnitOfWork
from tasks import process_tyre_impression_task
from werkzeug.datastructures import FileStorage
from database.models.data_types import FileType, FileModel
from services.file_service import FileService, FileSaveRequest
from domain import InvalidFileTypeError, FileSaveError, DatabaseError, InvalidFileError
from database.repositories import TyreImpressionRepository, TyreImpressionProcessingRepository


class TyreImpressionService:
    def __init__(self):
        self.tyre_impression_repository = TyreImpressionRepository()
        self.tyre_impression_processing_repository = TyreImpressionProcessingRepository()
        self.file_service = FileService()

    def get_all(self, page=1, page_size=20) -> (list[TyreImpression], int):
        return self.tyre_impression_repository.get_all(page=page, page_size=page_size)

    def upload_impression_image(self, file: FileStorage) -> TyreImpression:
        if not file:
            current_app.logger.error("No file provided in tyre impression service")
            raise InvalidFileTypeError("No file provided")

        try:
            validate_file(file, ["jpg", "jpeg", "png"])
        except InvalidFileTypeError as e:
            current_app.logger.error(f"Invalid file type error from validate_file in tyre impression service: {e}")
            raise InvalidFileTypeError(f"Invalid file type error from validate_file in tyre impression service: {e}")
        except InvalidFileError as e:
            current_app.logger.error(f"Invalid file error from validate_file in tyre impression service: {e}")
            raise InvalidFileTypeError(f"Invalid file error from validate_file in tyre impression service: {e}")

        uuid, filename = policies.uuid_filename(file)
        file.filename = filename

        with UnitOfWork(db.session):
            try:
                tyre_impression = self.tyre_impression_repository.create(uuid=uuid)
            except DatabaseError as e:
                current_app.logger.exception(f"Error creating tyre impression record in tyre impression service: {e}")
                raise DatabaseError(f"Error creating tyre impression record in tyre impression service: {e}")

            try:
                tyre_impression_processing = self.tyre_impression_processing_repository.create(
                    tyre_impression_id=tyre_impression.id,
                )
            except DatabaseError as e:
                current_app.logger.exception(
                    f"Error creating tyre impression processing record in tyre impression service: {e}")
                raise DatabaseError(f"Error creating tyre impression processing record in tyre impression service: {e}")

            try:
                self.file_service.handle_file(
                    FileSaveRequest(
                        file=file,
                        upload_directory=f"/tyre_match/files/tyre_impressions/{tyre_impression.id}/{FileType.original.value}",
                        valid_extensions=["png", "jpg", "jpeg", "webp"],
                        model=FileModel.tyre_model,
                        model_id=tyre_impression_processing.id,
                        file_type=FileType.original
                    )
                )
            except InvalidFileError as e:
                current_app.logger.error(f"Invalid file error from file service in tyre impression service: {e}")
                raise FileSaveError(f"Invalid file error from file service in tyre impression service: {e}")
            except InvalidFileTypeError as e:
                current_app.logger.error(f"Invalid file type error from file service in tyre impression service: {e}")
                raise FileSaveError(f"Invalid file type error from file service in tyre impression service: {e}")
            except (PermissionError, OSError) as e:
                current_app.logger.error(f"Permission or OS error from file service in tyre impression service: {e}")
                raise FileSaveError(f"Permission or OS error from file service in tyre impression service: {e}")
            except DatabaseError as e:
                current_app.logger.error(f"Database error from file service in tyre impression service: {e}")
                raise FileSaveError(f"Database error from file service in tyre impression service: {e}")

        # Trigger async processing task
        process_tyre_impression_task.delay(tyre_impression_processing.id)

        return tyre_impression
