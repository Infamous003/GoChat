CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    password_hash bytea NOT NULL,
    version integer NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ(0) NOT NULL DEFAULT NOW()
);