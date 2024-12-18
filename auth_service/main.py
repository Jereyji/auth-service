from fastapi import FastAPI
from handlers.handler import create_auth_router
from infrastructure.database import get_db_connection
from application.service import AuthService
from infrastructure.user_repository import UserRepository
from infrastructure.refresh_token_repository import RefreshTokenRepository
from pkg.secret_manager import SecretManager

secrets = SecretManager()

app = FastAPI()

def init_auth_service():
    connection = get_db_connection(secrets)
    user_repository = UserRepository(connection)
    refresh_token_repository = RefreshTokenRepository(connection)
    return AuthService(user_repository, refresh_token_repository, secrets.secret_key)

auth_service = init_auth_service()

auth_router = create_auth_router(auth_service)

app.include_router(auth_router, prefix="/auth")
