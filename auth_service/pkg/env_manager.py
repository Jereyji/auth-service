import os
from dotenv import load_dotenv

import os
from dotenv import load_dotenv

class SecretConfig:
    def __init__(self):
        self.secret_key = os.getenv("SECRET_KEY")
        self.access_token_expiration = int(os.getenv("ACCESS_TOKEN_EXPIRATION", 30))
        self.refresh_token_expiration = int(os.getenv("REFRESH_TOKEN_EXPIRATION", 14))

        if not all([self.secret_key, 
                    self.access_token_expiration, 
                    self.refresh_token_expiration]):
            raise ValueError("One or more secret-related environment variables are missing")

class PostgresConfig:
    def __init__(self):
        self.postgres_user = os.getenv("POSTGRES_USER")
        self.postgres_password = os.getenv("POSTGRES_PASSWORD")
        self.postgres_database = os.getenv("POSTGRES_DB")
        self.postgres_host = os.getenv("POSTGRES_HOST")
        self.postgres_port = os.getenv("POSTGRES_PORT")

        if not all([self.postgres_user, 
                    self.postgres_password, 
                    self.postgres_database, 
                    self.postgres_host,
                    self.postgres_port]):
            raise ValueError("One or more PostgreSQL environment variables are missing")

class Configs:
    def __init__(self):
        load_dotenv()
        self.secret_config = SecretConfig()
        self.postgres_config = PostgresConfig()
