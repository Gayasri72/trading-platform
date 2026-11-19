
package orderbook

import (
	"errors"
	"sort"
	"sync"

	"github.com/shopspring/decimal"
)

type Side string

const (
	Buy  Side = "BUY"
	Sell Side = "SELL"
)

var (
	ErrInvalidOrder = errors.New("invalid order")
	ErrDuplicateID  = errors.New("duplicate order ID")
)

type Order struct {
	ID        int64
	Side      Side
	Price     decimal.Decimal
	Quantity  decimal.Decimal
	Sequence  uint64
}

type Book struct {
	mu           sync.RWMutex
	bids         []Order
	asks         []Order
	nextSequence uint64
}

func NewBook() *Book {
	return &Book{
		bids: make([]Order, 0),
		asks: make([]Order, 0),
	}
}

func (b *Book) Add(order Order) error {
	if order.ID <= 0 ||
		(order.Side != Buy && order.Side != Sell) ||
		!order.Price.IsPositive() ||
		!order.Quantity.IsPositive() {
		return ErrInvalidOrder
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	for _, existing := range b.bids {
		if existing.ID == order.ID {
			return ErrDuplicateID
		}
	}
	for _, existing := range b.asks {
		if existing.ID == order.ID {
			return ErrDuplicateID
		}
	}

	b.nextSequence++
	order.Sequence = b.nextSequence

	if order.Side == Buy {
		b.bids = append(b.bids, order)

		sort.Slice(b.bids, func(i, j int) bool {
			a, c := b.bids[i], b.bids[j]
			if a.Price.Equal(c.Price) {
				return a.Sequence < c.Sequence
			}
			return a.Price.GreaterThan(c.Price)
		})
	} else {
		b.asks = append(b.asks, order)

		sort.Slice(b.asks, func(i, j int) bool {
			a, c := b.asks[i], b.asks[j]
			if a.Price.Equal(c.Price) {
				return a.Sequence < c.Sequence
			}
			return a.Price.LessThan(c.Price)
		})
	}
	return nil
}

func (b *Book) Snapshot() (bids, asks []Order) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	bids = append([]Order(nil), b.bids...)
	asks = append([]Order(nil), b.asks...)
	return
}

func (b *Book) BestBid() (Order, bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if len(b.bids) == 0 {
		return Order{}, false
	}
	return b.bids[0], true
}

func (b *Book) BestAsk() (Order, bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if len(b.asks) == 0 {
		return Order{}, false
	}
	return b.asks[0], true
}
