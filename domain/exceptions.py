class FileUploadError(Exception):
    pass


class FileReadError(Exception):
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


class ProcessorError(Exception):
    pass


class PipelineError(ProcessorError):
    pass