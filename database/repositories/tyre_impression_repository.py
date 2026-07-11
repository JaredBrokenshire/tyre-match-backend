from database.models.file import File
from database.models.data_types.files import FileModel
from database.models.tyre_impression import TyreImpression
from database.repositories.base_repository import BaseRepository


class TyreImpressionRepository(BaseRepository[TyreImpression]):

    def __init__(self):
        super().__init__(TyreImpression)

    def get_by_id(self, id_: int):
        tyre_impression = (
            self.db.query(self.model)
            .filter(self.model.id == id_)
            .first()
        )

        if not tyre_impression:
            return None

        # Load files
        files = (
            self.db.query(File)
            .filter(
                File.model == FileModel.tyre_impression,
                File.model_id == tyre_impression.id
            )
            .all()
        )
        tyre_impression.files = {f.file_type.value: f for f in files}

        return tyre_impression


