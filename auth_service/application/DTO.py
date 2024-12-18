from pydantic import BaseModel

class UserLoginRequest(BaseModel):
    username: str
    password: str

class TokensResponse(BaseModel):
    access_token: str
    refresh_token: str