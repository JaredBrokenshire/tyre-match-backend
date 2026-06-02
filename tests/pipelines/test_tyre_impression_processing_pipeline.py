import os
import pytest
import numpy as np
from unittest.mock import MagicMock, patch
from database.models.data_types.files import FileModel
from domain.exceptions import PipelineError, ProcessorError
from tests.helpers.factories.file_factory import FileFactory
from tests.helpers.factories.tyre_impression_factory import TyreImpressionFactory
from pipelines.tyre_impression_processing_pipeline import TyreImpressionProcessingPipeline
from tests.helpers.factories.tyre_impression_processing_factory import TyreImpressionProcessingFactory


@pytest.fixture
def pipeline():
    return TyreImpressionProcessingPipeline()


def test_tyre_impression_processing_pipeline_process_invalid_processing_id(pipeline):
    with pytest.raises(PipelineError, match="TyreImpressionProcessing with id 999 not found"):
        pipeline.process(999)


def test_tyre_impression_processing_pipeline_process_empty_original_file_location(pipeline):
    tyre_impression = TyreImpressionFactory.create()
    tyre_impression_processing = TyreImpressionProcessingFactory.create(tyre_impression.id)

    with pytest.raises(PipelineError, match="TyreImpressionProcessing has no original image"):
        pipeline.process(tyre_impression_processing.id)


def test_tyre_impression_processing_pipeline_process_invalid_original_file_location(pipeline):
    tyre_impression = TyreImpressionFactory.create()
    tyre_impression_processing = TyreImpressionProcessingFactory.create(tyre_impression.id)
    FileFactory.create(
        model=FileModel.tyre_impression,
        model_id=tyre_impression_processing.id,
        file_location="/invalid/file/location"
    )

    with pytest.raises(PipelineError, match="TyreImpressionProcessing original image location does not exist"):
        pipeline.process(tyre_impression_processing.id)


def test_process_invalid_file_read_from_cv2_imread(pipeline):
    tyre_impression = TyreImpressionFactory.create()
    tyre_impression_processing = TyreImpressionProcessingFactory.create(tyre_impression.id)
    os.makedirs("/files/test_directory", exist_ok=True)
    FileFactory.create(
        model=FileModel.tyre_impression,
        model_id=tyre_impression_processing.id,
        file_location="/files/test_directory"
    )

    with patch("cv2.imread", return_value=None):
        with pytest.raises(PipelineError, match="Unable to read image in tyre impression processing pipeline"):
            pipeline.process(tyre_impression_processing.id)


def test_tyre_impression_processing_pipeline_processor_error(pipeline):
    # Mock pipeline processor stages
    stage1 = MagicMock()
    stage1.name = "stage1"
    stage1.process.side_effect = ProcessorError("test error")

    stage2 = MagicMock()
    stage2.name = "stage2"
    stage2.process.return_value = "path2"

    pipeline.stages = [stage1, stage2]

    tyre_impression = TyreImpressionFactory.create()
    tyre_impression_processing = TyreImpressionProcessingFactory.create(tyre_impression.id)
    os.makedirs("/files/test_directory", exist_ok=True)
    FileFactory.create(
        model=FileModel.tyre_impression,
        model_id=tyre_impression_processing.id,
        file_location="/files/test_directory"
    )

    with patch("cv2.imread", return_value=np.zeros((10,10))):
        with pytest.raises(PipelineError, match="ProcessorError from stage1 processor: test error"):
            pipeline.process(tyre_impression_processing.id)


def test_tyre_impression_processing_pipeline_process(pipeline):
    # Mock pipeline processor stages
    stage1 = MagicMock()
    stage1.name = "stage1"
    stage1.process.return_value = np.ones((10,10))

    stage2 = MagicMock()
    stage2.name = "stage2"
    stage2.process.return_value = np.random.randint((10,10))

    pipeline.stages = [stage1, stage2]

    tyre_impression = TyreImpressionFactory.create()
    tyre_impression_processing = TyreImpressionProcessingFactory.create(tyre_impression.id)
    os.makedirs("/files/test_directory", exist_ok=True)
    original_file = FileFactory.create(
        model=FileModel.tyre_impression,
        model_id=tyre_impression_processing.id,
        file_location="/files/test_directory"
    )

    with patch("cv2.imread", return_value=np.zeros((10,10))):
        result = pipeline.process(tyre_impression_processing.id)

    stage1.process.assert_called_once()
    assert np.array_equal(stage1.process.call_args[0][0], np.zeros((10,10)))

    stage2.process.assert_called_once()
    assert np.array_equal(stage2.process.call_args[0][0], stage1.process.return_value)

    assert np.array_equal(stage2.process.return_value, result)
