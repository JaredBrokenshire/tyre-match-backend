from unittest.mock import patch

import cv2
import numpy as np
import pytest
from database.models.data_types.files import FileType
from domain.exceptions import ProcessorError
from pipelines.processors.normalisation_processor import NormalisationProcessor


@pytest.fixture
def processor():
    return NormalisationProcessor()


@pytest.fixture
def context():
    return {
        "target_width": 4096,
        "target_height": 4096,
        "skew_angle_threshold": 0.25,
        "output_directories": {"normalisation", "/files/test_directory/normalised"},
        "processing_id": 1,
        "file_types_on_completion": {"normalisation": FileType.normalised},
    }


def test_resize_error_from_cv2_resize(processor, context):
    image = np.zeros((10,10))

    with patch("cv2.resize", side_effect=Exception("test error")):
        with pytest.raises(ProcessorError, match="Error when resizing image: test error"):
            processor._resize(image, context)


def test_resize_increase_image_size(processor, context):
    image = np.zeros((10,10))

    result = processor._resize(image, context)

    assert result.shape[0] == context["target_height"]
    assert result.shape[1] == context["target_width"]

    original_aspect_ratio = image.shape[1] / image.shape[0]
    result_aspect_ratio = result.shape[1] / result.shape[0]

    assert original_aspect_ratio == result_aspect_ratio


def test_resize_decrease_image_size(processor, context):
    image = np.zeros((10000,10000))

    result = processor._resize(image, context)

    assert result.shape[0] == context["target_height"]
    assert result.shape[1] == context["target_width"]

    original_aspect_ratio = image.shape[1] / image.shape[0]
    result_aspect_ratio = result.shape[1] / result.shape[0]

    assert original_aspect_ratio == result_aspect_ratio


def test_resize_landscape_rectangle(processor, context):
    image = np.zeros((100,200))

    result = processor._resize(image, context)

    assert result.shape[0] < context["target_height"]
    assert result.shape[1] == context["target_width"]

    original_aspect_ratio = image.shape[1] / image.shape[0]
    result_aspect_ratio = result.shape[1] / result.shape[0]

    assert original_aspect_ratio == result_aspect_ratio


def test_resize_portrait_rectangle(processor, context):
    image = np.zeros((200,100))

    result = processor._resize(image, context)

    assert result.shape[0] == context["target_height"]
    assert result.shape[1] < context["target_width"]

    original_aspect_ratio = image.shape[1] / image.shape[0]
    result_aspect_ratio = result.shape[1] / result.shape[0]

    assert original_aspect_ratio == result_aspect_ratio


def test_resize_is_deterministic(processor, context):
    image = np.zeros((100,200))

    result_1 = processor._resize(image, context)
    result_2 = processor._resize(image, context)

    assert np.array_equal(result_1, result_2)


def test_correct_skew_error_from_cv2_threshold(processor, context):
    image = np.zeros((10,10))

    with patch("cv2.threshold", side_effect=Exception("test error")):
        with pytest.raises(ProcessorError, match="Error when building OTSU mask: test error"):
            processor._correct_skew(image, context)


def test_correct_skew_not_enough_structure_to_determine_skew(processor, context):
    image = np.zeros((10,10))

    with patch("cv2.threshold", return_value=(None, image)):
        result = processor._correct_skew(image, context)

    assert np.array_equal(image, result)


def test_correct_skew_horizontal_line_no_rotation(processor, context):
    image = np.zeros((200, 200), dtype=np.uint8)
    image[100, 20:180] = 255

    with patch("cv2.threshold", return_value=(None, image)):
        result = processor._correct_skew(image, context)

    assert result.shape == image.shape
    assert np.array_equal(result, image)


def test_correct_skew_vertical_line(processor, context):
    image = np.zeros((2000, 2000), dtype=np.uint8)
    image[20:180, :] = 255

    with patch("cv2.threshold", return_value=(None, image)):
        result = processor._correct_skew(image, context)

    assert result is not None
    assert result.ndim == 2
    assert not np.array_equal(result, image)


def test_correct_skew_diagonal_rotation(processor, context):
    image = np.zeros((2000, 2000), dtype=np.uint8)
    for i in range(100, 1900):
        image[i, i] = 255

    result = processor._correct_skew(image, context)

    # rotated image should differ
    assert not np.array_equal(result, image)
    assert result.shape[0] >= image.shape[0]
    assert result.shape[1] >= image.shape[1]


def test_correct_skew_no_structure_skips_rotation(processor, context):
    image = np.random.randint(0, 255, (2000, 2000), dtype=np.uint8)

    result = processor._correct_skew(image, context)

    # should safely return image or unchanged result
    assert np.array_equal(result, image)


def test_correct_skew_error_from_cv2_get_rotation_matrix(processor, context):
    image = np.zeros((2000, 2000), dtype=np.uint8)
    image[1000, :] = 255

    with patch("cv2.getRotationMatrix2D", side_effect=Exception("test error")):
        with pytest.raises(ProcessorError, match="Error when getting rotation matrix: test error"):
            processor._correct_skew(image, context)


def test_correct_skew_warp_affine_error(processor, context):
    image = np.zeros((2000, 2000), dtype=np.uint8)
    image[1000, :] = 255

    with patch("cv2.warpAffine", side_effect=Exception("test error")):
        with pytest.raises(ProcessorError, match="Error when applying rotation matrix: test error"):
            processor._correct_skew(image, context)


def test_correct_skew_is_deterministic(processor, context):
    image = np.zeros((2000, 2000), dtype=np.uint8)
    image[1000, :] = 255

    result_1 = processor._correct_skew(image, context)
    result_2 = processor._correct_skew(image, context)

    assert np.array_equal(result_1, result_2)