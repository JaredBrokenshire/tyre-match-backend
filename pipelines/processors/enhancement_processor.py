import cv2
import logging
import numpy as np
from domain.exceptions import ProcessorError
from pipelines.processors.base_processor import BaseProcessor

logger = logging.getLogger(__name__)


class EnhancementProcessor(BaseProcessor):
    name = "enhancement"

    def __init__(self):
        super().__init__(
            transform_steps = [
                self._denoise,
                self._apply_clahe,
                self._sharpen,
            ]
        )


    def _denoise(self, image, context: dict) -> np.ndarray:
        denoise_h = context.get("denoise_h", 10)
        template_window_size = context.get("denoise_template_window_size", 7)
        search_window_size = context.get("denoise_search_window_size", 21)

        try:
            result = cv2.fastNlMeansDenoising(
                image,
                None,
                denoise_h,
                template_window_size,
                search_window_size
            )
        except Exception as e:
            logger.error(f"Error when applying denoise: {e}")
            raise ProcessorError(f"Error when applying denoise: {e}") from e

        return result


    def _apply_clahe(self, image, context: dict) -> np.ndarray:
        clip_limit = context.get("clahe_clip_limit", 2.0)
        tile_grid_size = context.get("clahe_tile_grid_size", (4, 4))

        clahe = cv2.createCLAHE(clipLimit=clip_limit, tileGridSize=tile_grid_size)

        try:
            result = clahe.apply(image)
        except Exception as e:
            logger.error(f"Error when applying CLAHE: {e}")
            raise ProcessorError(f"Error when applying CLAHE: {e}") from e

        return result


    def _sharpen(self, image, context: dict) -> np.ndarray:
        strength = context.get("sharpen_strength", 1.2)
        blur_ksize = context.get("sharpen_blur_kernel_size", (0, 0))
        sigma = context.get("sharpen_sigma", 1.0)

        try:
            blurred = cv2.GaussianBlur(image, blur_ksize, sigma)
        except Exception as e:
            logger.error(f"Error when blurring during sharpen: {e}")
            raise ProcessorError(f"Error when blurring during sharpen: {e}") from e

        try:
            result = cv2.addWeighted(image, 1.0 + strength, blurred, -strength, 0)
        except Exception as e:
            logger.error(f"Error when applying sharpen: {e}")
            raise ProcessorError(f"Error when applying sharpen: {e}") from e

        return result
