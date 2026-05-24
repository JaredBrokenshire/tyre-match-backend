import pytest
from database.models import TyreModel
from domain import DatabaseError, ModelNotFoundError
from tests.helpers.factories import TyreModelFactory
from database.repositories.base_repository import BaseRepository


def test_get_all_returns_models():
    repo = BaseRepository(TyreModel)

    tyre_model_1 = TyreModelFactory.create()
    tyre_model_2 = TyreModelFactory.create()
    tyre_model_3 = TyreModelFactory.create()

    results, total_count = repo.get_all()
    # Ensure number of results is correct
    assert 3 == len(results) == total_count
    # Ensure correct models were returned
    assert tyre_model_1.id == results[0].id
    assert tyre_model_2.id == results[1].id
    assert tyre_model_3.id == results[2].id


def test_get_all_pagination_page_1():
    repo = BaseRepository(TyreModel)

    tyre_model_1 = TyreModelFactory.create()
    tyre_model_2 = TyreModelFactory.create()
    tyre_model_3 = TyreModelFactory.create()

    results, total_count = repo.get_all(page=1, page_size=1)

    assert 1 == len(results)
    assert 3 == total_count
    assert tyre_model_1.id == results[0].id
    assert tyre_model_2 not in results
    assert tyre_model_3 not in results


def test_get_all_pagination_page_2():
    repo = BaseRepository(TyreModel)

    tyre_model_1 = TyreModelFactory.create()
    tyre_model_2 = TyreModelFactory.create()
    tyre_model_3 = TyreModelFactory.create()

    results, total_count = repo.get_all(page=2, page_size=1)

    assert 1 == len(results)
    assert 3 == total_count
    assert tyre_model_2.id == results[0].id
    assert tyre_model_1 not in results
    assert tyre_model_3 not in results


def test_get_all_filtered():
    repo = BaseRepository(TyreModel)

    tyre_model_1 = TyreModelFactory.create(manufacturer="Test Manufacturer")
    tyre_model_2 = TyreModelFactory.create()
    tyre_model_3 = TyreModelFactory.create()

    search_term = f"%Test%"

    results, total_count = repo.get_all(filters=TyreModel.manufacturer.ilike(search_term))

    assert 1 == len(results) == total_count
    assert tyre_model_1.id == results[0].id
    assert tyre_model_2 not in results
    assert tyre_model_3 not in results


def test_get_by_id_invalid_id():
    repo = BaseRepository(TyreModel)

    result = repo.get_by_id(1)
    assert result is None


def test_get_by_id():
    repo = BaseRepository(TyreModel)

    tyre_model_1 = TyreModelFactory.create()

    result = repo.get_by_id(tyre_model_1.id)
    assert tyre_model_1.manufacturer == result.manufacturer
    assert tyre_model_1.model_name == result.model_name


def test_create_missing_required_fields():
    repo = BaseRepository(TyreModel)

    tyre_model = {}

    with pytest.raises(DatabaseError, match="Error inserting record into database: .* \"Column 'manufacturer' cannot be null\""):
        repo.create(**tyre_model)


def test_create():
    repo = BaseRepository(TyreModel)

    tyre_model = {
        "manufacturer": "Test Manufacturer",
        "model_name": "Test Model",
    }

    result = repo.create(**tyre_model)
    assert result.id != 0
    assert tyre_model["manufacturer"] == result.manufacturer
    assert tyre_model["model_name"] == "Test Model"


def test_update_invalid_data_type():
    repo = BaseRepository(TyreModel)

    tyre_model = TyreModelFactory.create()

    updated_model = {
        "width_mm": "This should be a number",
    }

    with pytest.raises(DatabaseError, match="Error updating record in database: .*Incorrect integer value: 'This should be a number' for column 'width_mm'"):
        repo.update(tyre_model, **updated_model)


def test_update():
    repo = BaseRepository(TyreModel)

    tyre_model = TyreModelFactory.create()

    updated_model = {
        "width_mm": 100,
    }

    result = repo.update(tyre_model, **updated_model)

    assert tyre_model.id == result.id
    assert 100 == result.width_mm


def test_delete_invalid_id():
    repo = BaseRepository(TyreModel)

    with pytest.raises(ModelNotFoundError, match="Model with id 1 not found"):
        repo.delete(1)


def test_delete():
    repo = BaseRepository(TyreModel)

    tyre_model = TyreModelFactory.create()

    res = repo.delete(tyre_model.id)

    assert True == res
    deleted_model = repo.get_by_id(tyre_model.id)
    assert None == deleted_model