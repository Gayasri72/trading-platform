
CREATE TABLE IF NOT EXISTS demo_users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS wallets (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES demo_users(id),
    asset_id BIGINT NOT NULL REFERENCES assets(id),
    available NUMERIC(30,10) NOT NULL DEFAULT 0,
    locked NUMERIC(30,10) NOT NULL DEFAULT 0,
    CHECK (available >= 0),
    CHECK (locked >= 0),
    UNIQUE (user_id, asset_id)
);

CREATE TABLE IF NOT EXISTS wallet_ledger (
    id BIGSERIAL PRIMARY KEY,
    wallet_id BIGINT NOT NULL REFERENCES wallets(id),
    amount NUMERIC(30,10) NOT NULL CHECK (amount <> 0),
    entry_type VARCHAR(30) NOT NULL,
    reference_id VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO demo_users (id, username)
VALUES (1, 'demo_trader')
ON CONFLICT (id) DO NOTHING;

INSERT INTO wallets (user_id, asset_id)
SELECT 1, id
FROM assets
WHERE symbol IN ('BTC', 'USDT')
ON CONFLICT (user_id, asset_id) DO NOTHING;
