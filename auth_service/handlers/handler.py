from fastapi import APIRouter, HTTPException, Response, Depends, Request
from fastapi.responses import JSONResponse
from application.service import AuthService
from application.DTO import UserLoginRequest, TokensResponse


def create_auth_router(auth_service: AuthService):
    router = APIRouter()

    @router.post("/register")
    def register(request: UserLoginRequest):
        try:
            auth_service.register_user(request.username, request.password)
            return {"detail": "User registered successfully"}
        except Exception as e:
            raise HTTPException(status_code=400, detail=str(e))

    @router.post("/login")
    def login(request: UserLoginRequest):
        try:
            tokens = auth_service.authenticate_user(request.username, request.password)

            response = JSONResponse(
                {"access_token": tokens.access_token, "token_type": "bearer"}
            )

            response.set_cookie(
                key="refresh_token",
                value=tokens.refresh_token,
                path="/auth",
                httponly=True,
                secure=True,
                samesite="Strict",
                expires=tokens.refresh_token_expiration  # Исправить!!!
            )

            return response
        except Exception as e:
            raise HTTPException(status_code=401, detail=str(e))

    def get_refresh_token_from_cookie(request: Request) -> str:
        refresh_token = request.cookies.get("refresh_token")
        if not refresh_token:
            raise HTTPException(status_code=401, detail="Refresh token is missing")
        return refresh_token
    
    @router.post("/refresh", response_model=TokensResponse)
    def refresh(response: Response, refresh_token: str = Depends(get_refresh_token_from_cookie)): # NEED FIX
        try:
            tokens = auth_service.refresh_tokens(refresh_token)

            response = JSONResponse(
                {"access_token": tokens.access_token, "token_type": "bearer"}
            )

            response.set_cookie(
                key="refresh_token",
                value=tokens.refresh_token,
                path="/auth",
                httponly=True,
                secure=True,
                samesite="Strict",
                expires=tokens.expired_at
            )

            return response
        except Exception as e:
            raise HTTPException(status_code=401, detail=str(e))

    return router
