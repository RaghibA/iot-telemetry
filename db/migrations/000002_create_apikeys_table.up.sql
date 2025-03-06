CREATE TABLE api_keys (
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    api_key TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (api_key)
);