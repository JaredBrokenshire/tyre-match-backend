import pytest
from unittest.mock import patch
from services.file_service import FileService
from werkzeug.datastructures import FileStorage
from database.models.data_types.files import FileModel
from services.tyre_model_service import TyreModelService
from tests.helpers.factories.tyre_model_factory import TyreModelFactory
from database.repositories.tyre_model_repository import TyreModelRepository
from domain.exceptions import ModelNotFoundError, DatabaseError, InvalidFileTypeError, InvalidFileError, FileSaveError


@pytest.fixture()
def service():
    return TyreModelService()


def test_get_all(service):
    tyre_model_1 = TyreModelFactory().create()
    tyre_model_2 = TyreModelFactory().create()
    tyre_model_3 = TyreModelFactory().create()

    results, total_count = service.get_all()

    assert 3 == len(results) == total_count
    assert tyre_model_1 in results
    assert tyre_model_2 in results
    assert tyre_model_3 in results


def test_get_all_pagination_page_1(service):
    tyre_model_1 = TyreModelFactory().create()
    tyre_model_2 = TyreModelFactory().create()
    tyre_model_3 = TyreModelFactory().create()

    results, total_count = service.get_all(page=1, page_size=1)

    assert 3 == total_count
    assert 1 == len(results)
    assert tyre_model_1 == results[0]
    assert tyre_model_2 not in results
    assert tyre_model_3 not in results


def test_get_all_pagination_page_2(service):
    tyre_model_1 = TyreModelFactory().create()
    tyre_model_2 = TyreModelFactory().create()
    tyre_model_3 = TyreModelFactory().create()

    results, total_count = service.get_all(page=2, page_size=1)

    assert 3 == total_count
    assert 1 == len(results)
    assert tyre_model_2 == results[0]
    assert tyre_model_1 not in results
    assert tyre_model_3 not in results


def test_get_all_search_by_manufacturer(service):
    tyre_model_1 = TyreModelFactory().create(manufacturer='Test Manufacturer')
    tyre_model_2 = TyreModelFactory().create()
    tyre_model_3 = TyreModelFactory().create()

    results, total_count = service.get_all(search="Test")

    assert 1 == len(results) == total_count
    assert tyre_model_1 == results[0]
    assert tyre_model_2 not in results
    assert tyre_model_3 not in results


def test_get_all_search_by_model_name(service):
    tyre_model_1 = TyreModelFactory().create(model_name='Test Model Name')
    tyre_model_2 = TyreModelFactory().create()
    tyre_model_3 = TyreModelFactory().create()

    results, total_count = service.get_all(search="Test")

    assert 1 == len(results) == total_count
    assert tyre_model_1 == results[0]
    assert tyre_model_2 not in results
    assert tyre_model_3 not in results


def test_get_by_id_invalid_id(service):
    with pytest.raises(ModelNotFoundError, match="TyreModel with id 1000 not found"):
        service.get_by_id(1000)


def test_get_by_id(service):
    tyre_model_1 = TyreModelFactory().create()
    tyre_model_2 = TyreModelFactory().create()

    result = service.get_by_id(tyre_model_1.id)

    assert tyre_model_1 == result
    assert tyre_model_2 != result


def test_create_database_error_from_tyre_model_repository(service):
    with patch.object(TyreModelRepository, 'create', side_effect=DatabaseError("test error")):
        with pytest.raises(DatabaseError, match="Error creating tyre model record: test error"):
            service.create({})


def test_create(service):
    dto = {
        "manufacturer": "Test Manufacturer",
        "model_name": "Test Model Name",
    }

    result = service.create(dto)
    assert result is not None
    assert result.id != 0
    assert result.manufacturer == "Test Manufacturer"
    assert result.model_name == "Test Model Name"


def test_upload_image_no_file(service):
    tyre_model = TyreModelFactory().create()

    with pytest.raises(InvalidFileTypeError, match="No file provided"):
        service.upload_image(tyre_model, None)


def test_upload_image_invalid_file_error_from_file_service(service):
    tyre_model = TyreModelFactory().create()

    file = FileStorage(filename="test-file.jpg")

    with patch.object(FileService, "handle_file", side_effect=InvalidFileError("test error")):
        with pytest.raises(FileSaveError, match="Invalid file error from file service in tyre model service: test error"):
            service.upload_image(tyre_model, file)


def test_upload_image_invalid_file_type_error_from_file_service(service):
    tyre_model = TyreModelFactory().create()

    file = FileStorage(filename="test-file.jpg")

    with patch.object(FileService, "handle_file", side_effect=InvalidFileTypeError("test error")):
        with pytest.raises(FileSaveError, match="Invalid file type error from file service in tyre model service: test error"):
            service.upload_image(tyre_model, file)


def test_upload_image_permission_error_from_file_service(service):
    tyre_model = TyreModelFactory().create()

    file = FileStorage(filename="test-file.jpg")

    with patch.object(FileService, "handle_file", side_effect=PermissionError("test error")):
        with pytest.raises(FileSaveError, match="Permission or OS error from file service in tyre model service: test error"):
            service.upload_image(tyre_model, file)


def test_upload_image_os_error_from_file_service(service):
    tyre_model = TyreModelFactory().create()

    file = FileStorage(filename="test-file.jpg")

    with patch.object(FileService, "handle_file", side_effect=OSError("test error")):
        with pytest.raises(FileSaveError, match="Permission or OS error from file service in tyre model service: test error"):
            service.upload_image(tyre_model, file)


def test_upload_image_database_error_from_file_service(service):
    tyre_model = TyreModelFactory().create()

    file = FileStorage(filename="test-file.jpg")

    with patch.object(FileService, "handle_file", side_effect=DatabaseError("test error")):
        with pytest.raises(FileSaveError, match="Database error from file service in tyre model service: test error"):
            service.upload_image(tyre_model, file)


def test_upload_image(service):
    tyre_model = TyreModelFactory().create()

    file = FileStorage(filename="test-file.jpg")

    res = service.upload_image(tyre_model, file)

    assert res is not None
    assert res.id == tyre_model.id
    assert 1 == len(res.files)
    assert res.files[0].model == FileModel.tyre_model
    assert res.files[0].model_id == tyre_model.id


def test_update_invalid_id(service):
    with pytest.raises(ModelNotFoundError, match="TyreModel with id 1000 not found"):
        service.get_by_id(1000)


def test_update_database_error_from_tyre_model_repository(service):
    tyre_model = TyreModelFactory().create()

    with patch.object(TyreModelRepository, 'update', side_effect=DatabaseError("test error")):
        with pytest.raises(DatabaseError, match="Error updating tyre model record: test error"):
            service.update(tyre_model.id, {})


def test_update(service):
    tyre_model = TyreModelFactory().create()

    updated_model = service.update(
        tyre_model.id,
        {
            "manufacturer": "Test Manufacturer",
            "model_name": "Test Model Name",
        }
    )

    assert tyre_model.id == updated_model.id
    assert "Test Manufacturer" == updated_model.manufacturer
    assert "Test Model Name" == updated_model.model_name


def test_delete_invalid_id(service):
    service = TyreModelService()

    with pytest.raises(ModelNotFoundError, match="Tyre model with id 1 not found"):
        result = service.delete(1)
        assert result == False


def test_delete_database_error_from_tyre_model_repository(service):
    with patch.object(TyreModelRepository, 'delete', side_effect=DatabaseError("test error")):
        with pytest.raises(DatabaseError, match="Error deleting tyre model record: test error"):
            result = service.delete(1)
            assert result == False


def test_delete(service):
    tyre_model = TyreModelFactory().create()

    result = service.delete(tyre_model.id)

    assert result == True
    with pytest.raises(ModelNotFoundError, match=f"TyreModel with id {tyre_model.id} not found"):
        service.get_by_id(tyre_model.id)
