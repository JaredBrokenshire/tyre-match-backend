from werkzeug.datastructures import FileStorage
from policies import uuid_filename, original_filename, prefixed_filename


def test_uuid_filename():
    file = FileStorage(
        filename="test-filename.png"
    )

    # Can generate uuid filename
    uuid_, filename = uuid_filename(file)
    assert filename == f"{uuid_}.png"


def test_original_filename():
    # Can generate secure filename
    risky_name_file = FileStorage(
        filename="Risky File Name (Spaces).png"
    )

    filename = original_filename(risky_name_file)
    secure_filename = "Risky_File_Name_Spaces.png"
    assert secure_filename == filename

def test_prefixed_filename():
    # Can generate prefixed filename
    risky_name_file = FileStorage(
        filename="Risky File Name (Spaces).png"
    )

    prefix = "test-prefix"

    filename = prefixed_filename(prefix, risky_name_file)

    secure_filename = "Risky_File_Name_Spaces.png"
    assert f"{prefix}_{secure_filename}" == filename
