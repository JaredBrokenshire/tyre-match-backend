import uuid
import pytest
import numpy as np
from unittest.mock import MagicMock, patch
from pipelines.processors.base_processor import BaseProcessor
from database.models.data_types.files import FileType, FileModel
from services.file_service import FileService, ProcessedImageRequest
from domain.exceptions import FileSaveError, DatabaseError, ProcessorError


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


def test_process_processor_error_from_transform(context, processor):
    image = MagicMock()

    with patch.object(processor, "transform", side_effect=ProcessorError("test error")):
        with pytest.raises(ProcessorError, match="ProcessorError when transforming image in base processor: test error"):
            processor.process(image, context)


def test_process_file_save_error_from_save_file(context, processor):
    image = MagicMock()

    with patch.object(processor, "transform", return_value=np.ndarray((100, 100))):
        with patch.object(processor, "save_file", side_effect=FileSaveError("test error")):
            with pytest.raises(ProcessorError, match="FileSaveError from base processor"):
                result = processor.process(image=image, context=context)
                assert result is None


def test_process(context, processor):
    image = np.zeros((10,10))

    with patch.object(processor, "transform", return_value=image):
        with patch.object(processor, "save_file"):
            result = processor.process(image=image, context=context)

            assert np.array_equal(result, image)


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


def test_save_file_permission_error_from_file_service(context, processor):
     with patch.object(FileService, "save_processed_image", side_effect=PermissionError("test error")):
        with pytest.raises(FileSaveError, match="PermissionError from file service in stage1 processor: test error"):
            result = processor.save_file(image="mock_image", context=context, stage_name="stage1")
            assert result is None


def test_save_file_os_error_from_file_service(context, processor):
    with patch.object(FileService, "save_processed_image", side_effect=OSError("test error")):
        with pytest.raises(FileSaveError, match="OSError from file service in stage1 processor: test error"):
            result = processor.save_file(image="mock_image", context=context, stage_name="stage1")
            assert result is None


def test_save_file_file_save_error_from_file_service(context, processor):
    with patch.object(FileService, "save_processed_image", side_effect=FileSaveError("test error")):
        with pytest.raises(FileSaveError, match="FileSaveError from file service in stage1 processor: test error"):
            result = processor.save_file(image="mock_image", context=context, stage_name="stage1")
            assert result is None


def test_save_file_database_error_from_file_service(context, processor):
    with patch.object(FileService, "save_processed_image", side_effect=DatabaseError("test error")):
        with pytest.raises(FileSaveError, match="DatabaseError from file service in stage1 processor: test error"):
            result = processor.save_file(image="mock_image", context=context, stage_name="stage1")
            assert result is None


def test_save_file(context, processor):
    mock_file_service = MagicMock()

    with patch.object(FileService, "save_processed_image", new=mock_file_service.save_processed_image):
        processor.save_file(image="mock_image", context=context, stage_name="stage1")

    mock_file_service.save_processed_image.assert_called_once()

    request = mock_file_service.save_processed_image.call_args.args[0]

    assert isinstance(request, ProcessedImageRequest)
    assert request.image == "mock_image"
    assert request.upload_directory == context["output_directories"]["stage1"]
    assert request.model == FileModel.tyre_impression
    assert request.model_id == context["processing_id"]
    assert request.file_type == context["file_types_on_completion"]["stage1"]
    assert request.extension == "png"

    # Validate generated filename
    assert request.file_name.endswith(".png")

    # Validate UUID filename
    uuid_part = request.file_name.removesuffix(".png")
    uuid.UUID(uuid_part)



