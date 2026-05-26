from tests.helpers.factories import TyreImpressionFactory


def test_repr():
    tyre_impression = TyreImpressionFactory.create()

    res = tyre_impression.__repr__()
    assert res == f'<TyreImpression {tyre_impression.id}>'
