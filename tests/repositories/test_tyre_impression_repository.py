from database.models.data_types.files import FileModel
from tests.helpers.factories.file_factory import FileFactory
from tests.helpers.factories.tyre_impression_factory import TyreImpressionFactory
from database.repositories.tyre_impression_repository import TyreImpressionRepository


def test_get_by_id_invalid_id():
    repo = TyreImpressionRepository()

    result = repo.get_by_id(1000)

    assert result is None


def test_get_by_id():
    repo = TyreImpressionRepository()

    tyre_impression = TyreImpressionFactory().create()

    file = FileFactory().create(
        model=FileModel.tyre_impression,
        model_id=tyre_impression.id,
    )

    result = repo.get_by_id(tyre_impression.id)
    assert tyre_impression == result

    # Ensure file was loaded correctly
    assert 1 == len(result.files)
    assert file == result.files[file.file_type.value]
