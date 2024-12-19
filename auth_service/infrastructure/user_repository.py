from psycopg2 import sql
from domain.models.user import User

class UserRepository:
    def __init__(self, db_connection):
        self.conn = db_connection

    def create_user(self, user: User) -> None:
        try:
            with self.conn.cursor() as cursor:
                query = sql.SQL("INSERT INTO users (id, username, hashed_password) VALUES (%s, %s, %s)")
                cursor.execute(query, (user.id, user.username, user.hashed_password))
                self.conn.commit()
        except Exception as e:
            self.conn.rollback()
            raise

    def get_user_by_username(self, username: str) -> User | None:
        with self.conn.cursor() as cursor:
            query = sql.SQL("SELECT id, username, hashed_password FROM users WHERE username = %s")
            cursor.execute(query, (username,))
            row = cursor.fetchone()
            if row:
                return User(user_id=row["id"], username=row["username"], hashed_password=row["hashed_password"])
        return None

    def get_user_by_id(self, id: str) -> User | None:
        with self.conn.cursor() as cursor:
            query = sql.SQL("SELECT id, username, hashed_password FROM users WHERE id = %s")
            cursor.execute(query, (id,))
            row = cursor.fetchone()
            if row:
                return User(user_id=row["id"], username=row["username"], hashed_password=row["hashed_password"])
        return None

    def update_user(self, user: User) -> None:
        with self.conn.cursor() as cursor:
            query = sql.SQL("UPDATE users SET username = %s, hashed_password = %s WHERE id = %s")
            cursor.execute(query, (user.username, user.hashed_password, user.id))
            self.conn.commit()

    def delete_user(self, username: str) -> None:
        with self.conn.cursor() as cursor:
            query = sql.SQL("DELETE FROM users WHERE username = %s")
            cursor.execute(query, (username,))
            self.conn.commit()
