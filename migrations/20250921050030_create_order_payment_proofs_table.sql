-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE order_payment_proofs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders(id),
    file_id UUID NOT NULL REFERENCES files(id),
    created_at TIMESTAMPTZ DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS order_payment_proofs;
-- +goose StatementEnd