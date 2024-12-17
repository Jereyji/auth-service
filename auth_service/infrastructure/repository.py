from psycopg2 import sql
from domain.models.user import User

class UserRepository:
    def __init__(self, db_connection):
        self.conn = db_connection

    def create_user(self, user: User) -> None:
        with self.conn.cursor() as cursor:
            query = sql.SQL("INSERT INTO users (id, username, password_hash) VALUES (%s, %s, %s)")
            cursor.execute(query, (user.id, user.username, user.password_hash))
            self.conn.commit()

    def get_user_by_username(self, username: str) -> User | None:
        with self.conn.cursor() as cursor:
            query = sql.SQL("SELECT id, username, password_hash FROM users WHERE username = %s")
            cursor.execute(query, (username,))
            row = cursor.fetchone()
            if row:
                return User(user_id=row["id"], username=row["username"], password_hash=row["password_hash"])
        return None
