CREATE TABLE IF NOT EXISTS users(
    uuid UUID PRIMARY KEY,
    email TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS used_refresh_tokens (
    signature varchar(100)
);