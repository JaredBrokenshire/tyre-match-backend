import logging
from database.extensions import db
from domain.exceptions import DatabaseError
from database.unit_of_work import UnitOfWork
from database.models.tyre_impression import TyreImpression
from database.models.data_types.tyre_impression_status import TyreImpressionStatus
from database.repositories.tyre_impression_repository import TyreImpressionRepository
from pipelines.tyre_impression_processing_pipeline import TyreImpressionProcessingPipeline

logger = logging.getLogger(__name__)


class TyreImpressionProcessingService:
    def __init__(self):
        self.tyre_impression_repository = TyreImpressionRepository()
        self.pipeline = TyreImpressionProcessingPipeline()


    def process_tyre_impression(self, tyre_impression: TyreImpression):
        # Set status -> processing
        tyre_impression.status = TyreImpressionStatus.processing

        with UnitOfWork(db.session):
            try:
                tyre_impression = self.tyre_impression_repository.update(tyre_impression)
            except DatabaseError as e:
                logger.error(f"Error setting tyre impression status `{TyreImpressionStatus.processing}` in processing service: {e}")
                raise DatabaseError(f"Error setting tyre impression status `{TyreImpressionStatus.processing}` in processing service: {e}")

            # Run pipeline
            try:
                self.pipeline.process(tyre_impression.id)
            except Exception as e:
                logger.error(f"Error processing tyre impression: {e}")
                raise e