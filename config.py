import os
from database.session import DATABASE_URL


class Config:
    SQLALCHEMY_DATABASE_URI = DATABASE_URL
    CELERY_BROKER_URL = os.getenv('CELERY_BROKER_URL')
    CELERY_RESULT_BACKEND = os.getenv('CELERY_BROKER_URL')