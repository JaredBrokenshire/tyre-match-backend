import re
import pytest
from unittest.mock import patch
from werkzeug.datastructures import FileStorage
from domain import InvalidFileError, InvalidFileTypeError
from utils import allowed_file, validate_file, make_directory


def test_validate_file_no_file():
    with pytest.raises(InvalidFileError, match="Empty file provided in validate_file"):
        validate_file(None, [])


def test_validate_file_no_filename():
    file = FileStorage(filename="")

    with pytest.raises(InvalidFileError, match="Empty file provided in validate_file"):
        validate_file(file, [])


def test_validate_file_invalid_extension():
    file = FileStorage(filename="invalid.extension")

    with patch("utils.allowed_file", return_value=False):
        with pytest.raises(InvalidFileTypeError, match=re.escape("File extension of invalid.extension is not in ['png']")):
            validate_file(file, ["png"])


def test_validate_file():
    file = FileStorage(filename="test.jpg")

    with patch("utils.allowed_file", return_value=True):
        assert validate_file(file, ["jpg"]) is None


valid_extensions = ["jpg", "jpeg", "png"]


def test_allowed_file_invalid_extension():
    response = allowed_file("test-file.txt", valid_extensions)
    assert response is False


def test_allowed_file_no_extension():
    response = allowed_file("test-file", valid_extensions)
    assert response is False


def test_allowed_file():
    response = allowed_file("test-file.jpg", valid_extensions)
    assert response is True


def test_make_directory_permission_error_from_makedirs():
    with patch("os.makedirs", side_effect=PermissionError("test error")):
        with pytest.raises(PermissionError, match="Permission denied creating directory `/test/file/path`: test error"):
            make_directory("/test/file/path")


def test_make_directory_os_error_from_makedirs():
    with patch("os.makedirs", side_effect=OSError("test error")):
        with pytest.raises(OSError, match="ailed to create directory `/test/file/path`: test error"):
            make_directory("/test/file/path")


def test_make_directory():
    with patch("os.makedirs") as mock_makedirs:
        make_directory("/test/file/path")

        mock_makedirs.assert_called_once_with(
            "/test/file/path",
            exist_ok=True
        )
