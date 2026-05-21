import http
from tests.mocks.data import MockFile
from tests.mocks.services import MockFileService
from database.repositories import TyreImpressionRepository


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
        "/tyre-impression/upload",
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
        "/tyre-impression/upload",
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
