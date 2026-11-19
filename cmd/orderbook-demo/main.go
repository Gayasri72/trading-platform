
package main

import (
	"fmt"

	"github.com/gayasri72/trading-platform/internal/orderbook"
	"github.com/shopspring/decimal"
)

func main() {
	book := orderbook.NewBook()

	orders := []orderbook.Order{
		{
			ID:       1,
			Side:     orderbook.Buy,
			Price:    decimal.RequireFromString("89900"),
			Quantity: decimal.RequireFromString("0.20"),
		},
		{
			ID:       2,
			Side:     orderbook.Buy,
			Price:    decimal.RequireFromString("90000"),
			Quantity: decimal.RequireFromString("0.10"),
		},
		{
			ID:       3,
			Side:     orderbook.Buy,
			Price:    decimal.RequireFromString("90000"),
			Quantity: decimal.RequireFromString("0.05"),
		},
		{
			ID:       4,
			Side:     orderbook.Sell,
			Price:    decimal.RequireFromString("90200"),
			Quantity: decimal.RequireFromString("0.30"),
		},
		{
			ID:       5,
			Side:     orderbook.Sell,
			Price:    decimal.RequireFromString("90100"),
			Quantity: decimal.RequireFromString("0.15"),
		},
	}

	for _, o := range orders {
		if err := book.Add(o); err != nil {
			panic(err)
		}
	}

	bids, asks := book.Snapshot()

	fmt.Println("BUY ORDERS (Highest first)")
	for _, o := range bids {
		fmt.Printf("ID=%d Price=%s Qty=%s Sequence=%d\n",
			o.ID, o.Price, o.Quantity, o.Sequence)
	}

	fmt.Println("\nSELL ORDERS (Lowest first)")
	for _, o := range asks {
		fmt.Printf("ID=%d Price=%s Qty=%s Sequence=%d\n",
			o.ID, o.Price, o.Quantity, o.Sequence)
	}
}
