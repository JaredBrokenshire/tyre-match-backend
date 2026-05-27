import os
from flask import Flask
from flask_cors import CORS
from database.extensions import db, migrate
from celery_config.config import init_celery
from logging_config.config import setup_logging


def create_app():
    setup_logging()

    app = Flask(__name__)

    CORS(app, origins=[
        "http://localtyrematch.com:8080",
        "http://localhost:8080",
    ])

    app.config.from_object("config.Config")
    app.config.from_mapping(
        CELERY=dict(
            broker_url=os.getenv("CELERY_BROKER_URL"),
            result_backend=os.getenv("CELERY_RESULT_BACKEND"),
            task_ignore_result=True,
        ),
    )

    db.init_app(app)
    migrate.init_app(app, db, directory="database/migrations")

    init_celery(app)

    from api.routes import register_blueprints
    register_blueprints(app)

    return app
