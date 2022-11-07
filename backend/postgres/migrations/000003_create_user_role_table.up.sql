CREATE TABLE IF NOT EXISTS user_role(
    id UUID PRIMARY KEY UNIQUE,
    user_id VARCHAR,
    role_id UUID,
    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
            REFERENCES users(id)
            ON DELETE CASCADE,
    CONSTRAINT fk_role
        FOREIGN KEY(role_id)
            REFERENCES roles(id)
            ON DELETE SET NULL
)
