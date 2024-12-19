from fastapi import FastAPI
from handlers.handler import create_auth_router
from infrastructure.database import get_db_connection
from application.service import AuthService
from infrastructure.user_repository import UserRepository
from infrastructure.refresh_token_repository import RefreshTokenRepository
from pkg.env_manager import Configs

configs = Configs()

app = FastAPI()

def init_auth_service():
    connection = get_db_connection(configs.postgres_config)
    user_repository = UserRepository(connection)
    refresh_token_repository = RefreshTokenRepository(connection)
    return AuthService(user_repository, refresh_token_repository, configs.secret_config)

auth_service = init_auth_service()

auth_router = create_auth_router(auth_service)

app.include_router(auth_router, prefix="/auth")
