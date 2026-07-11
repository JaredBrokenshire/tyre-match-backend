import os
import cv2
import logging
from database.models.data_types.files import FileType
from domain.exceptions import ProcessorError, PipelineError
from pipelines.processors.enhancement_processor import EnhancementProcessor
from pipelines.processors.normalisation_processor import NormalisationProcessor
from database.repositories.tyre_impression_repository import TyreImpressionRepository

logger = logging.getLogger(__name__)


class TyreImpressionProcessingPipeline:
    def __init__(self):
        self.stages = [
            NormalisationProcessor(),
            EnhancementProcessor(),
            # BinarisationProcessor(),
            # CleaningProcessor(),
            # SkeletonisationProcessor(),
        ]

        self.context = {
            # General
            "processing_id": None,

            # Normalisation
            "target_width": 4096,
            "target_height": 4096,
            "skew_anisotropy_threshold": 1.0,

            # Enhancement
            "denoise_h": 10,
            "denoise_template_window_size": 7,
            "denoise_search_window_size": 12,
            "clahe_clip_limit": 3.0,
            "clahe_tile_grid_size": (4,4),
            "sharpen_strength": 1.2,
            "sharpen_blur_kernel_size": (0,0),
            "sharpen_sigma": 1.0,

            "file_types_on_completion": {
                "normalisation": FileType.normalised,
                "enhancement": FileType.enhanced,
                "binary": FileType.binary,
                "clean": FileType.clean,
                "skeleton": FileType.skeleton,
            },
            "output_directories": {},
            "execution_trace": []
        }

        self.tyre_impression_repository = TyreImpressionRepository()

    def process(self, processing_id: int):
        tyre_impression = self.tyre_impression_repository.get_by_id(processing_id)
        if not tyre_impression:
            logger.error(f"TyreImpressionProcessing with id {processing_id} not found")
            raise PipelineError(f"TyreImpressionProcessing with id {processing_id} not found")

        original_file = tyre_impression.files.get(FileType.original.value)
        if not original_file or not original_file.file_location:
            logger.error("TyreImpressionProcessing has no original image")
            raise PipelineError("TyreImpressionProcessing has no original image")

        if not os.path.exists(original_file.file_location):
            logger.error("TyreImpressionProcessing original image location does not exist")
            raise PipelineError("TyreImpressionProcessing original image location does not exist")

        image = cv2.imread(f"{original_file.file_location}/{original_file.file_name}", cv2.IMREAD_GRAYSCALE)
        if image is None:
            logger.error(f"Unable to read image in tyre impression processing pipeline")
            raise PipelineError(f"Unable to read image in tyre impression processing pipeline")

        base_output_directory = f"/tyre_match/files/tyre_impressions/{tyre_impression.id}"

        self.context["processing_id"] = tyre_impression.id
        self.context["output_directories"] = {
                "normalisation": f"{base_output_directory}/normalised",
                "enhancement": f"{base_output_directory}/enhanced",
                "binary": f"{base_output_directory}/binary",
                "clean": f"{base_output_directory}/clean",
                "skeleton": f"{base_output_directory}/skeleton",
        }

        for stage in self.stages:
            self.context["execution_trace"].append(stage.name)
            try:
                image = stage.process(image, self.context)
            except ProcessorError as e:
                logger.error(f"ProcessorError from {stage.name} processor: {e}")
                raise PipelineError(f"ProcessorError from {stage.name} processor: {e}")

        return image
