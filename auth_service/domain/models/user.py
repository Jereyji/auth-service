import uuid
import bcrypt

class User:
    def __init__(self, username: str, password_hash: str, user_id: str = None):
        self.id = user_id if user_id else str(uuid.uuid4())
        self.username = username
        self.password_hash = password_hash

    def check_password(self, password: str) -> bool:
        return bcrypt.checkpw(password.encode("utf-8"), self.password_hash.encode("utf-8"))
