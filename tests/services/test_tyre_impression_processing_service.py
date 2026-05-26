import pytest
from unittest.mock import patch
from domain import ModelNotFoundError, DatabaseError
from pipelines import TyreImpressionProcessingPipeline
from database.repositories import TyreImpressionRepository
from database.models.data_types import TyreImpressionStatus
from tests.helpers.factories import TyreImpressionFactory, TyreImpressionProcessingFactory
from services.tyre_impression_processing_service import TyreImpressionProcessingService


def test_get_by_tyre_impression_id_invalid_id():
    service = TyreImpressionProcessingService()

    with pytest.raises(ModelNotFoundError, match="Error getting tyre impression processing by tyre impression id 999"):
        service.get_by_tyre_impression_id(999)


def test_get_by_tyre_impression_id():
    service = TyreImpressionProcessingService()

    tyre_impression_1 = TyreImpressionFactory().create()
    tyre_impression_2 = TyreImpressionFactory().create()
    tyre_impression_processing_1 = TyreImpressionProcessingFactory().create(tyre_impression_1.id)
    tyre_impression_processing_2 = TyreImpressionProcessingFactory().create(tyre_impression_2.id)

    result = service.get_by_tyre_impression_id(tyre_impression_1.id)

    assert result.id == tyre_impression_processing_1.id
    assert result.tyre_impression_id == tyre_impression_processing_1.tyre_impression_id
    assert result.id != tyre_impression_processing_2.id
    assert result.tyre_impression_id != tyre_impression_processing_2.tyre_impression_id


def test_process_tyre_impression_invalid_id():
    service = TyreImpressionProcessingService()


    with patch.object(TyreImpressionProcessingService, "_get_tyre_impression", side_effect=ModelNotFoundError("test error")):
        with pytest.raises(
                ModelNotFoundError,
                match="Error getting tyre impression in processing service with id 999: test error"):
            service.process_tyre_impression(999)


def test_process_tyre_impression_database_error_from_set_tyre_impression_status_processing():
    service = TyreImpressionProcessingService()

    tyre_impression = TyreImpressionFactory().create()

    with patch.object(
        TyreImpressionProcessingService,
        "_set_tyre_impression_status",
        side_effect=DatabaseError("test error")
    ):
        with pytest.raises(
            DatabaseError,
            match=f"Error setting tyre impression status `{TyreImpressionStatus.processing}` in processing service: test error"
        ):
            service.process_tyre_impression(tyre_impression.id)


# TODO: Update when exceptions have been defined in pipeline
def test_process_tyre_impression_error_from_pipeline():
    service = TyreImpressionProcessingService()

    tyre_impression = TyreImpressionFactory().create()

    with patch.object(
        TyreImpressionProcessingPipeline,
        "process",
        side_effect=Exception("test error")
    ):
        with pytest.raises(Exception, match="test error"):
            service.process_tyre_impression(tyre_impression.id)


def test_get_tyre_impression_invalid_id():
    service = TyreImpressionProcessingService()

    with pytest.raises(ModelNotFoundError, match="Error getting tyre impression with id 999"):
        service._get_tyre_impression(999)


def test_get_tyre_impression():
    service = TyreImpressionProcessingService()

    tyre_impression = TyreImpressionFactory().create()

    result = service._get_tyre_impression(tyre_impression.id)
    assert result is not None
    assert result.id == tyre_impression.id


def test_set_tyre_impression_status_database_error_from_repo():
    service = TyreImpressionProcessingService()

    tyre_impression = TyreImpressionFactory().create()

    with patch.object(TyreImpressionRepository, "update", side_effect=DatabaseError("test error")):
        with pytest.raises(DatabaseError, match="Error setting tyre impression status: test error"):
            service._set_tyre_impression_status(tyre_impression, TyreImpressionStatus.processing)


def test_set_tyre_impression_status():
    service = TyreImpressionProcessingService()

    tyre_impression = TyreImpressionFactory().create()

    result = service._set_tyre_impression_status(tyre_impression, TyreImpressionStatus.processing)
    assert result is not None
    assert result.id == tyre_impression.id
    assert result.status == TyreImpressionStatus.processing


