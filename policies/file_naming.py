import uuid
from werkzeug.utils import secure_filename


def uuid_filename(file) -> (str, str):
    id_ = uuid.uuid4()
    return id_, f"{id_}.{file.filename.split('.')[-1]}"


def original_filename(file) -> str:
    return secure_filename(file.filename)


def prefixed_filename(prefix: str, file) -> str:
    return f"{prefix}_{secure_filename(file.filename)}"