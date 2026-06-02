import cv2
import logging
import numpy as np
from domain.exceptions import ProcessorError
from pipelines.processors.base_processor import BaseProcessor

logger = logging.getLogger(__name__)


class NormalisationProcessor(BaseProcessor):
    name = "normalisation"

    def __init__(self):
        super().__init__(
            transform_steps=[
                self._correct_skew,
                self._resize,
            ]
        )


    def _resize(self, image, context: dict):
        target_width = context.get("target_width", 4096)
        target_height = context.get("target_height", 4096)

        original_height, original_width = image.shape[:2]

        scale = min(target_width / original_width, target_height / original_height)

        new_width = int(original_width * scale)
        new_height = int(original_height * scale)

        try:
            result = cv2.resize(image, (new_width, new_height), interpolation=cv2.INTER_LINEAR)
        except Exception as e:
            logger.error(f"Error when resizing image: {e}")
            raise ProcessorError(f"Error when resizing image: {e}") from e

        return result


    def _correct_skew(self, image, context: dict):
        try:
            _, binary = cv2.threshold(
                image,
                0,
                255,
                cv2.THRESH_BINARY + cv2.THRESH_OTSU
            )
        except Exception as e:
            logger.error(f"Error when building OTSU mask: {e}")
            raise ProcessorError(f"Error when building OTSU mask: {e}") from e

        coords = np.column_stack(np.where(binary > 0))

        # Not enough structure to estimate skew reliably
        if coords.shape[0] < 1000:
            return image

        # PCA for dominant orientation
        mean = np.mean(coords, axis=0)
        centered = coords - mean

        covariance_matrix = np.cov(centered, rowvar=False)
        eigen_values, eigen_vectors = np.linalg.eigh(covariance_matrix)

        skew_anisotropy_threshold = context.get("skew_anisotropy_threshold", 2.5)
        sorted_vals = np.sort(eigen_values)
        if sorted_vals[-1] / (sorted_vals[-2] + 1e-8) < skew_anisotropy_threshold:
            return image

        principal_axis = eigen_vectors[:, np.argmax(eigen_values)]

        angle = np.arctan2(principal_axis[1], principal_axis[0])
        angle_degrees = np.degrees(angle)


        # Compute expanded canvas size
        (height, width) = image.shape[:2]
        center = (width // 2, height // 2)

        # Normalize to minimal rotation
        angle_degrees = (angle_degrees + 90) % 180 - 180
        angle_radians = np.deg2rad(angle_degrees)
        cos = abs(np.cos(angle_radians))
        sin = abs(np.sin(angle_radians))

        new_width = int(height * sin + width * cos)
        new_height = int(height * cos + width * sin)

        # Rotation matrix with translation adjustment
        try:
            rotation_matrix = cv2.getRotationMatrix2D(center, angle_degrees, 1.0)
        except Exception as e:
            logger.error(f"Error when getting rotation matrix: {e}")
            raise ProcessorError(f"Error when getting rotation matrix: {e}") from e

        rotation_matrix[0, 2] += (new_width - width) / 2
        rotation_matrix[1, 2] += (new_height - height) / 2

        # Apply warp with clean border
        try:
            return cv2.warpAffine(
                image,
                rotation_matrix,
                (new_width, new_height),
                flags=cv2.INTER_LINEAR,
                borderMode=cv2.BORDER_CONSTANT,
                borderValue=0
            )
        except Exception as e:
            logger.error(f"Error when applying rotation matrix: {e}")
            raise ProcessorError(f"Error when applying rotation matrix: {e}") from e