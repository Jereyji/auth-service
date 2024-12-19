import jwt
from datetime import datetime, timedelta

class AccessToken:
    def __init__(self, username: str, secret_key: str, token_expiration: int):
        self.access_token_str = jwt.encode(
            {
                "username": username,
                "exp": datetime.now().replace(microsecond=0) + timedelta(minutes=token_expiration)
            },
            secret_key,
            algorithm="HS256"
        )