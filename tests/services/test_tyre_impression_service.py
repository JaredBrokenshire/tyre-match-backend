import pytest
from unittest.mock import patch
from tests.mocks.data import MockFile
from tests.mocks.services import MockFileService
from services import TyreImpressionService, FileService
from tests.helpers.factories import TyreImpressionFactory
from database.repositories import TyreImpressionRepository
from domain import InvalidFileTypeError, FileSaveError, DatabaseError


def test_get_all():
    service = TyreImpressionService()

    tyre_impression_1 = TyreImpressionFactory().create()
    tyre_impression_2 = TyreImpressionFactory().create()
    tyre_impression_3 = TyreImpressionFactory().create()

    results, total_count = service.get_all()

    assert 3 == len(results) == total_count
    assert tyre_impression_1 in results
    assert tyre_impression_2 in results
    assert tyre_impression_3 in results


def test_get_all_pagination_page_1():
    service = TyreImpressionService()

    tyre_impression_1 = TyreImpressionFactory().create()
    tyre_impression_2 = TyreImpressionFactory().create()
    tyre_impression_3 = TyreImpressionFactory().create()

    results, total_count = service.get_all(page=1, page_size=1)

    assert 3 == total_count
    assert 1 == len(results)
    assert tyre_impression_1 == results[0]
    assert tyre_impression_2 not in results
    assert tyre_impression_3 not in results


def test_get_all_pagination_page_2():
    service = TyreImpressionService()

    tyre_impression_1 = TyreImpressionFactory().create()
    tyre_impression_2 = TyreImpressionFactory().create()
    tyre_impression_3 = TyreImpressionFactory().create()

    results, total_count = service.get_all(page=2, page_size=1)

    assert 3 == total_count
    assert 1 == len(results)
    assert tyre_impression_2 == results[0]
    assert tyre_impression_1 not in results
    assert tyre_impression_3 not in results


def test_upload_impression_image_no_file():
    service = TyreImpressionService()

    with pytest.raises(InvalidFileTypeError, match="No file provided"):
        service.upload_impression_image(None)


def test_upload_impression_image_no_filename():
    service = TyreImpressionService()

    file = MockFile(filename="")

    with pytest.raises(InvalidFileTypeError, match="No filename provided"):
        service.upload_impression_image(file)


def test_upload_impression_image_invalid_file_type():
    service = TyreImpressionService()

    file = MockFile(filename="test-file.txt")

    with pytest.raises(InvalidFileTypeError, match="Error saving file: File type not allowed"):
        service.upload_impression_image(file)


def test_upload_impression_image_permission_error_from_file_service():
    service = TyreImpressionService()

    file = MockFile(filename="test-file.jpg")

    with patch.object(FileService, "save_file", side_effect=PermissionError("test error")):
        with pytest.raises(FileSaveError, match="Error saving file: test error"):
            service.upload_impression_image(file)


def test_upload_impression_image_os_error_from_file_service():
    service = TyreImpressionService()

    file = MockFile(filename="test-file.jpg")

    with patch.object(FileService, "save_file", side_effect=OSError("test error")):
        with pytest.raises(FileSaveError, match="Error saving file: test error"):
            service.upload_impression_image(file)


def test_upload_impression_image_database_error_from_file_impression_repository():
    service = TyreImpressionService()

    # Setup mock file service
    mock_file_service = MockFileService()
    mock_file_service.save_file_response = "/test/file/path"

    file = MockFile(filename="test-file.jpg")

    with patch.object(TyreImpressionRepository, "create", side_effect=DatabaseError("test error")):
        with patch.object(FileService, "save_file", new=mock_file_service.save_file):
            with pytest.raises(DatabaseError, match="Error uploading file: test error"):
                service.upload_impression_image(file)


def test_upload_impression_image():
    service = TyreImpressionService()

    # Setup mock file service
    mock_file_service = MockFileService()
    mock_file_service.save_file_response = "/test/file/path"

    file = MockFile(filename="test-file.jpg")

    with patch.object(FileService, "save_file", new=mock_file_service.save_file):
        result = service.upload_impression_image(file)

        # Ensure tyre impression model was returned
        assert result.id != 0
        assert "/test/file/path" == result.file_path
        # Ensure file service was called correctly
        assert 1 == len(mock_file_service.save_file_calls)
        assert "test-file.jpg" in mock_file_service.save_file_calls[0][0]
        assert "/tyre_match/files/tyre_impressions" == mock_file_service.save_file_calls[0][1]

