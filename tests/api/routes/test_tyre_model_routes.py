import http
from unittest.mock import patch
from domain import DatabaseError
from database.models import TyreModel
from services import TyreModelService
from database.repositories import TyreModelRepository
from tests.mocks.services import MockTyreModelService


def test_get_all(client, database_session):
    repo = TyreModelRepository()

    tyre_model_1 = repo.create(manufacturer='Michelin', model_name='Pilot Sport', category="Sport", vehicle_type="SUV")
    tyre_model_2 = repo.create(manufacturer='Pirelli', model_name='P Zero', category="Winter",
                               vehicle_type="Passenger Car")

    response = client.get("/tyre-models")

    # Can get status 200
    assert response.status_code == http.HTTPStatus.OK

    data = response.get_json()

    # Can return data and metadata
    assert "data" in data
    assert "total_count" in data
    assert data["total_count"] == 2

    tyre_models = data["data"]
    # Can return all items
    assert len(tyre_models) == 2

    first = tyre_models[0]
    second = tyre_models[1]

    # Can return all data
    assert first["id"] == tyre_model_1.id
    assert first["manufacturer"] == tyre_model_1.manufacturer
    assert first["model_name"] == tyre_model_1.model_name
    assert first["category"] == tyre_model_1.category
    assert first["vehicle_type"] == tyre_model_1.vehicle_type

    assert second["id"] == tyre_model_2.id
    assert second["manufacturer"] == tyre_model_2.manufacturer
    assert second["model_name"] == tyre_model_2.model_name
    assert second["category"] == tyre_model_2.category
    assert second["vehicle_type"] == tyre_model_2.vehicle_type

    # Can not return fields not included in slim response
    assert "notes" not in first
    assert "dataset_source" not in first
    assert "groove_count" not in first
    assert "width_mm" not in first


def test_get_all_empty(client):
    response = client.get("/tyre-models")

    # Can return status 200
    assert response.status_code == http.HTTPStatus.OK

    data = response.get_json()

    # Can return empty dataset
    assert data["total_count"] == 0
    assert data["data"] == []


def test_get_all_pagination(client, database_session):
    repo = TyreModelRepository()

    repo.create(manufacturer="A", model_name="A")
    tyre_model_2 = repo.create(manufacturer="B", model_name="B")
    repo.create(manufacturer="C", model_name="C")

    response = client.get("/tyre-models?page_size=1&page=2")

    data = response.get_json()

    # Can return correct number of items
    assert data["total_count"] == 3
    assert len(data["data"]) == 1

    # Can correctly offset response
    assert data["data"][0]["manufacturer"] == tyre_model_2.manufacturer


def test_get_all_search(client, database_session):
    repo = TyreModelRepository()

    tyre_model_1 = repo.create(manufacturer='Michelin', model_name='Pilot Sport')
    tyre_model_2 = repo.create(manufacturer='Pirelli', model_name='P Zero')
    tyre_model_3 = repo.create(manufacturer='Goodyear', model_name='Eagle F1')

    response = client.get(f"/tyre-models?search={tyre_model_1.manufacturer[:3]}")

    # Can search by manufacturer
    assert response.status_code == http.HTTPStatus.OK

    data = response.get_json()

    # Can get correct tyre model from search
    assert data["total_count"] == 1
    assert len(data["data"]) == 1
    assert data["data"][0]["id"] == tyre_model_1.id
    assert data["data"][0]["manufacturer"] == tyre_model_1.manufacturer
    assert tyre_model_2.manufacturer not in data["data"]
    assert tyre_model_3.manufacturer not in data["data"]

    response = client.get(f"/tyre-models?search={tyre_model_2.model_name[:3]}")

    # Can search by model name
    assert response.status_code == http.HTTPStatus.OK

    data = response.get_json()

    # Can get correct tyre model from search
    assert data["total_count"] == 1
    assert len(data["data"]) == 1
    assert data["data"][0]["id"] == tyre_model_2.id
    assert data["data"][0]["model_name"] == tyre_model_2.model_name
    assert tyre_model_1.model_name not in data["data"]
    assert tyre_model_3.model_name not in data["data"]


def test_get_by_id_not_exist(client):
    response = client.get("/tyre-models/test")

    # Can not get tyre_model with invalid id
    assert response.status_code == http.HTTPStatus.NOT_FOUND

    response = client.get("/tyre-models/1")

    # Can not get tyre_model with id that doesn't exist
    assert response.status_code == http.HTTPStatus.NOT_FOUND
    data = response.get_json()
    # Can return appropriate error message
    assert "error" in data
    assert data["error"] == "TyreModel with id 1 not found"


def test_get_by_id(client, database_session):
    repo = TyreModelRepository()

    tyre_model = repo.create(
        manufacturer="Michelin",
        model_name="Pilot Sport",
        category="Sport",
        vehicle_type="SUV",
        width_mm=185,
        aspect_ratio=55,
        rim_diameter_inches=16,
        groove_count=1,
        pattern_type="Symmetric",
        tread_pitch_length_mm=10,
        dataset_source="google.com",
        notes="This is a test"
    )

    response = client.get(f"/tyre-models/{tyre_model.id}")

    # Can get tyre_model with valid id
    assert response.status_code == http.HTTPStatus.OK
    data = response.get_json()

    # Can return all necessary fields
    assert "id" in data
    assert "manufacturer" in data
    assert "model_name" in data
    assert "category" in data
    assert "vehicle_type" in data
    assert "width_mm" in data
    assert "aspect_ratio" in data
    assert "rim_diameter_inches" in data
    assert "groove_count" in data
    assert "pattern_type" in data
    assert "tread_pitch_length_mm" in data
    assert "dataset_source" in data
    assert "notes" in data
    assert data["id"] == tyre_model.id
    assert data["manufacturer"] == tyre_model.manufacturer
    assert data["model_name"] == tyre_model.model_name
    assert data["category"] == tyre_model.category
    assert data["vehicle_type"] == tyre_model.vehicle_type
    assert data["width_mm"] == tyre_model.width_mm
    assert data["aspect_ratio"] == tyre_model.aspect_ratio
    assert data["rim_diameter_inches"] == tyre_model.rim_diameter_inches
    assert data["groove_count"] == tyre_model.groove_count
    assert data["pattern_type"] == tyre_model.pattern_type
    assert data["tread_pitch_length_mm"] == tyre_model.tread_pitch_length_mm
    assert data["dataset_source"] == tyre_model.dataset_source
    assert data["notes"] == tyre_model.notes


def test_create(client):
    # Can not create if there is an error from the tyre model service
    with patch.object(TyreModelService, "create", side_effect=DatabaseError("test error")):
        response = client.post(
            "/tyre-models",
            json={
                "manufacturer": "Michelin",
                "model_name": "Pilot Sport",
            }
        )
        # Ensure correct status code and error message are returned
        assert response.status_code == http.HTTPStatus.INTERNAL_SERVER_ERROR
        assert "Error creating tyre model record" == response.json["error"]

    # Setup mock tyre model service
    mock_tyre_model_service = MockTyreModelService()
    mock_tyre_model_service.create_response = TyreModel(manufacturer="Michelin", model_name="Pilot Sport")

    with patch.object(TyreModelService, "create", side_effect=mock_tyre_model_service.create):
        response = client.post(
            "/tyre-models",
            json={
                "manufacturer": "Michelin",
                "model_name": "Pilot Sport",
            }
        )

        # Ensure correct status code is returned
        assert response.status_code == http.HTTPStatus.CREATED
        # Ensure tyre model service was called correctly
        assert 1 == len(mock_tyre_model_service.create_calls)
        assert "Michelin" == mock_tyre_model_service.create_calls[0].get("manufacturer")
        assert "Pilot Sport" == mock_tyre_model_service.create_calls[0].get("model_name")

        # Ensure correct tyre model response is returned
        data = response.get_json()
        assert "id" in data
        assert "manufacturer" in data
        assert "model_name" in data
        assert "Michelin" == data["manufacturer"]
        assert "Pilot Sport" == data["model_name"]