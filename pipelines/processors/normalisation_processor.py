import logging
from pipelines.processors.base_processor import BaseProcessor

logger = logging.getLogger(__name__)


class NormalisationProcessor(BaseProcessor):
    name = "normalisation"

    def process(self, input_path: str, context: dict) -> str:
        return input_path


    def transform(self, image, context: dict) -> str:
        pass