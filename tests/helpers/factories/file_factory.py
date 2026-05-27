from utils import random_string
from database.repositories import FileRepository
from database.models.data_types import FileModel, FileType


class FileFactory:
    counter = 0

    @classmethod
    def create(cls, model: FileModel, model_id: int, **kwargs):
        repo = FileRepository()

        cls.counter += 1

        defaults = {
            "model": model,
            "model_id": model_id,
            "file_type": FileType.original,
            "file_name": f"{random_string()}.txt",
            "file_location": f"/test_directory/{random_string()}",
            "mime_type": "text/plain",
        }
        defaults.update(kwargs)

        return repo.create(**defaults)