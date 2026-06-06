from api.responses.tyre_impression_processing_responses import tyre_impression_processing_response
from database.models.tyre_impression import TyreImpression


def tyre_impression_response(tyre_impression: TyreImpression):
    return {
        "id": tyre_impression.id,
        "uuid": tyre_impression.uuid,
        "status": tyre_impression.status.value,
        "created_at": tyre_impression.created_at.isoformat() if tyre_impression.created_at else None,

        "processing": tyre_impression_processing_response(tyre_impression.processing) if tyre_impression.processing else None,
    }
