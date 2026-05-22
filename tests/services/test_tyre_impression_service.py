import pytest
from unittest.mock import patch
from tests.mocks.data import MockFile
from tests.mocks.services import MockFileService
from database.repositories import TyreImpressionRepository
from services.tyre_impression_service import TyreImpressionService
from domain import InvalidFileTypeError, FileSaveError, FileUploadError


def test_upload_impression_image(monkeypatch, database_session):
    service = TyreImpressionService()

    # Can not upload with empty file
    with pytest.raises(InvalidFileTypeError, match="No file provided"):
        service.upload_impression_image(None)

    # Can not upload with empty filename
    no_name_file = MockFile(filename="")
    with pytest.raises(InvalidFileTypeError, match="No filename provided"):
        service.upload_impression_image(no_name_file)

    # Setup mock file service
    mock_file_service = MockFileService()
    monkeypatch.setattr(
        "services.tyre_impression_service.save_file",
        mock_file_service.save_file,
    )

    mock_file_service.save_file_error = ValueError("Test File Service Error")

    file = MockFile(filename="test.jpg")

    # Can not upload if there is an error from the file service
    with pytest.raises(FileSaveError, match="Error saving file: Test File Service Error"):
        service.upload_impression_image(file)

    # Can not upload if there is an error from the database
    mock_file_service.reset()

    file = MockFile(filename="test.jpg")
    with patch.object(TyreImpressionRepository, "create", side_effect=OSError("write to db failed")):
        with pytest.raises(FileUploadError, match="Error uploading file: write to db failed"):
            service.upload_impression_image(file)

    # Can upload valid image file
    mock_file_service.reset()
    mock_file_service.save_file_response = "Test File Path"

    file = MockFile(filename="test.jpg")
    response = service.upload_impression_image(file)

    # Ensure file service was called correctly
    assert 1 == len(mock_file_service.save_file_calls)
    assert file.filename in mock_file_service.save_file_calls[0][0]
    assert "/images/tyre_impressions" == mock_file_service.save_file_calls[0][1]
    # Ensure DB record was created correctly
    tyre_impressions, total_count = service.repo.get_all()
    assert 1 == len(tyre_impressions)
    assert 1 == total_count
    assert "Test File Path" == tyre_impressions[0].file_path
    assert "uploaded" == tyre_impressions[0].status.value
    assert tyre_impressions[0].id == response.id