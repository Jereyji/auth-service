import jwt
import bcrypt
from infrastructure.repository import UserRepository
from domain.models.user import User

SECRET_KEY = "SECRET_KEY"
ALGORITHM = "HS256"

class AuthService:
    def __init__(self, user_repository: UserRepository):
        self.user_repository = user_repository

    def authenticate_user(self, username: str, password: str) -> str:
        user = self.user_repository.get_user_by_username(username)
        if not user:
            raise Exception("User not found")

        if not bcrypt.checkpw(password.encode("utf-8"), user.password_hash.encode("utf-8")):
            raise Exception("Invalid credentials")

        token = jwt.encode({"sub": user.username}, SECRET_KEY, algorithm=ALGORITHM)
        return token

    def register_user(self, username: str, password: str) -> None:
        hashed_password = bcrypt.hashpw(password.encode("utf-8"), bcrypt.gensalt()).decode("utf-8")
        new_user = User(username, hashed_password)
        self.user_repository.create_user(new_user)
