-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE IF NOT EXISTS users_beli (
  username TEXT PRIMARY KEY CHECK (char_length(username) BETWEEN 2 AND 30),
  password TEXT NOT NULL,
  email TEXT NOT NULL,
  isAdmin BOOLEAN NOT NULL DEFAULT FALSE
);

-- Unique constraint only for admin users
CREATE UNIQUE INDEX unique_admin_email ON users_beli (email) WHERE isAdmin = TRUE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP INDEX IF EXISTS unique_admin_email;
DROP TABLE IF EXISTS users_beli;
-- +goose StatementEnd
