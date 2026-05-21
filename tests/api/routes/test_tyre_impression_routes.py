import http
import uuid
from tests.mocks.data import MockFile
from tests.mocks.services import MockFileService
from database.repositories import TyreImpressionRepository


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

    _ = repo.create(uuid=uuid.uuid4(), file_path="/files/images/tyre_impressions")
    tyre_impression_2 = repo.create(uuid=uuid.uuid4(), file_path="/some/different/filepath")
    _ = repo.create(uuid=uuid.uuid4(), file_path="/another/unique/filepath")

    response = client.get("/tyre-impressions?page_size=1&page=2")

    data = response.get_json()

    # Can return correct number of items
    assert data["total_count"] == 3
    assert len(data["data"]) == 1

    # Can correctly offset response
    assert str(data["data"][0]["uuid"]) == str(tyre_impression_2.uuid)

def test_upload(client, database_session, monkeypatch):
    repo = TyreImpressionRepository()

    # Setup mock file service
    mock_file_service = MockFileService()
    monkeypatch.setattr(
        "api.routes.tyre_impression_routes.save_file",
        mock_file_service.save_file,
    )

    mock_file_service.save_file_error = ValueError("Test File Service Error")

    file = MockFile(filename="test.jpg")

    # Can not upload if there is an error from the file service
    response = client.post(
        "/tyre-impressions/upload",
        data={"file": (
            file.stream,
            file.filename,
        )},
        content_type="multipart/form-data",
    )
    assert response.status_code == http.HTTPStatus.BAD_REQUEST
    assert "File could not be saved" in response.json["error"]
    assert "Test File Service Error" in response.json["error"]

    # Can upload valid image file
    mock_file_service.reset()
    mock_file_service.save_file_response = "Test File Path"

    file = MockFile(filename="test.jpg")

    response = client.post(
        "/tyre-impressions/upload",
        data={"file": (
            file.stream,
            file.filename,
        )},
        content_type="multipart/form-data",
    )

    assert response.status_code == http.HTTPStatus.CREATED
    # Ensure file service was called correctly
    assert 1 == len(mock_file_service.save_file_calls)
    assert file.filename in mock_file_service.save_file_calls[0][0]
    assert "/files/images/tyre_impressions" == mock_file_service.save_file_calls[0][1]
    # Ensure DB record was created correctly
    tyre_impressions, total_count = repo.get_all()
    assert 1 == len(tyre_impressions)
    assert 1 == total_count
    assert "Test File Path" == tyre_impressions[0].file_path
    assert "uploaded" == tyre_impressions[0].status.value
