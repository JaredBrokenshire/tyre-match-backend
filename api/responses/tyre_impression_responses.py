from database.models import TyreImpression


def tyre_impression_response(tyre_impression: TyreImpression):
    return {
        "id": tyre_impression.id,
        "uuid": tyre_impression.uuid,
        "status": tyre_impression.status.value,
        "created_at": tyre_impression.created_at.isoformat() if tyre_impression.created_at else None,
    }
