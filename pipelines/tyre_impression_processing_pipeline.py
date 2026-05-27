import os
import logging
from database.models.data_types.files import FileType
from domain.exceptions import ProcessorError, PipelineError
from pipelines.processors.normalisation import NormalisationProcessor
from database.repositories.tyre_impression_processing_repository import TyreImpressionProcessingRepository

logger = logging.getLogger(__name__)

BASE_OUTPUT_DIRECTORY = "/tyre_match/tyre_impressions/"

class TyreImpressionProcessingPipeline:
    def __init__(self):
        self.stages = [
            NormalisationProcessor(),
            # ContrastProcessor(),
            # BinarisationProcessor(),
            # CleaningProcessor(),
            # SkeletonisationProcessor(),
        ]

        self.context = {
            # General
            "processing_id": None,

            # Normalisation
            "clahe_clip_limit": 2.5,
            "clahe_tile_grid_size": (4,4),

            "file_types_on_completion": {
                "normalisation": FileType.normalised,
                "enhanced": FileType.enhanced,
                "binary": FileType.binary,
                "clean": FileType.clean,
                "skeleton": FileType.skeleton,
            },
            "output_directories": {},
            "execution_trace": []
        }

        self.tyre_impression_processing_repository = TyreImpressionProcessingRepository()

    def process(self, processing_id: int):
        tyre_impression_processing = self.tyre_impression_processing_repository.get_by_id(processing_id)
        if not tyre_impression_processing:
            logger.error(f"TyreImpressionProcessing with id {processing_id} not found")
            raise PipelineError(f"TyreImpressionProcessing with id {processing_id} not found")

        original_file = tyre_impression_processing.files.get(FileType.original.value)
        if not original_file or not original_file.file_location:
            logger.error("TyreImpressionProcessing has no original image")
            raise PipelineError("TyreImpressionProcessing has no original image")

        current_path = original_file.file_location
        if not os.path.exists(current_path):
            logger.error("TyreImpressionProcessing original image location does not exist")
            raise PipelineError("TyreImpressionProcessing original image location does not exist")


        base_output_directory = f"/tyre_match/files/tyre_impressions/{tyre_impression_processing.tyre_impression_id}"

        self.context["processing_id"] = tyre_impression_processing.id
        self.context["output_directories"] = {
                "normalisation": f"{base_output_directory}/normalised",
                "enhanced": f"{base_output_directory}/enhanced",
                "binary": f"{base_output_directory}/binary",
                "clean": f"{base_output_directory}/clean",
                "skeleton": f"{base_output_directory}/skeleton",
        }

        for stage in self.stages:
            self.context["execution_trace"].append(stage.name)
            try:
                current_path = stage.process(current_path, self.context)
            except ProcessorError as e:
                logger.error(f"ProcessorError from {stage.name} processor: {e}")
                raise PipelineError(f"ProcessorError from {stage.name} processor: {e}")

        return current_path
