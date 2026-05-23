from database.extensions import db


class UnitOfWork:
    def __init__(self, session=None):
        self.session = session or db.session

    def __enter__(self):
        return self

    def __exit__(self, exc_type, exc, tb):
        if exc_type:
            self.session.rollback()
        else:
            self.session.commit()
