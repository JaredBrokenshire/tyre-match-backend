import os
from sqlalchemy.orm import declarative_base

DATABASE_URL = os.getenv("DATABASE_URL")

if os.getenv("TESTING") == "1":
    DATABASE_URL = os.getenv("TEST_DATABASE_URL")

Base = declarative_base()