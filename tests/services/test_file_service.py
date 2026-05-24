import os
import pytest
from unittest.mock import patch
from services import FileService
from tests.mocks.data import MockFile
from domain import InvalidFileTypeError


def test_save_file_no_file():
    service = FileService()

    with pytest.raises(InvalidFileTypeError, match="File cannot be empty"):
        service.save_file(None, "", [])


def test_save_file_no_filename():
    service = FileService()

    file = MockFile("")

    with pytest.raises(InvalidFileTypeError, match="File cannot be empty"):
        service.save_file(file, "", [])


def test_save_file_invalid_extension():
    service = FileService()

    file = MockFile("test-file.txt")
    # Can not save file with invalid extension
    with pytest.raises(InvalidFileTypeError, match="File type not allowed"):
        service.save_file(file, "", ["jpg"])


def test_save_file_permission_error():
    service = FileService()

    file = MockFile("test-file.jpg")
    # Can not save file without appropriate permissions
    with patch("os.makedirs", side_effect=PermissionError("no permission")):
        with pytest.raises(PermissionError, match="Permission denied creating directory"):
            service.save_file(file, "test_directory", ["jpg"])


def test_save_file_os_error():
    service = FileService()

    file = MockFile("test-file.jpg")
    # Can not save file if there is an error from the operating system
    with patch("os.makedirs", side_effect=OSError("disk error")):
        with pytest.raises(OSError, match="Failed to create directory"):
            service.save_file(file, "test_directory", ["jpg"])


def test_save_file_write_failure():
    service = FileService()

    file = MockFile("test-file.jpg")
    # Can not save file if there is a write failure
    with patch.object(MockFile, "save", side_effect=OSError("write failed")):
        with pytest.raises(OSError, match="Failed to save file"):
            service.save_file(file, "test_directory", ["jpg"])


def test_save_file():
    service = FileService()

    file = MockFile("test-file.jpg")
    response = service.save_file(file, "test_directory", ["jpg"])

    assert "/tyre_match/files/test_directory" in response
    assert len(os.listdir("/tyre_match/files/test_directory")) == 1
    assert os.listdir("/tyre_match/files/test_directory")[0] == file.filename
