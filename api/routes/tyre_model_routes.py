import http
import logging
from werkzeug.exceptions import BadRequest
from flask import Blueprint, jsonify, request
from services.tyre_model_service import TyreModelService
from api.responses.response_wrapper import paginated_response, error_response
from api.responses.tyre_model_responses import tyre_model_response, slim_tyre_model_response
from domain.exceptions import DatabaseError, ModelNotFoundError, InvalidFileTypeError, FileSaveError

logger = logging.getLogger(__name__)

tyre_model_blueprint = Blueprint('tyre_model', __name__)


@tyre_model_blueprint.route('/tyre-models', methods=['GET'])
def get_all():
    service = TyreModelService()

    # Query Parameters
    page_size = request.args.get("page_size", default=20, type=int)
    page = request.args.get("page", default=1, type=int)
    search = request.args.get("search", default="", type=str)

    tyre_models, total_count = service.get_all(
        page_size=page_size,
        page=page,
        search=search
    )

    res = paginated_response(
        [slim_tyre_model_response(t) for t in tyre_models],
        total_count
    )
    return jsonify(res)


@tyre_model_blueprint.route('/tyre-models/<int:id_>', methods=['GET'])
def get_by_id(id_):
    service = TyreModelService()

    try:
        tyre_model = service.get_by_id(id_)
    except ModelNotFoundError as e:
        logger.error(f"Tyre model with id {id_} not found: {e}")
        return error_response(http.HTTPStatus.NOT_FOUND, f"Tyre model with id {id_} not found")

    res = tyre_model_response(tyre_model)
    return jsonify(res)


@tyre_model_blueprint.route('/tyre-models', methods=['POST'])
def create():
    service = TyreModelService()

    try:
        dto = request.get_json(force=False, silent=False)
    except BadRequest as e:
        logger.error(f"Invalid JSON payload: {e}")
        return error_response(http.HTTPStatus.BAD_REQUEST, "Invalid JSON payload")

    if not dto:
        logger.error(f"Missing JSON payload: {dto}")
        return error_response(http.HTTPStatus.BAD_REQUEST, "Missing JSON payload")

    try:
        tyre_model = service.create(dto)
    except DatabaseError as e:
        logger.error(f"Error creating tyre model: {e}")
        return error_response(http.HTTPStatus.INTERNAL_SERVER_ERROR, "Error creating tyre model record")

    res = tyre_model_response(tyre_model)
    return jsonify(res), http.HTTPStatus.CREATED


@tyre_model_blueprint.route('/tyre-models/<int:id_>/upload', methods=['POST'])
def upload(id_: int):
    service = TyreModelService()

    # Extract file from request
    try:
        file = request.files['file']
    except KeyError as e:
        logger.error(f"No file provided: {e}")
        return error_response(http.HTTPStatus.BAD_REQUEST, "No file provided")

    if not file:
        logger.error(f"No filename provided: file={file}")
        return error_response(http.HTTPStatus.BAD_REQUEST, "No filename provided")

    try:
        tyre_model = service.get_by_id(id_)
    except ModelNotFoundError as e:
        logger.error(f"Tyre model with id {id_} not found: {e}")
        return error_response(http.HTTPStatus.NOT_FOUND, f"TyreModel with id {id_} not found")

    try:
        tyre_model = service.upload_image(tyre_model, file)
    except InvalidFileTypeError as e:
        logger.error(f"Invalid file type error: {e}")
        return error_response(http.HTTPStatus.BAD_REQUEST, "File type not supported")
    except FileSaveError as e:
        logger.error(f"File save error: {e}")
        return error_response(http.HTTPStatus.INTERNAL_SERVER_ERROR, "Error saving file to storage")
    except DatabaseError as e:
        logger.error(f"Database error: {e}")
        return error_response(http.HTTPStatus.INTERNAL_SERVER_ERROR, "Error saving file to database")

    res = tyre_model_response(tyre_model)
    return jsonify(res), http.HTTPStatus.CREATED


@tyre_model_blueprint.route('/tyre-models/<int:id_>', methods=['PATCH'])
def update(id_):
    service = TyreModelService()

    try:
        dto = request.get_json(force=False, silent=False)
    except BadRequest as e:
        logger.error(f"Invalid JSON payload: {e}")
        return error_response(http.HTTPStatus.BAD_REQUEST, "Invalid JSON payload")

    if not dto:
        logger.error(f"Missing JSON payload: {dto}")
        return error_response(http.HTTPStatus.BAD_REQUEST, "Missing JSON payload")

    try:
        tyre_model = service.update(id_, dto)
    except ModelNotFoundError as e:
        logger.error(f"Error fetching tyre model with id {id_} from database: {e}")
        return error_response(http.HTTPStatus.NOT_FOUND, "Tyre model could not be found")
    except DatabaseError as e:
        logger.error(f"Error updating tyre model: {e}")
        return error_response(http.HTTPStatus.INTERNAL_SERVER_ERROR, "Error updating tyre model record")

    res = tyre_model_response(tyre_model)
    return jsonify(res), http.HTTPStatus.OK


@tyre_model_blueprint.route('/tyre-models/<int:id_>', methods=['DELETE'])
def delete(id_):
    service = TyreModelService()

    try:
        service.delete(id_)
    except ModelNotFoundError as e:
        logger.error(f"Tyre model with id {id_} not found: {e}")
        return error_response(http.HTTPStatus.NOT_FOUND, "Tyre model could not be found")
    except DatabaseError as e:
        logger.error(f"Error deleting tyre model from database: {e}")
        return error_response(http.HTTPStatus.INTERNAL_SERVER_ERROR, "Error deleting tyre model from database")

    return "", http.HTTPStatus.NO_CONTENT

