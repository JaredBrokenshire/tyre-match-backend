import http
from flask import Blueprint, jsonify, request
from database.repositories import TyreModelRepository
from api.responses import paginated_response, slim_tyre_model_response, tyre_model_response, error_response

tyre_model_blueprint = Blueprint('tyre_model', __name__)

@tyre_model_blueprint.route('/tyre-models', methods=['GET'])
def get_all():
    repo = TyreModelRepository()

    # Query Parameters
    limit = request.args.get("limit", default=20, type=int)
    offset = request.args.get("offset", default=0, type=int)

    tyre_models, total_count = repo.get_all(limit, offset)

    res = paginated_response(
        [slim_tyre_model_response(t) for t in tyre_models],
        total_count
    )
    return jsonify(res)

@tyre_model_blueprint.route('/tyre-models/<int:id>', methods=['GET'])
def get_by_id(id):
    repo = TyreModelRepository()

    tyre_model = repo.get_by_id(id)
    if tyre_model is None:
        return error_response(http.HTTPStatus.NOT_FOUND, f"TyreModel with id {id} not found")

    res = tyre_model_response(tyre_model)
    return jsonify(res)