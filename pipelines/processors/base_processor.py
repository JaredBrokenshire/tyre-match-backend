import uuid
import logging
import numpy as np
from typing import Callable, Any
from database.models.data_types.files import FileModel
from services.file_service import ProcessedImageRequest, FileService
from domain.exceptions import FileSaveError, DatabaseError, ProcessorError

logger = logging.getLogger(__name__)

class BaseProcessor:
    name: str = "base"

    def __init__(self, transform_steps: list[Callable[[Any, dict], Any]]):
        self.file_service = FileService()
        self.transform_steps: list[Callable[[Any, dict], Any]] = transform_steps


    def process(self, image: np.ndarray, context: dict) -> np.ndarray:
        try:
            processed = self.transform(image, context)
        except ProcessorError as e:
            logger.error(f"ProcessorError when transforming image in {self.name} processor: {e}")
            raise ProcessorError(f"ProcessorError when transforming image in {self.name} processor: {e}") from e

        try:
            self.save_file(processed, context, self.name)
        except FileSaveError as e:
            logger.error(f"FileSaveError from {self.name} processor: {e}")
            raise ProcessorError(f"FileSaveError from {self.name} processor") from e

        return processed


    def transform(self, image, context: dict):
        for transform_step in self.transform_steps:
            try:
                image = transform_step(image, context)
            except ProcessorError as e:
                logger.error(f"ProcessorError from {self.name} processor ({transform_step.__name__}): {e}")
                raise ProcessorError(
                    f"ProcessorError from {self.name} processor ({transform_step.__name__}): {e}") from e

        return image


    def save_file(self, image, context: dict, stage_name: str):
        request = ProcessedImageRequest(
            image=image,
            file_name=f"{uuid.uuid4()}.png",
            upload_directory=context["output_directories"][stage_name],
            model=FileModel.tyre_impression,
            model_id=context["processing_id"],
            file_type=context["file_types_on_completion"][stage_name],
        )

        try:
            self.file_service.save_processed_image(request)
        except PermissionError as e:
            logger.error(f"PermissionError from file service in {stage_name} processor: {e}")
            raise FileSaveError(f"PermissionError from file service in {stage_name} processor: {e}") from e
        except OSError as e:
            logger.error(f"OSError from file service in {stage_name} processor: {e}")
            raise FileSaveError(f"OSError from file service in {stage_name} processor: {e}") from e
        except FileSaveError as e:
            logger.error(f"FileSaveError from file service in {stage_name} processor: {e}")
            raise FileSaveError(f"FileSaveError from file service in {stage_name} processor: {e}") from e
        except DatabaseError as e:
            logger.error(f"DatabaseError from file service in {stage_name} processor: {e}")
            raise FileSaveError(f"DatabaseError from file service in {stage_name} processor: {e}") from e