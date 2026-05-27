from database.models.data_types import FileModel
from database.models import TyreImpressionProcessing, File
from database.repositories.base_repository import BaseRepository


class TyreImpressionProcessingRepository(BaseRepository[TyreImpressionProcessing]):

    def __init__(self):
        super().__init__(TyreImpressionProcessing)


    def get_by_id(self, id_: int):
        tyre_impression_processing = (
            self.db.query(self.model)
            .filter(self.model.id == id_)
            .first()
        )

        if not tyre_impression_processing:
            return None

        files = (
            self.db.query(File)
            .filter(
                File.model == FileModel.tyre_impression,
                File.model_id == tyre_impression_processing.id
            )
            .all()
        )

        tyre_impression_processing.files = {f.file_type.value: f for f in files}

        return tyre_impression_processing


    def get_by_tyre_impression_id(self, tyre_impression_id: int) -> TyreImpressionProcessing:
        return self.db.query(self.model).filter(self.model.tyre_impression_id == tyre_impression_id).first()


