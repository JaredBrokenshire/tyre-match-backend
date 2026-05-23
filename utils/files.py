def allowed_file(filename: str, valid_extensions: list[str]) -> bool:
    return "." in filename and filename.rsplit(".", 1)[1] in valid_extensions
