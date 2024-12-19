import uuid
import bcrypt

class User:
    def __init__(self, username: str, hashed_password: str, user_id: str = None):
        self.id = user_id if user_id else str(uuid.uuid4())
        self.username = username
        self.hashed_password = hashed_password

    def generate_hashed_password(password: str) -> str:
        return bcrypt.hashpw(password.encode("utf-8"), bcrypt.gensalt()).decode("utf-8")

    def check_password(self, password: str) -> bool:
        return bcrypt.checkpw(password.encode("utf-8"), self.hashed_password.encode("utf-8"))
