package wallet

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (repo *Repository) GetBalances(
	ctx context.Context,
	userID int64,
) ([]Balance, error) {
	rows, err := repo.db.Query(ctx, `
		SELECT a.symbol,
		       w.available::text,
		       w.locked::text,
		       (w.available + w.locked)::text
		FROM wallets w
		JOIN assets a ON a.id = w.asset_id
		WHERE w.user_id = $1
		ORDER BY a.symbol
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	balances := make([]Balance, 0)
	for rows.Next() {
		var b Balance
		if err := rows.Scan(
			&b.Asset, &b.Available, &b.Locked, &b.Total,
		); err != nil {
			return nil, err
		}
		balances = append(balances, b)
	}
	return balances, rows.Err()
}

func (repo *Repository) GetLedger(
	ctx context.Context,
	userID int64,
) ([]LedgerEntry, error) {
	rows, err := repo.db.Query(ctx, `
		SELECT l.id, a.symbol, l.amount::text,
		       l.entry_type, l.reference_id, l.created_at
		FROM wallet_ledger l
		JOIN wallets w ON w.id = l.wallet_id
		JOIN assets a ON a.id = w.asset_id
		WHERE w.user_id = $1
		ORDER BY l.id DESC
		LIMIT 100
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entries := make([]LedgerEntry, 0)
	for rows.Next() {
		var e LedgerEntry
		if err := rows.Scan(
			&e.ID, &e.Asset, &e.Amount,
			&e.EntryType, &e.ReferenceID, &e.CreatedAt,
		); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, rows.Err()
}
