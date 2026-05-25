from datetime import datetime
from random import random, randint
from database.repositories import TyreImpressionProcessingRepository


class TyreImpressionProcessingFactory:
    counter = 0

    @classmethod
    def create(cls, tyre_impression_id: int, **kwargs):
        repo = TyreImpressionProcessingRepository()

        cls.counter += 1

        defaults = {
            "tyre_impression_id": tyre_impression_id,
            "grayscale_path": f"/files/tyre_impressions/grayscale/{cls.counter}",
            "binary_path": f"/files/tyre_impressions/binary/{cls.counter}",
            "skeleton_path": f"/files/tyre_impressions/skeleton/{cls.counter}",
            "edge_density": random(),
            "void_ratio": random(),
            "groove_count": randint(1,8),
            "preprocessing_version": randint(1,10),
            "created_at": datetime.now()
        }

        defaults.update(kwargs)

        return repo.create(**defaults)