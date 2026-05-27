from flask import Flask
from celery import Celery, Task

def init_celery(app: Flask):
    class ContextTask(Task):
        def __call__(self, *args, **kwargs):
            with app.app_context():
                return self.run(*args, **kwargs)

    celery_app = Celery(app.name, task_cls=ContextTask)
    celery_app.conf.update(app.config["CELERY"])
    celery_app.set_default()
    app.extensions["celery"] = celery_app

    return celery_app