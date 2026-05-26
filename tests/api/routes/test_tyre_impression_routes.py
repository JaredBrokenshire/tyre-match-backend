import http
from unittest.mock import patch
from database.models import TyreImpression
from services import TyreImpressionService
from werkzeug.datastructures import FileStorage
from tests.helpers.factories import TyreImpressionFactory
from database.models.data_types import TyreImpressionStatus
from domain import InvalidFileTypeError, DatabaseError, FileSaveError
from tests.helpers.assertions import assert_paginated_response, assert_tyre_impression_response, \
    assert_tyre_impression_not_in_response, assert_error_response


def test_get_all_empty_response(client):
    response = client.get("/tyre-impressions")

    # Ensure correct status_code is returned
    assert http.HTTPStatus.OK == response.status_code
    data = response.get_json()
    # Ensure returned data is empty
    assert_paginated_response(data, 0, 0)


def test_get_all_returns_models(client):
    tyre_impression_1 = TyreImpressionFactory.create()
    tyre_impression_2 = TyreImpressionFactory.create()
    tyre_impression_3 = TyreImpressionFactory.create()

    response = client.get("/tyre-impressions")

    # Ensure correct status code is returned
    assert http.HTTPStatus.OK == response.status_code
    data = response.get_json()

    assert_paginated_response(data, 3, 3)
    # Ensure correct tyre model information was returned
    tyre_models = data["data"]
    assert_tyre_impression_response(tyre_models[0], tyre_impression_1)
    assert_tyre_impression_response(tyre_models[1], tyre_impression_2)
    assert_tyre_impression_response(tyre_models[2], tyre_impression_3)


def test_get_all_pagination_get_page_1(client):
    tyre_impression_1 = TyreImpressionFactory.create()
    tyre_impression_2 = TyreImpressionFactory.create()
    tyre_impression_3 = TyreImpressionFactory.create()

    response = client.get("/tyre-impressions?page_size=1&page=1")

    # Ensure correct status code is returned
    assert http.HTTPStatus.OK == response.status_code
    data = response.get_json()

    assert_paginated_response(data, 1, 3)

    returned_impression = data["data"][0]
    # Ensure correct model was returned
    assert_tyre_impression_response(returned_impression, tyre_impression_1)
    # Ensure other models were not returned
    assert_tyre_impression_not_in_response(data, tyre_impression_2)
    assert_tyre_impression_not_in_response(data, tyre_impression_3)


def test_get_all_pagination_get_page_2(client):
    tyre_impression_1 = TyreImpressionFactory.create()
    tyre_impression_2 = TyreImpressionFactory.create()
    tyre_impression_3 = TyreImpressionFactory.create()

    response = client.get("/tyre-impressions?page_size=1&page=2")

    # Ensure correct status code is returned
    assert http.HTTPStatus.OK == response.status_code
    data = response.get_json()

    assert_paginated_response(data, 1, 3)

    returned_model = data["data"][0]
    # Ensure correct model was returned
    assert_tyre_impression_response(returned_model, tyre_impression_2)
    # Ensure other models were not returned
    assert_tyre_impression_not_in_response(data, tyre_impression_1)
    assert_tyre_impression_not_in_response(data, tyre_impression_3)


def test_upload_no_file(client):
    response = client.post(
        "/tyre-impressions/upload",
        data={"file": None},
        content_type="multipart/form-data",
    )

    # Ensure correct status code and error message were returned
    assert_error_response(response, http.HTTPStatus.BAD_REQUEST, "No file provided")


def test_upload_no_filename(client):
    file = FileStorage(filename="")

    response = client.post(
        "/tyre-impressions/upload",
        data={
            "file": (file.stream, file.filename)
        },
        content_type="multipart/form-data",
    )

    # Ensure correct status code and error message were returned
    assert_error_response(response, http.HTTPStatus.BAD_REQUEST, "No filename provided")


def test_upload_invalid_file_type_error_from_service(client):
    file = FileStorage(filename="test.jpg")

    with patch.object(TyreImpressionService, "upload_impression_image", side_effect=InvalidFileTypeError("test error")):
        response = client.post(
            "/tyre-impressions/upload",
            data={
                "file": (file.stream, file.filename)
            },
            content_type="multipart/form-data",
        )

        assert_error_response(response, http.HTTPStatus.BAD_REQUEST, "File type not supported")


def test_upload_file_save_error_from_service(client):
    file = FileStorage(filename="test.jpg")

    with patch.object(TyreImpressionService, "upload_impression_image", side_effect=FileSaveError("test error")):
        response = client.post(
            "/tyre-impressions/upload",
            data={
                "file": (file.stream, file.filename)
            },
            content_type="multipart/form-data",
        )

        assert_error_response(response, http.HTTPStatus.INTERNAL_SERVER_ERROR, "Error saving file to storage")


def test_upload_database_error_from_service(client):
    file = FileStorage(filename="test.jpg")

    with patch.object(TyreImpressionService, "upload_impression_image", side_effect=DatabaseError("test error")):
        response = client.post(
            "/tyre-impressions/upload",
            data={
                "file": (file.stream, file.filename)
            },
            content_type="multipart/form-data",
        )

        assert_error_response(response, http.HTTPStatus.INTERNAL_SERVER_ERROR, "Error saving file to database")


def test_upload(client):
    file = FileStorage(filename="test.jpg")

    mock_tyre_impression = TyreImpression(
        uuid="test-uuid",
        status=TyreImpressionStatus.uploaded
    )

    with patch.object(TyreImpressionService, "upload_impression_image", return_value=mock_tyre_impression):
        response = client.post(
            "/tyre-impressions/upload",
            data={
                "file": (file.stream, file.filename)
            },
            content_type="multipart/form-data",
        )

        # Ensure correct status code was returned
        assert http.HTTPStatus.CREATED == response.status_code

        data = response.get_json()
        assert_tyre_impression_response(data, mock_tyre_impression)
