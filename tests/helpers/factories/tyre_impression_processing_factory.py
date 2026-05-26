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

            "edge_density": random(),
            "void_ratio": random(),
            "groove_count": randint(1,10),

            "feature_vector_json": "{random: json}",
            "match_results_json": "{different_random: json}",

            "pipeline_version": 1,
        }

        defaults.update(kwargs)

        return repo.create(**defaults)