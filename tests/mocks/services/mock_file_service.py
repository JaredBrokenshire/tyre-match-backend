class MockFileService:
    def __init__(self):
        self.allowed_file_calls = []
        self.allowed_file_response = True

        self.save_file_calls = []
        self.save_file_response = ""
        self.save_file_error = None

        self.reset()

    def save_file(self, file, upload_dir: str, _: list[str]) -> str:
        self.save_file_calls.append((file.filename, upload_dir))

        if self.save_file_error:
            raise self.save_file_error

        return self.save_file_response

    def reset(self):
        self.allowed_file_calls = []
        self.allowed_file_response = True

        self.save_file_calls = []
        self.save_file_response = ""
        self.save_file_error = None
