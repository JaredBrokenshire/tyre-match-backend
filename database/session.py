import os
from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker, declarative_base

env_variable = "DATABASE_URL"
if os.getenv("TESTING") == "1":
    env_variable = "TEST_DATABASE_URL"

DATABASE_URL = os.getenv(env_variable)

engine = create_engine(
    DATABASE_URL,
    pool_pre_ping=True,
    pool_recycle=3600,
)

SessionLocal = sessionmaker(
    autocommit=False,
    autoflush=False,
    bind=engine
)

Base = declarative_base()