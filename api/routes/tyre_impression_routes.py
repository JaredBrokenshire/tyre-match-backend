import http
from flask import Blueprint, request, jsonify, current_app
from database.repositories import TyreImpressionRepository
from services.tyre_impression_service import TyreImpressionService
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
    except Exception:
        current_app.logger.exception("Failed to upload impression image")
        return error_response(http.HTTPStatus.BAD_REQUEST, f"Error uploading impression image")

    res = tyre_impression_response(tyre_impression)
    return jsonify(res), http.HTTPStatus.CREATED
