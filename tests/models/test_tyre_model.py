from database.models import TyreModel


def test_repr():
    tyre_model = TyreModel(manufacturer="Michelin", model_name="Pilot Sport")

    res = tyre_model.__repr__()
    assert res == f'<TyreModel {tyre_model.manufacturer} {tyre_model.model_name}>'