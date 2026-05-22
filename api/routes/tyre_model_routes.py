import http
from database.repositories import TyreModelRepository
from flask import Blueprint, jsonify, request, current_app
from api.responses import paginated_response, slim_tyre_model_response, tyre_model_response, error_response
from domain import DatabaseError
from services import TyreModelService

tyre_model_blueprint = Blueprint('tyre_model', __name__)

@tyre_model_blueprint.route('/tyre-models', methods=['GET'])
def get_all():
    repo = TyreModelRepository()

    # Query Parameters
    page_size = request.args.get("page_size", default=20, type=int)
    page = request.args.get("page", default=1, type=int)
    search = request.args.get("search", default="", type=str)

    tyre_models, total_count = repo.get_all(
        page_size=page_size,
        page=page,
        search=search
    )

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

@tyre_model_blueprint.route('/tyre-models', methods=['POST'])
def create():
    service = TyreModelService()

    request_json = request.get_json()
    try:
        tyre_model = service.create(request_json)
    except DatabaseError:
        return error_response(http.HTTPStatus.INTERNAL_SERVER_ERROR, "Error creating tyre model record")

    res = tyre_model_response(tyre_model)
    return jsonify(res), http.HTTPStatus.CREATED