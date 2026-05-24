from utils import allowed_file

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

