-- +goose Up
-- bcrypt hashes are always 60 chars, VARCHAR(50) rejects them
ALTER TABLE users
ALTER COLUMN password_hash TYPE VARCHAR(255);

-- +goose Down
ALTER TABLE users
ALTER COLUMN password_hash TYPE VARCHAR(50);
