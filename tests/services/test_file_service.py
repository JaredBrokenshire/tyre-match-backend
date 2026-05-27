import pytest
import numpy as np
from random import randint
from unittest.mock import patch
from werkzeug.datastructures import FileStorage
from services.file_service import FileSaveRequest
from utils.random_generators import random_string
from database.repositories.file_repository import FileRepository
from database.models.data_types.files import FileModel, FileType
from services.file_service import FileService, ProcessedImageRequest
from domain.exceptions import InvalidFileTypeError, InvalidFileError, DatabaseError


@pytest.fixture()
def service():
    return FileService()

@pytest.fixture()
def file_save_request():
    return FileSaveRequest(
        file=FileStorage(filename=f"{random_string()}.jpg",),
        upload_directory=f"/tyre_match/files/test_directory",
        valid_extensions=["png", "jpg", "jpeg"],
        model=FileModel.tyre_impression,
        model_id=randint(0, 1000),
        file_type=FileType.original,
    )

def test_handle_file_invalid_file_type_error_from_validate_file(service, file_save_request):
    with patch("services.file_service.validate_file", side_effect=InvalidFileTypeError("test error")):
        with pytest.raises(InvalidFileTypeError, match="Invalid file type error in file service: test error"):
            service.handle_file(file_save_request)


def test_handle_file_invalid_file_error_from_validate_file(service, file_save_request):
    with patch("services.file_service.validate_file", side_effect=InvalidFileError("test error")):
        with pytest.raises(InvalidFileError, match="Invalid file error in file service: test error"):
            service.handle_file(file_save_request)


def test_handle_file_permission_error_from_save_file(service, file_save_request):
   with patch.object(FileService, "_save_file", side_effect=PermissionError("test error")):
        with pytest.raises(PermissionError, match="Permission error when saving file in file service: test error"):
            service.handle_file(file_save_request)


def test_handle_file_os_error_from_save_file(service, file_save_request):
    with patch.object(FileService, "_save_file", side_effect=OSError("test error")):
        with pytest.raises(OSError, match="OS error when saving file in file service: test error"):
            service.handle_file(file_save_request)


def test_handle_file_database_error_from_repository(service, file_save_request):
    with patch.object(FileService, "_save_file", return_value="/test/file/path"):
        with patch.object(FileRepository, "create", side_effect=DatabaseError("test error")):
            with pytest.raises(DatabaseError, match="Database error when creating file in file service: test error"):
                service.handle_file(file_save_request)


def test_handle_file(service):
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


@pytest.fixture
def processed_image_request():
    return ProcessedImageRequest(
        image=np.zeros((100, 100, 3), dtype=np.uint8),
        file_name="test-file.jpg",
        upload_directory="/tyre_match/files/test_directory",
        model=FileModel.tyre_impression,
        model_id=1,
        file_type=FileType.original
    )


def test_save_processed_image_permission_error_from_make_directory(service, processed_image_request):
    with patch("services.file_service.make_directory", side_effect=PermissionError("test error")):
        with pytest.raises(PermissionError, match="PermissionError creating directory for processed image in file service: test error"):
            service.save_processed_image(processed_image_request)


def test_save_processed_image_os_error_from_make_directory(service, processed_image_request):
    with patch("services.file_service.make_directory", side_effect=OSError("test error")):
        with pytest.raises(PermissionError, match="OSError creating directory for processed image in file service: test error"):
            service.save_processed_image(processed_image_request)

def test_save_processed_image_exception_from_cv2_imwrite(service, processed_image_request):
    with patch("cv2.imwrite", side_effect=Exception("test error")):
        with pytest.raises(Exception, match="Error saving processed image in file service: test error"):
            service.save_processed_image(processed_image_request)


def test_save_processed_image_database_error_from_file_repository(service, processed_image_request):
    with patch("cv2.imwrite"):
        with patch.object(FileRepository, "create", side_effect=DatabaseError("test error")):
            with pytest.raises(DatabaseError, match="DatabaseError creating file for processed image in file service: test error"):
                service.save_processed_image(processed_image_request)


def test_save_processed_image(service, processed_image_request):
    with patch("cv2.imwrite"):
        result = service.save_processed_image(processed_image_request)

        assert result is not None
        # Ensure only one file db record was created
        files, total_count = FileRepository().get_all()
        assert 1 == len(files) == total_count
        assert files[0] == result
        assert result.file_location == f"{processed_image_request.upload_directory}/{processed_image_request.file_name}"
        assert result.file_name == processed_image_request.file_name
        assert result.model == processed_image_request.model
        assert result.file_type == processed_image_request.file_type


def test_save_file_permission_error_from_make_directory(service):
    file = FileStorage(filename="test-file.jpg")
    file_location = "/test/file/path"

    with patch("services.file_service.make_directory", side_effect=PermissionError("test error")):
        with pytest.raises(PermissionError, match=f"Permission denied making directory `{file_location}`: test error"):
            service._save_file(file, file_location)


def test_save_file_os_error_from_make_directory(service):
    file = FileStorage(filename="test-file.jpg")
    file_location = "/test/file/path"

    with patch("services.file_service.make_directory", side_effect=OSError("test error")):
        with pytest.raises(OSError, match=f"OS error when making directory `{file_location}`: test error"):
            service._save_file(file, file_location)


def test_save_file_os_error_from_save(service):
    file = FileStorage(filename="test-file.jpg")
    file_location = "/test/file/path"

    with patch.object(FileStorage, "save", side_effect=OSError("test error")):
        with pytest.raises(OSError, match=f"OS error when saving file: test error"):
            service._save_file(file, file_location)


def test_save_file(service):
    file = FileStorage(filename="test-file.jpg")
    file_location = "/test/file/path"

    with patch.object(FileStorage, "save"):
        result = service._save_file(file, file_location)

        assert result is not None
        assert f"{file_location}/{file.filename}" == result
