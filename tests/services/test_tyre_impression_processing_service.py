import pytest
from unittest.mock import patch
from domain.exceptions import DatabaseError
from tests.helpers.factories.tyre_impression_factory import TyreImpressionFactory
from database.models.data_types.tyre_impression_status import TyreImpressionStatus
from database.repositories.tyre_impression_repository import TyreImpressionRepository
from services.tyre_impression_processing_service import TyreImpressionProcessingService
from pipelines.tyre_impression_processing_pipeline import TyreImpressionProcessingPipeline


def test_process_tyre_impression_database_error_from_tyre_impression_repository_update():
    service = TyreImpressionProcessingService()

    tyre_impression = TyreImpressionFactory().create()

    with patch.object(
        TyreImpressionRepository,
        "update",
        side_effect=DatabaseError("test error")
    ):
        with pytest.raises(
            DatabaseError,
            match=f"Error setting tyre impression status `{TyreImpressionStatus.processing}` in processing service: test error"
        ):
            service.process_tyre_impression(tyre_impression)


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
            service.process_tyre_impression(tyre_impression)



