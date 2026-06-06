import http
from pathlib import Path
from database.models.data_types.files import FileModel
from tests.helpers.assertions import assert_error_response
from tests.helpers.factories.file_factory import FileFactory


def test_get_by_id_invalid_id(client):
    response = client.get("/files/1000")

    # Ensure correct status code and error message were returned
    assert_error_response(response, http.HTTPStatus.NOT_FOUND, "File with id 1000 not found")


def test_get_by_id(client, database_session):
    file_1 = FileFactory.create(FileModel.tyre_model, 1, True)
    file_2 = FileFactory.create(FileModel.tyre_model, 1, True)

    response = client.get(f"/files/{file_1.id}")

    assert http.HTTPStatus.OK == response.status_code
    expected = (
            Path(file_1.file_location)
            / file_1.file_name
    ).read_bytes()

    not_expected = (
        Path(file_2.file_location)
        / file_2.file_name
    ).read_bytes()

    assert response.data == expected
    assert response.data != not_expected
