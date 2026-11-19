package order

import (
	"context"
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(
	ctx context.Context,
	req CreateRequest,
) (Order, error) {

	req.Market = strings.ToUpper(strings.TrimSpace(req.Market))
	req.Side = strings.ToUpper(strings.TrimSpace(req.Side))
	req.ClientOrderID = strings.TrimSpace(req.ClientOrderID)

	if req.Market != "BTC-USDT" ||
		(req.Side != "BUY" && req.Side != "SELL") ||
		len(req.ClientOrderID) == 0 ||
		len(req.ClientOrderID) > 100 {
		return Order{}, ErrInvalidOrder
	}

	price, err := decimal.NewFromString(req.Price)
	if err != nil {
		return Order{}, ErrInvalidOrder
	}

	quantity, err := decimal.NewFromString(req.Quantity)
	if err != nil {
		return Order{}, ErrInvalidOrder
	}

	if !price.IsPositive() || !quantity.IsPositive() {
		return Order{}, ErrInvalidOrder
	}

	// BTC-USDT: 2 decimal places for price,
	// 8 decimal places for quantity.
	if price.Exponent() < -2 ||
		quantity.Exponent() < -8 {
		return Order{}, fmt.Errorf(
			"%w: price max 2 decimals, quantity max 8",
			ErrInvalidOrder,
		)
	}

	// Database NUMERIC(30,10) constraints.
	if price.Mul(quantity).GreaterThanOrEqual(
		decimal.RequireFromString("100000000000000000000"),
	) {
		return Order{}, ErrInvalidOrder
	}

	return s.repo.Create(ctx, ValidatedOrder{
		ClientOrderID: req.ClientOrderID,
		Market:        req.Market,
		Side:          req.Side,
		Price:         price,
		Quantity:      quantity,
	})
}

func (s *Service) Cancel(
	ctx context.Context,
	id int64,
) error {
	if id <= 0 {
		return ErrInvalidOrder
	}
	return s.repo.Cancel(ctx, id)
}
