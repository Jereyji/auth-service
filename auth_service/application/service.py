import jwt
import bcrypt
from datetime import datetime, timedelta
from infrastructure.user_repository import UserRepository
from infrastructure.refresh_token_repository import RefreshTokenRepository
from domain.models.user import User 
from domain.models.refresh_token import RefreshToken 
from application.DTO import TokensResponse

ALGORITHM = "HS256"

class AuthService:
    def __init__(self, user_repository: UserRepository, refresh_token_repository: RefreshTokenRepository, secret_key: str):
        self._user_repository = user_repository
        self._refresh_token_repository = refresh_token_repository
        self._secret_key = secret_key

    def authenticate_user(self, username: str, password: str) -> TokensResponse:
        user = self._user_repository.get_user_by_username(username)
        if not user:
            raise Exception("User not found")

        if not user.check_password(password):
            raise Exception("Invalid credentials")

        access_token = jwt.encode({"sub": user.username}, self._secret_key, algorithm=ALGORITHM)
        refresh_token = RefreshToken(user.id)

        self._refresh_token_repository.create_refresh_token(refresh_token)

        return TokensResponse(access_token, refresh_token.refresh_token)

    def register_user(self, username: str, password: str) -> None:
        hashed_password = bcrypt.hashpw(password.encode("utf-8"), bcrypt.gensalt()).decode("utf-8")
        new_user = User(username, hashed_password)
        self._user_repository.create_user(new_user)

    def refresh_tokens(self, refresh_token_str: str) -> TokensResponse:
        refresh_token = self._refresh_token_repository.get_refresh_token(refresh_token_str)
        if not refresh_token or refresh_token.expired_at < datetime.now():
            raise Exception("Invalid or expired refresh token")

        user = self._user_repository.get_user_by_id(refresh_token.user_id)
        if not user:
            raise Exception("User not found")

        new_access_token = jwt.encode({"sub": user.username}, self._secret_key, algorithm=ALGORITHM)
        new_refresh_token_str = RefreshToken(user.id).refresh_token
        new_expired_at = datetime.now() + timedelta(days=14)

        self._refresh_token_repository.update_refresh_token(refresh_token_str, new_refresh_token_str, new_expired_at)

        return TokensResponse(new_access_token, new_refresh_token_str)
