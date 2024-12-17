import uuid

class User:
    def __init__(self, username: str, password_hash: str):
        self.id = str(uuid.uuid4())
        self.username = username
        self.password_hash = password_hash
