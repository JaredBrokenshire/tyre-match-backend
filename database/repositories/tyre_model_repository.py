import logging
from database.models.file import File
from database.models.tyre_model import TyreModel
from domain.exceptions import ModelNotFoundError
from database.models.data_types.files import FileModel
from database.repositories.base_repository import BaseRepository

logger = logging.getLogger(__name__)


class TyreModelRepository(BaseRepository[TyreModel]):

    def __init__(self):
        super().__init__(TyreModel)


    def get_by_id(self, id_: int):
        tyre_model = (
            self.db.query(self.model)
            .filter(self.model.id == id_)
            .first()
        )

        if not tyre_model:
            logger.error(f"TyreModel with id {id_} not found")
            raise ModelNotFoundError(f"TyreModel with id {id_} not found")

        tyre_model.files = (
            self.db.query(File)
            .filter(
                File.model == FileModel.tyre_model,
                File.model_id == tyre_model.id
            )
            .all()
        )

        return tyre_model


