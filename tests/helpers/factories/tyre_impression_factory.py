import datetime
from database.repositories import TyreImpressionRepository
from database.models.data_types import TyreImpressionStatus


class TyreImpressionFactory:
    counter = 0

    @classmethod
    def create(cls, **kwargs):
        repo = TyreImpressionRepository()

        cls.counter += 1

        defaults = {
            "uuid": f"{cls.counter}-uuid",
            "file_path": f"{cls.counter}/file/path",
            "status": TyreImpressionStatus.uploaded,
            "created_at": datetime.datetime.now(),
        }

        defaults.update(kwargs)

        return repo.create(**defaults)