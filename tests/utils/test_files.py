from utils import allowed_file


def test_allowed_file():
    valid_extensions = ["jpg", "jpeg", "png"]

    # Can not use with file with invalid extension
    response = allowed_file("test-file.txt", valid_extensions)
    assert response is False

    # Can not use file with no extension
    response = allowed_file("test-file", valid_extensions)
    assert response is False

    # Can use file with valid extension
    response = allowed_file("test-file.jpg", valid_extensions)
    assert response is True

