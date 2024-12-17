from fastapi import FastAPI
from handlers.handler import create_auth_router
from application.service import AuthService
from infrastructure.database import get_db_connection
from infrastructure.repository import UserRepository
from pkg.secret_manager import SecretManager

secrets = SecretManager()

app = FastAPI()

def init_auth_service():
    connection = get_db_connection(secrets)
    user_repository = UserRepository(connection)
    return AuthService(user_repository)

auth_service = init_auth_service()

auth_router = create_auth_router(auth_service)

app.include_router(auth_router, prefix="/auth")
