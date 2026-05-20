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