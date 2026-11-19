
BEGIN;

WITH inserted_entry AS (
    INSERT INTO wallet_ledger (
        wallet_id,
        amount,
        entry_type,
        reference_id
    )
    SELECT w.id, 10000.00, 'DEMO_CREDIT', 'initial-usdt-1'
    FROM wallets w
    JOIN assets a ON a.id = w.asset_id
    WHERE w.user_id = 1 AND a.symbol = 'USDT'
    ON CONFLICT (reference_id) DO NOTHING
    RETURNING wallet_id, amount
)
UPDATE wallets w
SET available = w.available + e.amount
FROM inserted_entry e
WHERE w.id = e.wallet_id;

COMMIT;
