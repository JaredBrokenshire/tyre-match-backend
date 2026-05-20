from flask import Blueprint, jsonify, request
from database.repositories import TyreModelRepository
from api.responses import paginated_response, slim_tyre_model_response

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