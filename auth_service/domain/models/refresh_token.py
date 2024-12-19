import uuid
import secrets
from datetime import datetime, timedelta

class RefreshToken:
    def __init__(self, user_id: str, adding_time: int):
        self.id = str(uuid.uuid4())
        self.user_id = user_id
        self.refresh_token_str: str
        self.expired_at: datetime

        self.generate_refresh_token()
        self.generate_expired_at(adding_time)

    def generate_refresh_token(self):
        self.refresh_token_str = secrets.token_urlsafe(64)

    def generate_expired_at(self, adding_time: int):
        self.expired_at = datetime.now().replace(microsecond=0) + timedelta(days=adding_time)