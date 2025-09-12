-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TYPE product_category AS ENUM (
  'Food', 
  'Tools', 
  'Clothes', 
  'Beverage', 
  'Furniture'
);

CREATE TABLE products (
  id SERIAL PRIMARY KEY,
  name VARCHAR(32) NOT NULL DEFAULT '',
  user_id INT REFERENCES users(id) ON DELETE CASCADE,
  file_id INT REFERENCES files(id),
  sku VARCHAR(32) NOT NULL DEFAULT '',
  qty INT NOT NULL DEFAULT 0,
  price INT NOT NULL DEFAULT 0,
  file_uri TEXT NOT NULL DEFAULT '',
  file_thumbnail_uri TEXT NOT NULL DEFAULT '',
  category product_category NOT NULL DEFAULT 'Food',
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_products_user_id ON products(user_id);
CREATE INDEX idx_products_file_id ON products(file_id);
CREATE INDEX idx_products_price ON products(price);
CREATE INDEX idx_products_category ON products(category);
CREATE INDEX idx_products_created_at ON products(created_at);
CREATE UNIQUE INDEX idx_products_sku ON products(sku);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DROP INDEX IF EXISTS idx_products_sku;
DROP INDEX IF EXISTS idx_products_price;
DROP INDEX IF EXISTS idx_products_category;
DROP INDEX IF EXISTS idx_products_created_at;
DROP INDEX IF EXISTS idx_products_file_id;
DROP INDEX IF EXISTS idx_products_user_id;

DROP TABLE IF EXISTS products;
DROP TYPE  IF EXISTS product_category;
-- +goose StatementEnd
