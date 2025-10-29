CREATE TABLE users(
    id SERIAL PRIMARY KEY,
    login VARCHAR(64) pass_hash VARCHAR(255)
);
CREATE INDEX idx_user_id ON users(id);