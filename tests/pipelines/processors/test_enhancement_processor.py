import cv2
import pytest
import numpy as np
from unittest.mock import patch, MagicMock
from domain.exceptions import ProcessorError
from database.models.data_types.files import FileType
from pipelines.processors.enhancement_processor import EnhancementProcessor


@pytest.fixture
def processor():
    return EnhancementProcessor()


@pytest.fixture
def context():
    return {
        "clahe_clip_limit": 2.5,
        "clahe_tile_grid_size": (4,4),
        "output_directories": {"enhancement", "/files/test_directory/enhanced"},
        "processing_id": 1,
        "file_types_on_completion": {"enhancement": FileType.enhanced},
    }

@pytest.fixture
def noisy_image_with_grid():
    image = np.zeros((200, 200), dtype=np.uint8)
    image[50:150, 50:150] = 128
    image += np.random.randint(0, 20, (200, 200), dtype=np.uint8)
    return image


def test_denoise_exception_from_fast_nl_means_denoising(processor, context):
    image = np.zeros((10, 10), np.uint8)

    with patch("cv2.fastNlMeansDenoising", side_effect=Exception("test error")):
        with pytest.raises(ProcessorError, match="Error when applying denoise: test error"):
            processor._denoise(image, context)


def test_denoise_calls_opencv_with_correct_params(processor):
    image = np.zeros((10, 10), dtype=np.uint8)

    context = {
        "denoise_h": 15,
        "denoise_template_window_size": 7,
        "denoise_search_window_size": 21,
    }

    with patch("cv2.fastNlMeansDenoising", return_value=image) as mock_fn:
        processor._denoise(image, context)

        mock_fn.assert_called_once_with(
            image,
            None,
            15,
            7,
            21
        )


def test_denoise_is_deterministic(processor, context):
    image = np.random.randint(0, 256, (100, 100), dtype=np.uint8)

    result_1 = processor._denoise(image, context)
    result_2 = processor._denoise(image, context)

    assert np.array_equal(result_1, result_2)


def test_denoise_reduces_noise(processor, noisy_image_with_grid):
    np.random.seed(0)

    context = {
        "denoise_h": 10,
        "denoise_template_window_size": 7,
        "denoise_search_window_size": 21,
    }

    denoised = processor._denoise(noisy_image_with_grid, context)

    noisy_var = np.var(noisy_image_with_grid)
    denoised_var = np.var(denoised)

    assert denoised_var < noisy_var


def test_denoise_preserves_shape_and_dtype(processor, context):
    image = np.random.randint(0, 256, (200, 300), dtype=np.uint8)

    result = processor._denoise(image, context)

    assert isinstance(result, np.ndarray)
    assert result.shape == image.shape
    assert result.dtype == image.dtype


def test_denoise_does_not_destroy_strong_edges(processor, context):
    image = np.zeros((100, 100), dtype=np.uint8)
    image[:, 50:] = 255  # sharp vertical edge

    result = processor._denoise(image, context)

    input_edge = np.diff(image.mean(axis=0))
    output_edge = np.diff(result.mean(axis=0))

    # edge location should remain stable
    assert np.argmax(input_edge) == np.argmax(output_edge)


def test_denoise_changes_with_h(processor, noisy_image_with_grid):
    low_h_context = {"denoise_h": 3}
    high_h_context = {"denoise_h": 20}

    low_h_result = processor._denoise(noisy_image_with_grid, low_h_context)
    high_h_result = processor._denoise(noisy_image_with_grid, high_h_context)

    assert not np.array_equal(low_h_result, high_h_result)


def test_denoise_changes_with_template_window_size(processor, noisy_image_with_grid):
    low_template_window_size_context = {"denoise_template_window_size": 3}
    high_template_window_size_context = {"denoise_template_window_size": 20}

    low_template_window_size_result = processor._denoise(noisy_image_with_grid, low_template_window_size_context)
    high_template_window_size_result = processor._denoise(noisy_image_with_grid, high_template_window_size_context)

    assert not np.array_equal(low_template_window_size_result, high_template_window_size_result)


def test_denoise_changes_with_search_window_size(processor, noisy_image_with_grid):
    low_search_window_size_context = {"denoise_search_window_size": 3}
    high_search_window_size_context = {"denoise_search_window_size": 20}

    low_search_window_size_result = processor._denoise(noisy_image_with_grid, low_search_window_size_context)
    high_search_window_size_result = processor._denoise(noisy_image_with_grid, high_search_window_size_context)

    assert not np.array_equal(low_search_window_size_result, high_search_window_size_result)


def test_denoise_keeps_valid_pixel_range(processor, context):
    image = np.random.randint(0, 256, (100, 100), dtype=np.uint8)

    result = processor._denoise(image, context)

    assert result.min() >= 0
    assert result.max() <= 255


def test_denoise_reduces_high_frequency_noise_but_preserves_structure(processor, noisy_image_with_grid):
    context = {
        "denoise_h": 10,
        "denoise_template_window_size": 7,
        "denoise_search_window_size": 21,
    }

    result = processor._denoise(noisy_image_with_grid, context)

    # structural signal should remain detectable
    input_energy = np.sum(np.abs(np.diff(noisy_image_with_grid.astype(np.int16))))
    output_energy = np.sum(np.abs(np.diff(result.astype(np.int16))))

    assert output_energy < input_energy


def test_denoise_stability_on_reapplication(processor, context):
    image = np.random.randint(0, 256, (200, 200), dtype=np.uint8)

    once = processor._denoise(image, context)
    twice = processor._denoise(once, context)

    diff = np.mean(np.abs(once.astype(np.int16) - twice.astype(np.int16)))

    # sharpening twice should not explode differences uncontrollably
    assert diff < 1


def test_apply_clahe_exception_from_clahe_apply(processor, context):
    image = np.zeros((10, 10), np.uint8)

    mock_clahe = MagicMock()
    mock_clahe.apply.side_effect = Exception("test error")

    with patch("cv2.createCLAHE", return_value=mock_clahe):
        with pytest.raises(ProcessorError, match="Error when applying CLAHE: test error"):
            processor._apply_clahe(image, context)


def test_apply_clahe_calls_clahe_with_correct_params(processor):
    image = np.zeros((10, 10), dtype=np.uint8)

    context = {
        "clahe_clip_limit": 3.5,
        "clahe_tile_grid_size": (4, 4),
    }

    mock_clahe = MagicMock()
    mock_clahe.apply.return_value = image

    with patch("cv2.createCLAHE", return_value=mock_clahe) as mock_create:
        processor._apply_clahe(image, context)

        mock_create.assert_called_once_with(
            clipLimit=3.5,
            tileGridSize=(4, 4),
        )
        mock_clahe.apply.assert_called_once()


def test_apply_clahe_is_deterministic(context, processor):
    image = np.random.randint(0, 256, (100, 100), dtype=np.uint8)

    result_1 = processor._apply_clahe(image, context)
    result_2 = processor._apply_clahe(image, context)

    assert np.array_equal(result_1, result_2)


def test_apply_clahe_changes_with_tile_grid_size(processor):
    image = np.full((200, 200), 120, dtype=np.uint8)

    small_tile_grid_size_context = {"clahe_clip_limit": 2.0, "clahe_tile_grid_size": (4, 4)}
    large_tile_grid_size_context = {"clahe_clip_limit": 2.0, "clahe_tile_grid_size": (16, 16)}

    small_tile_grid_size_result = processor._apply_clahe(image, small_tile_grid_size_context)
    large_tile_grid_size_result = processor._apply_clahe(image, large_tile_grid_size_context)

    assert not np.array_equal(small_tile_grid_size_result, large_tile_grid_size_result)


def test_apply_clahe_changes_with_clip_limit(processor):
    image = np.full((200, 200), 120, dtype=np.uint8)

    low_clip_limit_context = {"clahe_clip_limit": 2.0, "clahe_tile_grid_size": (4, 4)}
    high_clip_limit_context = {"clahe_clip_limit": 4.0, "clahe_tile_grid_size": (4, 4)}

    low_clip_limit_result = processor._apply_clahe(image, low_clip_limit_context)
    high_clip_limit_result = processor._apply_clahe(image, high_clip_limit_context)

    assert not np.array_equal(low_clip_limit_result, high_clip_limit_result)


def test_apply_clahe_contrast_preserved_in_flat_image(processor, context):
    image = np.full((100, 100), 100, dtype=np.uint8)

    result = processor._apply_clahe(image, context)

    assert 0.0 == result.std() == image.std()


def test_apply_clahe_higher_clip_limit_increases_contrast_in_noisy_image(processor):
    np.random.seed(0)

    noise = np.random.randint(0, 200, (100,100), dtype=np.uint8)
    image = np.full((100, 100), 55, dtype=np.uint8)
    image = np.clip(image + noise, 0, 255).astype(np.uint8)

    low_clip_limit_context = {"clahe_clip_limit": 1.0, "clahe_tile_grid_size": (4, 4)}
    high_clip_limit_context = {"clahe_clip_limit": 10.0, "clahe_tile_grid_size": (4, 4)}

    high_clip_limit_result = processor._apply_clahe(image, high_clip_limit_context)
    low_clip_limit_result = processor._apply_clahe(image, low_clip_limit_context)

    high_contrast = cv2.Laplacian(high_clip_limit_result, cv2.CV_64F).var()
    low_contrast = cv2.Laplacian(low_clip_limit_result, cv2.CV_64F).var()

    assert high_contrast > low_contrast


def test_apply_clahe_keeps_valid_pixel_range(context, processor):
    image = np.random.randint(0, 256, (100, 100), dtype=np.uint8)

    result = processor._apply_clahe(image, context)

    assert result.min() >= 0
    assert result.max() <= 255


def test_apply_clahe_output_shape_and_data_type_matches_input(context, processor):
    image = np.random.randint(0, 256, (100,100), dtype=np.uint8)

    result = processor._apply_clahe(image, context)

    assert isinstance(result, np.ndarray)
    assert result.shape == image.shape
    assert result.dtype == image.dtype


def test_apply_clahe_preserves_structural_edges(context, processor):
    image = np.zeros((100, 100), dtype=np.uint8)
    image[:, 50:] = 255  # sharp vertical edge

    result = processor._apply_clahe(image, context)

    input_edge = np.diff(image.mean(axis=0))
    output_edge = np.diff(result.mean(axis=0))

    assert np.argmax(input_edge) == np.argmax(output_edge)


def test_apply_clahe_changes_histogram(context, processor):
    image = np.zeros((100,100), dtype=np.uint8)
    image[10:50, 10:50] = 128

    result = processor._apply_clahe(image, context)

    # Calculate histograms for input and output image
    input_image_histogram = cv2.calcHist([image], [0], None, [256], [0, 256])
    result_histogram = cv2.calcHist([result], [0], None, [256], [0, 256])

    # Calculate Euclidean Distance between histograms
    assert np.linalg.norm(input_image_histogram - result_histogram) > 0


def test_apply_clahe_stability_on_reapplication(context, processor):
    image = np.random.randint(0, 256, (100, 100), dtype=np.uint8)

    once = processor._apply_clahe(image, context)
    twice = processor._apply_clahe(once, context)

    # Not identical, but difference should be small
    diff = np.mean(np.abs(once.astype(np.int16) - twice.astype(np.int16)))
    assert diff < 1


def test_apply_clahe_extreme_parameters(processor):
    image = np.random.randint(0, 256, (10000, 10000), dtype=np.uint8)

    context = {
        "clahe_clip_limit": 100.0,
        "clahe_tile_grid_size": (1, 1),
    }

    result = processor._apply_clahe(image, context)

    assert result.shape == image.shape
    assert result.dtype == np.uint8
    assert result.min() >= 0
    assert result.max() <= 255


def test_apply_clahe_requires_grayscale_image(processor, context):
    colour_image = np.random.randint(0, 256, (10, 10, 3), dtype=np.uint8)

    with pytest.raises(ProcessorError, match=r"Error when applying CLAHE:.*"):
        processor._apply_clahe(colour_image, context)


def test_sharpen_exception_from_gaussian_blur(processor, context):
    image = np.zeros((10, 10), np.uint8)

    with patch("cv2.GaussianBlur", side_effect=Exception("test error")):
        with pytest.raises(ProcessorError, match="Error when blurring during sharpen: test error"):
            processor._sharpen(image, context)


def test_sharpen_exception_from_add_weighted(processor, context):
    image = np.zeros((10, 10), np.uint8)

    with patch("cv2.addWeighted", side_effect=Exception("test error")):
        with pytest.raises(ProcessorError, match="Error when applying sharpen: test error"):
            processor._sharpen(image, context)


def test_sharpen_calls_gaussian_blur_with_correct_params(processor):
    image = np.zeros((10, 10), dtype=np.uint8)

    context = {
        "sharpen_blur_kernel_size": (5, 5),
        "sharpen_sigma": 2.0,
    }

    with patch("cv2.GaussianBlur", return_value=image) as mock_blur:
        processor._sharpen(image, context)

        mock_blur.assert_called_once_with(
            image,
            (5, 5),
            2.0
        )


def test_sharpen_is_deterministic(processor, context):
    image = np.random.randint(0, 256, (100, 100), dtype=np.uint8)

    result_1 = processor._sharpen(image, context)
    result_2 = processor._sharpen(image, context)

    assert np.array_equal(result_1, result_2)


def test_sharpen_changes_with_strength(processor, noisy_image_with_grid):
    low_strength_context = {
        "sharpen_strength": 0.2,
    }

    high_strength_context = {
        "sharpen_strength": 2.0,
    }

    low_strength_result = processor._sharpen(noisy_image_with_grid, low_strength_context)
    high_strength_result = processor._sharpen(noisy_image_with_grid, high_strength_context)

    assert not np.array_equal(low_strength_result, high_strength_result)


def test_sharpen_changes_with_sigma(processor, noisy_image_with_grid):
    low_sigma_context = {
        "sharpen_sigma": 0.5,
    }

    high_sigma_context = {
        "sharpen_sigma": 3.0,
    }

    low_sigma_result = processor._sharpen(noisy_image_with_grid, low_sigma_context)
    high_sigma_result = processor._sharpen(noisy_image_with_grid, high_sigma_context)

    assert not np.array_equal(low_sigma_result, high_sigma_result)


def test_sharpen_changes_with_kernel_size(processor, noisy_image_with_grid):
    small_kernel_context = {
        "sharpen_strength": 1.2,
        "sharpen_blur_kernel_size": (3, 3),
        "sharpen_sigma": 1.0,
    }

    large_kernel_context = {
        "sharpen_strength": 1.2,
        "sharpen_blur_kernel_size": (9, 9),
        "sharpen_sigma": 1.0,
    }

    small_kernel_result = processor._sharpen(noisy_image_with_grid, small_kernel_context)
    large_kernel_result = processor._sharpen(noisy_image_with_grid, large_kernel_context)

    assert not np.array_equal(small_kernel_result, large_kernel_result)


def test_sharpen_preserves_shape_and_dtype(processor, context):
    image = np.random.randint(0, 256, (120, 180), dtype=np.uint8)

    result = processor._sharpen(image, context)

    assert isinstance(result, np.ndarray)
    assert result.shape == image.shape
    assert result.dtype == image.dtype


def test_sharpen_keeps_valid_pixel_range(processor, context):
    image = np.random.randint(0, 256, (100, 100), dtype=np.uint8)

    result = processor._sharpen(image, context)

    assert result.min() >= 0
    assert result.max() <= 255


def test_sharpen_increases_edge_contrast(processor, context):
    image = np.zeros((100, 100), dtype=np.uint8)
    image[:, 50:] = 255  # strong vertical edge

    result = processor._sharpen(image, context)

    input_edge_strength = np.var(np.diff(image.mean(axis=0)))
    output_edge_strength = np.var(np.diff(result.mean(axis=0)))

    assert output_edge_strength >= input_edge_strength


def test_sharpen_reduces_local_blur_effect(processor, context):
    image = np.full((100, 100), 120, dtype=np.uint8)
    image[40:60, 40:60] = 200  # soft square signal

    blurred = cv2.GaussianBlur(image, (5, 5), 1.0)
    sharpened = processor._sharpen(blurred, context)

    assert np.std(sharpened) > np.std(blurred)


def test_sharpen_with_zero_strength_is_nearly_identical(processor):
    image = np.random.randint(0, 256, (100, 100), dtype=np.uint8)

    context = {
        "sharpen_strength": 0.0,
        "sharpen_blur_kernel_size": (3, 3),
        "sharpen_sigma": 1.0,
    }

    result = processor._sharpen(image, context)

    diff = np.mean(np.abs(result.astype(np.int16) - image.astype(np.int16)))
    assert diff < 1


def test_sharpen_extreme_strength_does_not_break_image(processor):
    image = np.random.randint(0, 256, (200, 200), dtype=np.uint8)

    context = {
        "sharpen_strength": 10.0,
        "sharpen_blur_kernel_size": (3, 3),
        "sharpen_sigma": 1.0,
    }

    result = processor._sharpen(image, context)

    assert result.shape == image.shape
    assert result.dtype == np.uint8
    assert result.min() >= 0
    assert result.max() <= 255
