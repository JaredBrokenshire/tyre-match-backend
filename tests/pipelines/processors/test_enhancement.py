import cv2
import pytest
import numpy as np
from unittest.mock import patch, MagicMock
from database.models.data_types.files import FileType
from domain.exceptions import FileReadError, FileSaveError, ProcessorError
from pipelines.processors.enhancement_processor import EnhancementProcessor


@pytest.fixture
def processor():
    return EnhancementProcessor()


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
        with patch.object(EnhancementProcessor, "transform", return_value=np.ndarray((100, 100))):
            with patch.object(EnhancementProcessor, "save_file", side_effect=FileSaveError("test error")):
                with pytest.raises(ProcessorError, match="FileSaveError from enhancement processor"):
                    result = processor.process(input_path="/files/test_directory", context=context)
                    assert result is None


def test_process(context, processor):
    with patch("cv2.imread", return_value=np.ndarray((100, 100))):
        with patch.object(EnhancementProcessor, "transform", return_value=np.ndarray((100, 100))):
            with patch.object(EnhancementProcessor, "save_file", return_value="/files/test_directory/normalised"):
                result = processor.process(input_path="/files/test_directory", context=context)

                assert result is not None
                assert result == "/files/test_directory/normalised"


def test_transform_exception_from_clahe_apply(processor, context):
    image = np.zeros((10, 10), np.uint8)

    mock_clahe = MagicMock()
    mock_clahe.apply.side_effect = Exception("test error")

    with patch("cv2.createCLAHE", return_value=mock_clahe):
        with pytest.raises(ProcessorError, match="Exception from apply CLAHE in transform: test error"):
            processor.transform(image, context)


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


def test_transform_is_deterministic(context, processor):
    image = np.random.randint(0, 256, (100, 100), dtype=np.uint8)

    result_1 = processor.transform(image, context)
    result_2 = processor.transform(image, context)

    assert np.array_equal(result_1, result_2)


def test_transform_changes_with_tile_grid_size(processor):
    image = np.full((200, 200), 120, dtype=np.uint8)

    small_tile_grid_size_context = {"clahe_clip_limit": 2.0, "clahe_tile_grid_size": (4, 4)}
    large_tile_grid_size_context = {"clahe_clip_limit": 2.0, "clahe_tile_grid_size": (16, 16)}

    small_tile_grid_size_result = processor.transform(image, small_tile_grid_size_context)
    large_tile_grid_size_result = processor.transform(image, large_tile_grid_size_context)

    assert not np.array_equal(small_tile_grid_size_result, large_tile_grid_size_result)


def test_transform_changes_with_clip_limit(processor):
    image = np.full((200, 200), 120, dtype=np.uint8)

    low_clip_limit_context = {"clahe_clip_limit": 2.0, "clahe_tile_grid_size": (4, 4)}
    high_clip_limit_context = {"clahe_clip_limit": 4.0, "clahe_tile_grid_size": (4, 4)}

    low_clip_limit_result = processor.transform(image, low_clip_limit_context)
    high_clip_limit_result = processor.transform(image, high_clip_limit_context)

    assert not np.array_equal(low_clip_limit_result, high_clip_limit_result)


def test_transform_contrast_preserved_in_flat_image(processor, context):
    image = np.full((100, 100), 100, dtype=np.uint8)

    result = processor.transform(image, context)

    assert 0.0 == result.std() == image.std()


def test_transform_higher_clip_limit_increases_contrast_in_noisy_image(processor):
    np.random.seed(0)

    noise = np.random.randint(0, 200, (100,100), dtype=np.uint8)
    image = np.full((100, 100), 55, dtype=np.uint8)
    image = np.clip(image + noise, 0, 255).astype(np.uint8)

    low_clip_limit_context = {"clahe_clip_limit": 1.0, "clahe_tile_grid_size": (4, 4)}
    high_clip_limit_context = {"clahe_clip_limit": 10.0, "clahe_tile_grid_size": (4, 4)}

    high_clip_limit_result = processor.transform(image, high_clip_limit_context)
    low_clip_limit_result = processor.transform(image, low_clip_limit_context)

    high_contrast = cv2.Laplacian(high_clip_limit_result, cv2.CV_64F).var()
    low_contrast = cv2.Laplacian(low_clip_limit_result, cv2.CV_64F).var()

    assert high_contrast > low_contrast


def test_transform_keeps_valid_pixel_range(context, processor):
    image = np.random.randint(0, 256, (100, 100), dtype=np.uint8)

    result = processor.transform(image, context)

    assert result.min() >= 0
    assert result.max() <= 255


def test_transform_output_shape_and_data_type_matches_input(context, processor):
    image = np.random.randint(0, 256, (100,100), dtype=np.uint8)

    result = processor.transform(image, context)

    assert isinstance(result, np.ndarray)
    assert result.shape == image.shape
    assert result.dtype == image.dtype


def test_transform_preserves_structural_edges(context, processor):
    image = np.zeros((100, 100), dtype=np.uint8)
    image[:, 50:] = 255  # sharp vertical edge

    result = processor.transform(image, context)

    input_edge = np.diff(image.mean(axis=0))
    output_edge = np.diff(result.mean(axis=0))

    assert np.argmax(input_edge) == np.argmax(output_edge)


def test_transform_changes_histogram(context, processor):
    image = np.zeros((100,100), dtype=np.uint8)
    image[10:50, 10:50] = 128

    result = processor.transform(image, context)

    # Calculate histograms for input and output image
    input_image_histogram = cv2.calcHist([image], [0], None, [256], [0, 256])
    result_histogram = cv2.calcHist([result], [0], None, [256], [0, 256])

    # Calculate Euclidean Distance between histograms
    assert np.linalg.norm(input_image_histogram - result_histogram) > 0


def test_transform_stability_on_reapplication(context, processor):
    image = np.random.randint(0, 256, (100, 100), dtype=np.uint8)

    once = processor.transform(image, context)
    twice = processor.transform(once, context)

    # Not identical, but difference should be small
    diff = np.mean(np.abs(once.astype(np.int16) - twice.astype(np.int16)))
    assert diff < 1


def test_transform_extreme_parameters(processor):
    image = np.random.randint(0, 256, (10000, 10000), dtype=np.uint8)

    context = {
        "clahe_clip_limit": 100.0,
        "clahe_tile_grid_size": (1, 1),
    }

    result = processor.transform(image, context)

    assert result.shape == image.shape
    assert result.dtype == np.uint8
    assert result.min() >= 0
    assert result.max() <= 255


def test_transform_requires_grayscale_image(processor, context):
    colour_image = np.random.randint(0, 256, (10, 10, 3), dtype=np.uint8)

    with pytest.raises(ProcessorError, match=r"Exception from apply CLAHE in transform:.*"):
        processor.transform(colour_image, context)