CREATE TABLE IF NOT EXISTS users
(
    id           UUID PRIMARY KEY,
    username     VARCHAR(100) NOT NULL UNIQUE,
    password     VARCHAR(100) NOT NULL,
    access_level INT CHECK (access_level IN (1, 2)) -- 1: Moderator, 2: Client
);

CREATE TABLE IF NOT EXISTS refresh_sessions
(
    id            UUID PRIMARY KEY,
    user_id       UUID         NOT NULL REFERENCES users (id),
    refresh_token VARCHAR(100) NOT NULL UNIQUE,
    created_at    TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expired_at    TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);
