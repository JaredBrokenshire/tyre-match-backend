from tests.mocks.data import MockFile
from policies import uuid_filename, original_filename, prefixed_filename


def test_uuid_filename():
    file = MockFile(
        filename="test-filename.png"
    )

    # Can generate uuid filename
    id, filename = uuid_filename(file)
    assert str(id) in filename
    assert file.filename in filename

    # Can generate secure filename
    risky_name_file = MockFile(
        filename="Risky File Name (Spaces).png"
    )

    id, filename = uuid_filename(risky_name_file)
    secure_filename = "Risky_File_Name_Spaces.png"
    assert str(id) in filename
    assert secure_filename in filename

def test_original_filename():
    # Can generate secure filename
    risky_name_file = MockFile(
        filename="Risky File Name (Spaces).png"
    )

    filename = original_filename(risky_name_file)
    secure_filename = "Risky_File_Name_Spaces.png"
    assert secure_filename == filename

def test_prefixed_filename():
    # Can generate prefixed filename
    risky_name_file = MockFile(
        filename="Risky File Name (Spaces).png"
    )

    prefix = "test-prefix"

    filename = prefixed_filename(prefix, risky_name_file)

    secure_filename = "Risky_File_Name_Spaces.png"
    assert f"{prefix}_{secure_filename}" == filename
