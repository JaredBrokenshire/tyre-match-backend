import os
from flask import current_app
from utils import allowed_file
from domain import InvalidFileTypeError

BASE_FILE_DIRECTORY = "/tyre_match/files"

class FileService:

    def __init__(self):
        self.base_directory = BASE_FILE_DIRECTORY

    def save_file(self, file, upload_dir: str, valid_extensions: list[str]) -> str:
        if not file or file.filename == "":
            current_app.logger.error("File cannot be empty")
            raise InvalidFileTypeError("File cannot be empty")

        if not allowed_file(file.filename, valid_extensions):
            current_app.logger.error("File extension must be one of {}".format(valid_extensions))
            raise InvalidFileTypeError("File type not allowed")

        dir_path = os.path.join(self.base_directory, upload_dir)

        current_app.logger.info("Saving file: {}".format(dir_path))

        try:
            os.makedirs(dir_path, exist_ok=True)
        except PermissionError as e:
            raise PermissionError(f"Permission denied creating directory {dir_path}") from e
        except OSError as e:
            raise OSError(f"Failed to create directory {dir_path}: {str(e)}") from e

        path = os.path.join(dir_path, file.filename)

        try:
            file.save(path)
        except OSError as e:
            raise OSError(f"Failed to save file to {path}: {str(e)}") from e

        return path