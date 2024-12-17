import os
from dotenv import load_dotenv

class SecretManager:
    def __init__(self):
        load_dotenv()

        self.postgres_user = os.getenv("POSTGRES_USER")
        self.postgres_password = os.getenv("POSTGRES_PASSWORD")
        self.postgres_database = os.getenv("POSTGRES_DB")
        self.postgres_host = os.getenv("POSTGRES_HOST")
        self.postgres_port = os.getenv("POSTGRES_PORT")
        self.secret_key = os.getenv("SECRET_KEY")

        if not all([self.postgres_user, self.postgres_password, self.postgres_database, self.secret_key]):
            raise ValueError("One or more required environment variables are missing")
