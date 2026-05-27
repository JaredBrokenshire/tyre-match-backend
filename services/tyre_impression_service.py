import logging
from database.extensions import db
from utils.files import validate_file
from database.unit_of_work import UnitOfWork
from policies.file_naming import uuid_filename
from werkzeug.datastructures import FileStorage
from database.models.tyre_impression import TyreImpression
from services.file_service import FileService, FileSaveRequest
from database.models.data_types.files import FileType, FileModel
from tasks.tyre_impression_tasks import process_tyre_impression_task
from database.repositories.tyre_impression_repository import TyreImpressionRepository
from domain.exceptions import InvalidFileTypeError, FileSaveError, DatabaseError, InvalidFileError
from database.repositories.tyre_impression_processing_repository import TyreImpressionProcessingRepository

logger = logging.getLogger(__name__)


class TyreImpressionService:
    def __init__(self):
        self.tyre_impression_repository = TyreImpressionRepository()
        self.tyre_impression_processing_repository = TyreImpressionProcessingRepository()
        self.file_service = FileService()

    def get_all(self, page=1, page_size=20) -> (list[TyreImpression], int):
        return self.tyre_impression_repository.get_all(page=page, page_size=page_size)

    def upload_impression_image(self, file: FileStorage) -> TyreImpression:
        if not file:
            logger.error("No file provided in tyre impression service")
            raise InvalidFileTypeError("No file provided")

        try:
            validate_file(file, ["jpg", "jpeg", "png"])
        except InvalidFileTypeError as e:
            logger.error(f"Invalid file type error from validate_file in tyre impression service: {e}")
            raise InvalidFileTypeError(f"Invalid file type error from validate_file in tyre impression service: {e}")
        except InvalidFileError as e:
            logger.error(f"Invalid file error from validate_file in tyre impression service: {e}")
            raise InvalidFileTypeError(f"Invalid file error from validate_file in tyre impression service: {e}")

        uuid, filename = uuid_filename(file)
        file.filename = filename

        with UnitOfWork(db.session):
            try:
                tyre_impression = self.tyre_impression_repository.create(uuid=uuid)
            except DatabaseError as e:
                logger.exception(f"Error creating tyre impression record in tyre impression service: {e}")
                raise DatabaseError(f"Error creating tyre impression record in tyre impression service: {e}")

            try:
                tyre_impression_processing = self.tyre_impression_processing_repository.create(
                    tyre_impression_id=tyre_impression.id,
                )
            except DatabaseError as e:
                logger.exception(
                    f"Error creating tyre impression processing record in tyre impression service: {e}")
                raise DatabaseError(f"Error creating tyre impression processing record in tyre impression service: {e}")

            try:
                self.file_service.handle_file(
                    FileSaveRequest(
                        file=file,
                        upload_directory=f"/tyre_match/files/tyre_impressions/{tyre_impression.id}/{FileType.original.value}",
                        valid_extensions=["png", "jpg", "jpeg", "webp"],
                        model=FileModel.tyre_impression,
                        model_id=tyre_impression_processing.id,
                        file_type=FileType.original
                    )
                )
            except InvalidFileError as e:
                logger.error(f"Invalid file error from file service in tyre impression service: {e}")
                raise FileSaveError(f"Invalid file error from file service in tyre impression service: {e}")
            except InvalidFileTypeError as e:
                logger.error(f"Invalid file type error from file service in tyre impression service: {e}")
                raise FileSaveError(f"Invalid file type error from file service in tyre impression service: {e}")
            except (PermissionError, OSError) as e:
                logger.error(f"Permission or OS error from file service in tyre impression service: {e}")
                raise FileSaveError(f"Permission or OS error from file service in tyre impression service: {e}")
            except DatabaseError as e:
                logger.error(f"Database error from file service in tyre impression service: {e}")
                raise FileSaveError(f"Database error from file service in tyre impression service: {e}")

        # Trigger async processing task
        process_tyre_impression_task.delay(tyre_impression_processing.id)

        return tyre_impression
