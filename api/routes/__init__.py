from .tyre_model_routes import tyre_model_blueprint
from .tyre_impression_routes import tyre_impression_blueprint

def register_blueprints(app):
    app.register_blueprint(tyre_model_blueprint)
    app.register_blueprint(tyre_impression_blueprint)
