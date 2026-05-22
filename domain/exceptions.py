class FileUploadError(Exception):
    pass


class InvalidFileTypeError(FileUploadError):
    pass


class FileSaveError(FileUploadError):
    pass


class DatabaseError(Exception):
    pass
