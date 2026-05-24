import http
from unittest.mock import patch
from tests.mocks.data import MockFile
from database.models import TyreImpression
from tests.mocks.services import MockFileService
from domain import InvalidFileTypeError, DatabaseError
from services import TyreImpressionService, FileService
from tests.helpers.factories import TyreImpressionFactory
from database.models.data_types import TyreImpressionStatus
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


def test_upload_impression_image_no_file(client):
    response = client.post(
        "/tyre-impressions/upload",
        data={"file": None},
        content_type="multipart/form-data",
    )

    # Ensure correct status code and error message were returned
    assert_error_response(response, http.HTTPStatus.BAD_REQUEST, "No file provided")


def test_upload_impression_image_no_filename(client):
    file = MockFile(filename="")

    response = client.post(
        "/tyre-impressions/upload",
        data={
            "file": (file.stream, file.filename)
        },
        content_type="multipart/form-data",
    )

    # Ensure correct status code and error message were returned
    assert_error_response(response, http.HTTPStatus.BAD_REQUEST, "No filename provided")


def test_upload_impression_image_invalid_file_type(client):
    with patch.object(FileService, "save_file", side_effect=InvalidFileTypeError("test error")):
        file = MockFile(filename="test.txt")

        response = client.post(
            "/tyre-impressions/upload",
            data={
                "file": (file.stream, file.filename)
            },
            content_type="multipart/form-data",
        )

        # Ensure correct status code and error message were returned
        assert_error_response(response, http.HTTPStatus.BAD_REQUEST, "File type not supported")


def test_upload_impression_image_permission_error_from_file_service(client):
    with patch.object(FileService, "save_file", side_effect=PermissionError("test error")):
        file = MockFile(filename="test.jpg")

        response = client.post(
            "/tyre-impressions/upload",
            data={
                "file": (file.stream, file.filename)
            },
            content_type="multipart/form-data",
        )

        # Ensure correct status code and error message were returned
        assert_error_response(response, http.HTTPStatus.INTERNAL_SERVER_ERROR, "Error saving file to storage")


def test_upload_impression_image_os_error_from_file_service(client):
    with patch.object(FileService, "save_file", side_effect=OSError("test error")):
        file = MockFile(filename="test.jpg")

        response = client.post(
            "/tyre-impressions/upload",
            data={
                "file": (file.stream, file.filename)
            },
            content_type="multipart/form-data",
        )

        # Ensure correct status code and error message were returned
        assert_error_response(response, http.HTTPStatus.INTERNAL_SERVER_ERROR, "Error saving file to storage")


def test_upload_impression_image_database_error_from_tyre_impression_service(client):
    with patch.object(TyreImpressionService, "upload_impression_image", side_effect=DatabaseError("test error")):
        file = MockFile(filename="test.jpg")

        response = client.post(
            "/tyre-impressions/upload",
            data={
                "file": (file.stream, file.filename)
            },
            content_type="multipart/form-data",
        )

        # Ensure correct status code and error message were returned
        assert_error_response(response, http.HTTPStatus.INTERNAL_SERVER_ERROR, "Error uploading file to database")

def test_upload_impression_image(client):
    mock_file_service = MockFileService()
    mock_file_service.save_file_response = "/test/path"

    with patch.object(FileService, "save_file", side_effect=mock_file_service.save_file):
        file = MockFile(filename="test.jpg")

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
        # Ensure tyre impression model was returned
        assert_tyre_impression_response(
            data,
            TyreImpression(
                file_path="/test/path",
                status=TyreImpressionStatus.uploaded,
            )
        )
