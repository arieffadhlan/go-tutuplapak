-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE order_items (
  id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  order_id UUID REFERENCES orders(id),
  product_id UUID REFERENCES products(id),
  name VARCHAR(32) NOT NULL DEFAULT '',
  sku VARCHAR(32) NOT NULL DEFAULT '',
  qty INT NOT NULL DEFAULT 0,
  price INT NOT NULL DEFAULT 0,
  file_uri TEXT NOT NULL DEFAULT '',
  file_thumbnail_uri TEXT NOT NULL DEFAULT '',
  category VARCHAR(10) NOT NULL,
  purchase_qty INT NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_order_items_product_id ON order_items(product_id)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS order_items;
-- +goose StatementEnd