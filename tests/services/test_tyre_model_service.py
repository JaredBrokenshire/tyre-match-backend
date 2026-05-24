import pytest
from unittest.mock import patch
from services import TyreModelService
from domain import ModelNotFoundError, DatabaseError
from tests.helpers.factories import TyreModelFactory
from database.repositories import TyreModelRepository


def test_get_all():
    service = TyreModelService()

    tyre_model_1 = TyreModelFactory().create()
    tyre_model_2 = TyreModelFactory().create()
    tyre_model_3 = TyreModelFactory().create()

    results, total_count = service.get_all()

    assert 3 == len(results) == total_count
    assert tyre_model_1 in results
    assert tyre_model_2 in results
    assert tyre_model_3 in results


def test_get_all_pagination_page_1():
    service = TyreModelService()

    tyre_model_1 = TyreModelFactory().create()
    tyre_model_2 = TyreModelFactory().create()
    tyre_model_3 = TyreModelFactory().create()

    results, total_count = service.get_all(page=1, page_size=1)

    assert 3 == total_count
    assert 1 == len(results)
    assert tyre_model_1 == results[0]
    assert tyre_model_2 not in results
    assert tyre_model_3 not in results


def test_get_all_pagination_page_2():
    service = TyreModelService()

    tyre_model_1 = TyreModelFactory().create()
    tyre_model_2 = TyreModelFactory().create()
    tyre_model_3 = TyreModelFactory().create()

    results, total_count = service.get_all(page=2, page_size=1)

    assert 3 == total_count
    assert 1 == len(results)
    assert tyre_model_2 == results[0]
    assert tyre_model_1 not in results
    assert tyre_model_3 not in results


def test_get_all_search_by_manufacturer():
    service = TyreModelService()

    tyre_model_1 = TyreModelFactory().create(manufacturer='Test Manufacturer')
    tyre_model_2 = TyreModelFactory().create()
    tyre_model_3 = TyreModelFactory().create()

    results, total_count = service.get_all(search="Test")

    assert 1 == len(results) == total_count
    assert tyre_model_1 == results[0]
    assert tyre_model_2 not in results
    assert tyre_model_3 not in results


def test_get_all_search_by_model_name():
    service = TyreModelService()

    tyre_model_1 = TyreModelFactory().create(model_name='Test Model Name')
    tyre_model_2 = TyreModelFactory().create()
    tyre_model_3 = TyreModelFactory().create()

    results, total_count = service.get_all(search="Test")

    assert 1 == len(results) == total_count
    assert tyre_model_1 == results[0]
    assert tyre_model_2 not in results
    assert tyre_model_3 not in results


def test_get_by_id_invalid_id():
    service = TyreModelService()

    with pytest.raises(ModelNotFoundError, match="Error getting tyre model by id: 1"):
        service.get_by_id(1)


def test_get_by_id():
    service = TyreModelService()

    tyre_model_1 = TyreModelFactory().create()
    tyre_model_2 = TyreModelFactory().create()
    tyre_model_3 = TyreModelFactory().create()

    result = service.get_by_id(tyre_model_1.id)

    assert tyre_model_1 == result
    assert tyre_model_2 != result
    assert tyre_model_3 != result


def test_create_database_error_from_tyre_model_repository():
    service = TyreModelService()

    with patch.object(TyreModelRepository, 'create', side_effect=DatabaseError("test error")):
        with pytest.raises(DatabaseError, match="Error creating tyre model record: test error"):
            service.create({})


def test_create():
    service = TyreModelService()

    dto = {
        "manufacturer": "Test Manufacturer",
        "model_name": "Test Model Name",
    }

    result = service.create(dto)
    assert result is not None
    assert result.id != 0
    assert result.manufacturer == "Test Manufacturer"
    assert result.model_name == "Test Model Name"


def test_update_invalid_id():
    service = TyreModelService()

    with pytest.raises(ModelNotFoundError, match="Error getting tyre model by id: 1"):
        service.get_by_id(1)


def test_update_database_error_from_tyre_model_repository():
    service = TyreModelService()

    tyre_model = TyreModelFactory().create()

    with patch.object(TyreModelRepository, 'update', side_effect=DatabaseError("test error")):
        with pytest.raises(DatabaseError, match="Error updating tyre model record: test error"):
            service.update(tyre_model.id, {})


def test_update():
    service = TyreModelService()

    tyre_model = TyreModelFactory().create()

    updated_model = service.update(
        tyre_model.id,
        {
            "manufacturer": "Test Manufacturer",
            "model_name": "Test Model Name",
        }
    )

    assert tyre_model.id == updated_model.id
    assert "Test Manufacturer" == updated_model.manufacturer
    assert "Test Model Name" == updated_model.model_name


def test_delete_invalid_id():
    service = TyreModelService()

    with pytest.raises(ModelNotFoundError, match="Tyre model with id 1 not found"):
        result = service.delete(1)
        assert result == False


def test_delete_database_error_from_tyre_model_repository():
    service = TyreModelService()

    with patch.object(TyreModelRepository, 'delete', side_effect=DatabaseError("test error")):
        with pytest.raises(DatabaseError, match="Error deleting tyre model record: test error"):
            result = service.delete(1)
            assert result == False


def test_delete():
    service = TyreModelService()

    tyre_model = TyreModelFactory().create()

    result = service.delete(tyre_model.id)

    assert result == True
    with pytest.raises(ModelNotFoundError, match=f"Error getting tyre model by id: {tyre_model.id}"):
        service.get_by_id(tyre_model.id)
