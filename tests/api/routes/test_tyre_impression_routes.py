import http
import uuid
from unittest.mock import patch
from tests.mocks.data import MockFile
from database.models import TyreImpression
from services import TyreImpressionService
from database.repositories import TyreImpressionRepository
from database.models.data_types import TyreImpressionStatus
from domain import InvalidFileTypeError, FileSaveError, DatabaseError
from tests.mocks.services.mock_tyre_impression_service import MockTyreImpressionService


def test_get_all(client, database_session):
    repo = TyreImpressionRepository()

    tyre_impression_1 = repo.create(uuid=uuid.uuid4(), file_path="/files/images/tyre_impressions")
    tyre_impression_2 = repo.create(uuid=uuid.uuid4(), file_path="/some/different/filepath")

    response = client.get("/tyre-impressions")

    # Can get status 200
    assert response.status_code == http.HTTPStatus.OK

    data = response.get_json()

    # Can return data and metadata
    assert "data" in data
    assert "total_count" in data
    assert data["total_count"] == 2

    tyre_impressions = data["data"]
    # Can return all items
    assert len(tyre_impressions) == 2

    first = tyre_impressions[0]
    second = tyre_impressions[1]

    # Can return all data
    assert first["id"] == tyre_impression_1.id
    assert str(first["uuid"]) == str(tyre_impression_1.uuid)
    assert first["file_path"] == tyre_impression_1.file_path
    assert first["status"] == tyre_impression_1.status.value
    assert first["created_at"] == tyre_impression_1.created_at.isoformat()

    assert second["id"] == tyre_impression_2.id
    assert str(second["uuid"]) == str(tyre_impression_2.uuid)
    assert second["file_path"] == tyre_impression_2.file_path
    assert second["status"] == tyre_impression_2.status.value
    assert second["created_at"] == tyre_impression_2.created_at.isoformat()


def test_get_all_empty(client):
    response = client.get("/tyre-impressions")

    # Can return status 200
    assert response.status_code == http.HTTPStatus.OK

    data = response.get_json()

    # Can return empty dataset
    assert data["total_count"] == 0
    assert data["data"] == []


def test_get_all_pagination(client, database_session):
    repo = TyreImpressionRepository()

    repo.create(uuid=uuid.uuid4(), file_path="/files/images/tyre_impressions")
    tyre_impression_2 = repo.create(uuid=uuid.uuid4(), file_path="/some/different/filepath")
    repo.create(uuid=uuid.uuid4(), file_path="/another/unique/filepath")

    response = client.get("/tyre-impressions?page_size=1&page=2")

    data = response.get_json()

    # Can return correct number of items
    assert data["total_count"] == 3
    assert len(data["data"]) == 1

    # Can correctly offset response
    assert str(data["data"][0]["uuid"]) == str(tyre_impression_2.uuid)


def test_upload(client, database_session, monkeypatch):
    # Can not upload if there is an invalid file type error from the tyre impression service
    with patch.object(TyreImpressionService, "upload_impression_image", side_effect=InvalidFileTypeError("test error")):
        file = MockFile(filename="test.jpg")
        response = client.post(
            "/tyre-impressions/upload",
            data={"file": (
                file.stream,
                file.filename,
            )},
            content_type="multipart/form-data",
        )
        # Ensure correct status code and error message are returned
        assert response.status_code == http.HTTPStatus.BAD_REQUEST
        assert "File type not supported" == response.json["error"]

    # Can not upload if there is a file save error from the tyre impression service
    with patch.object(TyreImpressionService, "upload_impression_image",
                      side_effect=FileSaveError("test error")):
        file = MockFile(filename="test.jpg")
        response = client.post(
            "/tyre-impressions/upload",
            data={"file": (
                file.stream,
                file.filename,
            )},
            content_type="multipart/form-data",
        )
        # Ensure correct status code and error message are returned
        assert response.status_code == http.HTTPStatus.INTERNAL_SERVER_ERROR
        assert "Error saving file to storage" == response.json["error"]

    # Can not upload if there is a database error from the tyre impression service
    with patch.object(TyreImpressionService, "upload_impression_image",
                      side_effect=DatabaseError("test error")):
        file = MockFile(filename="test.jpg")
        response = client.post(
            "/tyre-impressions/upload",
            data={"file": (
                file.stream,
                file.filename,
            )},
            content_type="multipart/form-data",
        )
        # Ensure correct status code and error message are returned
        assert response.status_code == http.HTTPStatus.INTERNAL_SERVER_ERROR
        assert "Error uploading file to database" == response.json["error"]

    # Setup mock tyre impression service
    mock_tyre_impression_service = MockTyreImpressionService()

    # Can upload tyre impression image
    mock_tyre_impression_service.upload_impression_image_response = TyreImpression(uuid="test-uuid",
                                                                                   file_path="/test/file/path",
                                                                                   status=TyreImpressionStatus.uploaded)

    with patch.object(TyreImpressionService, "upload_impression_image",
                      mock_tyre_impression_service.upload_impression_image):
        file = MockFile(filename="test.jpg")
        response = client.post(
            "/tyre-impressions/upload",
            data={"file": (
                file.stream,
                file.filename,
            )},
            content_type="multipart/form-data",
        )
        # Ensure tyre impression service was called with the correct parameters
        assert 1 == len(mock_tyre_impression_service.upload_impression_image_calls)
        assert file.filename == mock_tyre_impression_service.upload_impression_image_calls[0].filename
        # Ensure correct status code is returned
        assert response.status_code == http.HTTPStatus.CREATED
        data = response.get_json()
        assert "test-uuid" == data["uuid"]
        assert "/test/file/path" == data["file_path"]
        assert TyreImpressionStatus.uploaded.value == data["status"]
