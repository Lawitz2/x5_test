-- +goose Up
-- +goose StatementBegin
CREATE TYPE order_status AS ENUM (
    'NEW',
    'PROCESSING',
    'FULFILLED',
    'FAILED'
);

CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY,
    customer_id VARCHAR(255) NOT NULL,
    items       JSONB NOT NULL,
    status      order_status NOT NULL DEFAULT 'NEW',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- индексы для ускорения получения индивидуальных заказов и списков заказов
CREATE INDEX idx_orders_customer_id ON orders(customer_id);
CREATE INDEX idx_orders_status ON orders(status);

-- автообновление updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_orders_updated_at
    BEFORE UPDATE ON orders
    FOR EACH ROW
    EXECUTE PROCEDURE update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
DROP TYPE IF EXISTS order_status;
-- +goose StatementEnd
