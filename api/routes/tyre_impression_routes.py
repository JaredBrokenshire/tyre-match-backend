import http

import policies
from services import save_file
from flask import Blueprint, request, jsonify
from database.repositories import TyreImpressionRepository
from api.responses import tyre_impression_response, error_response

tyre_impression_blueprint = Blueprint('tyre_impression', __name__)

@tyre_impression_blueprint.route('/tyre-impression/upload', methods=['POST'])
def upload():
    repo = TyreImpressionRepository()

    # Extract file from request
    file = request.files['file']

    valid_file_extensions = ["png", "jpg", "jpeg"]

    # Generate uuid and filename
    uuid, uuid_file_name = policies.uuid_filename(file)
    file.filename = uuid_file_name
    # Save file locally
    try:
        path = save_file(file, '/files/images/tyre_impressions', valid_file_extensions)
    except ValueError as e:
        return error_response(400, "File could not be saved: {}".format(e))

    # Create DB record
    tyre_impression = repo.create(
        uuid=uuid,
        file_path=path,
    )

    res = tyre_impression_response(tyre_impression)
    return jsonify(res), http.HTTPStatus.CREATED
