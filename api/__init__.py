from flask import Flask
from database.session import DATABASE_URL
from database.extensions import db, migrate


def create_app():
    app = Flask(__name__)

    app.config.from_mapping(
        SQLALCHEMY_DATABASE_URI = DATABASE_URL,
    )

    db.init_app(app)
    migrate.init_app(app, db, directory="database/migrations")

    return app
