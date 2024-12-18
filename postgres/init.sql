CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    refresh_token VARCHAR(255) NOT NULL,
    expired_at TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_expired_at ON refresh_tokens (expired_at);
