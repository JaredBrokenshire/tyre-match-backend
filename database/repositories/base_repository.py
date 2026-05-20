from sqlalchemy.orm import Session
from sqlalchemy.exc import IntegrityError
from typing import TypeVar, Generic, Type, List, Optional

T = TypeVar('T')

class BaseRepository(Generic[T]):
    """
    Generic repository providing basic CRUD operations
    """

    def __init__(self, db: Session, model: Type[T]):
        self.db = db
        self.model = model

    def get_all(self) -> List[T]:
        return self.db.query(self.model).all()

    def get_by_id(self, entity_id: int) -> Optional[T]:
        return self.db.query(self.model).filter(self.model.id == entity_id).first()

    def create(self, **kwargs) -> T:
        entity = self.model(**kwargs)
        self.db.add(entity)

        try:
            self.db.flush()
        except IntegrityError:
            self.db.rollback()
            raise

        return entity

    def delete(self, entity_id: int) -> bool:
        entity = self.get_by_id(entity_id)
        if not entity:
            return False

        self.db.delete(entity)
        self.db.flush()
        return True