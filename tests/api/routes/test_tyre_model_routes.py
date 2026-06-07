import http
from unittest.mock import patch
from werkzeug.datastructures import FileStorage
from database.models.tyre_model import TyreModel
from services.tyre_model_service import TyreModelService
from database.repositories.file_repository import FileRepository
from tests.helpers.factories.tyre_model_factory import TyreModelFactory
from domain.exceptions import DatabaseError, InvalidFileTypeError, FileSaveError
from tests.helpers.assertions import assert_paginated_response, assert_slim_tyre_model_response, \
    assert_tyre_model_not_in_response, assert_error_response, assert_tyre_model_response, \
    assert_is_not_tyre_model, assert_file_response


def test_get_all_empty_response(client):
    response = client.get("/tyre-models")

    # Ensure correct status_code is returned
    assert http.HTTPStatus.OK == response.status_code
    data = response.get_json()
    # Ensure returned data is empty
    assert_paginated_response(data, 0, 0)


def test_get_all_returns_models(client, database_session):
    tyre_model_1 = TyreModelFactory.create()
    tyre_model_2 = TyreModelFactory.create()
    tyre_model_3 = TyreModelFactory.create()

    response = client.get("/tyre-models")

    # Ensure correct status code is returned
    assert http.HTTPStatus.OK == response.status_code
    data = response.get_json()

    assert_paginated_response(data, 3, 3)
    # Ensure correct tyre model information was returned
    tyre_models = data["data"]
    assert_slim_tyre_model_response(tyre_models[0], tyre_model_1)
    assert_slim_tyre_model_response(tyre_models[1], tyre_model_2)
    assert_slim_tyre_model_response(tyre_models[2], tyre_model_3)


def test_get_all_pagination_get_page_1(client, database_session):
    tyre_model_1 = TyreModelFactory.create()
    tyre_model_2 = TyreModelFactory.create()
    tyre_model_3 = TyreModelFactory.create()

    response = client.get("/tyre-models?page_size=1&page=1")

    # Ensure correct status code is returned
    assert http.HTTPStatus.OK == response.status_code
    data = response.get_json()

    assert_paginated_response(data, 1, 3)

    returned_model = data["data"][0]
    # Ensure correct model was returned
    assert_slim_tyre_model_response(returned_model, tyre_model_1)
    # Ensure other models were not returned
    assert_tyre_model_not_in_response(data, tyre_model_2)
    assert_tyre_model_not_in_response(data, tyre_model_3)


def test_get_all_pagination_get_page_2(client, database_session):
    tyre_model_1 = TyreModelFactory.create()
    tyre_model_2 = TyreModelFactory.create()
    tyre_model_3 = TyreModelFactory.create()

    response = client.get("/tyre-models?page_size=1&page=2")

    # Ensure correct status code is returned
    assert http.HTTPStatus.OK == response.status_code
    data = response.get_json()

    assert_paginated_response(data, 1, 3)

    returned_model = data["data"][0]
    # Ensure correct model was returned
    assert_slim_tyre_model_response(returned_model, tyre_model_2)
    # Ensure other models were not returned
    assert_tyre_model_not_in_response(data, tyre_model_1)
    assert_tyre_model_not_in_response(data, tyre_model_3)


def test_get_all_filter_by_manufacturer(client, database_session):
    tyre_model_1 = TyreModelFactory.create(manufacturer="Test Manufacturer")
    tyre_model_2 = TyreModelFactory.create()
    tyre_model_3 = TyreModelFactory.create()

    response = client.get(f"/tyre-models?search=Test")

    # Ensure correct status code was returned
    assert http.HTTPStatus.OK == response.status_code
    data = response.get_json()
    # Ensure correct tyre model was returned
    assert_paginated_response(data, 1, 1)

    returned_model = data["data"][0]
    # Ensure correct model was returned
    assert_slim_tyre_model_response(returned_model, tyre_model_1)
    # Ensure other models were not returned
    assert_tyre_model_not_in_response(data, tyre_model_2)
    assert_tyre_model_not_in_response(data, tyre_model_3)


def test_get_all_filter_by_model_name(client, database_session):
    tyre_model_1 = TyreModelFactory.create(model_name="Test Model Name")
    tyre_model_2 = TyreModelFactory.create()
    tyre_model_3 = TyreModelFactory.create()

    response = client.get(f"/tyre-models?search=Test")

    # Ensure correct status code was returned
    assert http.HTTPStatus.OK == response.status_code
    data = response.get_json()
    # Ensure correct tyre model was returned
    assert_paginated_response(data, 1, 1)

    returned_model = data["data"][0]
    # Ensure correct model was returned
    assert_slim_tyre_model_response(returned_model, tyre_model_1)
    # Ensure other models were not returned
    assert_tyre_model_not_in_response(data, tyre_model_2)
    assert_tyre_model_not_in_response(data, tyre_model_3)


def test_get_by_id_invalid_id(client):
    response = client.get(f"/tyre-models/1000")

    # Ensure correct status code and error message were returned
    assert_error_response(response, http.HTTPStatus.NOT_FOUND, "Tyre model with id 1000 not found")


def test_get_by_id(client, database_session):
    tyre_model_1 = TyreModelFactory.create()
    tyre_model_2 = TyreModelFactory.create()

    response = client.get(f"/tyre-models/{tyre_model_1.id}")

    assert http.HTTPStatus.OK == response.status_code
    data = response.get_json()
    assert_tyre_model_response(data, tyre_model_1)
    assert_is_not_tyre_model(data, tyre_model_2)


def test_create_invalid_json(client):
    response = client.post(
        "/tyre-models",
        data="{invalid json",
        headers={"Content-Type": "application/json"},
    )

    assert_error_response(response, http.HTTPStatus.BAD_REQUEST, "Invalid JSON payload")


def test_create_missing_json(client):
    response = client.post(
        "/tyre-models",
        json="",
        headers={"Content-Type": "application/json"},
    )

    assert_error_response(response, http.HTTPStatus.BAD_REQUEST, "Missing JSON payload")


def test_create_error_from_tyre_model_service(client):
    with patch.object(TyreModelService, "create", side_effect=DatabaseError("test error")):
        response = client.post(
            "/tyre-models",
            json={
                "manufacturer": "Test Manufacturer",
                "model_name": "Test Model Name",
            },
            headers={"Content-Type": "application/json"},
        )
        # Ensure correct status code and error message were returned
        assert_error_response(response, http.HTTPStatus.INTERNAL_SERVER_ERROR, "Error creating tyre model record")


def test_create(client, database_session):
    response = client.post(
        "/tyre-models",
        json={
            "manufacturer": "Test Manufacturer",
            "model_name": "Test Model Name",
        },
        headers={"Content-Type": "application/json"},
    )

    # Ensure correct status code was returned
    assert http.HTTPStatus.CREATED == response.status_code
    data = response.get_json()
    assert_tyre_model_response(data, TyreModel(manufacturer="Test Manufacturer", model_name="Test Model Name"))


def test_upload_no_file(client):
    response = client.post(
        "/tyre-models/1000/upload",
        data={"file": None},
        content_type="multipart/form-data",
    )

    # Ensure correct status code and error message were returned
    assert_error_response(response, http.HTTPStatus.BAD_REQUEST, "No file provided")


def test_upload_no_filename(client):
    file = FileStorage(filename="")

    response = client.post(
        "/tyre-models/1000/upload",
        data={
            "file": (file.stream, file.filename)
        },
        content_type="multipart/form-data",
    )

    # Ensure correct status code and error message were returned
    assert_error_response(response, http.HTTPStatus.BAD_REQUEST, "No filename provided")


def test_upload_model_not_found_error_from_tyre_model_service(client):
    file = FileStorage(filename="test-file.png")

    response = client.post(
        "/tyre-models/1000/upload",
        data={
            "file": (file.stream, file.filename)
        },
        content_type="multipart/form-data",
    )

    assert_error_response(response, http.HTTPStatus.NOT_FOUND, "TyreModel with id 1000 not found")


def test_upload_invalid_file_type_error_from_file_service(client):
    tyre_model = TyreModelFactory.create()

    file = FileStorage(filename="test.jpg")

    with patch.object(TyreModelService, "upload_image", side_effect=InvalidFileTypeError("test error")):
        response = client.post(
            f"/tyre-models/{tyre_model.id}/upload",
            data={
                "file": (file.stream, file.filename)
            },
            content_type="multipart/form-data",
        )

        assert_error_response(response, http.HTTPStatus.BAD_REQUEST, "File type not supported")


def test_upload_file_save_error_from_file_service(client):
    tyre_model = TyreModelFactory.create()

    file = FileStorage(filename="test.jpg")

    with patch.object(TyreModelService, "upload_image", side_effect=FileSaveError("test error")):
        response = client.post(
            f"/tyre-models/{tyre_model.id}/upload",
            data={
                "file": (file.stream, file.filename)
            },
            content_type="multipart/form-data",
        )

        assert_error_response(response, http.HTTPStatus.INTERNAL_SERVER_ERROR, "Error saving file to storage")


def test_upload_database_error_from_file_service(client):
    tyre_model = TyreModelFactory.create()

    file = FileStorage(filename="test.jpg")

    with patch.object(TyreModelService, "upload_image", side_effect=DatabaseError("test error")):
        response = client.post(
            f"/tyre-models/{tyre_model.id}/upload",
            data={
                "file": (file.stream, file.filename)
            },
            content_type="multipart/form-data",
        )

        assert_error_response(response, http.HTTPStatus.INTERNAL_SERVER_ERROR, "Error saving file to database")


def test_upload(client):
    tyre_model = TyreModelFactory.create()

    file = FileStorage(filename="test.jpg")

    response = client.post(
        f"/tyre-models/{tyre_model.id}/upload",
        data={
            "file": (file.stream, file.filename)
        },
        content_type="multipart/form-data",
    )

    assert http.HTTPStatus.CREATED == response.status_code
    data = response.get_json()
    assert_tyre_model_response(data, tyre_model)
    assert 1 == len(data.get("files"))

    files, total_count = FileRepository().get_all()
    assert 1 == len(files) == total_count
    assert_file_response(data.get("files")[0], files[0])


def test_update_invalid_json(client):
    response = client.patch(
        "/tyre-models/1",
        data="{invalid json",
        headers={"Content-Type": "application/json"},
    )

    # Ensure correct status code and error message was returned
    assert_error_response(response, http.HTTPStatus.BAD_REQUEST, "Invalid JSON payload")


def test_update_missing_json(client):
    response = client.patch(
        "/tyre-models/1",
        json="",
        headers={"Content-Type": "application/json"},
    )

    # Ensure correct status code and error message was returned
    assert_error_response(response, http.HTTPStatus.BAD_REQUEST, "Missing JSON payload")


def test_update_invalid_id(client, database_session):
    response = client.patch(
        "/tyre-models/1",
        json={
            "manufacturer": "Test Manufacturer",
        },
        headers={"Content-Type": "application/json"},
    )

    # Ensure correct status code and error message was returned
    assert_error_response(response, http.HTTPStatus.NOT_FOUND, "Tyre model could not be found")


def test_update_database_error_from_tyre_model_service(client, database_session):
    tyre_model = TyreModelFactory.create()

    with patch.object(TyreModelService, "update", side_effect=DatabaseError("test error")):
        response = client.patch(
            f"/tyre-models/{tyre_model.id}",
            json={
                "manufacturer": "Test Manufacturer",
            },
            headers={"Content-Type": "application/json"},
        )

        # Ensure correct status code and error message was returned
        assert_error_response(response, http.HTTPStatus.INTERNAL_SERVER_ERROR, "Error updating tyre model record")


def test_update(client, database_session):
    tyre_model = TyreModelFactory.create()

    dto = {
        "manufacturer": "Test Manufacturer",
        "model_name": "Test Model Name",
        "category": "Test Category",
        "vehicle_type": "Test Vehicle Type",
        "width_mm": 1000,
        "aspect_ratio": 1000,
        "rim_diameter_inches": 1000,
        "groove_count": 1000,
        "tread_pitch_length_mm": 1000,
        "dataset_source": "Test Dataset Source",
        "notes": "Test Notes",
    }

    response = client.patch(
        f"/tyre-models/{tyre_model.id}",
        json=dto,
        headers={"Content-Type": "application/json"},
    )

    # Ensure correct status code was returned
    assert http.HTTPStatus.OK == response.status_code
    data = response.get_json()
    assert_tyre_model_response(data, TyreModel(**dto))


def test_delete_invalid_id(client, database_session):
    response = client.delete("/tyre-models/1")

    # Ensure correct status code and error message were returned
    assert_error_response(response, http.HTTPStatus.NOT_FOUND, "Tyre model could not be found")


def test_delete_database_error_from_tyre_model_service(client, database_session):
    tyre_model = TyreModelFactory.create()

    with patch.object(TyreModelService, "delete", side_effect=DatabaseError("test error")):
        response = client.delete(f"/tyre-models/{tyre_model.id}")

        # Ensure correct status code and error message were returned
        assert_error_response(response, http.HTTPStatus.INTERNAL_SERVER_ERROR, "Error deleting tyre model from database")


def test_delete(client, database_session):
    tyre_model = TyreModelFactory.create()

    response = client.delete(f"/tyre-models/{tyre_model.id}")

    # Ensure correct status code was returned
    assert http.HTTPStatus.NO_CONTENT == response.status_code


