import psycopg2
from psycopg2.extras import RealDictCursor
from pkg.env_manager import PostgresConfig

def get_db_connection(secrets: PostgresConfig):
    DATABASE_URL = (
        f"postgresql://{secrets.postgres_user}:"
        f"{secrets.postgres_password}@"
        f"{secrets.postgres_host}:"
        f"{secrets.postgres_port}/"
        f"{secrets.postgres_database}"
    )

    conn = psycopg2.connect(DATABASE_URL, cursor_factory=RealDictCursor)
    return conn
