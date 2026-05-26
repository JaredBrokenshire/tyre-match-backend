from database.repositories import TyreImpressionProcessingRepository
from tests.helpers.factories import TyreImpressionFactory, TyreImpressionProcessingFactory


def test_get_by_tyre_impression_id_invalid_id(database_session):
    repo = TyreImpressionProcessingRepository()

    result = repo.get_by_tyre_impression_id(1)
    assert result is None


def test_get_by_tyre_impression_id(database_session):
    repo = TyreImpressionProcessingRepository()

    tyre_impression = TyreImpressionFactory.create()
    tyre_impression_processing = TyreImpressionProcessingFactory.create(tyre_impression.id)

    result = repo.get_by_tyre_impression_id(tyre_impression.id)
    assert result == tyre_impression_processing

