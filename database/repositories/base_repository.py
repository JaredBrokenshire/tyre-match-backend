from flask import current_app
from database.extensions import db
from domain import DatabaseError, ModelNotFoundError
from sqlalchemy.exc import IntegrityError, DataError
from typing import TypeVar, Generic, Type, List, Optional

T = TypeVar('T')


class BaseRepository(Generic[T]):
    """
    Generic repository providing basic CRUD operations
    """

    def __init__(self, model: Type[T]):
        self.db = db.session
        self.model = model

    def get_all(self, page: int = 1, page_size: int = 20, filters=None) -> (List[T], int):
        query = self.db.query(self.model)

        if filters is not None:
            query = query.filter(filters)

        query = query.order_by(self.model.id)

        query = query.paginate(page=page, per_page=page_size, error_out=False, count=True)

        total_count = query.total
        res = query.items

        return res, total_count

    def get_by_id(self, entity_id: int) -> Optional[T]:
        return self.db.query(self.model).filter(self.model.id == entity_id).first()

    def create(self, **kwargs) -> T:
        entity = self.model(**kwargs)

        try:
            self.db.add(entity)
            self.db.flush()
        except IntegrityError as e:
            self.db.rollback()
            current_app.logger.error(f"Error creating record in database: {e}")
            raise DatabaseError(f"Error inserting record into database: {e}")

        return entity

    def update(self, entity: T, **kwargs) -> T:
        try:
            for key, value in kwargs.items():
                setattr(entity, key, value)

            self.db.flush()

            return entity
        except (IntegrityError, DataError) as e:
            self.db.rollback()
            current_app.logger.error(f"Error updating record in database: {e}")
            raise DatabaseError(f"Error updating record in database: {e}")


    def delete(self, entity_id: int) -> bool:
        entity = self.get_by_id(entity_id)
        if not entity:
            raise ModelNotFoundError(f"Model with id {entity_id} not found")

        self.db.delete(entity)
        self.db.flush()
        return True
