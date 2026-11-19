package order

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

const demoUserID int64 = 1

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

const orderColumns = `
	o.id, o.client_order_id, tp.symbol, o.side,
	o.order_type, o.price::text, o.quantity::text,
	o.remaining_quantity::text, o.reserved_amount::text,
	o.status, o.created_at
`

func scanOrder(row pgx.Row) (Order, error) {
	var o Order
	err := row.Scan(
		&o.ID, &o.ClientOrderID, &o.Market, &o.Side,
		&o.OrderType, &o.Price, &o.Quantity,
		&o.RemainingQuantity, &o.ReservedAmount,
		&o.Status, &o.CreatedAt,
	)
	return o, err
}

func (repo *Repository) Create(
	ctx context.Context,
	input ValidatedOrder,
) (Order, error) {
	tx, err := repo.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return Order{}, err
	}
	defer tx.Rollback(ctx)

	var pairID int64
	var assetID int64
	var pricePrecision, quantityPrecision int32

	err = tx.QueryRow(ctx, `
		SELECT tp.id,
		       CASE WHEN $2 = 'BUY'
		            THEN tp.quote_asset_id
		            ELSE tp.base_asset_id END,
		       tp.price_precision,
		       tp.quantity_precision
		FROM trading_pairs tp
		JOIN assets base ON base.id = tp.base_asset_id
		JOIN assets quote ON quote.id = tp.quote_asset_id
		WHERE tp.symbol = $1
		  AND tp.is_active
		  AND base.is_active
		  AND quote.is_active
	`, input.Market, input.Side).Scan(
		&pairID, &assetID,
		&pricePrecision, &quantityPrecision,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return Order{}, ErrInvalidOrder
	}
	if err != nil {
		return Order{}, err
	}

	if input.Price.Exponent() < -pricePrecision ||
		input.Quantity.Exponent() < -quantityPrecision {
		return Order{}, ErrInvalidOrder
	}

	reserve := input.Quantity
	if input.Side == "BUY" {
		reserve = input.Price.Mul(input.Quantity)
	}

	// NUMERIC(30,10) allows at most 20 integer digits.
	maxAmount := decimal.RequireFromString("100000000000000000000")
	if !reserve.IsPositive() ||
		input.Price.GreaterThanOrEqual(maxAmount) ||
		input.Quantity.GreaterThanOrEqual(maxAmount) ||
		reserve.GreaterThanOrEqual(maxAmount) {
		return Order{}, ErrInvalidOrder
	}

	// Conditional UPDATE obtains a row lock and prevents
	// two orders from spending the same available funds.
	var walletID int64
	err = tx.QueryRow(ctx, `
		UPDATE wallets
		SET available = available - $1::numeric,
		    locked = locked + $1::numeric
		WHERE user_id = $2
		  AND asset_id = $3
		  AND available >= $1::numeric
		RETURNING id
	`, reserve.String(), demoUserID, assetID).Scan(&walletID)

	if errors.Is(err, pgx.ErrNoRows) {
		return Order{}, ErrInsufficientFunds
	}
	if err != nil {
		return Order{}, err
	}

	var orderID int64
	err = tx.QueryRow(ctx, `
		INSERT INTO orders (
			user_id, trading_pair_id, side,
			price, quantity, remaining_quantity,
			reserved_amount, client_order_id
		)
		VALUES (
			$1, $2, $3,
			$4::numeric, $5::numeric, $5::numeric,
			$6::numeric, $7
		)
		RETURNING id
	`,
		demoUserID, pairID, input.Side,
		input.Price.String(), input.Quantity.String(),
		reserve.String(), input.ClientOrderID,
	).Scan(&orderID)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return Order{}, ErrDuplicateOrder
		}
		return Order{}, err
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO wallet_movements (
			wallet_id, order_id, movement_type,
			available_delta, locked_delta
		)
		VALUES (
			$1, $2, 'RESERVE',
			-$3::numeric, $3::numeric
		)
	`, walletID, orderID, reserve.String())
	if err != nil {
		return Order{}, err
	}

	o, err := scanOrder(tx.QueryRow(ctx, `
		SELECT `+orderColumns+`
		FROM orders o
		JOIN trading_pairs tp ON tp.id = o.trading_pair_id
		WHERE o.id = $1
	`, orderID))
	if err != nil {
		return Order{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return Order{}, fmt.Errorf("commit order: %w", err)
	}
	return o, nil
}

func (repo *Repository) List(
	ctx context.Context,
) ([]Order, error) {
	rows, err := repo.db.Query(ctx, `
		SELECT `+orderColumns+`
		FROM orders o
		JOIN trading_pairs tp ON tp.id = o.trading_pair_id
		WHERE o.user_id = $1
		ORDER BY o.created_at DESC, o.id DESC
		LIMIT 100
	`, demoUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]Order, 0)
	for rows.Next() {
		o, err := scanOrder(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, o)
	}
	return result, rows.Err()
}

func (repo *Repository) Cancel(
	ctx context.Context,
	id int64,
) error {
	tx, err := repo.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var walletID int64
	var reservedText string
	var status string

	// Lock the order so cancellation and future matching
	// cannot update it simultaneously.
	err = tx.QueryRow(ctx, `
		SELECT w.id, o.reserved_amount::text, o.status
		FROM orders o
		JOIN trading_pairs tp ON tp.id = o.trading_pair_id
		JOIN wallets w
		  ON w.user_id = o.user_id
		 AND w.asset_id = CASE
		     WHEN o.side = 'BUY' THEN tp.quote_asset_id
		     ELSE tp.base_asset_id END
		WHERE o.id = $1 AND o.user_id = $2
		FOR UPDATE OF o
	`, id, demoUserID).Scan(
		&walletID, &reservedText, &status,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}
	if status != "OPEN" && status != "PARTIALLY_FILLED" {
		return ErrNotCancellable
	}

	// Lock and update the wallet using a guarded UPDATE.
	var updatedWalletID int64
	err = tx.QueryRow(ctx, `
		UPDATE wallets
		SET available = available + $1::numeric,
		    locked = locked - $1::numeric
		WHERE id = $2 AND locked >= $1::numeric
		RETURNING id
	`, reservedText, walletID).Scan(&updatedWalletID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		UPDATE orders
		SET status = 'CANCELLED',
		    reserved_amount = 0,
		    updated_at = NOW()
		WHERE id = $1
	`, id)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO wallet_movements (
			wallet_id, order_id, movement_type,
			available_delta, locked_delta
		)
		VALUES (
			$1, $2, 'RELEASE',
			$3::numeric, -$3::numeric
		)
	`, walletID, id, reservedText)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
