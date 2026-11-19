
CREATE TABLE IF NOT EXISTS orders (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES demo_users(id),
    trading_pair_id BIGINT NOT NULL REFERENCES trading_pairs(id),

    side VARCHAR(4) NOT NULL
        CHECK (side IN ('BUY', 'SELL')),

    order_type VARCHAR(10) NOT NULL DEFAULT 'LIMIT'
        CHECK (order_type = 'LIMIT'),

    price NUMERIC(30,10) NOT NULL CHECK (price > 0),
    quantity NUMERIC(30,10) NOT NULL CHECK (quantity > 0),

    remaining_quantity NUMERIC(30,10) NOT NULL
        CHECK (remaining_quantity >= 0),

    status VARCHAR(20) NOT NULL DEFAULT 'OPEN'
        CHECK (status IN (
            'OPEN', 'PARTIALLY_FILLED',
            'FILLED', 'CANCELLED'
        )),

    reserved_amount NUMERIC(30,10) NOT NULL
        CHECK (reserved_amount >= 0),

    client_order_id VARCHAR(100) NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (user_id, client_order_id),
    CHECK (remaining_quantity <= quantity)
);

CREATE INDEX IF NOT EXISTS idx_orders_open
ON orders (trading_pair_id, side, price, created_at, id)
WHERE status IN ('OPEN', 'PARTIALLY_FILLED');

CREATE TABLE IF NOT EXISTS wallet_movements (
    id BIGSERIAL PRIMARY KEY,
    wallet_id BIGINT NOT NULL REFERENCES wallets(id),
    order_id BIGINT NOT NULL REFERENCES orders(id),

    movement_type VARCHAR(20) NOT NULL
        CHECK (movement_type IN ('RESERVE', 'RELEASE')),

    available_delta NUMERIC(30,10) NOT NULL,
    locked_delta NUMERIC(30,10) NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (order_id, movement_type),
    CHECK (available_delta + locked_delta = 0)
);
