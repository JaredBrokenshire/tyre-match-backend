from database.models.data_types import FileModel, FileType
from database.repositories import TyreImpressionProcessingRepository
from tests.helpers.factories import TyreImpressionFactory, TyreImpressionProcessingFactory, FileFactory


def test_get_by_id_invalid_id():
    repo = TyreImpressionProcessingRepository()

    result = repo.get_by_id(1000)
    assert result is None


def test_get_by_id_valid_id():
    repo = TyreImpressionProcessingRepository()

    tyre_impression = TyreImpressionFactory().create()
    tyre_impression_processing = TyreImpressionProcessingFactory().create(tyre_impression.id)

    result = repo.get_by_id(tyre_impression_processing.id)
    assert tyre_impression_processing == result


def test_get_by_id_load_files():
    repo = TyreImpressionProcessingRepository()

    tyre_impression = TyreImpressionFactory().create()
    tyre_impression_processing = TyreImpressionProcessingFactory().create(tyre_impression.id)

    # Create files
    original_file = FileFactory().create(
        model=FileModel.tyre_impression,
        model_id=tyre_impression_processing.id,
        file_type=FileType.original
    )
    normalised_file = FileFactory().create(
        model=FileModel.tyre_impression,
        model_id=tyre_impression_processing.id,
        file_type=FileType.normalised
    )
    enhanced_file = FileFactory().create(
        model=FileModel.tyre_impression,
        model_id=tyre_impression_processing.id,
        file_type=FileType.enhanced
    )
    binary_file = FileFactory().create(
        model=FileModel.tyre_impression,
        model_id=tyre_impression_processing.id,
        file_type=FileType.binary
    )
    clean_file = FileFactory().create(
        model=FileModel.tyre_impression,
        model_id=tyre_impression_processing.id,
        file_type=FileType.clean
    )
    skeleton_file = FileFactory().create(
        model=FileModel.tyre_impression,
        model_id=tyre_impression_processing.id,
        file_type=FileType.skeleton
    )

    result = repo.get_by_id(tyre_impression_processing.id)
    
    assert tyre_impression_processing == result
    assert result.files is not None
    assert original_file == result.files.get(FileType.original.value)
    assert normalised_file == result.files.get(FileType.normalised.value)
    assert enhanced_file == result.files.get(FileType.enhanced.value)
    assert binary_file == result.files.get(FileType.binary.value)
    assert clean_file == result.files.get(FileType.clean.value)
    assert skeleton_file == result.files.get(FileType.skeleton.value)


def test_get_by_tyre_impression_id_invalid_id(database_session):
    repo = TyreImpressionProcessingRepository()

    result = repo.get_by_tyre_impression_id(1)
    assert result is None


def test_get_by_tyre_impression_id(database_session):
    repo = TyreImpressionProcessingRepository()

    tyre_impression = TyreImpressionFactory.create()
    tyre_impression_processing = TyreImpressionProcessingFactory.create(tyre_impression.id)

    result = repo.get_by_tyre_impression_id(tyre_impression.id)
    assert result == tyre_impression_processing

