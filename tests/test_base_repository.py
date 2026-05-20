import pytest
from database.models import TyreModel
from sqlalchemy.exc import IntegrityError
from database.repositories.base_repository import BaseRepository


def test_create(app_ctx, database_session):
    repo = BaseRepository(TyreModel)

    # Can not create entity with required fields missing
    with pytest.raises(IntegrityError):
        repo.create()

    # Can create entity
    tyre_model = repo.create(manufacturer='Michelin', model_name='Test Tyre Model')

    assert tyre_model.manufacturer is not None
    assert tyre_model.manufacturer == 'Michelin'
    assert tyre_model.model_name is not None
    assert tyre_model.model_name == 'Test Tyre Model'

def test_get_all(app_ctx, database_session):
    repo = BaseRepository(TyreModel)

    tyre_model_1 = repo.create(manufacturer='Michelin', model_name='Test Tyre Model')
    tyre_model_2 = repo.create(manufacturer='Not Michelin', model_name='Different Model Name')
    tyre_model_3 = repo.create(manufacturer='Pirelli', model_name='Third Model')

    # Can get all results with no pagination
    results, total_count = repo.get_all()

    assert len(results) == 3
    assert total_count == 3
    assert results[0].id == tyre_model_1.id
    assert results[1].id == tyre_model_2.id
    assert results[2].id == tyre_model_3.id

    # Can limit number of results returned with pagination
    limit_results, total_count = repo.get_all(limit=2)

    assert len(limit_results) == 2
    assert total_count == 3
    assert limit_results[0].id == tyre_model_1.id
    assert limit_results[1].id == tyre_model_2.id

    # Can offset start point of pagination
    offset_results, total_count = repo.get_all(limit=2, offset=1)

    assert len(offset_results) == 2
    assert total_count == 3
    assert offset_results[0].id == tyre_model_2.id
    assert offset_results[1].id == tyre_model_3.id

    # Can not return results with invalid offset
    empty_results, total_count = repo.get_all(offset=10)

    assert len(empty_results) == 0
    assert total_count == 3

def test_get_by_id(app_ctx, database_session):
    repo = BaseRepository(TyreModel)

    # Can get by ID
    created = repo.create(manufacturer='Michelin', model_name='Test Tyre Model')
    found = repo.get_by_id(created.id)

    assert found.id is not None
    assert found.id == created.id
    assert found.manufacturer == 'Michelin'
    assert found.model_name == 'Test Tyre Model'


    # Can not get by invalid ID
    invalid_id = 999999

    found = repo.get_by_id(invalid_id)

    assert found is None


def test_delete(app_ctx, database_session):
    repo = BaseRepository(TyreModel)

    # Can delete model that exists
    tyre_model = repo.create(manufacturer='Delete Me', model_name='Delete Me')

    result = repo.delete(tyre_model.id)

    assert result is True
    assert repo.get_by_id(tyre_model.id) is None

    # Can't delete model that doesn't exist
    result = repo.delete(999999)

    assert result is False