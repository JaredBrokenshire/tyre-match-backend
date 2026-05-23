import pytest
from unittest.mock import patch
from domain import DatabaseError
from services import TyreModelService
from database.repositories import TyreModelRepository


def test_create():
    service = TyreModelService()

    # Can not create is there is an error from the repo
    with patch.object(TyreModelRepository, "create", side_effect=DatabaseError("test repo error")):
        with pytest.raises(DatabaseError, match="Error creating tyre_model record: test repo error"):
            service.create({
                "manufacturer": "Michelin",
                "model_name": "Pilot Sport"
            })

    # Can create tyre model
    dto = {
        "manufacturer": "Michelin",
        "model_name": "Pilot Sport",
        "category": "All-Season",
        "vehicle_type": "Passenger Car",
        "width_mm": 205,
        "aspect_ratio": 55,
        "rim_diameter_inches": 16,
        "groove_count": 3,
        "pattern_type": "Symmetrical",
        "tread_pitch_length_mm": 10,
        "dataset_source": "Google",
        "notes": "Test Notes"
    }

    response = service.create(dto)

    # Ensure record was created in db
    tyre_models, total_count = service.repo.get_all()
    assert 1 == len(tyre_models)
    assert 1 == total_count
    db_model = tyre_models[0]
    assert db_model.id == response.id
    assert db_model.manufacturer == response.manufacturer == "Michelin"
    assert db_model.model_name == response.model_name == "Pilot Sport"
    assert db_model.category == response.category == "All-Season"
    assert db_model.vehicle_type == response.vehicle_type == "Passenger Car"
    assert db_model.width_mm == response.width_mm == 205
    assert db_model.aspect_ratio == response.aspect_ratio == 55
    assert db_model.rim_diameter_inches == response.rim_diameter_inches == 16
    assert db_model.groove_count == response.groove_count == 3
    assert db_model.pattern_type == response.pattern_type == "Symmetrical"
    assert db_model.tread_pitch_length_mm == response.tread_pitch_length_mm == 10
    assert db_model.dataset_source == response.dataset_source == "Google"
    assert db_model.notes == response.notes == "Test Notes"
