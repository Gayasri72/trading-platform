package market

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

func (repo *Repository) ListAssets(ctx context.Context) ([]Asset, error) {
	rows, err := repo.db.Query(ctx, `
		SELECT id, symbol, name, decimal_places, is_active
		FROM assets
		ORDER BY symbol
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	assets := make([]Asset, 0)

	for rows.Next() {
		var asset Asset

		err := rows.Scan(
			&asset.ID,
			&asset.Symbol,
			&asset.Name,
			&asset.DecimalPlaces,
			&asset.IsActive,
		)
		if err != nil {
			return nil, err
		}

		assets = append(assets, asset)
	}

	return assets, rows.Err()
}

func (repo *Repository) GetAsset(
	ctx context.Context,
	symbol string,
) (Asset, error) {
	var asset Asset

	err := repo.db.QueryRow(ctx, `
		SELECT id, symbol, name, decimal_places, is_active
		FROM assets
		WHERE symbol = $1
	`, symbol).Scan(
		&asset.ID,
		&asset.Symbol,
		&asset.Name,
		&asset.DecimalPlaces,
		&asset.IsActive,
	)

	return asset, err
}

const marketQuery = `
	SELECT
		tp.id,
		tp.symbol,
		base.symbol,
		quote.symbol,
		tp.price_precision,
		tp.quantity_precision,
		mp.reference_price::text,
		mp.source,
		CASE
			WHEN tp.is_active THEN 'active'
			ELSE 'inactive'
		END
	FROM trading_pairs tp
	JOIN assets base ON base.id = tp.base_asset_id
	JOIN assets quote ON quote.id = tp.quote_asset_id
	JOIN market_prices mp ON mp.trading_pair_id = tp.id
`

func (repo *Repository) ListMarkets(ctx context.Context) ([]Market, error) {
	rows, err := repo.db.Query(ctx, marketQuery+" ORDER BY tp.symbol")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	markets := make([]Market, 0)

	for rows.Next() {
		var item Market

		err := rows.Scan(
			&item.ID,
			&item.Symbol,
			&item.BaseAsset,
			&item.QuoteAsset,
			&item.PricePrecision,
			&item.QuantityPrecision,
			&item.ReferencePrice,
			&item.PriceSource,
			&item.Status,
		)
		if err != nil {
			return nil, err
		}

		markets = append(markets, item)
	}

	return markets, rows.Err()
}

func (repo *Repository) GetMarket(
	ctx context.Context,
	symbol string,
) (Market, error) {
	var item Market

	err := repo.db.QueryRow(ctx,
		marketQuery+" WHERE tp.symbol = $1",
		symbol,
	).Scan(
		&item.ID,
		&item.Symbol,
		&item.BaseAsset,
		&item.QuoteAsset,
		&item.PricePrecision,
		&item.QuantityPrecision,
		&item.ReferencePrice,
		&item.PriceSource,
		&item.Status,
	)

	return item, err
}
