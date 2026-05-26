import pytest
from unittest.mock import patch
from services import FileService
from werkzeug.datastructures import FileStorage
from database.repositories import FileRepository
from services.file_service import FileSaveRequest
from database.models.data_types import FileModel, FileType
from tests.helpers.requests import generic_file_save_request
from domain import InvalidFileTypeError, InvalidFileError, DatabaseError


def test_handle_file_invalid_file_type_error_from_validate_file():
    service = FileService()

    request = generic_file_save_request()

    with patch("services.file_service.validate_file", side_effect=InvalidFileTypeError("test error")):
        with pytest.raises(InvalidFileTypeError, match="Invalid file type error in file service: test error"):
            service.handle_file(request)


def test_handle_file_invalid_file_error_from_validate_file():
    service = FileService()

    request = generic_file_save_request()

    with patch("services.file_service.validate_file", side_effect=InvalidFileError("test error")):
        with pytest.raises(InvalidFileError, match="Invalid file error in file service: test error"):
            service.handle_file(request)


def test_handle_file_permission_error_from_save_file():
    service = FileService()

    request = generic_file_save_request()

    with patch.object(FileService, "_save_file", side_effect=PermissionError("test error")):
        with pytest.raises(PermissionError, match="Permission error when saving file in file service: test error"):
            service.handle_file(request)


def test_handle_file_os_error_from_save_file():
    service = FileService()

    request = generic_file_save_request()

    with patch.object(FileService, "_save_file", side_effect=OSError("test error")):
        with pytest.raises(OSError, match="OS error when saving file in file service: test error"):
            service.handle_file(request)


def test_handle_file_database_error_from_repository():
    service = FileService()

    request = generic_file_save_request()

    with patch.object(FileService, "_save_file", return_value="/test/file/path"):
        with patch.object(FileRepository, "create", side_effect=DatabaseError("test error")):
            with pytest.raises(DatabaseError, match="Database error when creating file in file service: test error"):
                service.handle_file(request)


def test_handle_file():
    service = FileService()

    request = FileSaveRequest(
        file=FileStorage(filename="test-file.jpg"),
        upload_directory="/test/file/path",
        valid_extensions=["jpg"],
        model=FileModel.tyre_model,
        model_id=1,
        file_type=FileType.original
    )

    with patch.object(FileService, "_save_file", return_value=request.upload_directory):
        result = service.handle_file(request)

        # Ensure only one file record was created
        files, total_count = FileRepository().get_all()
        assert 1 == len(files) == total_count

        assert result is not None
        assert files[0] == result
        assert result.file_location == request.upload_directory
        assert result.file_name == request.file.filename


def test_save_file_permission_error_from_make_directory():
    service = FileService()

    file = FileStorage(filename="test-file.jpg")
    file_location = "/test/file/path"

    with patch("services.file_service.make_directory", side_effect=PermissionError("test error")):
        with pytest.raises(PermissionError, match=f"Permission denied making directory `{file_location}`: test error"):
            service._save_file(file, file_location)


def test_save_file_os_error_from_make_directory():
    service = FileService()

    file = FileStorage(filename="test-file.jpg")
    file_location = "/test/file/path"

    with patch("services.file_service.make_directory", side_effect=OSError("test error")):
        with pytest.raises(OSError, match=f"OS error when making directory `{file_location}`: test error"):
            service._save_file(file, file_location)


def test_save_file_os_error_from_save():
    service = FileService()

    file = FileStorage(filename="test-file.jpg")
    file_location = "/test/file/path"

    with patch.object(FileStorage, "save", side_effect=OSError("test error")):
        with pytest.raises(OSError, match=f"OS error when saving file: test error"):
            service._save_file(file, file_location)


def test_save_file():
    service = FileService()

    file = FileStorage(filename="test-file.jpg")
    file_location = "/test/file/path"

    with patch.object(FileStorage, "save"):
        result = service._save_file(file, file_location)

        assert result is not None
        assert f"{file_location}/{file.filename}" == result
