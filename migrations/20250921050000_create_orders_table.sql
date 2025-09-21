-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_type WHERE typname = 'sender_type'
    ) THEN
        CREATE TYPE sender_type AS ENUM ('email', 'phone');
    END IF;
END$$;

CREATE TABLE orders (
  id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  total_price INT NOT NULL DEFAULT 0,
  sender_name VARCHAR(55) NOT NULL,
  sender_contact_detail VARCHAR(55) NOT NULL,
  sender_contact_type sender_type NOT NULL DEFAULT 'email',
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd