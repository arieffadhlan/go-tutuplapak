-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    public_id UUID NOT NULL UNIQUE,
    email VARCHAR(255) UNIQUE,
    phone VARCHAR(255) UNIQUE,
    password VARCHAR(255) NOT NULL,
    file_id VARCHAR(255),
    file_uri VARCHAR(255),
    file_thumbnail_uri VARCHAR(255),
    bank_account_name VARCHAR(255),
    bank_account_holder VARCHAR(255),
    bank_account_number VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users (email);
CREATE INDEX idx_users_phone ON users (phone);
CREATE INDEX idx_users_public_id ON users (public_id);
CREATE INDEX idx_users_file_id ON users (file_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists users;
SELECT 'down SQL query';
-- +goose StatementEnd
