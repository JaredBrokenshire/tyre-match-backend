import pytest
from unittest.mock import patch
from database.models import File
from werkzeug.datastructures import FileStorage
from services import TyreImpressionService, FileService
from tests.helpers.factories import TyreImpressionFactory
from database.models.data_types import TyreImpressionStatus, FileModel, FileType
from domain import InvalidFileTypeError, FileSaveError, DatabaseError, InvalidFileError
from database.repositories import TyreImpressionRepository, TyreImpressionProcessingRepository


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


def test_upload_impression_image_invalid_file_type_error_from_validate_file():
    service = TyreImpressionService()

    file = FileStorage(filename="test-file.txt")

    with patch("services.tyre_impression_service.validate_file", side_effect=InvalidFileTypeError("test error")):
        with pytest.raises(InvalidFileTypeError, match="Invalid file type error from validate_file in tyre impression service: test error"):
            service.upload_impression_image(file)


def test_upload_impression_image_invalid_file_error_from_validate_file():
    service = TyreImpressionService()

    file = FileStorage(filename="test-file.txt")

    with patch("services.tyre_impression_service.validate_file", side_effect=InvalidFileError("test error")):
        with pytest.raises(InvalidFileTypeError, match="Invalid file error from validate_file in tyre impression service: test error"):
            service.upload_impression_image(file)


def test_upload_impression_image_database_error_from_tyre_impression_repository_create():
    service = TyreImpressionService()

    file = FileStorage(filename="test-file.jpg")

    with patch.object(TyreImpressionRepository, "create", side_effect=DatabaseError("test error")):
        with patch.object(FileService, "handle_file"):
            with pytest.raises(DatabaseError, match="Error creating tyre impression record in tyre impression service: test error"):
                service.upload_impression_image(file)


def test_upload_impression_image_invalid_file_error_from_file_service():
    service = TyreImpressionService()

    file = FileStorage(filename="test-file.jpg")

    with patch.object(FileService, "handle_file", side_effect=InvalidFileError("test error")):
        with pytest.raises(FileSaveError, match="Invalid file error from file service in tyre impression service: test error"):
            service.upload_impression_image(file)


def test_upload_impression_image_invalid_file_type_error_from_file_service():
    service = TyreImpressionService()

    file = FileStorage(filename="test-file.jpg")

    with patch.object(FileService, "handle_file", side_effect=InvalidFileTypeError("test error")):
        with pytest.raises(FileSaveError, match="Invalid file type error from file service in tyre impression service: test error"):
            service.upload_impression_image(file)


def test_upload_impression_image_permission_error_from_file_service():
    service = TyreImpressionService()

    file = FileStorage(filename="test-file.jpg")

    with patch.object(FileService, "handle_file", side_effect=PermissionError("test error")):
        with pytest.raises(FileSaveError, match="Permission or OS error from file service in tyre impression service: test error"):
            service.upload_impression_image(file)


def test_upload_impression_image_os_error_from_file_service():
    service = TyreImpressionService()

    file = FileStorage(filename="test-file.jpg")

    with patch.object(FileService, "handle_file", side_effect=OSError("test error")):
        with pytest.raises(FileSaveError, match="Permission or OS error from file service in tyre impression service: test error"):
            service.upload_impression_image(file)


def test_upload_impression_image_database_error_from_file_service():
    service = TyreImpressionService()

    file = FileStorage(filename="test-file.jpg")

    with patch.object(FileService, "handle_file", side_effect=DatabaseError("test error")):
        with pytest.raises(FileSaveError, match="Database error from file service in tyre impression service: test error"):
            service.upload_impression_image(file)


def test_upload_impression_image_database_error_from_tyre_impression_processing_repository():
    service = TyreImpressionService()

    file = FileStorage(filename="test-file.jpg")

    with patch.object(FileService, "handle_file", return_value=File(
        model=FileModel.tyre_impression,
        file_type=FileType.original,
        file_name="test-file.jpg",
        file_location="/tyre_match/files/test_directory",
        mime_type="image/jpeg",
    )):
        with patch.object(TyreImpressionProcessingRepository, "create", side_effect=DatabaseError("test error")):
            with pytest.raises(DatabaseError, match="Error creating tyre impression processing record in tyre impression service: test error"):
                service.upload_impression_image(file)



def test_upload_impression_image():
    service = TyreImpressionService()

    file = FileStorage(filename="test-file.jpg")

    with patch.object(FileService, "handle_file"):
        with patch("services.tyre_impression_service.process_tyre_impression_task"):
            with patch.object(FileService, "handle_file", return_value=File(
                    model=FileModel.tyre_impression,
                    file_type=FileType.original,
                    file_name="test-file.jpg",
                    file_location="/tyre_match/files/test_directory",
                    mime_type="image/jpeg",
            )):
                result = service.upload_impression_image(file)

                # Ensure tyre impression model was returned
                assert result is not None
                assert result.status == TyreImpressionStatus.uploaded

                # Ensure only one tyre impression record was created
                tyre_impressions, total_count = TyreImpressionRepository().get_all()
                assert 1 == len(tyre_impressions) == total_count
                assert tyre_impressions[0] == result

                # Ensure only one tyre impression processing record was created
                tyre_impression_processing_records, total_count = TyreImpressionProcessingRepository().get_all()
                assert 1 == len(tyre_impression_processing_records) == total_count
                assert result.id == tyre_impression_processing_records[0].tyre_impression_id
