import os
from flask import g
from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker, declarative_base

env_variable = "DATABASE_URL"
if os.getenv("TESTING") == "1":
    env_variable = "TEST_DATABASE_URL"

DATABASE_URL = os.getenv(env_variable)

def get_engine():
    return create_engine(
        DATABASE_URL,
        pool_pre_ping=True,
        pool_recycle=3600,
    )

SessionLocal = sessionmaker(
    autocommit=False,
    autoflush=False
)

def get_db():
    if "db" not in g:
        g.db = SessionLocal(bind=get_engine())
    return g.db

def close_db(exception=None):
    db = g.pop("db", None)
    if db is not None:
        if exception:
            db.rollback()
        else:
            db.commit()
        db.close()

Base = declarative_base()