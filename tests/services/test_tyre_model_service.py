import pytest
from unittest.mock import patch
from domain import DatabaseError
from database.models import TyreModel
from services import TyreModelService
from database.repositories import TyreModelRepository
from tests.mocks.database.repositories import MockBaseRepository


def test_create():
    service = TyreModelService()

    # Setup mock repository
    mock_repo = MockBaseRepository()

    # Can not create is there is an error from the repo
    mock_repo.create_error = DatabaseError("test repo error")
    with patch.object(TyreModelRepository, "create", new=mock_repo.create):
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

    # Setup mock tyre model repository
    mock_repo.reset()
    mock_repo.create_response = TyreModel(
        manufacturer=dto.get("manufacturer"),
        model_name=dto.get("model_name"),
        category=dto.get("category"),
        vehicle_type=dto.get("vehicle_type"),
        width_mm=dto.get("width_mm"),
        aspect_ratio=dto.get("aspect_ratio"),
        rim_diameter_inches=dto.get("rim_diameter_inches"),
        groove_count=dto.get("groove_count"),
        pattern_type=dto.get("pattern_type"),
        tread_pitch_length_mm=dto.get("tread_pitch_length_mm"),
        dataset_source=dto.get("dataset_source"),
        notes=dto.get("notes"),
    )


    with patch.object(TyreModelRepository, "create", new=mock_repo.create):
        response = service.create(dto)

        # Ensure repo was called correctly
        assert 1 == len(mock_repo.create_calls)
        assert dto.get("manufacturer") == mock_repo.create_calls[0].get("manufacturer")
        assert dto.get("model_name") == mock_repo.create_calls[0].get("model_name")
        assert dto.get("category") == mock_repo.create_calls[0].get("category")
        assert dto.get("vehicle_type") == mock_repo.create_calls[0].get("vehicle_type")
        assert dto.get("width_mm") == mock_repo.create_calls[0].get("width_mm")
        assert dto.get("aspect_ratio") == mock_repo.create_calls[0].get("aspect_ratio")
        assert dto.get("rim_diameter_inches") == mock_repo.create_calls[0].get("rim_diameter_inches")
        assert dto.get("groove_count") == mock_repo.create_calls[0].get("groove_count")
        assert dto.get("pattern_type") == mock_repo.create_calls[0].get("pattern_type")
        assert dto.get("tread_pitch_length_mm") == mock_repo.create_calls[0].get("tread_pitch_length_mm")
        assert dto.get("dataset_source") == mock_repo.create_calls[0].get("dataset_source")
        assert dto.get("notes") == mock_repo.create_calls[0].get("notes")

        # Ensure created model matches input
        assert response.manufacturer == "Michelin"
        assert response.model_name == "Pilot Sport"
        assert response.category == "All-Season"
        assert response.vehicle_type == "Passenger Car"
        assert response.width_mm == 205
        assert response.aspect_ratio == 55
        assert response.rim_diameter_inches == 16
        assert response.groove_count == 3
        assert response.pattern_type == "Symmetrical"
        assert response.tread_pitch_length_mm == 10
        assert response.dataset_source == "Google"
        assert response.notes == "Test Notes"
