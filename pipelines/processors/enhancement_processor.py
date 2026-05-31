import logging
from domain.exceptions import ProcessorError
from pipelines.processors.base_processor import BaseProcessor

logger = logging.getLogger(__name__)


class EnhancementProcessor(BaseProcessor):
    name = "enhancement"

    def __init__(self):
        super().__init__(
            transform_steps = [
                self._apply_clahe,
            ]
        )


    def _apply_clahe(self, image, context: dict):
        import cv2

        clip_limit = context.get("clahe_clip_limit", 2.0)
        tile_grid_size = context.get("clahe_tile_grid_size", (4, 4))

        clahe = cv2.createCLAHE(clipLimit=clip_limit, tileGridSize=tile_grid_size)

        try:
            result = clahe.apply(image)
        except Exception as e:
            logger.error(f"Error when applying CLAHE: {e}")
            raise ProcessorError(f"Error when applying CLAHE: {e}") from e

        return result
