import pytest
from domain.exceptions import ModelNotFoundError
from database.models.data_types.files import FileModel
from tests.helpers.factories.file_factory import FileFactory
from tests.helpers.factories.tyre_model_factory import TyreModelFactory
from database.repositories.tyre_model_repository import TyreModelRepository


@pytest.fixture()
def repo():
    return TyreModelRepository()


def test_get_by_id_invalid_id(repo):
    with pytest.raises(ModelNotFoundError, match="TyreModel with id 1000 not found"):
        repo.get_by_id(1000)


def test_get_by_id(repo):
    tyre_model_1 = TyreModelFactory().create()
    file_1 = FileFactory().create(FileModel.tyre_model, tyre_model_1.id)

    tyre_model_2 = TyreModelFactory().create()
    file_2 = FileFactory().create(FileModel.tyre_model, tyre_model_2.id)

    res = repo.get_by_id(tyre_model_1.id)

    assert res.id == tyre_model_1.id
    assert res.id != tyre_model_2.id
    assert res.model_name == tyre_model_1.model_name
    assert res.manufacturer == tyre_model_1.manufacturer

    assert 1 == len(res.files)
    assert res.files[0] == file_1
    assert file_2 not in res.files