from database.models import TyreModel, TyreImpression


# API Response Objects

def assert_paginated_response(data: dict, data_length: int, total_count: int):
    assert "data" in data
    assert "total_count" in data
    assert total_count == data["total_count"]
    assert data_length == len(data["data"])


def assert_error_response(response, error_code: int, error_message: str):
    assert response.status_code == error_code

    data = response.get_json(silent=True)
    assert data is not None, "Response is not JSON"

    assert "error" in data
    assert error_message == data["error"]


# Tyre Models

def assert_tyre_model_response(data: dict, tyre_model: TyreModel):
    if tyre_model.id:
        assert tyre_model.id == data["id"]
    if tyre_model.manufacturer:
        assert tyre_model.manufacturer == data["manufacturer"]
    if tyre_model.model_name:
        assert tyre_model.model_name == data["model_name"]
    if tyre_model.category:
        assert tyre_model.category == data["category"]
    if tyre_model.vehicle_type:
        assert tyre_model.vehicle_type == data["vehicle_type"]
    if tyre_model.width_mm:
        assert tyre_model.width_mm == data["width_mm"]
    if tyre_model.aspect_ratio:
        assert tyre_model.aspect_ratio == data["aspect_ratio"]
    if tyre_model.rim_diameter_inches:
        assert tyre_model.rim_diameter_inches == data["rim_diameter_inches"]
    if tyre_model.groove_count:
        assert tyre_model.groove_count == data["groove_count"]
    if tyre_model.pattern_type:
        assert tyre_model.pattern_type == data["pattern_type"]
    if tyre_model.tread_pitch_length_mm:
        assert tyre_model.tread_pitch_length_mm == data["tread_pitch_length_mm"]
    if tyre_model.dataset_source:
        assert tyre_model.dataset_source == data["dataset_source"]
    if tyre_model.notes:
        assert tyre_model.notes == data["notes"]


def assert_slim_tyre_model_response(data: dict, tyre_model: TyreModel):
    assert tyre_model.id == data["id"]
    assert tyre_model.manufacturer == data["manufacturer"]
    assert tyre_model.model_name == data["model_name"]
    assert tyre_model.category == data["category"]
    assert tyre_model.vehicle_type == data["vehicle_type"]

    assert "width_mm" not in data
    assert "aspect_ratio" not in data
    assert "rim_diameter_inches" not in data
    assert "groove_count" not in data
    assert "pattern_type" not in data
    assert "tread_pitch_length_mm" not in data
    assert "dataset_source" not in data
    assert "notes" not in data

def assert_tyre_model_not_in_response(data: dict, tyre_model: TyreModel):
    matching = [
        item for item in data["data"]
        if item.get("id") == tyre_model.id
    ]

    assert not matching, f"Tyre model unexpectedly found: {matching}"

def assert_is_not_tyre_model(data: dict, tyre_model: TyreModel):
    assert data.get("id") != tyre_model.id

# Tyre Impressions

def assert_tyre_impression_response(data: dict, tyre_impression: TyreImpression):
    if tyre_impression.id:
        assert tyre_impression.id == data["id"]
    if tyre_impression.uuid:
        assert tyre_impression.uuid == data["uuid"]
    if tyre_impression.status:
        assert tyre_impression.status.value == data["status"]

def assert_tyre_impression_not_in_response(data: dict, tyre_impression: TyreImpression):
    matching = [
        item for item in data["data"]
        if item.get("id") == tyre_impression.id
    ]

    assert not matching, f"Tyre impression unexpectedly found: {matching}"