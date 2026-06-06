from database.models.file import File


def file_response(file: File):
    return {
        "id": file.id,
        "model": file.model.value,
        "model_id": file.model_id,
        "file_type": file.file_type.value,
        "file_name": file.file_name,
        "file_location": file.file_location,
        "mime_type": file.mime_type,
        "created_at": file.created_at,
        "updated_at": file.updated_at,
    }