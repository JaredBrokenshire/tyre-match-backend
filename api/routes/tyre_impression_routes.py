import http
from services import TyreImpressionService
from flask import Blueprint, request, jsonify, current_app
from database.repositories import TyreImpressionRepository
from domain import InvalidFileTypeError, FileSaveError, DatabaseError
from api.responses import tyre_impression_response, error_response, paginated_response

tyre_impression_blueprint = Blueprint('tyre_impression', __name__)

@tyre_impression_blueprint.route('/tyre-impressions', methods=['GET'])
def get_all():
    repo = TyreImpressionRepository()

    # Query Parameters
    page_size = request.args.get("page_size", default=20, type=int)
    page = request.args.get("page", default=1, type=int)

    tyre_impressions, total_count = repo.get_all(
        page_size=page_size,
        page=page,
    )

    res = paginated_response(
        [tyre_impression_response(t) for t in tyre_impressions],
        total_count
    )
    return jsonify(res)

@tyre_impression_blueprint.route('/tyre-impressions/upload', methods=['POST'])
def upload():
    service = TyreImpressionService()

    # Extract file from request
    file = request.files['file']
    try:
        tyre_impression = service.upload_impression_image(file)
    except InvalidFileTypeError as e:
        current_app.logger.exception(f"Invalid file type error: {e}")
        return error_response(http.HTTPStatus.BAD_REQUEST, "File type not supported")
    except FileSaveError as e:
        current_app.logger.exception(f"File save error: {e}")
        return error_response(http.HTTPStatus.INTERNAL_SERVER_ERROR, "Error saving file to storage")
    except DatabaseError as e:
        current_app.logger.exception(f"Database error: {e}")
        return error_response(http.HTTPStatus.INTERNAL_SERVER_ERROR, "Error uploading file to database")

    res = tyre_impression_response(tyre_impression)
    return jsonify(res), http.HTTPStatus.CREATED
