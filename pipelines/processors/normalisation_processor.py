import logging
from pipelines.processors.base_processor import BaseProcessor

logger = logging.getLogger(__name__)


class NormalisationProcessor(BaseProcessor):
    name = "normalisation"

    def __init__(self):
        super().__init__(
            transform_steps=[
                # TODO: Write normalisation steps
            ]
        )