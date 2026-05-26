import enum


class FileModel(enum.Enum):
    tyre_model = "tyre_model"
    tyre_impression = "tyre_impression"

class FileType(enum.Enum):
    original = "original"
    normalised = "normalised"
    enhanced = "enhanced"
    binary = "binary"
    clean = "clean"
    skeleton = "skeleton"
