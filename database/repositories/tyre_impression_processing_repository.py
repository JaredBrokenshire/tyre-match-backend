from database.models import TyreImpressionProcessing
from database.repositories.base_repository import BaseRepository


class TyreImpressionProcessingRepository(BaseRepository[TyreImpressionProcessing]):

    def __init__(self):
        super().__init__(TyreImpressionProcessing)


