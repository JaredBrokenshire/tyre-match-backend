import random
from utils.random_generators import random_string
from database.repositories.tyre_model_repository import TyreModelRepository


class TyreModelFactory:
    counter = 0

    @classmethod
    def create(cls, **kwargs):
        repo = TyreModelRepository()

        cls.counter += 1

        defaults = {
            "manufacturer": random_string(10),
            "model_name": random_string(10),
            "category": random_string(10),
            "vehicle_type": random_string(10),
            "width_mm": random.randint(170, 220),
            "aspect_ratio": random.randint(1, 100),
            "rim_diameter_inches": random.randint(12, 24),
            "groove_count": random.randint(1, 8),
            "pattern_type": random_string(10),
            "tread_pitch_length_mm": random.randint(5, 30),
            "dataset_source": random_string(100),
            "notes": random_string(100),
        }

        defaults.update(kwargs)

        return repo.create(**defaults)
