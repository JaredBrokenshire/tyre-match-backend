import os
from flask import current_app
from werkzeug.datastructures import FileStorage
from domain import InvalidFileError, InvalidFileTypeError


def validate_file(file: FileStorage, valid_extensions: list[str]):
    if not file or file.filename == "":
        current_app.logger.error("Empty file provided in validate_file")
        raise InvalidFileError("Empty file provided in validate_file")

    if not allowed_file(file.filename, valid_extensions):
        current_app.logger.error(f"File extension of {file.filename} is not in {valid_extensions}")
        raise InvalidFileTypeError(f"File extension of {file.filename} is not in {valid_extensions}")


def allowed_file(filename: str, valid_extensions: list[str]) -> bool:
    return "." in filename and filename.rsplit(".", 1)[1] in valid_extensions


def make_directory(directory_path: str):
    try:
        os.makedirs(directory_path, exist_ok=True)
    except PermissionError as e:
        raise PermissionError(f"Permission denied creating directory `{directory_path}`: {e}")
    except OSError as e:
        raise OSError(f"Failed to create directory `{directory_path}`: {e}")
