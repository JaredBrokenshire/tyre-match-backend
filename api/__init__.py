import logging

from flask import Flask
from flask_cors import CORS
from database.session import DATABASE_URL
from database.extensions import db, migrate


def create_app():
    app = Flask(__name__)

    CORS(app, origins=[
        "http://localtyrematch.com:8080",
        "http://localhost:8080",
    ])

    app.config.from_mapping(
        SQLALCHEMY_DATABASE_URI = DATABASE_URL,
    )
    app.logger.setLevel(logging.INFO)

    db.init_app(app)
    migrate.init_app(app, db, directory="database/migrations")

    from api.routes import register_blueprints
    register_blueprints(app)

    return app
