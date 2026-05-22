import os
import pytest
from unittest.mock import patch
from services import file_service
from tests.mocks.data import MockFile

def test_allowed_file():
    valid_filename = "test-file.jpg"
    invalid_filename = "test-file.txt"
    valid_extensions = ["jpg", "jpeg", "png"]

    # Can use file with valid extension
    response = file_service.allowed_file(valid_filename, valid_extensions)
    assert response is True

    # Can not use with file with invalid extension
    response = file_service.allowed_file(invalid_filename, valid_extensions)
    assert response is False

def test_save_file():
    file = MockFile("")

    # Can not save file with no file
    with pytest.raises(ValueError, match="File cannot be empty"):
        file_service.save_file(None, "", [])

    # Can not save file with no filename
    with pytest.raises(ValueError, match="File cannot be empty"):
        file_service.save_file(file, "", [])

    file = MockFile("test-file.txt")
    # Can not save file with invalid extension
    with pytest.raises(ValueError, match="File type not allowed"):
        file_service.save_file(file, "", ["jpg"])

    file = MockFile("test-file.jpg")
    # Can not save file without appropriate permissions
    with patch("os.makedirs", side_effect=PermissionError("no permission")):
        with pytest.raises(PermissionError, match="Permission denied creating directory"):
            file_service.save_file(file, "images", ["jpg"])

    file = MockFile("test-file.jpg")
    # Can not save file if there is an error from the operating system
    with patch("os.makedirs", side_effect=OSError("disk error")):
        with pytest.raises(OSError, match="Failed to create directory"):
            file_service.save_file(file, "images", ["jpg"])

    file = MockFile("test-file.jpg")
    # Can not save file if there is a write failure
    with patch.object(MockFile, "save", side_effect=OSError("write failed")):
        with pytest.raises(OSError, match="Failed to save file"):
            file_service.save_file(file, "images", ["jpg"])

    file = MockFile("test-file.jpg")
    response = file_service.save_file(file, "images", ["jpg"])
    # Can save file
    assert "/files/images" in response
    assert len(os.listdir("/files/images")) == 1
    assert os.listdir("/files/images")[0] == file.filename
