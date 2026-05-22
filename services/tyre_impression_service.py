import policies
from flask import current_app
from services import save_file
from database.extensions import db
from database.models import TyreImpression
from database.unit_of_work import UnitOfWork
from database.repositories import TyreImpressionRepository
from domain import InvalidFileTypeError, FileSaveError, FileUploadError


class TyreImpressionService:
    def __init__(self):
        self.repo = TyreImpressionRepository()

    def upload_impression_image(self, file) -> TyreImpression:
        if not file:
            raise InvalidFileTypeError("No file provided")

        if not file.filename:
            raise InvalidFileTypeError("No filename provided")

        uuid, filename = policies.uuid_filename(file)
        file.filename = filename

        try:
            path = save_file(
                file,
                "/images/tyre_impressions",
                ["png", "jpg", "jpeg", "webp"]
            )
        except ValueError as e:
            current_app.logger.exception(e)
            raise FileSaveError(f"Error saving file: {str(e)}")

        with UnitOfWork(db.session):
            try:
                tyre_impression = self.repo.create(
                    uuid=uuid,
                    file_path=path,
                )

                return tyre_impression
            except Exception as e:
                current_app.logger.exception(e)
                raise FileUploadError(f"Error uploading file: {str(e)}")
