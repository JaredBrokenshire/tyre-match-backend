import http
from services import TyreImpressionService
from flask import Blueprint, request, jsonify, current_app
from domain import InvalidFileTypeError, FileSaveError, DatabaseError
from api.responses import tyre_impression_response, error_response, paginated_response

tyre_impression_blueprint = Blueprint('tyre_impression', __name__)

@tyre_impression_blueprint.route('/tyre-impressions', methods=['GET'])
def get_all():
    service = TyreImpressionService()

    # Query Parameters
    page_size = request.args.get("page_size", default=20, type=int)
    page = request.args.get("page", default=1, type=int)

    tyre_impressions, total_count = service.get_all(
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
    try:
        file = request.files['file']
    except KeyError as e:
        current_app.logger.error(f"No file provided: {e}")
        return error_response(http.HTTPStatus.BAD_REQUEST, "No file provided")

    if not file:
        current_app.logger.error(f"No filename provided: file={file}")
        return error_response(http.HTTPStatus.BAD_REQUEST, "No filename provided")

    try:
        tyre_impression = service.upload_impression_image(file)
    except InvalidFileTypeError as e:
        current_app.logger.error(f"Invalid file type error: {e}")
        return error_response(http.HTTPStatus.BAD_REQUEST, "File type not supported")
    except FileSaveError as e:
        current_app.logger.error(f"File save error: {e}")
        return error_response(http.HTTPStatus.INTERNAL_SERVER_ERROR, "Error saving file to storage")
    except DatabaseError as e:
        current_app.logger.error(f"Database error: {e}")
        return error_response(http.HTTPStatus.INTERNAL_SERVER_ERROR, "Error uploading file to database")

    res = tyre_impression_response(tyre_impression)
    return jsonify(res), http.HTTPStatus.CREATED
