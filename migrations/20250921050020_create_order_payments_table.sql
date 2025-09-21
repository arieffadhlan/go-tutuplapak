-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE order_payments (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID REFERENCES orders(id),
    seller_id UUID REFERENCES users(id),
    bank_account_name VARCHAR(128) NOT NULL,
    bank_account_holder VARCHAR(128) NOT NULL,
    bank_account_number VARCHAR(64) NOT NULL,
    total_price INT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_order_payments_order_id ON order_payments(order_id);
CREATE INDEX idx_order_payments_seller_id ON order_payments(seller_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DROP INDEX IF EXISTS idx_order_payments_seller_id;
DROP INDEX IF EXISTS idx_order_payments_order_id;

DROP TABLE IF EXISTS order_payments;
-- +goose StatementEnd