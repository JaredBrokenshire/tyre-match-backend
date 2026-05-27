import cv2
import pytest
import numpy as np
from unittest.mock import patch, MagicMock
from database.models.data_types.files import FileType
from pipelines.processors.normalisation import NormalisationProcessor
from domain.exceptions import FileReadError, FileSaveError, ProcessorError


@pytest.fixture
def processor():
    return NormalisationProcessor()

@pytest.fixture
def context():
    return {
        "clahe_clip_limit": 2.5,
        "clahe_tile_grid_size": (4,4),
        "output_directories": {"normalisation", "/files/test_directory/normalised"},
        "processing_id": 1,
        "file_types_on_completion": {"normalisation": FileType.normalised},
    }


def test_process_no_image_read(context, processor):
    with patch("cv2.imread", return_value=None):
        with pytest.raises(FileReadError, match="Unable to read image in NormalisationProcessor"):
            result = processor.process(input_path="/files/test_directory", context=context)
            assert result is None


def test_process_file_save_error_from_save_file(context, processor):
    with patch("cv2.imread", return_value=np.ndarray((100, 100))):
        with patch.object(NormalisationProcessor, "transform", return_value=np.ndarray((100, 100))):
            with patch.object(NormalisationProcessor, "save_file", side_effect=FileSaveError("test error")):
                with pytest.raises(ProcessorError, match="FileSaveError from normalisation processor"):
                    result = processor.process(input_path="/files/test_directory", context=context)
                    assert result is None


def test_process(context, processor):
    with patch("cv2.imread", return_value=np.ndarray((100, 100))):
        with patch.object(NormalisationProcessor, "transform", return_value=np.ndarray((100, 100))):
            with patch.object(NormalisationProcessor, "save_file", return_value="/files/test_directory/normalised"):
                result = processor.process(input_path="/files/test_directory", context=context)

                assert result is not None
                assert result == "/files/test_directory/normalised"


def test_transform_calls_clahe_with_correct_params(processor):
    image = np.zeros((10, 10), dtype=np.uint8)

    context = {
        "clahe_clip_limit": 3.5,
        "clahe_tile_grid_size": (4, 4),
    }

    mock_clahe = MagicMock()
    mock_clahe.apply.return_value = image

    with patch("cv2.createCLAHE", return_value=mock_clahe) as mock_create:
        processor.transform(image, context)

        mock_create.assert_called_once_with(
            clipLimit=3.5,
            tileGridSize=(4, 4),
        )
        mock_clahe.apply.assert_called_once()


def test_transform_output_shape_and_data_type_matches_input(context, processor):
    image = np.random.randint(0, 256, (100,100), dtype=np.uint8)

    result = processor.transform(image, context)

    assert isinstance(result, np.ndarray)
    assert result.shape == image.shape
    assert result.dtype == image.dtype


def test_transform_changes_contrast(context, processor):
    # Low contrast image, all values similar
    image = np.full((100,100), 120, dtype=np.uint8)

    result = processor.transform(image, context)

    # CLAHE should introduce variation
    assert not np.array_equal(result, image)
    print(">>> ", image, result)


def test_transform_changes_histogram(context, processor):
    image = np.zeros((100,100), dtype=np.uint8)
    image[10:50, 10:50] = 128

    result = processor.transform(image, context)

    # Calculate histograms for input and output image
    input_image_histogram = cv2.calcHist([image], [0], None, [256], [0, 256])
    result_histogram = cv2.calcHist([result], [0], None, [256], [0, 256])

    # Calculate Euclidean Distance between histograms
    assert np.linalg.norm(input_image_histogram - result_histogram) > 0


def test_transform_extreme_parameters(processor):
    image = np.random.randint(0, 256, (1000, 1000), dtype=np.uint8)

    context = {
        "clahe_clip_limit": 100.0,
        "clahe_tile_grid_size": (1, 1),
    }

    result = processor.transform(image, context)

    assert result.shape == image.shape
    assert result.dtype == np.uint8
    assert result.min() >= 0
    assert result.max() <= 255