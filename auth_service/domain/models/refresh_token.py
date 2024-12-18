import uuid
import secrets
from datetime import datetime, timedelta

ADDING_TIME=14

class RefreshToken:
    def __init__(self, user_id: str):
        self.id = str(uuid.uuid4())
        self.user_id = user_id
        self.refresh_token = self._generate_refresh_token()
        self.expired_at = self._generate_expired_at()

    def _generate_refresh_token(self) -> str:
        return secrets.token_urlsafe(64)

    def _generate_expired_at(self) -> datetime:
        return datetime.now().replace(microsecond=0) + timedelta(days=ADDING_TIME)