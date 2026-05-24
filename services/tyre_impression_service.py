import policies
from flask import current_app
from services import FileService
from database.extensions import db
from database.models import TyreImpression
from database.unit_of_work import UnitOfWork
from database.repositories import TyreImpressionRepository
from domain import InvalidFileTypeError, FileSaveError, DatabaseError


class TyreImpressionService:
    def __init__(self):
        self.repo = TyreImpressionRepository()
        self.file_service = FileService()

    def get_all(self, page=1, page_size=20) -> (list[TyreImpression], int):
        return self.repo.get_all(page=page, page_size=page_size)

    def upload_impression_image(self, file) -> TyreImpression:
        if not file:
            raise InvalidFileTypeError("No file provided")

        if not file.filename:
            raise InvalidFileTypeError("No filename provided")

        uuid, filename = policies.uuid_filename(file)
        file.filename = filename

        try:
            path = self.file_service.save_file(
                file,
                "/tyre_match/files/tyre_impressions",
                ["png", "jpg", "jpeg", "webp"]
            )
        except InvalidFileTypeError as e:
            current_app.logger.error(f"Invalid file type error: {e}")
            raise InvalidFileTypeError(f"Error saving file: {str(e)}")
        except PermissionError as e:
            current_app.logger.error(f"Permission error: {e}")
            raise FileSaveError(f"Error saving file: {str(e)}")
        except OSError as e:
            current_app.logger.error(f"OS error: {e}")
            raise FileSaveError(f"Error saving file: {str(e)}")

        with UnitOfWork(db.session):
            try:
                tyre_impression = self.repo.create(
                    uuid=uuid,
                    file_path=path,
                )
            except DatabaseError as e:
                current_app.logger.exception(e)
                raise DatabaseError(f"Error uploading file: {str(e)}")

        return tyre_impression
