from fastapi import APIRouter, HTTPException
from application.service import AuthService
from application.DTO import UserLoginRequest, TokenResponse

def create_auth_router(auth_service: AuthService):
    router = APIRouter()

    @router.post("/login", response_model=TokenResponse)
    def login(request: UserLoginRequest):
        try:
            token = auth_service.authenticate_user(request.username, request.password)
            return {"access_token": token, "token_type": "bearer"}
        except Exception as e:
            raise HTTPException(status_code=401, detail=str(e))

    @router.post("/register")
    def register(request: UserLoginRequest):
        try:
            auth_service.register_user(request.username, request.password)
            return {"detail": "User registered successfully"}
        except Exception as e:
            raise HTTPException(status_code=400, detail=str(e))

    return router
