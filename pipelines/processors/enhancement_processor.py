import logging
from pipelines.processors.base_processor import BaseProcessor
from domain.exceptions import FileReadError, FileSaveError, ProcessorError

logger = logging.getLogger(__name__)

class EnhancementProcessor(BaseProcessor):
    name = "enhancement"

    def process(self, input_path: str, context: dict) -> str:
        import cv2

        image = cv2.imread(input_path, cv2.IMREAD_GRAYSCALE)
        if image is None:
            logger.error("Unable to read image in NormalisationProcessor")
            raise FileReadError("Unable to read image in NormalisationProcessor")

        processed = self.transform(image, context)

        try:
            file_location = self.save_file(processed, context, self.name)
        except FileSaveError as e:
            logger.error(f"FileSaveError from {self.name} processor: {e}")
            raise ProcessorError(f"FileSaveError from {self.name} processor") from e

        return file_location


    def transform(self, image, context: dict):
        import cv2

        clip_limit = context.get("clahe_clip_limit", 2.0)
        tile_grid_size = context.get("clahe_tile_grid_size", (4,4))

        clahe = cv2.createCLAHE(clipLimit=clip_limit, tileGridSize=tile_grid_size)

        try:
            result = clahe.apply(image)
        except Exception as e:
            logger.error(f"Exception from apply CLAHE in transform: {e}")
            raise ProcessorError(f"Exception from apply CLAHE in transform: {e}") from e

        return result