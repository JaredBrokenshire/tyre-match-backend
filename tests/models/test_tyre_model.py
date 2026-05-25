from tests.helpers.factories import TyreModelFactory


def test_repr():
    tyre_model = TyreModelFactory().create()

    res = tyre_model.__repr__()
    assert res == f'<TyreModel {tyre_model.manufacturer} {tyre_model.model_name}>'