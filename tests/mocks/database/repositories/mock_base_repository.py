
class MockBaseRepository:
    def __init__(self):
        self.create_calls = []
        self.create_error = None
        self.create_response = None

    def create(self, **kwargs):
        self.create_calls.append(kwargs)

        if self.create_error:
            raise self.create_error

        return self.create_response

    def reset(self):
        self.create_calls = []
        self.create_error = None
        self.create_response = None