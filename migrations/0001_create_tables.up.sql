CREATE TABLE IF NOT EXISTS users(
    uuid UUID PRIMARY KEY,
    email TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS sessions(
  id SERIAL PRIMARY KEY,
  user_id UUID references users (uuid),
  refresh_token VARCHAR(80),
  expired_at_refresh TIMESTAMP without time zone,
  ip VARCHAR(15) NOT NULL
);
