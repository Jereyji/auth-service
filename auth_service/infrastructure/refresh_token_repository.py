from psycopg2 import sql
from domain.models.refresh_token import RefreshToken
from datetime import datetime

class RefreshTokenRepository:
    def __init__(self, db_connection):
        self.conn = db_connection

    def create_refresh_token(self, token: RefreshToken) -> None:
        try:
            with self.conn.cursor() as cursor:
                query = sql.SQL("""
                    INSERT INTO refresh_tokens (id, user_id, refresh_token, expired_at)
                    VALUES (%s, %s, %s, %s)
                """)
                cursor.execute(query, (token.id, token.user_id, token.refresh_token, token.expired_at))
                self.conn.commit()
        except Exception as e:
            self.conn.rollback()
            raise

    def get_refresh_token(self, refresh_token: str) -> RefreshToken | None:
        with self.conn.cursor() as cursor:
            query = sql.SQL("""
                SELECT id, user_id, refresh_token, expired_at
                FROM refresh_tokens
                WHERE refresh_token = %s
            """)
            cursor.execute(query, (refresh_token,))
            row = cursor.fetchone()
            if row:
                return RefreshToken(
                    user_id=row["user_id"],
                    refresh_token=row["refresh_token"],
                    expired_at=row["expired_at"]
                )
        return None
    
    def update_refresh_token(self, old_token: str, new_token: str, new_expired_at: datetime) -> None:
        with self.conn.cursor() as cursor:
            query = """
                UPDATE refresh_tokens
                SET refresh_token = %s, expired_at = %s
                WHERE refresh_token = %s
            """
            cursor.execute(query, (new_token, new_expired_at, old_token))
            self.conn.commit()

    def delete_refresh_token(self, refresh_token: str) -> None:
        with self.conn.cursor() as cursor:
            query = sql.SQL("DELETE FROM refresh_tokens WHERE refresh_token = %s")
            cursor.execute(query, (refresh_token,))
            self.conn.commit()

    def delete_expired_tokens(self) -> None:
        with self.conn.cursor() as cursor:
            query = sql.SQL("DELETE FROM refresh_tokens WHERE expired_at < %s")
            cursor.execute(query, (datetime.now(),))
            self.conn.commit()
