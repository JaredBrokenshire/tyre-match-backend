import pytest
from unittest.mock import patch
from tests.mocks.data import MockFile
from database.models import TyreImpression
from services import TyreImpressionService, FileService
from tests.mocks.services import MockFileService
from database.repositories import TyreImpressionRepository
from tests.mocks.database.repositories import MockBaseRepository
from domain import InvalidFileTypeError, FileSaveError, DatabaseError


def test_upload_impression_image(database_session):
    service = TyreImpressionService()

    # Setup mocks
    mock_repo = MockBaseRepository()
    mock_file_service = MockFileService()

    # Can not upload with empty file
    with pytest.raises(InvalidFileTypeError, match="No file provided"):
        service.upload_impression_image(None)

    # Can not upload with empty filename
    no_name_file = MockFile(filename="")

    with pytest.raises(InvalidFileTypeError, match="No filename provided"):
        service.upload_impression_image(no_name_file)

    # Can not upload if there is an invalid file type error from the file service
    mock_file_service.save_file_error = InvalidFileTypeError("invalid file type error")
    file = MockFile(filename="test.jpg")

    with patch.object(FileService, "save_file", new=mock_file_service.save_file):
        with pytest.raises(FileSaveError, match="Error saving file: invalid file type error"):
            service.upload_impression_image(file)

    # Can not upload if there is a permission error from the file service
    mock_file_service.reset()
    mock_file_service.save_file_error = PermissionError("permission error")
    file = MockFile(filename="test.jpg")

    with patch.object(FileService, "save_file", new=mock_file_service.save_file):
        with pytest.raises(FileSaveError, match="Error saving file: permission error"):
            service.upload_impression_image(file)

    # Can not upload if there is a OS error from the file service
    mock_file_service.reset()
    mock_file_service.save_file_error = OSError("os error")
    file = MockFile(filename="test.jpg")

    with patch.object(FileService, "save_file", new=mock_file_service.save_file):
        with pytest.raises(FileSaveError, match="Error saving file: os error"):
            service.upload_impression_image(file)

    # Can not upload if there is an error from the database
    mock_repo.reset()
    mock_file_service.reset()

    mock_repo.create_error = DatabaseError("test repo error")

    file = MockFile(filename="test.jpg")

    with patch.object(TyreImpressionRepository, "create", new=mock_repo.create):
        with patch.object(FileService, "save_file", mock_file_service.save_file):
            with pytest.raises(DatabaseError, match="Error uploading file: test repo error"):
                service.upload_impression_image(file)

    # Can upload valid image file
    mock_repo.reset()
    mock_repo.create_response = TyreImpression(
        uuid="test-uuid",
        file_path="Test File Path",
    )
    mock_file_service.reset()
    mock_file_service.save_file_response = "Test File Path"

    file = MockFile(filename="test.jpg")

    with patch.object(FileService, "save_file", new=mock_file_service.save_file):
        with patch.object(TyreImpressionRepository, "create", new=mock_repo.create):
            response = service.upload_impression_image(file)

            # Ensure file service was called correctly
            assert 1 == len(mock_file_service.save_file_calls)
            assert file.filename in mock_file_service.save_file_calls[0][0]
            assert "/tyre_match/files/tyre_impressions" == mock_file_service.save_file_calls[0][1]
            # Ensure repo was called correctly
            assert 1 == len(mock_repo.create_calls)
            assert response.uuid == "test-uuid"
            assert response.file_path == "Test File Path"
