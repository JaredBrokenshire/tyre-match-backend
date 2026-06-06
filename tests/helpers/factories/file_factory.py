from pathlib import Path
from utils.random_generators import random_string
from database.repositories.file_repository import FileRepository
from database.models.data_types.files import FileModel, FileType


class FileFactory:
    counter = 0

    @classmethod
    def create(cls, model: FileModel, model_id: int, create_file_on_disk: bool = False, **kwargs):
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

        if create_file_on_disk:
            # Create directory and file on disk
            directory = Path(defaults["file_location"])
            directory.mkdir(parents=True, exist_ok=True)

            file_path = directory / defaults["file_name"]
            file_path.write_text(
                f"test content for file {cls.counter}",
                encoding="utf-8",
            )

        return repo.create(**defaults)