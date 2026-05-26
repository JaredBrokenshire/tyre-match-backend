import os
import pytest
import shutil
import logging
import sqlalchemy
from main import create_app


os.environ["TESTING"] = "1"

@pytest.fixture(scope="session")
def app():
    app = create_app()
    app.config["TESTING"] = True
    return app

@pytest.fixture(autouse=True)
def app_ctx(app):
    with app.app_context():
        yield

@pytest.fixture(autouse=True)
def database_session(app):
    from database.extensions import db

    connection = db.engine.connect()
    transaction = connection.begin()

    db.session.remove()
    db.session.bind = connection

    nested = connection.begin_nested()

    @sqlalchemy.event.listens_for(db.session, "after_transaction_end")
    def restart_savepoint(session, trans):
        nonlocal nested

        # only restart savepoint when the outer transaction ends
        if trans.nested and not trans._parent.nested:
            nested = connection.begin_nested()

    try:
        yield db.session
    finally:
        db.session.remove()
        transaction.rollback()
        connection.close()

@pytest.fixture
def client(app):
    with app.test_client() as client:
        with app.app_context():
            yield client

@pytest.fixture(autouse=True, scope="session")
def disable_logging():
    logging.disable(logging.CRITICAL)
    yield
    logging.disable(logging.NOTSET)

@pytest.fixture(autouse=True)
def cleanup():
    yield
    shutil.rmtree("/files/test_directory", ignore_errors=True)
