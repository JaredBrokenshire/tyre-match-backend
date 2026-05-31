import pytest
import numpy as np
from unittest.mock import MagicMock, patch
from services.file_service import FileService
from database.models.data_types.files import FileType
from pipelines.processors.base_processor import BaseProcessor
from domain.exceptions import FileSaveError, DatabaseError, FileReadError, ProcessorError


@pytest.fixture
def context():
    return {
        "output_directories": {"stage1": "/files/test_directory/stage1"},
        "processing_id": 1,
        "file_types_on_completion": {"stage1": FileType.normalised},
    }

@pytest.fixture
def processor():
    return BaseProcessor([])

@pytest.fixture
def mock_file_record():
    mock_file_record = MagicMock()
    mock_file_record.file_location = "/files/test_directory/test-file.jpg"
    return mock_file_record


def test_process_no_image_read(context, processor):
    with patch("cv2.imread", return_value=None):
        with pytest.raises(FileReadError, match="Unable to read image in base processor"):
            result = processor.process(input_path="/files/test_directory", context=context)
            assert result is None


def test_process_file_save_error_from_save_file(context, processor):
    with patch("cv2.imread", return_value=np.ndarray((100, 100))):
        with patch.object(processor, "transform", return_value=np.ndarray((100, 100))):
            with patch.object(processor, "save_file", side_effect=FileSaveError("test error")):
                with pytest.raises(ProcessorError, match="FileSaveError from base processor"):
                    result = processor.process(input_path="/files/test_directory", context=context)
                    assert result is None


def test_process(context, processor):
    with patch("cv2.imread", return_value=np.ndarray((100, 100))):
        with patch.object(processor, "transform", return_value=np.ndarray((100, 100))):
            with patch.object(processor, "save_file", return_value="/files/test_directory/normalised"):
                result = processor.process(input_path="/files/test_directory", context=context)

                assert result is not None
                assert result == "/files/test_directory/normalised"


def test_transform_processor_error_from_transform_step(processor, context):
    mock_transform_step = MagicMock()
    mock_transform_step.__name__ = "test_transform_step"
    mock_transform_step.side_effect = ProcessorError("test error")

    processor.transform_steps = [mock_transform_step]

    image = np.zeros((10,10))

    with pytest.raises(ProcessorError, match=r"ProcessorError from base processor \(test_transform_step\): test error"):
        processor.transform(image, context)


def test_transform(processor, context):
    mock_transform_step_1 = MagicMock()
    mock_transform_step_1.return_value = np.full((10,10), 1)
    mock_transform_step_2 = MagicMock()
    mock_transform_step_2.return_value = np.full((10,10), 2)

    processor.transform_steps = [mock_transform_step_1, mock_transform_step_2]

    image = np.zeros((10,10))

    result = processor.transform(image, context)

    # Ensure final step result is returned
    assert np.array_equal(result, mock_transform_step_2.return_value)
    # Ensure all steps were called with the correct parameters
    mock_transform_step_1.assert_called_once_with(image, context)
    mock_transform_step_2.assert_called_once_with(mock_transform_step_1.return_value, context)


def test_save_file_permission_error_from_file_service(context, processor, mock_file_record):
     with patch.object(FileService, "save_processed_image", side_effect=PermissionError("test error")):
        with pytest.raises(FileSaveError, match="PermissionError from file service in stage1 processor: test error"):
            result = processor.save_file(image="mock_image", context=context, stage_name="stage1")
            assert result is None


def test_save_file_os_error_from_file_service(context, processor, mock_file_record):
    with patch.object(FileService, "save_processed_image", side_effect=OSError("test error")):
        with pytest.raises(FileSaveError, match="OSError from file service in stage1 processor: test error"):
            result = processor.save_file(image="mock_image", context=context, stage_name="stage1")
            assert result is None


def test_save_file_file_save_error_from_file_service(context, processor, mock_file_record):
    with patch.object(FileService, "save_processed_image", side_effect=FileSaveError("test error")):
        with pytest.raises(FileSaveError, match="FileSaveError from file service in stage1 processor: test error"):
            result = processor.save_file(image="mock_image", context=context, stage_name="stage1")
            assert result is None


def test_save_file_database_error_from_file_service(context, processor, mock_file_record):
    with patch.object(FileService, "save_processed_image", side_effect=DatabaseError("test error")):
        with pytest.raises(FileSaveError, match="DatabaseError from file service in stage1 processor: test error"):
            result = processor.save_file(image="mock_image", context=context, stage_name="stage1")
            assert result is None


def test_save_file(context, processor, mock_file_record):
    with patch.object(FileService, "save_processed_image", return_value=mock_file_record):
        result = processor.save_file(image="mock_image", context=context, stage_name="stage1")

        assert result is not None
        assert result == mock_file_record.file_location

