import logging
from flask import Flask
from flask_cors import CORS
from celery_app import init_celery
from database.extensions import db, migrate


def create_app():
    app = Flask(__name__)

    CORS(app, origins=[
        "http://localtyrematch.com:8080",
        "http://localhost:8080",
    ])

    app.config.from_object("config.Config")
    app.logger.setLevel(logging.INFO)

    db.init_app(app)
    migrate.init_app(app, db, directory="database/migrations")

    init_celery(app)

    from api.routes import register_blueprints
    register_blueprints(app)

    return app
