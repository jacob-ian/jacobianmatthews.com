CREATE TABLE IF NOT EXISTS users(
    id VARCHAR PRIMARY KEY UNIQUE,
    name VARCHAR NOT NULL,
    email VARCHAR NOT NULL,
    email_verified BOOLEAN NOT NULL,
    image_url VARCHAR,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    deleted_at timestamptz
)