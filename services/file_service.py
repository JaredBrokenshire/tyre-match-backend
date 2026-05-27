import os
import logging
import numpy as np
from dataclasses import dataclass
from database.models.file import File
from werkzeug.datastructures import FileStorage
from utils.files import validate_file, make_directory
from database.repositories.file_repository import FileRepository
from database.models.data_types.files import FileModel, FileType
from domain.exceptions import InvalidFileTypeError, InvalidFileError, DatabaseError, FileSaveError

logger = logging.getLogger(__name__)

BASE_FILE_DIRECTORY = "/tyre_match/files"

@dataclass
class FileSaveRequest:
    file: FileStorage
    upload_directory: str
    valid_extensions: list[str]

    model: FileModel
    model_id: int
    file_type: FileType

@dataclass
class ProcessedImageRequest:
    image: np.ndarray
    file_name: str
    upload_directory: str
    model: FileModel
    model_id: int
    file_type: FileType
    extension: str = "jpg"


class FileService:
    def __init__(self):
        self.base_directory = BASE_FILE_DIRECTORY
        self.file_repository = FileRepository()


    def handle_file(self, request: FileSaveRequest) -> File:
        file = request.file

        # Validate file
        try:
            validate_file(file, request.valid_extensions)
        except InvalidFileTypeError as e:
            logger.error(f"Invalid file type error in file service: {e}")
            raise InvalidFileTypeError(f"Invalid file type error in file service: {e}")
        except InvalidFileError as e:
            logger.error(f"Invalid file error in file service: {e}")
            raise InvalidFileError(f"Invalid file error in file service: {e}")

        # Save file locally
        try:
            file_path = self._save_file(file, request.upload_directory)
        except PermissionError as e:
            logger.error(f"Permission error when saving file in file service: {e}")
            raise PermissionError(f"Permission error when saving file in file service: {e}")
        except OSError as e:
            logger.error(f"OS error when saving file in file service: {e}")
            raise PermissionError(f"OS error when saving file in file service: {e}")

        # Create DB record
        try:
            file = self.file_repository.create(
                model=request.model,
                model_id=request.model_id,
                file_type=request.file_type,
                file_name=file.filename,
                file_location=file_path,
                mime_type=file.mimetype,
            )
        except DatabaseError as e:
            logger.error(f"Database error when creating file in file service: {e}")
            raise DatabaseError(f"Database error when creating file in file service: {e}")

        return file

    def save_processed_image(self, request: ProcessedImageRequest) -> File:
        import cv2

        directory_path = os.path.join(self.base_directory, request.upload_directory)

        try:
            make_directory(directory_path)
        except PermissionError as e:
            logger.error(f"PermissionError creating directory for processed image in file service: {e}")
            raise PermissionError(f"PermissionError creating directory for processed image in file service: {e}")
        except OSError as e:
            logger.error(f"OSError creating directory for processed image in file service: {e}")
            raise PermissionError(f"OSError creating directory for processed image in file service: {e}")

        file_path = os.path.join(directory_path, request.file_name)

        try:
            cv2.imwrite(file_path, request.image)
        except Exception as e:
            logger.error(f"Error saving processed image in file service: {e}")
            raise FileSaveError(f"Error saving processed image in file service: {e}")

        try:
            file_record = self.file_repository.create(
                model=request.model,
                model_id=request.model_id,
                file_type=request.file_type,
                file_name=request.file_name,
                file_location=file_path,
                mime_type="image/jpeg",
            )
        except DatabaseError as e:
            logger.error(f"DatabaseError creating file for processed image in file service: {e}")
            raise DatabaseError(f"DatabaseError creating file for processed image in file service: {e}")

        return file_record


    def _save_file(self, file: FileStorage, location: str):
        directory_path = os.path.join(self.base_directory, location)

        logger.info("Saving file: {}".format(directory_path))

        try:
            make_directory(directory_path)
        except PermissionError as e:
            logger.error(f"Permission denied making directory `{directory_path}`: {e}")
            raise PermissionError(f"Permission denied making directory `{directory_path}`: {e}")
        except OSError as e:
            logger.error(f"OS error when making directory `{directory_path}`: {e}")
            raise OSError(f"OS error when making directory `{directory_path}`: {e}")

        path = os.path.join(directory_path, file.filename)

        try:
            file.save(path)
        except OSError as e:
            logger.error(f"OS error when saving file: {e}")
            raise OSError(f"OS error when saving file: {e}")

        return path


