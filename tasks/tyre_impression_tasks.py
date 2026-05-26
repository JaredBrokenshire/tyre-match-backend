from celery_app import celery
from services.tyre_impression_processing_service import TyreImpressionProcessingService


@celery.task(autoretry_for=(Exception,), retry_backoff=True, max_retries=3)
def process_tyre_impression_task(tyre_impression_id: int):
    service = TyreImpressionProcessingService()
    return service.process_tyre_impression(tyre_impression_id)