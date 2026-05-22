import os
from flask import current_app

BASE_FILE_DIRECTORY = "/files"

def allowed_file(filename: str, valid_extensions: list[str]) -> bool:
    return "." in filename and filename.rsplit(".", 1)[1] in valid_extensions

def save_file(file, upload_dir: str, valid_extensions: list[str]) -> str:
    if not file or file.filename == "":
        current_app.logger.error("File cannot be empty")
        raise ValueError("File cannot be empty")

    if not allowed_file(file.filename, valid_extensions):
        current_app.logger.error("File extension must be one of {}".format(valid_extensions))
        raise ValueError("File type not allowed")

    dir_path = os.path.join(BASE_FILE_DIRECTORY, upload_dir)

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