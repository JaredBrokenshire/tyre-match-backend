from database.models import File
from database.repositories.base_repository import BaseRepository


class FileRepository(BaseRepository):
    def __init__(self):
        super().__init__(File)