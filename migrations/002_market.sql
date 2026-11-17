
-- 1. Supported crypto assets
CREATE TABLE IF NOT EXISTS assets (
    id BIGSERIAL PRIMARY KEY,
    symbol VARCHAR(10) NOT NULL UNIQUE,
    name VARCHAR(50) NOT NULL,
    decimal_places INTEGER NOT NULL CHECK (decimal_places BETWEEN 0 AND 18),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 2. Trading pairs
CREATE TABLE IF NOT EXISTS trading_pairs (
    id BIGSERIAL PRIMARY KEY,
    symbol VARCHAR(25) NOT NULL UNIQUE,
    base_asset_id BIGINT NOT NULL REFERENCES assets(id),
    quote_asset_id BIGINT NOT NULL REFERENCES assets(id),
    price_precision INTEGER NOT NULL CHECK (price_precision BETWEEN 0 AND 10),
    quantity_precision INTEGER NOT NULL CHECK (quantity_precision BETWEEN 0 AND 18),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    CHECK (base_asset_id <> quote_asset_id)
);

-- 3. Simulated reference market prices
CREATE TABLE IF NOT EXISTS market_prices (
    trading_pair_id BIGINT PRIMARY KEY
        REFERENCES trading_pairs(id),
    reference_price NUMERIC(30, 10) NOT NULL
        CHECK (reference_price > 0),
    source VARCHAR(30) NOT NULL DEFAULT 'SIMULATED',
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Insert supported assets
INSERT INTO assets (symbol, name, decimal_places)
VALUES
    ('BTC', 'Bitcoin', 8),
    ('USDT', 'Tether', 6)
ON CONFLICT (symbol) DO NOTHING;

-- Create BTC/USDT trading pair
INSERT INTO trading_pairs (
    symbol, base_asset_id, quote_asset_id,
    price_precision, quantity_precision
)
SELECT
    'BTC-USDT', btc.id, usdt.id, 2, 8
FROM assets btc
CROSS JOIN assets usdt
WHERE btc.symbol = 'BTC'
  AND usdt.symbol = 'USDT'
ON CONFLICT (symbol) DO NOTHING;

-- Initial simulated reference price
INSERT INTO market_prices (trading_pair_id, reference_price)
SELECT id, 90000.00
FROM trading_pairs
WHERE symbol = 'BTC-USDT'
ON CONFLICT (trading_pair_id) DO NOTHING;
