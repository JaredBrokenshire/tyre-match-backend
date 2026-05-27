import logging
from pipelines.processors.base_processor import BaseProcessor
from domain.exceptions import FileReadError, FileSaveError, ProcessorError

logger = logging.getLogger(__name__)

class NormalisationProcessor(BaseProcessor):
    name = "normalisation"

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

        clip_limit = context["clahe_clip_limit"]
        tile_grid_size = context["clahe_tile_grid_size"]

        clahe = cv2.createCLAHE(clipLimit=clip_limit, tileGridSize=tile_grid_size)

        return clahe.apply(image)