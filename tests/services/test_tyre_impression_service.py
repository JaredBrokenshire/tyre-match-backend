import pytest
from unittest.mock import patch
from tests.mocks.data import MockFile
from tests.mocks.services import MockFileService
from database.repositories import TyreImpressionRepository
from domain import InvalidFileTypeError, FileSaveError, DatabaseError
from services import TyreImpressionService, FileService


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

    # Can not upload if there is an invalid file type error from the file service
    mock_file_service.save_file_error = InvalidFileTypeError("invalid file type error")
    file = MockFile(filename="test.jpg")

    with patch.object(service.file_service, "save_file", new=mock_file_service.save_file):
        with pytest.raises(FileSaveError, match="Error saving file: invalid file type error"):
            service.upload_impression_image(file)

    # Can not upload if there is a permission error from the file service
    mock_file_service.reset()
    mock_file_service.save_file_error = PermissionError("permission error")
    file = MockFile(filename="test.jpg")

    with patch.object(service.file_service, "save_file", new=mock_file_service.save_file):
        with pytest.raises(FileSaveError, match="Error saving file: permission error"):
            service.upload_impression_image(file)

    # Can not upload if there is a OS error from the file service
    mock_file_service.reset()
    mock_file_service.save_file_error = OSError("os error")
    file = MockFile(filename="test.jpg")

    with patch.object(service.file_service, "save_file", new=mock_file_service.save_file):
        with pytest.raises(FileSaveError, match="Error saving file: os error"):
            service.upload_impression_image(file)

    # Can not upload if there is an error from the database
    mock_file_service.reset()
    file = MockFile(filename="test.jpg")

    with patch.object(service.repo, "create", side_effect=DatabaseError("database error")):
        with pytest.raises(DatabaseError, match="Error uploading file: database error"):
            service.upload_impression_image(file)

    # Can upload valid image file
    mock_file_service.reset()
    mock_file_service.save_file_response = "Test File Path"

    file = MockFile(filename="test.jpg")

    with patch.object(service.file_service, "save_file", new=mock_file_service.save_file):
        response = service.upload_impression_image(file)

        # Ensure file service was called correctly
        assert 1 == len(mock_file_service.save_file_calls)
        assert file.filename in mock_file_service.save_file_calls[0][0]
        assert "/tyre_match/files/tyre_impressions" == mock_file_service.save_file_calls[0][1]
        # Ensure DB record was created correctly
        tyre_impressions, total_count = service.repo.get_all()
        assert 1 == len(tyre_impressions)
        assert 1 == total_count
        assert "Test File Path" == tyre_impressions[0].file_path
        assert "uploaded" == tyre_impressions[0].status.value
        assert tyre_impressions[0].id == response.id