import pytest
from unittest.mock import MagicMock, patch
from services.file_service import FileService
from database.models.data_types.files import FileType
from domain.exceptions import FileSaveError, DatabaseError
from pipelines.processors.base_processor import BaseProcessor


@pytest.fixture
def context():
    return {
        "output_directories": {"stage1": "/files/test_directory/stage1"},
        "processing_id": 1,
        "file_types_on_completion": {"stage1": FileType.normalised},
    }

@pytest.fixture
def processor():
    return BaseProcessor()

@pytest.fixture
def mock_file_record():
    mock_file_record = MagicMock()
    mock_file_record.file_location = "/files/test_directory/test-file.jpg"
    return mock_file_record


def test_base_processor_save_file_permission_error_from_file_service(context, processor, mock_file_record):
     with patch.object(FileService, "save_processed_image", side_effect=PermissionError("test error")):
        with pytest.raises(FileSaveError, match="PermissionError from file service in stage1 processor: test error"):
            result = processor.save_file(image="mock_image", context=context, stage_name="stage1")
            assert result is None


def test_base_processor_save_file_os_error_from_file_service(context, processor, mock_file_record):
    with patch.object(FileService, "save_processed_image", side_effect=OSError("test error")):
        with pytest.raises(FileSaveError, match="OSError from file service in stage1 processor: test error"):
            result = processor.save_file(image="mock_image", context=context, stage_name="stage1")
            assert result is None


def test_base_processor_save_file_file_save_error_from_file_service(context, processor, mock_file_record):
    with patch.object(FileService, "save_processed_image", side_effect=FileSaveError("test error")):
        with pytest.raises(FileSaveError, match="FileSaveError from file service in stage1 processor: test error"):
            result = processor.save_file(image="mock_image", context=context, stage_name="stage1")
            assert result is None


def test_base_processor_save_file_database_error_from_file_service(context, processor, mock_file_record):
    with patch.object(FileService, "save_processed_image", side_effect=DatabaseError("test error")):
        with pytest.raises(FileSaveError, match="DatabaseError from file service in stage1 processor: test error"):
            result = processor.save_file(image="mock_image", context=context, stage_name="stage1")
            assert result is None


def test_base_processor_save_file(context, processor, mock_file_record):
    with patch.object(FileService, "save_processed_image", return_value=mock_file_record):
        result = processor.save_file(image="mock_image", context=context, stage_name="stage1")

        assert result is not None
        assert result == mock_file_record.file_location

