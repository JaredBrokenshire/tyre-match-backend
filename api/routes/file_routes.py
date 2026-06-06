import http
import logging
from services.file_service import FileService
from flask import Blueprint, send_from_directory
from domain.exceptions import ModelNotFoundError
from api.responses.response_wrapper import error_response

logger = logging.getLogger(__name__)

file_blueprint = Blueprint('file', __name__)


@file_blueprint.route('/files/<int:id_>', methods=['GET'])
def get(id_):
    service = FileService()

    logger.info(f">>> Getting file by id: {id_}")

    try:
        file_record = service.get_by_id(id_)
        logger.info(f">>> FILE RECORD: {file_record}")
    except ModelNotFoundError as e:
        logger.error(f"File with id {id_} not found: {e}")
        return error_response(http.HTTPStatus.NOT_FOUND, f"File with id {id_} not found")

    return send_from_directory(
        file_record.file_location,
        file_record.file_name,
    )