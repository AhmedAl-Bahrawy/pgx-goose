-- +goose Up
ALTER TABLE users
ADD COLUMN password_hash VARCHAR(50);

-- +goose Down
ALTER TABLE users
DROP COLUMN password_hash;
