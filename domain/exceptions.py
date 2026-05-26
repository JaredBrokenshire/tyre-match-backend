class FileUploadError(Exception):
    pass


class InvalidFileError(FileUploadError):
    pass


class InvalidFileTypeError(FileUploadError):
    pass


class FileSaveError(FileUploadError):
    pass


class DatabaseError(Exception):
    pass


class ModelNotFoundError(DatabaseError):
    pass
