from api.responses.file_responses import file_response
from database.models.tyre_impression_processing import TyreImpressionProcessing


def tyre_impression_processing_response(tyre_impression_processing: TyreImpressionProcessing):
    res = {
        "id": tyre_impression_processing.id,
        "tyre_impression_id": tyre_impression_processing.tyre_impression_id,
        "created_at": tyre_impression_processing.created_at,
        "files": {},
    }

    if tyre_impression_processing.files:
        res["files"] = {
            key: file_response(file)
            for key, file in tyre_impression_processing.files.items()
        }

    return res