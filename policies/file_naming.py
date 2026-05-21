import uuid
from werkzeug.utils import secure_filename


def uuid_filename(file) -> (str, str):
    id = uuid.uuid4()
    return id, f"{id}_{secure_filename(file.filename)}"


def original_filename(file) -> str:
    return secure_filename(file.filename)


def prefixed_filename(prefix: str, file) -> str:
    return f"{prefix}_{secure_filename(file.filename)}"