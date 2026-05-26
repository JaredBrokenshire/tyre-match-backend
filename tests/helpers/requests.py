from random import randint
from utils import random_string
from werkzeug.datastructures import FileStorage
from services.file_service import FileSaveRequest
from database.models.data_types import FileModel, FileType


def generic_file_save_request(extension:str="jpg") -> FileSaveRequest:
    return FileSaveRequest(
        file=FileStorage(filename=f"{random_string()}.{extension}",),
        upload_directory=f"/tyre_match/files/test_directory",
        valid_extensions=["png", "jpg", "jpeg"],
        model=FileModel.tyre_impression,
        model_id=randint(0, 1000),
        file_type=FileType.original,
    )