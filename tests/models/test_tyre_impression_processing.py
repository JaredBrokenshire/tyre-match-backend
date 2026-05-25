from tests.helpers.factories import TyreImpressionProcessingFactory, TyreImpressionFactory


def test_repr(database_session):
    tyre_impression = TyreImpressionFactory().create()
    tyre_impression_processing = TyreImpressionProcessingFactory().create(tyre_impression.id)

    res = tyre_impression_processing.__repr__()
    assert res == f"<TyreImpressionProcessing {tyre_impression_processing.id} v{tyre_impression_processing.preprocessing_version}>"