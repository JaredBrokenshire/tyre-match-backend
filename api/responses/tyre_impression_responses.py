from api.responses.file_responses import file_response
from database.models.tyre_impression import TyreImpression


def tyre_impression_response(tyre_impression: TyreImpression):
    res = {
        "id": tyre_impression.id,
        "uuid": tyre_impression.uuid,
        "status": tyre_impression.status.value,

        "edge_density": tyre_impression.edge_density,
        "void_ratio": tyre_impression.void_ratio,
        "groove_count": tyre_impression.groove_count,

        "files": {},

        "created_at": tyre_impression.created_at,
        "updated_at": tyre_impression.updated_at,
    }

    if tyre_impression.files:
        res["files"] = {
            key: file_response(file)
            for key, file in tyre_impression.files.items()
        }

    return res
