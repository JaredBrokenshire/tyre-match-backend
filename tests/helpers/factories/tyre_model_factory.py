import random
from database.repositories import TyreModelRepository


class TyreModelFactory:
    counter = 0

    @classmethod
    def create(cls, **kwargs):
        repo = TyreModelRepository()

        cls.counter += 1

        defaults = {
            "manufacturer": f"Manufacturer {cls.counter}",
            "model_name": f"Model {cls.counter}",
            "category": "Summer",
            "vehicle_type": "Passenger Car",
            "width_mm": random.randint(170, 220),
            "aspect_ratio": random.randint(1, 100),
            "rim_diameter_inches": random.randint(12, 24),
            "groove_count": random.randint(1, 8),
            "pattern_type": "Symmetric",
            "tread_pitch_length_mm": random.randint(5, 30),
            "dataset_source": "test",
            "notes": "test notes",
        }

        defaults.update(kwargs)

        return repo.create(**defaults)
