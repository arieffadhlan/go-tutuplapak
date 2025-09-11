-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS files (
  id bigserial primary key,
  url text NOT NULL,
  thumbnail_url text NOT NULL,
  created_at timestamptz NOT NULL default current_timestamp
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS files
-- +goose StatementEnd
