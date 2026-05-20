import http
from database.repositories import TyreModelRepository


def test_get_all(client, database_session):
    repo = TyreModelRepository()

    tyre_model_1 = repo.create(manufacturer='Michelin', model_name='Pilot Sport', category="Sport", vehicle_type="SUV")
    tyre_model_2 = repo.create(manufacturer='Pirelli', model_name='P Zero', category="Winter", vehicle_type="Passenger Car")

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
    repo.create(manufacturer="B", model_name="B")
    repo.create(manufacturer="C", model_name="C")

    response = client.get("/tyre-models?limit=2&offset=1")

    data = response.get_json()

    # Can return correct number of items
    assert data["total_count"] == 3
    assert len(data["data"]) == 2

    # Can correctly offset response
    assert data["data"][0]["manufacturer"] == "B"
    assert data["data"][1]["manufacturer"] == "C"

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


