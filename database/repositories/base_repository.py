from flask import current_app
from sqlalchemy import or_
from domain import DatabaseError
from database.extensions import db
from sqlalchemy.exc import IntegrityError
from typing import TypeVar, Generic, Type, List, Optional

T = TypeVar('T')


class BaseRepository(Generic[T]):
    """
    Generic repository providing basic CRUD operations
    """

    def __init__(self, model: Type[T]):
        self.db = db.session
        self.model = model

    def get_all(self, page_size: int = 20, page: int = 1, search: str = "") -> (List[T], int):
        query = self.db.query(self.model)

        if search != "":
            search_term = f"%{search}%"
            query = query.filter(
                or_(
                    self.model.manufacturer.ilike(search_term),
                    self.model.model_name.ilike(search_term),
                )
            )

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
            current_app.logger.error(f"Error creating record in db: {e}")
            raise DatabaseError("Error inserting record into DB")

        return entity

    def delete(self, entity_id: int) -> bool:
        entity = self.get_by_id(entity_id)
        if not entity:
            return False

        self.db.delete(entity)
        self.db.flush()
        return True
