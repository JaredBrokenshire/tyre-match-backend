from database.models import TyreImpression


def test_repr():
    tyre_impression = TyreImpression(file_path="/test/file/path", status="test status")

    res = tyre_impression.__repr__()
    assert res == f'<TyreImpression {tyre_impression.file_path} {tyre_impression.status}>'
