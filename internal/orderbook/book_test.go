
package orderbook

import (
	"errors"
	"testing"

	"github.com/shopspring/decimal"
)

func makeOrder(id int64, side Side, price string) Order {
	return Order{
		ID:       id,
		Side:     side,
		Price:    decimal.RequireFromString(price),
		Quantity: decimal.RequireFromString("0.1"),
	}
}

func TestBuyPricePriority(t *testing.T) {
	book := NewBook()

	for _, o := range []Order{
		makeOrder(1, Buy, "89000"),
		makeOrder(2, Buy, "90000"),
		makeOrder(3, Buy, "89500"),
	} {
		if err := book.Add(o); err != nil {
			t.Fatal(err)
		}
	}

	best, ok := book.BestBid()
	if !ok || best.ID != 2 {
		t.Fatalf("expected best bid ID 2, got %+v", best)
	}
}

func TestSellPricePriority(t *testing.T) {
	book := NewBook()

	for _, o := range []Order{
		makeOrder(1, Sell, "90200"),
		makeOrder(2, Sell, "90100"),
		makeOrder(3, Sell, "90300"),
	} {
		if err := book.Add(o); err != nil {
			t.Fatal(err)
		}
	}

	best, ok := book.BestAsk()
	if !ok || best.ID != 2 {
		t.Fatalf("expected best ask ID 2, got %+v", best)
	}
}

func TestTimePriority(t *testing.T) {
	book := NewBook()

	book.Add(makeOrder(10, Buy, "90000"))
	book.Add(makeOrder(11, Buy, "90000"))

	bids, _ := book.Snapshot()

	if bids[0].ID != 10 || bids[1].ID != 11 {
		t.Fatal("time priority violated")
	}
}

func TestDuplicateOrder(t *testing.T) {
	book := NewBook()

	book.Add(makeOrder(1, Buy, "90000"))

	err := book.Add(makeOrder(1, Sell, "90100"))
	if !errors.Is(err, ErrDuplicateID) {
		t.Fatalf("expected duplicate ID error, got %v", err)
	}
}

func TestInvalidOrder(t *testing.T) {
	book := NewBook()

	err := book.Add(makeOrder(1, Buy, "0"))
	if !errors.Is(err, ErrInvalidOrder) {
		t.Fatalf("expected invalid order error, got %v", err)
	}
}
