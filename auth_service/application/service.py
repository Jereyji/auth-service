from datetime import datetime, timezone
from infrastructure.user_repository import UserRepository
from infrastructure.refresh_token_repository import RefreshTokenRepository
from domain.models.user import User 
from domain.models.access_token import AccessToken
from domain.models.refresh_token import RefreshToken
from application.DTO import TokensResponse
from pkg.env_manager import SecretConfig

class AuthService:
    def __init__(self, user_repository: UserRepository, refresh_token_repository: RefreshTokenRepository, secret_config: SecretConfig):
        self._user_repository = user_repository
        self._refresh_token_repository = refresh_token_repository
        self._secret_config = secret_config

    def register_user(self, username: str, password: str) -> None:
        new_user = User(username, User.generate_hashed_password(password))
        self._user_repository.create_user(new_user)

    def authenticate_user(self, username: str, password: str) -> TokensResponse:
        user = self._user_repository.get_user_by_username(username)
        if not user:
            raise Exception("User not found")

        if not user.check_password(password):
            raise Exception("Invalid credentials")

        access_token = AccessToken(user.username, self._secret_config.secret_key, self._secret_config.access_token_expiration)

        refresh_token = RefreshToken(user.id, self._secret_config.refresh_token_expiration)

        self._refresh_token_repository.create_refresh_token(refresh_token)

        return TokensResponse(
            access_token = access_token.access_token_str,
            refresh_token = refresh_token.refresh_token_str,
            refresh_token_expiration = refresh_token.expired_at.astimezone(timezone.utc)
        )

    def refresh_tokens(self, refresh_token_str: str) -> TokensResponse:
        refresh_token = self._refresh_token_repository.get_refresh_token(refresh_token_str)
        if not refresh_token or refresh_token.expired_at < datetime.now():
            raise Exception("Invalid or expired refresh token")

        user = self._user_repository.get_user_by_id(refresh_token.user_id)
        if not user:
            raise Exception("User not found")

        new_access_token = jwt.encode({"username": user.username}, self._secret_key, algorithm=ALGORITHM)
        refresh_token.refresh_token_str = refresh_token.new_refresh_token()
        refresh_token.expired_at = refresh_token.generate_expired_at()

        self._refresh_token_repository.update_refresh_token(refresh_token_str, refresh_token)

        return TokensResponse(
            access_token = new_access_token, 
            refresh_token = refresh_token.refresh_token_str)
