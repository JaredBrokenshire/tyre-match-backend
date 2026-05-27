import os
import pytest
from unittest.mock import MagicMock
from domain.exceptions import PipelineError, ProcessorError
from database.models.data_types.files import FileModel
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

    with pytest.raises(
            PipelineError,
            match="TyreImpressionProcessing original image location does not exist"
    ):
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

    with pytest.raises(PipelineError, match="ProcessorError from stage1 processor: test error"):
        result = pipeline.process(tyre_impression_processing.id)
        assert result is None


def test_tyre_impression_processing_pipeline_process(pipeline):
    # Mock pipeline processor stages
    stage1 = MagicMock()
    stage1.name = "stage1"
    stage1.process.return_value = "path1"

    stage2 = MagicMock()
    stage2.name = "stage2"
    stage2.process.return_value = "path2"

    pipeline.stages = [stage1, stage2]

    tyre_impression = TyreImpressionFactory.create()
    tyre_impression_processing = TyreImpressionProcessingFactory.create(tyre_impression.id)
    os.makedirs("/files/test_directory", exist_ok=True)
    original_file = FileFactory.create(
        model=FileModel.tyre_impression,
        model_id=tyre_impression_processing.id,
        file_location="/files/test_directory"
    )

    result = pipeline.process(tyre_impression_processing.id)

    assert stage1.process.called
    assert stage2.process.called

    assert stage1.process.call_args[0][0] == original_file.file_location
    assert stage2.process.call_args[0][0] == "path1"

    assert "path2" == result
