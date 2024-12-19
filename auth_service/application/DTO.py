from pydantic import BaseModel
from datetime import datetime

class UserLoginRequest(BaseModel):
    username: str
    password: str

class TokensResponse(BaseModel):
    access_token: str
    refresh_token: str
    refresh_token_expiration: datetime

    class Config:
        arbitrary_types_allowed = True