package order

import (
	"errors"
	"time"

	"github.com/shopspring/decimal"
)

var (
	ErrInvalidOrder      = errors.New("invalid order")
	ErrInsufficientFunds = errors.New("insufficient available balance")
	ErrNotFound          = errors.New("order not found")
	ErrNotCancellable    = errors.New("order cannot be cancelled")
	ErrDuplicateOrder    = errors.New("client order ID already exists")
)

type CreateRequest struct {
	ClientOrderID string `json:"client_order_id"`
	Market        string `json:"market"`
	Side          string `json:"side"`
	Price         string `json:"price"`
	Quantity      string `json:"quantity"`
}

type Order struct {
	ID                int64     `json:"id"`
	ClientOrderID     string    `json:"client_order_id"`
	Market            string    `json:"market"`
	Side              string    `json:"side"`
	OrderType         string    `json:"order_type"`
	Price             string    `json:"price"`
	Quantity          string    `json:"quantity"`
	RemainingQuantity string    `json:"remaining_quantity"`
	ReservedAmount    string    `json:"reserved_amount"`
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"created_at"`
}

type ValidatedOrder struct {
	ClientOrderID string
	Market        string
	Side          string
	Price         decimal.Decimal
	Quantity      decimal.Decimal
}
